package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"cloud.google.com/go/pubsub"
	telemetryAPI "github.com/capsule8/capsule8/api/v0"
	"github.com/capsule8/capsule8/pkg/sensor"
	"github.com/golang/glog"
	"github.com/golang/protobuf/jsonpb"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
)

var (
	subscription telemetryAPI.Subscription
	nEvents      int
)

var flags struct {
	pubsubPath       string
	credentialsFile  string
	subscriptionFile string
	dumpConfig       bool
	useEmulator      bool
}

func main() {
	var err error
	ctx, cancel := context.WithCancel(context.Background())

	flag.StringVar(&flags.pubsubPath, "pubsub", "",
		"Cloud Pubsub path (e.g. projects/PROJECT_ID/topics/TOPIC_ID)")

	// Config file can contain pubsub topic, creds, and subscription spec
	flag.StringVar(&flags.credentialsFile, "credentials", "",
		"path to Google Cloud credentials file")

	flag.StringVar(&flags.subscriptionFile, "subscription", "",
		"path to JSON subscription file")

	flag.BoolVar(&flags.dumpConfig, "p", false,
		"print subscription as JSON")

	flag.BoolVar(&flags.useEmulator, "emulator", false,
		"use local Pub/Sub emulator")

	// Configure glog
	flag.Set("logtostderr", "true")
	flag.Parse()

	if len(flags.subscriptionFile) > 0 {
		f, err := os.Open(flags.subscriptionFile)
		if err != nil {
			glog.Fatal("couldn't open subscription JSON file: ", err)
		}

		err = jsonpb.Unmarshal(f, &subscription)
		if err != nil {
			glog.Fatal("couldn't parse subscription JSON: ", err)
		}
	} else {
		subscription = createSubscription()
	}

	if flags.dumpConfig {
		m := jsonpb.Marshaler{}
		jsonString, err := m.MarshalToString(&subscription)
		if err != nil {
			glog.Fatal(err)
		}

		fmt.Println(string(jsonString))
		os.Exit(0)
	}

	// Try to connect to PubSub
	if len(flags.pubsubPath) == 0 {
		flag.Usage()
		os.Exit(0)
	}

	topic, err := connectToPubsubTopic(ctx, flags.pubsubPath)
	if err != nil {
		glog.Fatal(err)
	}

	//
	// Create sensor based on given or default subscription
	//

	s, err := sensor.NewSensor()
	if err != nil {
		glog.Fatal(err)
	}

	err = s.Start()
	if err != nil {
		glog.Fatal(err)
	}

	sub := s.NewSubscription()
	sub.ProcessTelemetryServiceSubscription(&subscription)

	errors, err := sub.Run(ctx, func(event sensor.TelemetryEvent) {
		nEvents++

		//
		// The output formatting could use some work
		//
		var e struct {
			EventType      string
			TelemetryEvent *sensor.TelemetryEvent
		}

		e.EventType = fmt.Sprintf("%T", event)
		e.TelemetryEvent = &event

		// Need an event type name
		jsonString, err := json.Marshal(&e)
		if err != nil {
			glog.Warning(err)
			return
		}

		m := pubsub.Message{Data: jsonString}
		topic.Publish(ctx, &m)
	})

	if err != nil {
		glog.Fatal(err)
	}

	if len(errors) > 0 {
		glog.Fatal(errors)
	}

	// Trap Control-C
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	glog.Info("Received interrupt signal, exiting...")
	cancel()
	s.Stop()

	glog.Infof("Received %d events", nEvents)
	os.Exit(1)
}

func connectToPubsubTopic(ctx context.Context, path string) (*pubsub.Topic, error) {
	parts := strings.Split(path, "/")
	if parts[0] != "projects" || parts[2] != "topics" {
		return nil, fmt.Errorf("could not parse pubsub path %s", path)
	}

	project := parts[1]
	topic := parts[3]

	var options []option.ClientOption
	if flags.useEmulator {
		options = append(options, option.WithoutAuthentication())
		options = append(options, option.WithGRPCDialOption(grpc.WithInsecure()))
	} else if len(flags.credentialsFile) > 0 {
		options = append(options, option.WithCredentialsFile(flags.credentialsFile))
	}

	c, err := pubsub.NewClient(ctx, project, options...)
	if err != nil {
		return nil, err
	}

	t := c.Topic(topic)
	ok, err := t.Exists(ctx)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, nil
	}

	return t, nil
}

func createSubscription() telemetryAPI.Subscription {
	processEvents := []*telemetryAPI.ProcessEventFilter{
		&telemetryAPI.ProcessEventFilter{
			Type: telemetryAPI.ProcessEventType_PROCESS_EVENT_TYPE_FORK,
		},
		&telemetryAPI.ProcessEventFilter{
			Type: telemetryAPI.ProcessEventType_PROCESS_EVENT_TYPE_EXEC,
		},
		&telemetryAPI.ProcessEventFilter{
			Type: telemetryAPI.ProcessEventType_PROCESS_EVENT_TYPE_EXIT,
		},
	}

	fileEvents := []*telemetryAPI.FileEventFilter{
		&telemetryAPI.FileEventFilter{
			Type: telemetryAPI.FileEventType_FILE_EVENT_TYPE_OPEN,
		},
	}

	networkEvents := []*telemetryAPI.NetworkEventFilter{
		&telemetryAPI.NetworkEventFilter{
			Type: telemetryAPI.NetworkEventType_NETWORK_EVENT_TYPE_ACCEPT_ATTEMPT,
		},

		&telemetryAPI.NetworkEventFilter{
			Type: telemetryAPI.NetworkEventType_NETWORK_EVENT_TYPE_ACCEPT_RESULT,
		},

		&telemetryAPI.NetworkEventFilter{
			Type: telemetryAPI.NetworkEventType_NETWORK_EVENT_TYPE_CONNECT_ATTEMPT,
		},

		&telemetryAPI.NetworkEventFilter{
			Type: telemetryAPI.NetworkEventType_NETWORK_EVENT_TYPE_CONNECT_RESULT,
		},

		&telemetryAPI.NetworkEventFilter{
			Type: telemetryAPI.NetworkEventType_NETWORK_EVENT_TYPE_BIND_ATTEMPT,
		},

		&telemetryAPI.NetworkEventFilter{
			Type: telemetryAPI.NetworkEventType_NETWORK_EVENT_TYPE_BIND_RESULT,
		},
	}

	containerEvents := []*telemetryAPI.ContainerEventFilter{
		&telemetryAPI.ContainerEventFilter{
			Type: telemetryAPI.ContainerEventType_CONTAINER_EVENT_TYPE_CREATED,
		},
		&telemetryAPI.ContainerEventFilter{
			Type: telemetryAPI.ContainerEventType_CONTAINER_EVENT_TYPE_RUNNING,
		},
		&telemetryAPI.ContainerEventFilter{
			Type: telemetryAPI.ContainerEventType_CONTAINER_EVENT_TYPE_EXITED,
		},
		&telemetryAPI.ContainerEventFilter{
			Type: telemetryAPI.ContainerEventType_CONTAINER_EVENT_TYPE_DESTROYED,
		},
	}

	eventFilter := &telemetryAPI.EventFilter{
		ProcessEvents:   processEvents,
		FileEvents:      fileEvents,
		NetworkEvents:   networkEvents,
		ContainerEvents: containerEvents,
	}

	sub := telemetryAPI.Subscription{
		EventFilter: eventFilter,
	}

	return sub
}
