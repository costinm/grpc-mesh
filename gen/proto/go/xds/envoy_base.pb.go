// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        (unknown)
// source: xds/envoy_base.proto

package xds

import (
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

type SocketAddress_Protocol int32

const (
	SocketAddress_TCP SocketAddress_Protocol = 0
	// [#not-implemented-hide:]
	SocketAddress_UDP SocketAddress_Protocol = 1
)

// Enum value maps for SocketAddress_Protocol.
var (
	SocketAddress_Protocol_name = map[int32]string{
		0: "TCP",
		1: "UDP",
	}
	SocketAddress_Protocol_value = map[string]int32{
		"TCP": 0,
		"UDP": 1,
	}
)

func (x SocketAddress_Protocol) Enum() *SocketAddress_Protocol {
	p := new(SocketAddress_Protocol)
	*p = x
	return p
}

func (x SocketAddress_Protocol) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (SocketAddress_Protocol) Descriptor() protoreflect.EnumDescriptor {
	return file_xds_envoy_base_proto_enumTypes[0].Descriptor()
}

func (SocketAddress_Protocol) Type() protoreflect.EnumType {
	return &file_xds_envoy_base_proto_enumTypes[0]
}

func (x SocketAddress_Protocol) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use SocketAddress_Protocol.Descriptor instead.
func (SocketAddress_Protocol) EnumDescriptor() ([]byte, []int) {
	return file_xds_envoy_base_proto_rawDescGZIP(), []int{1, 0}
}

// Identifies location of where either Envoy runs or where upstream hosts run.
type Locality struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Region this :ref:`zone <envoy_api_field_core.Locality.zone>` belongs to.
	Region string `protobuf:"bytes,1,opt,name=region,proto3" json:"region,omitempty"`
	// Defines the local service zone where Envoy is running. Though optional, it
	// should be set if discovery service routing is used and the discovery
	// service exposes :ref:`zone data <config_cluster_manager_sds_api_host_az>`,
	// either in this message or via :option:`--service-zone`. The meaning of zone
	// is context dependent, e.g. `Availability Zone (AZ)
	// <https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/using-regions-availability-zones.html>`_
	// on AWS, `Zone <https://cloud.google.com/compute/docs/regions-zones/>`_ on
	// GCP, etc.
	Zone string `protobuf:"bytes,2,opt,name=zone,proto3" json:"zone,omitempty"`
	// When used for locality of upstream hosts, this field further splits zone
	// into smaller chunks of sub-zones so they can be load balanced
	// independently.
	SubZone string `protobuf:"bytes,3,opt,name=sub_zone,json=subZone,proto3" json:"sub_zone,omitempty"`
}

func (x *Locality) Reset() {
	*x = Locality{}
	if protoimpl.UnsafeEnabled {
		mi := &file_xds_envoy_base_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Locality) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Locality) ProtoMessage() {}

func (x *Locality) ProtoReflect() protoreflect.Message {
	mi := &file_xds_envoy_base_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Locality.ProtoReflect.Descriptor instead.
func (*Locality) Descriptor() ([]byte, []int) {
	return file_xds_envoy_base_proto_rawDescGZIP(), []int{0}
}

func (x *Locality) GetRegion() string {
	if x != nil {
		return x.Region
	}
	return ""
}

func (x *Locality) GetZone() string {
	if x != nil {
		return x.Zone
	}
	return ""
}

func (x *Locality) GetSubZone() string {
	if x != nil {
		return x.SubZone
	}
	return ""
}

type SocketAddress struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Protocol SocketAddress_Protocol `protobuf:"varint,1,opt,name=protocol,proto3,enum=xds.SocketAddress_Protocol" json:"protocol,omitempty"`
	// The address for this socket. :ref:`Listeners <config_listeners>` will bind
	// to the address or outbound connections will be made. An empty address is
	// not allowed, specify ``0.0.0.0`` or ``::`` to bind any. It's still possible to
	// distinguish on an address via the prefix/suffix matching in
	// FilterChainMatch after connection. For :ref:`clusters
	// <config_cluster_manager_cluster>`, an address may be either an IP or
	// hostname to be resolved via DNS. If it is a hostname, :ref:`resolver_name
	// <envoy_api_field_core.SocketAddress.resolver_name>` should be set unless default
	// (i.e. DNS) resolution is expected.
	Address string `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
	// Types that are assignable to PortSpecifier:
	//	*SocketAddress_PortValue
	//	*SocketAddress_NamedPort
	PortSpecifier isSocketAddress_PortSpecifier `protobuf_oneof:"port_specifier"`
	// The name of the resolver. This must have been registered with Envoy. If this is
	// empty, a context dependent default applies. If address is a hostname this
	// should be set for resolution other than DNS. If the address is a concrete
	// IP address, no resolution will occur.
	ResolverName string `protobuf:"bytes,5,opt,name=resolver_name,json=resolverName,proto3" json:"resolver_name,omitempty"`
	// When binding to an IPv6 address above, this enables `IPv4 compatibity
	// <https://tools.ietf.org/html/rfc3493#page-11>`_. Binding to ``::`` will
	// allow both IPv4 and IPv6 connections, with peer IPv4 addresses mapped into
	// IPv6 space as ``::FFFF:<IPv4-address>``.
	Ipv4Compat bool `protobuf:"varint,6,opt,name=ipv4_compat,json=ipv4Compat,proto3" json:"ipv4_compat,omitempty"`
}

func (x *SocketAddress) Reset() {
	*x = SocketAddress{}
	if protoimpl.UnsafeEnabled {
		mi := &file_xds_envoy_base_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SocketAddress) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SocketAddress) ProtoMessage() {}

func (x *SocketAddress) ProtoReflect() protoreflect.Message {
	mi := &file_xds_envoy_base_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SocketAddress.ProtoReflect.Descriptor instead.
func (*SocketAddress) Descriptor() ([]byte, []int) {
	return file_xds_envoy_base_proto_rawDescGZIP(), []int{1}
}

func (x *SocketAddress) GetProtocol() SocketAddress_Protocol {
	if x != nil {
		return x.Protocol
	}
	return SocketAddress_TCP
}

func (x *SocketAddress) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (m *SocketAddress) GetPortSpecifier() isSocketAddress_PortSpecifier {
	if m != nil {
		return m.PortSpecifier
	}
	return nil
}

func (x *SocketAddress) GetPortValue() uint32 {
	if x, ok := x.GetPortSpecifier().(*SocketAddress_PortValue); ok {
		return x.PortValue
	}
	return 0
}

func (x *SocketAddress) GetNamedPort() string {
	if x, ok := x.GetPortSpecifier().(*SocketAddress_NamedPort); ok {
		return x.NamedPort
	}
	return ""
}

func (x *SocketAddress) GetResolverName() string {
	if x != nil {
		return x.ResolverName
	}
	return ""
}

func (x *SocketAddress) GetIpv4Compat() bool {
	if x != nil {
		return x.Ipv4Compat
	}
	return false
}

type isSocketAddress_PortSpecifier interface {
	isSocketAddress_PortSpecifier()
}

type SocketAddress_PortValue struct {
	PortValue uint32 `protobuf:"varint,3,opt,name=port_value,json=portValue,proto3,oneof"`
}

type SocketAddress_NamedPort struct {
	// This is only valid if :ref:`resolver_name
	// <envoy_api_field_core.SocketAddress.resolver_name>` is specified below and the
	// named resolver is capable of named port resolution.
	NamedPort string `protobuf:"bytes,4,opt,name=named_port,json=namedPort,proto3,oneof"`
}

func (*SocketAddress_PortValue) isSocketAddress_PortSpecifier() {}

func (*SocketAddress_NamedPort) isSocketAddress_PortSpecifier() {}

// Metadata provides additional inputs to filters based on matched listeners,
// filter chains, routes and endpoints. It is structured as a map from filter
// name (in reverse DNS format) to metadata specific to the filter. Metadata
// key-values for a filter are merged as connection and request handling occurs,
// with later values for the same key overriding earlier values.
//
// An example use of metadata is providing additional values to
// http_connection_manager in the envoy.http_connection_manager.access_log
// namespace.
//
// For load balancing, Metadata provides a means to subset cluster endpoints.
// Endpoints have a Metadata object associated and routes contain a Metadata
// object to match against. There are some well defined metadata used today for
// this purpose:
//
// * ``{"envoy.lb": {"canary": <bool> }}`` This indicates the canary status of an
//   endpoint and is also used during header processing
//   (x-envoy-upstream-canary) and for stats purposes.
type Metadata struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Key is the reverse DNS filter name, e.g. com.acme.widget. The envoy.*
	// namespace is reserved for Envoy's built-in filters.
	FilterMetadata map[string]*Struct `protobuf:"bytes,1,rep,name=filter_metadata,json=filterMetadata,proto3" json:"filter_metadata,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *Metadata) Reset() {
	*x = Metadata{}
	if protoimpl.UnsafeEnabled {
		mi := &file_xds_envoy_base_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Metadata) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Metadata) ProtoMessage() {}

func (x *Metadata) ProtoReflect() protoreflect.Message {
	mi := &file_xds_envoy_base_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Metadata.ProtoReflect.Descriptor instead.
func (*Metadata) Descriptor() ([]byte, []int) {
	return file_xds_envoy_base_proto_rawDescGZIP(), []int{2}
}

func (x *Metadata) GetFilterMetadata() map[string]*Struct {
	if x != nil {
		return x.FilterMetadata
	}
	return nil
}

// Addresses specify either a logical or physical address and port, which are
// used to tell Envoy where to bind/listen, connect to upstream and find
// management servers.
type Address struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Address:
	//	*Address_SocketAddress
	//	*Address_Pipe
	Address isAddress_Address `protobuf_oneof:"address"`
}

func (x *Address) Reset() {
	*x = Address{}
	if protoimpl.UnsafeEnabled {
		mi := &file_xds_envoy_base_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Address) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Address) ProtoMessage() {}

func (x *Address) ProtoReflect() protoreflect.Message {
	mi := &file_xds_envoy_base_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Address.ProtoReflect.Descriptor instead.
func (*Address) Descriptor() ([]byte, []int) {
	return file_xds_envoy_base_proto_rawDescGZIP(), []int{3}
}

func (m *Address) GetAddress() isAddress_Address {
	if m != nil {
		return m.Address
	}
	return nil
}

func (x *Address) GetSocketAddress() *SocketAddress {
	if x, ok := x.GetAddress().(*Address_SocketAddress); ok {
		return x.SocketAddress
	}
	return nil
}

func (x *Address) GetPipe() *Pipe {
	if x, ok := x.GetAddress().(*Address_Pipe); ok {
		return x.Pipe
	}
	return nil
}

type isAddress_Address interface {
	isAddress_Address()
}

type Address_SocketAddress struct {
	SocketAddress *SocketAddress `protobuf:"bytes,1,opt,name=socket_address,json=socketAddress,proto3,oneof"`
}

type Address_Pipe struct {
	Pipe *Pipe `protobuf:"bytes,2,opt,name=pipe,proto3,oneof"`
}

func (*Address_SocketAddress) isAddress_Address() {}

func (*Address_Pipe) isAddress_Address() {}

type Pipe struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Unix Domain Socket path. On Linux, paths starting with '@' will use the
	// abstract namespace. The starting '@' is replaced by a null byte by Envoy.
	// Paths starting with '@' will result in an error in environments other than
	// Linux.
	Path string `protobuf:"bytes,1,opt,name=path,proto3" json:"path,omitempty"`
}

func (x *Pipe) Reset() {
	*x = Pipe{}
	if protoimpl.UnsafeEnabled {
		mi := &file_xds_envoy_base_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Pipe) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Pipe) ProtoMessage() {}

func (x *Pipe) ProtoReflect() protoreflect.Message {
	mi := &file_xds_envoy_base_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Pipe.ProtoReflect.Descriptor instead.
func (*Pipe) Descriptor() ([]byte, []int) {
	return file_xds_envoy_base_proto_rawDescGZIP(), []int{4}
}

func (x *Pipe) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

type BindConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The address to bind to when creating a socket.
	SourceAddress *SocketAddress `protobuf:"bytes,1,opt,name=source_address,json=sourceAddress,proto3" json:"source_address,omitempty"`
}

func (x *BindConfig) Reset() {
	*x = BindConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_xds_envoy_base_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BindConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BindConfig) ProtoMessage() {}

func (x *BindConfig) ProtoReflect() protoreflect.Message {
	mi := &file_xds_envoy_base_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BindConfig.ProtoReflect.Descriptor instead.
func (*BindConfig) Descriptor() ([]byte, []int) {
	return file_xds_envoy_base_proto_rawDescGZIP(), []int{5}
}

func (x *BindConfig) GetSourceAddress() *SocketAddress {
	if x != nil {
		return x.SourceAddress
	}
	return nil
}

// CidrRange specifies an IP Address and a prefix length to construct
// the subnet mask for a `CIDR <https://tools.ietf.org/html/rfc4632>`_ range.
type CidrRange struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// IPv4 or IPv6 address, e.g. ``192.0.0.0`` or ``2001:db8::``.
	AddressPrefix string `protobuf:"bytes,1,opt,name=address_prefix,json=addressPrefix,proto3" json:"address_prefix,omitempty"`
	// Length of prefix, e.g. 0, 32.
	PrefixLen *UInt32Value `protobuf:"bytes,2,opt,name=prefix_len,json=prefixLen,proto3" json:"prefix_len,omitempty"`
}

func (x *CidrRange) Reset() {
	*x = CidrRange{}
	if protoimpl.UnsafeEnabled {
		mi := &file_xds_envoy_base_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CidrRange) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CidrRange) ProtoMessage() {}

func (x *CidrRange) ProtoReflect() protoreflect.Message {
	mi := &file_xds_envoy_base_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CidrRange.ProtoReflect.Descriptor instead.
func (*CidrRange) Descriptor() ([]byte, []int) {
	return file_xds_envoy_base_proto_rawDescGZIP(), []int{6}
}

func (x *CidrRange) GetAddressPrefix() string {
	if x != nil {
		return x.AddressPrefix
	}
	return ""
}

func (x *CidrRange) GetPrefixLen() *UInt32Value {
	if x != nil {
		return x.PrefixLen
	}
	return nil
}

var File_xds_envoy_base_proto protoreflect.FileDescriptor

var file_xds_envoy_base_proto_rawDesc = []byte{
	0x0a, 0x14, 0x78, 0x64, 0x73, 0x2f, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x5f, 0x62, 0x61, 0x73, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x78, 0x64, 0x73, 0x1a, 0x0e, 0x78, 0x64, 0x73,
	0x2f, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x51, 0x0a, 0x08, 0x4c,
	0x6f, 0x63, 0x61, 0x6c, 0x69, 0x74, 0x79, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x67, 0x69, 0x6f,
	0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x12,
	0x12, 0x0a, 0x04, 0x7a, 0x6f, 0x6e, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x7a,
	0x6f, 0x6e, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x73, 0x75, 0x62, 0x5f, 0x7a, 0x6f, 0x6e, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x73, 0x75, 0x62, 0x5a, 0x6f, 0x6e, 0x65, 0x22, 0x9a,
	0x02, 0x0a, 0x0d, 0x53, 0x6f, 0x63, 0x6b, 0x65, 0x74, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x12, 0x37, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x1b, 0x2e, 0x78, 0x64, 0x73, 0x2e, 0x53, 0x6f, 0x63, 0x6b, 0x65, 0x74, 0x41,
	0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x2e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x52,
	0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64,
	0x72, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72,
	0x65, 0x73, 0x73, 0x12, 0x1f, 0x0a, 0x0a, 0x70, 0x6f, 0x72, 0x74, 0x5f, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x48, 0x00, 0x52, 0x09, 0x70, 0x6f, 0x72, 0x74, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x12, 0x1f, 0x0a, 0x0a, 0x6e, 0x61, 0x6d, 0x65, 0x64, 0x5f, 0x70, 0x6f,
	0x72, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x09, 0x6e, 0x61, 0x6d, 0x65,
	0x64, 0x50, 0x6f, 0x72, 0x74, 0x12, 0x23, 0x0a, 0x0d, 0x72, 0x65, 0x73, 0x6f, 0x6c, 0x76, 0x65,
	0x72, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x72, 0x65,
	0x73, 0x6f, 0x6c, 0x76, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x69, 0x70,
	0x76, 0x34, 0x5f, 0x63, 0x6f, 0x6d, 0x70, 0x61, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x0a, 0x69, 0x70, 0x76, 0x34, 0x43, 0x6f, 0x6d, 0x70, 0x61, 0x74, 0x22, 0x1c, 0x0a, 0x08, 0x50,
	0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x12, 0x07, 0x0a, 0x03, 0x54, 0x43, 0x50, 0x10, 0x00,
	0x12, 0x07, 0x0a, 0x03, 0x55, 0x44, 0x50, 0x10, 0x01, 0x42, 0x10, 0x0a, 0x0e, 0x70, 0x6f, 0x72,
	0x74, 0x5f, 0x73, 0x70, 0x65, 0x63, 0x69, 0x66, 0x69, 0x65, 0x72, 0x22, 0xa6, 0x01, 0x0a, 0x08,
	0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x4a, 0x0a, 0x0f, 0x66, 0x69, 0x6c, 0x74,
	0x65, 0x72, 0x5f, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x21, 0x2e, 0x78, 0x64, 0x73, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0x2e, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x52, 0x0e, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x4d, 0x65, 0x74, 0x61,
	0x64, 0x61, 0x74, 0x61, 0x1a, 0x4e, 0x0a, 0x13, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x4d, 0x65,
	0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x21, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x78,
	0x64, 0x73, 0x2e, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x3a, 0x02, 0x38, 0x01, 0x22, 0x72, 0x0a, 0x07, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12,
	0x3b, 0x0a, 0x0e, 0x73, 0x6f, 0x63, 0x6b, 0x65, 0x74, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73,
	0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x78, 0x64, 0x73, 0x2e, 0x53, 0x6f,
	0x63, 0x6b, 0x65, 0x74, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x48, 0x00, 0x52, 0x0d, 0x73,
	0x6f, 0x63, 0x6b, 0x65, 0x74, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x1f, 0x0a, 0x04,
	0x70, 0x69, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x09, 0x2e, 0x78, 0x64, 0x73,
	0x2e, 0x50, 0x69, 0x70, 0x65, 0x48, 0x00, 0x52, 0x04, 0x70, 0x69, 0x70, 0x65, 0x42, 0x09, 0x0a,
	0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x22, 0x1a, 0x0a, 0x04, 0x50, 0x69, 0x70, 0x65,
	0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x74, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x70, 0x61, 0x74, 0x68, 0x22, 0x47, 0x0a, 0x0a, 0x42, 0x69, 0x6e, 0x64, 0x43, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x12, 0x39, 0x0a, 0x0e, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x61, 0x64, 0x64,
	0x72, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x78, 0x64, 0x73,
	0x2e, 0x53, 0x6f, 0x63, 0x6b, 0x65, 0x74, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x52, 0x0d,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x22, 0x63, 0x0a,
	0x09, 0x43, 0x69, 0x64, 0x72, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x12, 0x25, 0x0a, 0x0e, 0x61, 0x64,
	0x64, 0x72, 0x65, 0x73, 0x73, 0x5f, 0x70, 0x72, 0x65, 0x66, 0x69, 0x78, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0d, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x50, 0x72, 0x65, 0x66, 0x69,
	0x78, 0x12, 0x2f, 0x0a, 0x0a, 0x70, 0x72, 0x65, 0x66, 0x69, 0x78, 0x5f, 0x6c, 0x65, 0x6e, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x78, 0x64, 0x73, 0x2e, 0x55, 0x49, 0x6e, 0x74,
	0x33, 0x32, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x09, 0x70, 0x72, 0x65, 0x66, 0x69, 0x78, 0x4c,
	0x65, 0x6e, 0x42, 0x32, 0x5a, 0x2d, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x63, 0x6f, 0x73, 0x74, 0x69, 0x6e, 0x6d, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2d, 0x6d, 0x65,
	0x73, 0x68, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x2f,
	0x78, 0x64, 0x73, 0x88, 0x01, 0x01, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_xds_envoy_base_proto_rawDescOnce sync.Once
	file_xds_envoy_base_proto_rawDescData = file_xds_envoy_base_proto_rawDesc
)

func file_xds_envoy_base_proto_rawDescGZIP() []byte {
	file_xds_envoy_base_proto_rawDescOnce.Do(func() {
		file_xds_envoy_base_proto_rawDescData = protoimpl.X.CompressGZIP(file_xds_envoy_base_proto_rawDescData)
	})
	return file_xds_envoy_base_proto_rawDescData
}

var file_xds_envoy_base_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_xds_envoy_base_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_xds_envoy_base_proto_goTypes = []interface{}{
	(SocketAddress_Protocol)(0), // 0: xds.SocketAddress.Protocol
	(*Locality)(nil),            // 1: xds.Locality
	(*SocketAddress)(nil),       // 2: xds.SocketAddress
	(*Metadata)(nil),            // 3: xds.Metadata
	(*Address)(nil),             // 4: xds.Address
	(*Pipe)(nil),                // 5: xds.Pipe
	(*BindConfig)(nil),          // 6: xds.BindConfig
	(*CidrRange)(nil),           // 7: xds.CidrRange
	nil,                         // 8: xds.Metadata.FilterMetadataEntry
	(*UInt32Value)(nil),         // 9: xds.UInt32Value
	(*Struct)(nil),              // 10: xds.Struct
}
var file_xds_envoy_base_proto_depIdxs = []int32{
	0,  // 0: xds.SocketAddress.protocol:type_name -> xds.SocketAddress.Protocol
	8,  // 1: xds.Metadata.filter_metadata:type_name -> xds.Metadata.FilterMetadataEntry
	2,  // 2: xds.Address.socket_address:type_name -> xds.SocketAddress
	5,  // 3: xds.Address.pipe:type_name -> xds.Pipe
	2,  // 4: xds.BindConfig.source_address:type_name -> xds.SocketAddress
	9,  // 5: xds.CidrRange.prefix_len:type_name -> xds.UInt32Value
	10, // 6: xds.Metadata.FilterMetadataEntry.value:type_name -> xds.Struct
	7,  // [7:7] is the sub-list for method output_type
	7,  // [7:7] is the sub-list for method input_type
	7,  // [7:7] is the sub-list for extension type_name
	7,  // [7:7] is the sub-list for extension extendee
	0,  // [0:7] is the sub-list for field type_name
}

func init() { file_xds_envoy_base_proto_init() }
func file_xds_envoy_base_proto_init() {
	if File_xds_envoy_base_proto != nil {
		return
	}
	file_xds_base_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_xds_envoy_base_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Locality); i {
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
		file_xds_envoy_base_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SocketAddress); i {
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
		file_xds_envoy_base_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Metadata); i {
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
		file_xds_envoy_base_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Address); i {
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
		file_xds_envoy_base_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Pipe); i {
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
		file_xds_envoy_base_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BindConfig); i {
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
		file_xds_envoy_base_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CidrRange); i {
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
	file_xds_envoy_base_proto_msgTypes[1].OneofWrappers = []interface{}{
		(*SocketAddress_PortValue)(nil),
		(*SocketAddress_NamedPort)(nil),
	}
	file_xds_envoy_base_proto_msgTypes[3].OneofWrappers = []interface{}{
		(*Address_SocketAddress)(nil),
		(*Address_Pipe)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_xds_envoy_base_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_xds_envoy_base_proto_goTypes,
		DependencyIndexes: file_xds_envoy_base_proto_depIdxs,
		EnumInfos:         file_xds_envoy_base_proto_enumTypes,
		MessageInfos:      file_xds_envoy_base_proto_msgTypes,
	}.Build()
	File_xds_envoy_base_proto = out.File
	file_xds_envoy_base_proto_rawDesc = nil
	file_xds_envoy_base_proto_goTypes = nil
	file_xds_envoy_base_proto_depIdxs = nil
}
