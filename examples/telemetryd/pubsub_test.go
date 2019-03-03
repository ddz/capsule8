package main

import (
	"context"
	"testing"

	"cloud.google.com/go/pubsub"
)

func getOrCreateTopic(ctx context.Context, c *pubsub.Client, id string) (*pubsub.Topic, error) {
	t := c.Topic(id)
	ok, err := t.Exists(ctx)
	if err != nil {
		return nil, err
	}

	if ok {
		return t, nil
	}

	return c.CreateTopic(ctx, id)
}

func getOrCreateSub(ctx context.Context, c *pubsub.Client, id string, t *pubsub.Topic) (*pubsub.Subscription, error) {
	s := c.Subscription(id)
	ok, err := s.Exists(ctx)
	if err != nil {
		return nil, err
	}

	if ok {
		return s, nil
	}

	return c.CreateSubscription(ctx, id, pubsub.SubscriptionConfig{Topic: t})
}

func TestPubSub(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// publish to fully-qualified pubsub path: projects/PROJECT_ID/topics/TOPIC_ID

	// Auth via:
	// WithAPIKey(apiKey)
	// WithCredentialsFile(filename)

	pubsubClient, err := pubsub.NewClient(ctx, "project-id")
	if err != nil {
		t.Fatal(err)
	}

	// Create a test topic
	topic, err := getOrCreateTopic(ctx, pubsubClient, "topic-name")
	if err != nil {
		t.Fatal(err)
	}

	// can adjust topic.PublishSettings here

	sub, err := getOrCreateSub(ctx, pubsubClient, "sub-name", topic)
	if err != nil {
		t.Fatal(err)
	}

	// Make this random
	m := pubsub.Message{Data: []byte("foo")}
	res := topic.Publish(ctx, &m)

	defer topic.Stop()

	serverID, err := res.Get(ctx)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Sent message as serverID: ", serverID)

	err = sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
		t.Log("Received: ", m)
		m.Ack()

		// We got one, cancel the context to stop blocking on sub.Receive()
		cancel()
	})

	if err != nil {
		t.Fatal(err)
	}
}
