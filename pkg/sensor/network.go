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
	"github.com/capsule8/capsule8/pkg/expression"
	"github.com/capsule8/capsule8/pkg/sys/perf"

	"golang.org/x/sys/unix"
)

// NetworkAttemptEventTypes defines the field types that can be used with
// filters on network attempt telemetry events that do not have more specific
// fields of their own.
var NetworkAttemptEventTypes = expression.FieldTypeMap{
	"fd": expression.ValueTypeUnsignedInt64,
}

// NetworkAttemptWithAddressEventTypes defines the field types that can be
// used with filters on network attempt telemetry events that include address
// information but do not have more specific fields of their own.
var NetworkAttemptWithAddressEventTypes = expression.FieldTypeMap{
	"fd":             expression.ValueTypeUnsignedInt64,
	"sa_family":      expression.ValueTypeUnsignedInt16,
	"sin_port":       expression.ValueTypeUnsignedInt16,
	"sin_addr":       expression.ValueTypeUnsignedInt32,
	"sun_path":       expression.ValueTypeString,
	"sin6_port":      expression.ValueTypeUnsignedInt16,
	"sin6_addr_high": expression.ValueTypeUnsignedInt64,
	"sin6_addr_low":  expression.ValueTypeUnsignedInt64,
}

// NetworkListenAttemptEventTypes defines the field types that can be used with
// filters on network listen attempt telemetry events.
var NetworkListenAttemptEventTypes = expression.FieldTypeMap{
	"fd":      expression.ValueTypeUnsignedInt64,
	"backlog": expression.ValueTypeUnsignedInt64,
}

// NetworkResultEventTypes defines the field types that can be used with
// filters on network result telemetry events.
var NetworkResultEventTypes = expression.FieldTypeMap{
	"ret": expression.ValueTypeSignedInt64,
}

// NetworkAttemptTelemetryEventData is the data common to all network attempt
// telemetry events.
type NetworkAttemptTelemetryEventData struct {
	FD uint64
}

func (ted *NetworkAttemptTelemetryEventData) initWithSample(sample *perf.Sample) {
	ted.FD, _ = sample.GetUnsignedInt64("fd")
}

// NetworkAddressTelemetryEventData is the data common to all network attempt
// telemetry events that have addresses.
type NetworkAddressTelemetryEventData struct {
	Family          uint16
	UnixPath        string
	IPv4Address     uint32
	IPv4Port        uint16
	IPv6AddressHigh uint64
	IPv6AddressLow  uint64
	IPv6Port        uint16
}

func (ted *NetworkAddressTelemetryEventData) initWithSample(sample *perf.Sample) {
	ted.Family, _ = sample.GetUnsignedInt16("sa_family")
	switch ted.Family {
	case unix.AF_LOCAL:
		ted.UnixPath, _ = sample.GetString("sun_path")
	case unix.AF_INET:
		ted.IPv4Address, _ = sample.GetUnsignedInt32("sin_addr")
		ted.IPv4Port, _ = sample.GetUnsignedInt16("sin_port")
	case unix.AF_INET6:
		ted.IPv6AddressHigh, _ = sample.GetUnsignedInt64("sin6_addr_high")
		ted.IPv6AddressLow, _ = sample.GetUnsignedInt64("sin6_addr_low")
		ted.IPv6Port, _ = sample.GetUnsignedInt16("sin6_port")
	}
}

// NetworkResultTelemetryEventData is the data common to all network result
// telemetry events.
type NetworkResultTelemetryEventData struct {
	Return int64
}

func (ted *NetworkResultTelemetryEventData) initWithSample(sample *perf.Sample) {
	ted.Return, _ = sample.GetSignedInt64("ret")
}

// NetworkAcceptAttemptTelemetryEvent is a telemetry event generated by the
// network accept attempt event source.
type NetworkAcceptAttemptTelemetryEvent struct {
	TelemetryEventData
	NetworkAttemptTelemetryEventData
}

// CommonTelemetryEventData returns the telemtry event data common to all
// telemetry events for a chargen telemetry event.
func (e NetworkAcceptAttemptTelemetryEvent) CommonTelemetryEventData() TelemetryEventData {
	return e.TelemetryEventData
}

// NetworkAcceptResultTelemetryEvent is a telemetry event generated by the
// network accept result event source.
type NetworkAcceptResultTelemetryEvent struct {
	TelemetryEventData
	NetworkResultTelemetryEventData
}

// CommonTelemetryEventData returns the telemtry event data common to all
// telemetry events for a chargen telemetry event.
func (e NetworkAcceptResultTelemetryEvent) CommonTelemetryEventData() TelemetryEventData {
	return e.TelemetryEventData
}

// NetworkBindAttemptTelemetryEvent is a telemetry event generated by the
// network bind attempt event source.
type NetworkBindAttemptTelemetryEvent struct {
	TelemetryEventData
	NetworkAttemptTelemetryEventData
	NetworkAddressTelemetryEventData
}

// CommonTelemetryEventData returns the telemtry event data common to all
// telemetry events for a chargen telemetry event.
func (e NetworkBindAttemptTelemetryEvent) CommonTelemetryEventData() TelemetryEventData {
	return e.TelemetryEventData
}

// NetworkBindResultTelemetryEvent is a telemetry event generated by the
// network bind result event source.
type NetworkBindResultTelemetryEvent struct {
	TelemetryEventData
	NetworkResultTelemetryEventData
}

// CommonTelemetryEventData returns the telemtry event data common to all
// telemetry events for a chargen telemetry event.
func (e NetworkBindResultTelemetryEvent) CommonTelemetryEventData() TelemetryEventData {
	return e.TelemetryEventData
}

// NetworkConnectAttemptTelemetryEvent is a telemetry event generated by the
// network connect attempt event source.
type NetworkConnectAttemptTelemetryEvent struct {
	TelemetryEventData
	NetworkAttemptTelemetryEventData
	NetworkAddressTelemetryEventData
}

// CommonTelemetryEventData returns the telemtry event data common to all
// telemetry events for a chargen telemetry event.
func (e NetworkConnectAttemptTelemetryEvent) CommonTelemetryEventData() TelemetryEventData {
	return e.TelemetryEventData
}

// NetworkConnectResultTelemetryEvent is a telemetry event generated by the
// network connect result event source.
type NetworkConnectResultTelemetryEvent struct {
	TelemetryEventData
	NetworkResultTelemetryEventData
}

// CommonTelemetryEventData returns the telemtry event data common to all
// telemetry events for a chargen telemetry event.
func (e NetworkConnectResultTelemetryEvent) CommonTelemetryEventData() TelemetryEventData {
	return e.TelemetryEventData
}

// NetworkListenAttemptTelemetryEvent is a telemetry event generated by the
// network listen attempt event source.
type NetworkListenAttemptTelemetryEvent struct {
	TelemetryEventData
	NetworkAttemptTelemetryEventData

	Backlog uint64
}

// CommonTelemetryEventData returns the telemtry event data common to all
// telemetry events for a chargen telemetry event.
func (e NetworkListenAttemptTelemetryEvent) CommonTelemetryEventData() TelemetryEventData {
	return e.TelemetryEventData
}

// NetworkListenResultTelemetryEvent is a telemetry event generated by the
// network listen result event source.
type NetworkListenResultTelemetryEvent struct {
	TelemetryEventData
	NetworkResultTelemetryEventData
}

// CommonTelemetryEventData returns the telemtry event data common to all
// telemetry events for a chargen telemetry event.
func (e NetworkListenResultTelemetryEvent) CommonTelemetryEventData() TelemetryEventData {
	return e.TelemetryEventData
}

// NetworkRecvfromAttemptTelemetryEvent is a telemetry event generated by the
// network recvfrom attempt event source.
type NetworkRecvfromAttemptTelemetryEvent struct {
	TelemetryEventData
	NetworkAttemptTelemetryEventData
}

// CommonTelemetryEventData returns the telemtry event data common to all
// telemetry events for a chargen telemetry event.
func (e NetworkRecvfromAttemptTelemetryEvent) CommonTelemetryEventData() TelemetryEventData {
	return e.TelemetryEventData
}

// NetworkRecvfromResultTelemetryEvent is a telemetry event generated by the
// network recvfrom result event source.
type NetworkRecvfromResultTelemetryEvent struct {
	TelemetryEventData
	NetworkResultTelemetryEventData
}

// CommonTelemetryEventData returns the telemtry event data common to all
// telemetry events for a chargen telemetry event.
func (e NetworkRecvfromResultTelemetryEvent) CommonTelemetryEventData() TelemetryEventData {
	return e.TelemetryEventData
}

// NetworkSendtoAttemptTelemetryEvent is a telemetry event generated by the
// network sendto attempt event source.
type NetworkSendtoAttemptTelemetryEvent struct {
	TelemetryEventData
	NetworkAttemptTelemetryEventData
	NetworkAddressTelemetryEventData
}

// CommonTelemetryEventData returns the telemtry event data common to all
// telemetry events for a chargen telemetry event.
func (e NetworkSendtoAttemptTelemetryEvent) CommonTelemetryEventData() TelemetryEventData {
	return e.TelemetryEventData
}

// NetworkSendtoResultTelemetryEvent is a telemetry event generated by the
// network sendto result event source.
type NetworkSendtoResultTelemetryEvent struct {
	TelemetryEventData
	NetworkResultTelemetryEventData
}

// CommonTelemetryEventData returns the telemtry event data common to all
// telemetry events for a chargen telemetry event.
func (e NetworkSendtoResultTelemetryEvent) CommonTelemetryEventData() TelemetryEventData {
	return e.TelemetryEventData
}

const (
	networkKprobeBindSymbol    = "sys_bind"
	networkKprobeBindFetchargs = "fd=%di sa_family=+0(%si):u16 " +
		"sin_port=+2(%si):u16 sin_addr=+4(%si):u32 " +
		"sun_path=+2(%si):string " +
		"sin6_port=+2(%si):u16 sin6_addr_high=+8(%si):u64 sin6_addr_low=+16(%si):u64"

	networkKprobeConnectSymbol    = "sys_connect"
	networkKprobeConnectFetchargs = "fd=%di sa_family=+0(%si):u16 " +
		"sin_port=+2(%si):u16 sin_addr=+4(%si):u32 " +
		"sun_path=+2(%si):string " +
		"sin6_port=+2(%si):u16 sin6_addr_high=+8(%si):u64 sin6_addr_low=+16(%si):u64"

	networkKprobeSendmsgSymbol    = "sys_sendmsg"
	networkKprobeSendmsgFetchargs = "fd=%di sa_family=+0(+0(%si)):u16 " +
		"sin_port=+2(+0(%si)):u16 sin_addr=+4(+0(%si)):u32 " +
		"sun_path=+2(+0(%si)):string " +
		"sin6_port=+2(+0(%si)):u16 sin6_addr_high=+8(+0(%si)):u64 sin6_addr_low=+16(+0(%si)):u64"

	networkKprobeSendtoSymbol    = "sys_sendto"
	networkKprobeSendtoFetchargs = "fd=%di sa_family=+0(%r8):u16 " +
		"sin_port=+2(%r8):u16 sin_addr=+4(%r8):u32 " +
		"sun_path=+2(%r8):string " +
		"sin6_port=+2(%r8):u16 sin6_addr_high=+8(%r8):u64 sin6_addr_low=+16(%r8):u64"
)

func (s *Subscription) handleSysEnterAccept(eventid uint64, sample *perf.Sample) {
	var e NetworkAcceptAttemptTelemetryEvent
	if e.InitWithSample(s.sensor, sample) {
		e.NetworkAttemptTelemetryEventData.initWithSample(sample)
		s.DispatchEvent(eventid, e, nil)
	}
}

func (s *Subscription) handleSysExitAccept(eventid uint64, sample *perf.Sample) {
	var e NetworkAcceptResultTelemetryEvent
	if e.InitWithSample(s.sensor, sample) {
		e.NetworkResultTelemetryEventData.initWithSample(sample)
		s.DispatchEvent(eventid, e, nil)
	}
}

func (s *Subscription) handleSysBind(eventid uint64, sample *perf.Sample) {
	var e NetworkBindAttemptTelemetryEvent
	if e.InitWithSample(s.sensor, sample) {
		e.NetworkAttemptTelemetryEventData.initWithSample(sample)
		e.NetworkAddressTelemetryEventData.initWithSample(sample)
		s.DispatchEvent(eventid, e, nil)
	}
}

func (s *Subscription) handleSysExitBind(eventid uint64, sample *perf.Sample) {
	var e NetworkBindResultTelemetryEvent
	if e.InitWithSample(s.sensor, sample) {
		e.NetworkResultTelemetryEventData.initWithSample(sample)
		s.DispatchEvent(eventid, e, nil)
	}
}

func (s *Subscription) handleSysConnect(eventid uint64, sample *perf.Sample) {
	var e NetworkConnectAttemptTelemetryEvent
	if e.InitWithSample(s.sensor, sample) {
		e.NetworkAttemptTelemetryEventData.initWithSample(sample)
		e.NetworkAddressTelemetryEventData.initWithSample(sample)
		s.DispatchEvent(eventid, e, nil)
	}
}

func (s *Subscription) handleSysExitConnect(eventid uint64, sample *perf.Sample) {
	var e NetworkConnectResultTelemetryEvent
	if e.InitWithSample(s.sensor, sample) {
		e.NetworkResultTelemetryEventData.initWithSample(sample)
		s.DispatchEvent(eventid, e, nil)
	}
}

func (s *Subscription) handleSysEnterListen(eventid uint64, sample *perf.Sample) {
	var e NetworkListenAttemptTelemetryEvent
	if e.InitWithSample(s.sensor, sample) {
		e.NetworkAttemptTelemetryEventData.initWithSample(sample)
		e.Backlog, _ = sample.GetUnsignedInt64("backlog")
		s.DispatchEvent(eventid, e, nil)
	}
}

func (s *Subscription) handleSysExitListen(eventid uint64, sample *perf.Sample) {
	var e NetworkListenResultTelemetryEvent
	if e.InitWithSample(s.sensor, sample) {
		e.NetworkResultTelemetryEventData.initWithSample(sample)
		s.DispatchEvent(eventid, e, nil)
	}
}

func (s *Subscription) handleSysEnterRecvfrom(eventid uint64, sample *perf.Sample) {
	var e NetworkRecvfromAttemptTelemetryEvent
	if e.InitWithSample(s.sensor, sample) {
		e.NetworkAttemptTelemetryEventData.initWithSample(sample)
		s.DispatchEvent(eventid, e, nil)
	}
}

func (s *Subscription) handleSysExitRecvfrom(eventid uint64, sample *perf.Sample) {
	var e NetworkRecvfromResultTelemetryEvent
	if e.InitWithSample(s.sensor, sample) {
		e.NetworkResultTelemetryEventData.initWithSample(sample)
		s.DispatchEvent(eventid, e, nil)
	}
}

func (s *Subscription) handleSysSendto(eventid uint64, sample *perf.Sample) {
	var e NetworkSendtoAttemptTelemetryEvent
	if e.InitWithSample(s.sensor, sample) {
		e.NetworkAttemptTelemetryEventData.initWithSample(sample)
		e.NetworkAddressTelemetryEventData.initWithSample(sample)
		s.DispatchEvent(eventid, e, nil)
	}
}

func (s *Subscription) handleSysExitSendto(eventid uint64, sample *perf.Sample) {
	var e NetworkSendtoResultTelemetryEvent
	if e.InitWithSample(s.sensor, sample) {
		e.NetworkResultTelemetryEventData.initWithSample(sample)
		s.DispatchEvent(eventid, e, nil)
	}
}

// RegisterNetworkAcceptAttemptEventFilter registers a network accept attempt
// event filter with a subscription.
func (s *Subscription) RegisterNetworkAcceptAttemptEventFilter(expr *expression.Expression) {
	s.registerTracepoint("syscalls/sys_enter_accept",
		s.handleSysEnterAccept, expr)
	s.registerTracepoint("syscalls/sys_enter_accept4",
		s.handleSysEnterAccept, expr)
}

// RegisterNetworkAcceptResultEventFilter registers a network accept result
// event filter with a subscription.
func (s *Subscription) RegisterNetworkAcceptResultEventFilter(expr *expression.Expression) {
	s.registerTracepoint("syscalls/sys_exit_accept",
		s.handleSysExitAccept, expr)
	s.registerTracepoint("syscalls/sys_exit_accept4",
		s.handleSysExitAccept, expr)
}

// RegisterNetworkBindAttemptEventFilter registers a network bind attempt event
// filter with a subscription.
func (s *Subscription) RegisterNetworkBindAttemptEventFilter(expr *expression.Expression) {
	s.registerKprobe(networkKprobeBindSymbol, false,
		networkKprobeBindFetchargs, s.handleSysBind, expr, false)
}

// RegisterNetworkBindResultEventFilter registers a network bind result event
// filter with a subscription.
func (s *Subscription) RegisterNetworkBindResultEventFilter(expr *expression.Expression) {
	s.registerTracepoint("syscalls/sys_exit_bind",
		s.handleSysExitBind, expr)
}

// RegisterNetworkConnectAttemptEventFilter registers a network connect attempt
// event filter with a subscription.
func (s *Subscription) RegisterNetworkConnectAttemptEventFilter(expr *expression.Expression) {
	s.registerKprobe(networkKprobeConnectSymbol, false,
		networkKprobeConnectFetchargs, s.handleSysConnect, expr, false)
}

// RegisterNetworkConnectResultEventFilter registers a network connect result
// event filter with a subscription.
func (s *Subscription) RegisterNetworkConnectResultEventFilter(expr *expression.Expression) {
	s.registerTracepoint("syscalls/sys_exit_connect",
		s.handleSysExitConnect, expr)
}

// RegisterNetworkListenAttemptEventFilter registers a network listen attempt
// event filter with a subscription.
func (s *Subscription) RegisterNetworkListenAttemptEventFilter(expr *expression.Expression) {
	s.registerTracepoint("syscalls/sys_enter_listen",
		s.handleSysEnterListen, expr)
}

// RegisterNetworkListenResultEventFilter registers a network listen result
// event filter with a subscription.
func (s *Subscription) RegisterNetworkListenResultEventFilter(expr *expression.Expression) {
	s.registerTracepoint("syscalls/sys_exit_listen",
		s.handleSysExitListen, expr)
}

// RegisterNetworkRecvfromAttemptEventFilter registers a network recvfrom
// attempt event filter with a subscription.
func (s *Subscription) RegisterNetworkRecvfromAttemptEventFilter(expr *expression.Expression) {
	s.registerTracepoint("syscalls/sys_enter_recvfrom",
		s.handleSysEnterRecvfrom, expr)
	s.registerTracepoint("syscalls/sys_enter_recvmsg",
		s.handleSysEnterRecvfrom, expr)
}

// RegisterNetworkRecvfromResultEventFilter registers a network recvfrom result
// event filter with a subscription.
func (s *Subscription) RegisterNetworkRecvfromResultEventFilter(expr *expression.Expression) {
	s.registerTracepoint("syscalls/sys_exit_recvfrom",
		s.handleSysExitRecvfrom, expr)
	s.registerTracepoint("syscalls/sys_exit_recvmsg",
		s.handleSysExitRecvfrom, expr)
}

// RegisterNetworkSendtoAttemptEventFilter registers a network sendto attempt
// event filter with a subscription.
func (s *Subscription) RegisterNetworkSendtoAttemptEventFilter(expr *expression.Expression) {
	s.registerKprobe(networkKprobeSendmsgSymbol, false,
		networkKprobeSendmsgFetchargs, s.handleSysSendto, expr, false)
	s.registerKprobe(networkKprobeSendtoSymbol, false,
		networkKprobeSendtoFetchargs, s.handleSysSendto, expr, false)
}

// RegisterNetworkSendtoResultEventFilter registers a network sendto result
// event filter with a subscription.
func (s *Subscription) RegisterNetworkSendtoResultEventFilter(expr *expression.Expression) {
	s.registerTracepoint("syscalls/sys_exit_sendmsg",
		s.handleSysExitSendto, expr)
	s.registerTracepoint("syscalls/sys_exit_sendto",
		s.handleSysExitSendto, expr)
}
