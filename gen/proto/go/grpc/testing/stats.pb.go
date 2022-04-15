// Copyright 2015 gRPC authors.
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

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        (unknown)
// source: grpc/testing/stats.proto

package testing

import (
	core "github.com/costinm/grpc-mesh/proto/grpc/core"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ServerStats struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// wall clock time change in seconds since last reset
	TimeElapsed float64 `protobuf:"fixed64,1,opt,name=time_elapsed,json=timeElapsed,proto3" json:"time_elapsed,omitempty"`
	// change in user time (in seconds) used by the server since last reset
	TimeUser float64 `protobuf:"fixed64,2,opt,name=time_user,json=timeUser,proto3" json:"time_user,omitempty"`
	// change in server time (in seconds) used by the server process and all
	// threads since last reset
	TimeSystem float64 `protobuf:"fixed64,3,opt,name=time_system,json=timeSystem,proto3" json:"time_system,omitempty"`
	// change in total cpu time of the server (data from proc/stat)
	TotalCpuTime uint64 `protobuf:"varint,4,opt,name=total_cpu_time,json=totalCpuTime,proto3" json:"total_cpu_time,omitempty"`
	// change in idle time of the server (data from proc/stat)
	IdleCpuTime uint64 `protobuf:"varint,5,opt,name=idle_cpu_time,json=idleCpuTime,proto3" json:"idle_cpu_time,omitempty"`
	// Number of polls called inside completion queue
	CqPollCount uint64 `protobuf:"varint,6,opt,name=cq_poll_count,json=cqPollCount,proto3" json:"cq_poll_count,omitempty"`
	// Core library stats
	CoreStats *core.Stats `protobuf:"bytes,7,opt,name=core_stats,json=coreStats,proto3" json:"core_stats,omitempty"`
}

func (x *ServerStats) Reset() {
	*x = ServerStats{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_testing_stats_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServerStats) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServerStats) ProtoMessage() {}

func (x *ServerStats) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_testing_stats_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServerStats.ProtoReflect.Descriptor instead.
func (*ServerStats) Descriptor() ([]byte, []int) {
	return file_grpc_testing_stats_proto_rawDescGZIP(), []int{0}
}

func (x *ServerStats) GetTimeElapsed() float64 {
	if x != nil {
		return x.TimeElapsed
	}
	return 0
}

func (x *ServerStats) GetTimeUser() float64 {
	if x != nil {
		return x.TimeUser
	}
	return 0
}

func (x *ServerStats) GetTimeSystem() float64 {
	if x != nil {
		return x.TimeSystem
	}
	return 0
}

func (x *ServerStats) GetTotalCpuTime() uint64 {
	if x != nil {
		return x.TotalCpuTime
	}
	return 0
}

func (x *ServerStats) GetIdleCpuTime() uint64 {
	if x != nil {
		return x.IdleCpuTime
	}
	return 0
}

func (x *ServerStats) GetCqPollCount() uint64 {
	if x != nil {
		return x.CqPollCount
	}
	return 0
}

func (x *ServerStats) GetCoreStats() *core.Stats {
	if x != nil {
		return x.CoreStats
	}
	return nil
}

// Histogram params based on grpc/support/histogram.c
type HistogramParams struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Resolution  float64 `protobuf:"fixed64,1,opt,name=resolution,proto3" json:"resolution,omitempty"`                      // first bucket is [0, 1 + resolution)
	MaxPossible float64 `protobuf:"fixed64,2,opt,name=max_possible,json=maxPossible,proto3" json:"max_possible,omitempty"` // use enough buckets to allow this value
}

func (x *HistogramParams) Reset() {
	*x = HistogramParams{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_testing_stats_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HistogramParams) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HistogramParams) ProtoMessage() {}

func (x *HistogramParams) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_testing_stats_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HistogramParams.ProtoReflect.Descriptor instead.
func (*HistogramParams) Descriptor() ([]byte, []int) {
	return file_grpc_testing_stats_proto_rawDescGZIP(), []int{1}
}

func (x *HistogramParams) GetResolution() float64 {
	if x != nil {
		return x.Resolution
	}
	return 0
}

func (x *HistogramParams) GetMaxPossible() float64 {
	if x != nil {
		return x.MaxPossible
	}
	return 0
}

// Histogram data based on grpc/support/histogram.c
type HistogramData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Bucket       []uint32 `protobuf:"varint,1,rep,packed,name=bucket,proto3" json:"bucket,omitempty"`
	MinSeen      float64  `protobuf:"fixed64,2,opt,name=min_seen,json=minSeen,proto3" json:"min_seen,omitempty"`
	MaxSeen      float64  `protobuf:"fixed64,3,opt,name=max_seen,json=maxSeen,proto3" json:"max_seen,omitempty"`
	Sum          float64  `protobuf:"fixed64,4,opt,name=sum,proto3" json:"sum,omitempty"`
	SumOfSquares float64  `protobuf:"fixed64,5,opt,name=sum_of_squares,json=sumOfSquares,proto3" json:"sum_of_squares,omitempty"`
	Count        float64  `protobuf:"fixed64,6,opt,name=count,proto3" json:"count,omitempty"`
}

func (x *HistogramData) Reset() {
	*x = HistogramData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_testing_stats_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HistogramData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HistogramData) ProtoMessage() {}

func (x *HistogramData) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_testing_stats_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HistogramData.ProtoReflect.Descriptor instead.
func (*HistogramData) Descriptor() ([]byte, []int) {
	return file_grpc_testing_stats_proto_rawDescGZIP(), []int{2}
}

func (x *HistogramData) GetBucket() []uint32 {
	if x != nil {
		return x.Bucket
	}
	return nil
}

func (x *HistogramData) GetMinSeen() float64 {
	if x != nil {
		return x.MinSeen
	}
	return 0
}

func (x *HistogramData) GetMaxSeen() float64 {
	if x != nil {
		return x.MaxSeen
	}
	return 0
}

func (x *HistogramData) GetSum() float64 {
	if x != nil {
		return x.Sum
	}
	return 0
}

func (x *HistogramData) GetSumOfSquares() float64 {
	if x != nil {
		return x.SumOfSquares
	}
	return 0
}

func (x *HistogramData) GetCount() float64 {
	if x != nil {
		return x.Count
	}
	return 0
}

type RequestResultCount struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StatusCode int32 `protobuf:"varint,1,opt,name=status_code,json=statusCode,proto3" json:"status_code,omitempty"`
	Count      int64 `protobuf:"varint,2,opt,name=count,proto3" json:"count,omitempty"`
}

func (x *RequestResultCount) Reset() {
	*x = RequestResultCount{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_testing_stats_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestResultCount) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestResultCount) ProtoMessage() {}

func (x *RequestResultCount) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_testing_stats_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestResultCount.ProtoReflect.Descriptor instead.
func (*RequestResultCount) Descriptor() ([]byte, []int) {
	return file_grpc_testing_stats_proto_rawDescGZIP(), []int{3}
}

func (x *RequestResultCount) GetStatusCode() int32 {
	if x != nil {
		return x.StatusCode
	}
	return 0
}

func (x *RequestResultCount) GetCount() int64 {
	if x != nil {
		return x.Count
	}
	return 0
}

type ClientStats struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Latency histogram. Data points are in nanoseconds.
	Latencies *HistogramData `protobuf:"bytes,1,opt,name=latencies,proto3" json:"latencies,omitempty"`
	// See ServerStats for details.
	TimeElapsed float64 `protobuf:"fixed64,2,opt,name=time_elapsed,json=timeElapsed,proto3" json:"time_elapsed,omitempty"`
	TimeUser    float64 `protobuf:"fixed64,3,opt,name=time_user,json=timeUser,proto3" json:"time_user,omitempty"`
	TimeSystem  float64 `protobuf:"fixed64,4,opt,name=time_system,json=timeSystem,proto3" json:"time_system,omitempty"`
	// Number of failed requests (one row per status code seen)
	RequestResults []*RequestResultCount `protobuf:"bytes,5,rep,name=request_results,json=requestResults,proto3" json:"request_results,omitempty"`
	// Number of polls called inside completion queue
	CqPollCount uint64 `protobuf:"varint,6,opt,name=cq_poll_count,json=cqPollCount,proto3" json:"cq_poll_count,omitempty"`
	// Core library stats
	CoreStats *core.Stats `protobuf:"bytes,7,opt,name=core_stats,json=coreStats,proto3" json:"core_stats,omitempty"`
}

func (x *ClientStats) Reset() {
	*x = ClientStats{}
	if protoimpl.UnsafeEnabled {
		mi := &file_grpc_testing_stats_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClientStats) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClientStats) ProtoMessage() {}

func (x *ClientStats) ProtoReflect() protoreflect.Message {
	mi := &file_grpc_testing_stats_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClientStats.ProtoReflect.Descriptor instead.
func (*ClientStats) Descriptor() ([]byte, []int) {
	return file_grpc_testing_stats_proto_rawDescGZIP(), []int{4}
}

func (x *ClientStats) GetLatencies() *HistogramData {
	if x != nil {
		return x.Latencies
	}
	return nil
}

func (x *ClientStats) GetTimeElapsed() float64 {
	if x != nil {
		return x.TimeElapsed
	}
	return 0
}

func (x *ClientStats) GetTimeUser() float64 {
	if x != nil {
		return x.TimeUser
	}
	return 0
}

func (x *ClientStats) GetTimeSystem() float64 {
	if x != nil {
		return x.TimeSystem
	}
	return 0
}

func (x *ClientStats) GetRequestResults() []*RequestResultCount {
	if x != nil {
		return x.RequestResults
	}
	return nil
}

func (x *ClientStats) GetCqPollCount() uint64 {
	if x != nil {
		return x.CqPollCount
	}
	return 0
}

func (x *ClientStats) GetCoreStats() *core.Stats {
	if x != nil {
		return x.CoreStats
	}
	return nil
}

var File_grpc_testing_stats_proto protoreflect.FileDescriptor

var file_grpc_testing_stats_proto_rawDesc = []byte{
	0x0a, 0x18, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x74, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x2f, 0x73,
	0x74, 0x61, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x67, 0x72, 0x70, 0x63,
	0x2e, 0x74, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x1a, 0x15, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x63,
	0x6f, 0x72, 0x65, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x8d, 0x02, 0x0a, 0x0b, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x73, 0x12,
	0x21, 0x0a, 0x0c, 0x74, 0x69, 0x6d, 0x65, 0x5f, 0x65, 0x6c, 0x61, 0x70, 0x73, 0x65, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0b, 0x74, 0x69, 0x6d, 0x65, 0x45, 0x6c, 0x61, 0x70, 0x73,
	0x65, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x08, 0x74, 0x69, 0x6d, 0x65, 0x55, 0x73, 0x65, 0x72, 0x12,
	0x1f, 0x0a, 0x0b, 0x74, 0x69, 0x6d, 0x65, 0x5f, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x01, 0x52, 0x0a, 0x74, 0x69, 0x6d, 0x65, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d,
	0x12, 0x24, 0x0a, 0x0e, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x63, 0x70, 0x75, 0x5f, 0x74, 0x69,
	0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0c, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x43,
	0x70, 0x75, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x22, 0x0a, 0x0d, 0x69, 0x64, 0x6c, 0x65, 0x5f, 0x63,
	0x70, 0x75, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x69,
	0x64, 0x6c, 0x65, 0x43, 0x70, 0x75, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x22, 0x0a, 0x0d, 0x63, 0x71,
	0x5f, 0x70, 0x6f, 0x6c, 0x6c, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x0b, 0x63, 0x71, 0x50, 0x6f, 0x6c, 0x6c, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x2f,
	0x0a, 0x0a, 0x63, 0x6f, 0x72, 0x65, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x73, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x10, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x53,
	0x74, 0x61, 0x74, 0x73, 0x52, 0x09, 0x63, 0x6f, 0x72, 0x65, 0x53, 0x74, 0x61, 0x74, 0x73, 0x22,
	0x54, 0x0a, 0x0f, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x67, 0x72, 0x61, 0x6d, 0x50, 0x61, 0x72, 0x61,
	0x6d, 0x73, 0x12, 0x1e, 0x0a, 0x0a, 0x72, 0x65, 0x73, 0x6f, 0x6c, 0x75, 0x74, 0x69, 0x6f, 0x6e,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0a, 0x72, 0x65, 0x73, 0x6f, 0x6c, 0x75, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x21, 0x0a, 0x0c, 0x6d, 0x61, 0x78, 0x5f, 0x70, 0x6f, 0x73, 0x73, 0x69, 0x62,
	0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0b, 0x6d, 0x61, 0x78, 0x50, 0x6f, 0x73,
	0x73, 0x69, 0x62, 0x6c, 0x65, 0x22, 0xab, 0x01, 0x0a, 0x0d, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x67,
	0x72, 0x61, 0x6d, 0x44, 0x61, 0x74, 0x61, 0x12, 0x16, 0x0a, 0x06, 0x62, 0x75, 0x63, 0x6b, 0x65,
	0x74, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0d, 0x52, 0x06, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x12,
	0x19, 0x0a, 0x08, 0x6d, 0x69, 0x6e, 0x5f, 0x73, 0x65, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x01, 0x52, 0x07, 0x6d, 0x69, 0x6e, 0x53, 0x65, 0x65, 0x6e, 0x12, 0x19, 0x0a, 0x08, 0x6d, 0x61,
	0x78, 0x5f, 0x73, 0x65, 0x65, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x07, 0x6d, 0x61,
	0x78, 0x53, 0x65, 0x65, 0x6e, 0x12, 0x10, 0x0a, 0x03, 0x73, 0x75, 0x6d, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x01, 0x52, 0x03, 0x73, 0x75, 0x6d, 0x12, 0x24, 0x0a, 0x0e, 0x73, 0x75, 0x6d, 0x5f, 0x6f,
	0x66, 0x5f, 0x73, 0x71, 0x75, 0x61, 0x72, 0x65, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x01, 0x52,
	0x0c, 0x73, 0x75, 0x6d, 0x4f, 0x66, 0x53, 0x71, 0x75, 0x61, 0x72, 0x65, 0x73, 0x12, 0x14, 0x0a,
	0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x22, 0x4b, 0x0a, 0x12, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x52, 0x65,
	0x73, 0x75, 0x6c, 0x74, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x22, 0xc9, 0x02, 0x0a, 0x0b, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x74, 0x73,
	0x12, 0x39, 0x0a, 0x09, 0x6c, 0x61, 0x74, 0x65, 0x6e, 0x63, 0x69, 0x65, 0x73, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x69,
	0x6e, 0x67, 0x2e, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x67, 0x72, 0x61, 0x6d, 0x44, 0x61, 0x74, 0x61,
	0x52, 0x09, 0x6c, 0x61, 0x74, 0x65, 0x6e, 0x63, 0x69, 0x65, 0x73, 0x12, 0x21, 0x0a, 0x0c, 0x74,
	0x69, 0x6d, 0x65, 0x5f, 0x65, 0x6c, 0x61, 0x70, 0x73, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x01, 0x52, 0x0b, 0x74, 0x69, 0x6d, 0x65, 0x45, 0x6c, 0x61, 0x70, 0x73, 0x65, 0x64, 0x12, 0x1b,
	0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x01, 0x52, 0x08, 0x74, 0x69, 0x6d, 0x65, 0x55, 0x73, 0x65, 0x72, 0x12, 0x1f, 0x0a, 0x0b, 0x74,
	0x69, 0x6d, 0x65, 0x5f, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x18, 0x04, 0x20, 0x01, 0x28, 0x01,
	0x52, 0x0a, 0x74, 0x69, 0x6d, 0x65, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x12, 0x49, 0x0a, 0x0f,
	0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x5f, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x18,
	0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x74, 0x65, 0x73,
	0x74, 0x69, 0x6e, 0x67, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x52, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x0e, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x12, 0x22, 0x0a, 0x0d, 0x63, 0x71, 0x5f, 0x70, 0x6f,
	0x6c, 0x6c, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b,
	0x63, 0x71, 0x50, 0x6f, 0x6c, 0x6c, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x2f, 0x0a, 0x0a, 0x63,
	0x6f, 0x72, 0x65, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x10, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x53, 0x74, 0x61, 0x74,
	0x73, 0x52, 0x09, 0x63, 0x6f, 0x72, 0x65, 0x53, 0x74, 0x61, 0x74, 0x73, 0x42, 0x50, 0x0a, 0x0f,
	0x69, 0x6f, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x42,
	0x0a, 0x53, 0x74, 0x61, 0x74, 0x73, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x2f, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x6f, 0x73, 0x74, 0x69, 0x6e,
	0x6d, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2d, 0x6d, 0x65, 0x73, 0x68, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x74, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_grpc_testing_stats_proto_rawDescOnce sync.Once
	file_grpc_testing_stats_proto_rawDescData = file_grpc_testing_stats_proto_rawDesc
)

func file_grpc_testing_stats_proto_rawDescGZIP() []byte {
	file_grpc_testing_stats_proto_rawDescOnce.Do(func() {
		file_grpc_testing_stats_proto_rawDescData = protoimpl.X.CompressGZIP(file_grpc_testing_stats_proto_rawDescData)
	})
	return file_grpc_testing_stats_proto_rawDescData
}

var file_grpc_testing_stats_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_grpc_testing_stats_proto_goTypes = []interface{}{
	(*ServerStats)(nil),        // 0: grpc.testing.ServerStats
	(*HistogramParams)(nil),    // 1: grpc.testing.HistogramParams
	(*HistogramData)(nil),      // 2: grpc.testing.HistogramData
	(*RequestResultCount)(nil), // 3: grpc.testing.RequestResultCount
	(*ClientStats)(nil),        // 4: grpc.testing.ClientStats
	(*core.Stats)(nil),         // 5: grpc.core.Stats
}
var file_grpc_testing_stats_proto_depIdxs = []int32{
	5, // 0: grpc.testing.ServerStats.core_stats:type_name -> grpc.core.Stats
	2, // 1: grpc.testing.ClientStats.latencies:type_name -> grpc.testing.HistogramData
	3, // 2: grpc.testing.ClientStats.request_results:type_name -> grpc.testing.RequestResultCount
	5, // 3: grpc.testing.ClientStats.core_stats:type_name -> grpc.core.Stats
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_grpc_testing_stats_proto_init() }
func file_grpc_testing_stats_proto_init() {
	if File_grpc_testing_stats_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_grpc_testing_stats_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServerStats); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_grpc_testing_stats_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HistogramParams); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_grpc_testing_stats_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HistogramData); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_grpc_testing_stats_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestResultCount); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_grpc_testing_stats_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClientStats); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_grpc_testing_stats_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_grpc_testing_stats_proto_goTypes,
		DependencyIndexes: file_grpc_testing_stats_proto_depIdxs,
		MessageInfos:      file_grpc_testing_stats_proto_msgTypes,
	}.Build()
	File_grpc_testing_stats_proto = out.File
	file_grpc_testing_stats_proto_rawDesc = nil
	file_grpc_testing_stats_proto_goTypes = nil
	file_grpc_testing_stats_proto_depIdxs = nil
}
