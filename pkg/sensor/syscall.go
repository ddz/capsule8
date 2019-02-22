// Copyright 2017 Capsule8, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sensor

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/capsule8/capsule8/pkg/expression"
	"github.com/capsule8/capsule8/pkg/sys"
	"github.com/capsule8/capsule8/pkg/sys/perf"

	"github.com/golang/glog"
)

// SyscallEnterEventTypes defines the field types that can be used with filters
// on syscall enter telemetry events.
var SyscallEnterEventTypes = expression.FieldTypeMap{
	"id":   expression.ValueTypeSignedInt64,
	"arg0": expression.ValueTypeUnsignedInt64,
	"arg1": expression.ValueTypeUnsignedInt64,
	"arg2": expression.ValueTypeUnsignedInt64,
	"arg3": expression.ValueTypeUnsignedInt64,
	"arg4": expression.ValueTypeUnsignedInt64,
	"arg5": expression.ValueTypeUnsignedInt64,
}

// SyscallExitEventTypes defines the field types that can be used with filters
// on syscall exit telemetry events.
var SyscallExitEventTypes = expression.FieldTypeMap{
	"id":  expression.ValueTypeSignedInt64,
	"ret": expression.ValueTypeSignedInt64,
}

func handleDummySysEnter(_ uint64, _ *perf.Sample) {}

// SyscallEnterTelemetryEvent is a telemetry event generated by the syscall
// enter event source.
type SyscallEnterTelemetryEvent struct {
	TelemetryEventData

	ID        int64
	Arguments [6]uint64
}

// CommonTelemetryEventData returns the telemtry event data common to all
// telemetry events for a chargen telemetry event.
func (e SyscallEnterTelemetryEvent) CommonTelemetryEventData() TelemetryEventData {
	return e.TelemetryEventData
}

// SyscallExitTelemetryEvent is a telemetry event generated by the syscall
// exit event source.
type SyscallExitTelemetryEvent struct {
	TelemetryEventData

	ID     int64
	Return int64
}

// CommonTelemetryEventData returns the telemtry event data common to all
// telemetry events for a chargen telemetry event.
func (e SyscallExitTelemetryEvent) CommonTelemetryEventData() TelemetryEventData {
	return e.TelemetryEventData
}

var argKeys = []string{"arg0", "arg1", "arg2", "arg3", "arg4", "arg5"}

func (s *Subscription) handleSyscallTraceEnter(eventid uint64, sample *perf.Sample) {
	var e SyscallEnterTelemetryEvent
	if e.InitWithSample(s.sensor, sample) {
		e.ID, _ = sample.GetSignedInt64("id")
		for i := 0; i < 6; i++ {
			e.Arguments[i], _ = sample.GetUnsignedInt64(argKeys[i])
		}
		s.DispatchEvent(eventid, e, nil)
	}
}

func (s *Subscription) handleSysExit(eventid uint64, sample *perf.Sample) {
	var e SyscallExitTelemetryEvent
	if e.InitWithSample(s.sensor, sample) {
		e.ID, _ = sample.GetSignedInt64("id")
		e.Return, _ = sample.GetSignedInt64("ret")
		s.DispatchEvent(eventid, e, nil)
	}
}

const (
	syscallNewEnterKprobeAddress string = "syscall_trace_enter_phase1"
	syscallOldEnterKprobeAddress string = "syscall_trace_enter"

	// These offsets index into the x86_64 version of struct pt_regs
	// in the kernel. This is a stable structure.
	syscallEnterKprobeFetchargs string = "id=+120(%di):s64 " + // orig_ax
		"arg0=+112(%di):u64 " + // di
		"arg1=+104(%di):u64 " + // si
		"arg2=+96(%di):u64 " + // dx
		"arg3=+56(%di):u64 " + // r10
		"arg4=+72(%di):u64 " + // r8
		"arg5=+64(%di):u64" // r9
)

var (
	syscallOnce      sync.Once
	syscallEnterName string
	syscallExitName  string
)

func (s *Subscription) initSyscallNames() {
	// Different sensor instances could have different tracing dirs, but
	// ultimately they're all still coming from the same running kernel,
	// so the tracepoint names are not going to change. So while this might
	// appear to be a bit shady, it's perfectly safe.
	syscallEnterName = "raw_syscalls/sys_enter"
	if s.sensor.Monitor().DoesTracepointExist(syscallEnterName) {
		syscallExitName = "raw_syscalls/sys_exit"
	} else {
		syscallEnterName = "syscalls/sys_enter"
		if !s.sensor.Monitor().DoesTracepointExist(syscallEnterName) {
			glog.Fatal("No syscall sys_enter tracepoint exists!")
		}
		syscallExitName = "syscalls/sys_exit"
	}
}

func (s *Subscription) registerGlobalDummySyscallEvent() bool {
	if atomic.AddInt64(&s.sensor.dummySyscallEventCount, 1) == 1 {
		// Use the container cache event group ID here. No events will
		// ever be generated, but we need something that won't ever be
		// removed. It doesn't make sense to create a whole new group
		// that'll never hold anything.
		eventID, err := s.sensor.Monitor().RegisterTracepoint(
			syscallEnterName, handleDummySysEnter,
			s.sensor.ContainerCache.EventGroupID,
			perf.WithFilter("id == 0x7fffffff"))
		if err != nil {
			s.logStatus(
				fmt.Sprintf("Could not register dummy syscall event %s: %v",
					syscallEnterName, err))
			atomic.AddInt64(&s.sensor.dummySyscallEventCount, -1)
			return false
		}
		s.sensor.dummySyscallEventID = eventID
		return true
	}
	return true
}

func (s *Subscription) registerLocalDummySyscallEvent() bool {
	var err error
	if err = s.createEventGroup(); err != nil {
		s.logStatus(
			fmt.Sprintf("Could not create subscription event group: %v", err))
		return false
	}
	_, err = s.sensor.Monitor().RegisterTracepoint(
		syscallEnterName, handleDummySysEnter,
		s.eventGroupID,
		perf.WithFilter("id == 0x7fffffff"))
	if err != nil {
		s.logStatus(
			fmt.Sprintf("Could not register dummy syscall event %s: %v",
				syscallEnterName, err))
		return false
	}
	return true
}

// RegisterSyscallEnterEventFilter registers a syscall enter event filter with
// a subscription.
func (s *Subscription) RegisterSyscallEnterEventFilter(
	filter *expression.Expression,
) {
	syscallOnce.Do(s.initSyscallNames)

	// Create the dummy syscall event. This event is needed to put
	// the kernel into a mode where it'll make the function calls
	// needed to make the kprobe we'll add fire. Add the tracepoint,
	// but make sure it never adds events into the ringbuffer by
	// using a filter that will never evaluate true. It also never
	// gets enabled, but just creating it is enough.
	//
	// For kernels older than 3.x, create this dummy event in all
	// event groups, because we cannot remove it when we don't need
	// it anymore due to bugs in CentOS 6.x kernels (2.6.32).
	var (
		result     bool
		unregister func(*eventSink)
	)
	if major, _, _ := sys.KernelVersion(); major < 3 {
		result = s.registerLocalDummySyscallEvent()
	} else {
		result = s.registerGlobalDummySyscallEvent()
		unregister = func(es *eventSink) {
			eventID := s.sensor.dummySyscallEventID
			if atomic.AddInt64(&s.sensor.dummySyscallEventCount, -1) == 0 {
				s.sensor.Monitor().UnregisterEvent(eventID)
			}
		}
	}
	if result {
		// There are two possible kprobes. Newer kernels (>= 4.1) have
		// refactored syscall entry code, so syscall_trace_enter_phase1
		// is the right one, but for older kernels syscall_trace_enter
		// is the right one. Both have the same signature, so the
		// fetchargs doesn't have to change. Try the new probe first,
		// because the old probe will also set in the newer kernels,
		// but it won't fire.
		kprobeSymbol := syscallNewEnterKprobeAddress
		if !s.sensor.IsKernelSymbolAvailable(kprobeSymbol) {
			kprobeSymbol = syscallOldEnterKprobeAddress
		}

		es, err := s.registerKprobe(kprobeSymbol, false,
			syscallEnterKprobeFetchargs, s.handleSyscallTraceEnter,
			filter, false)
		if err != nil {
			if unregister != nil {
				unregister(nil)
			}
		} else {
			es.unregister = unregister
		}
	}
}

// RegisterSyscallExitEventFilter registers a syscall exit event filter with a
// subscription.
func (s *Subscription) RegisterSyscallExitEventFilter(
	filter *expression.Expression,
) {
	syscallOnce.Do(s.initSyscallNames)

	s.registerTracepoint(syscallExitName, s.handleSysExit, filter)
}
