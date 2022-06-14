// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: m/v1/m.proto

package mv1connect

import (
	context "context"
	errors "errors"
	connect_go "github.com/bufbuild/connect-go"
	v1 "github.com/catouc/m/internal/m/v1"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect_go.IsAtLeastVersion0_1_0

const (
	// MServiceName is the fully-qualified name of the MService service.
	MServiceName = "m.v1.MService"
)

// MServiceClient is a client for the m.v1.MService service.
type MServiceClient interface {
	ListVideosForChannel(context.Context, *connect_go.Request[v1.YoutubeChanneListRequest]) (*connect_go.Response[v1.YoutubeChannelListResponse], error)
}

// NewMServiceClient constructs a client for the m.v1.MService service. By default, it uses the
// Connect protocol with the binary Protobuf Codec, asks for gzipped responses, and sends
// uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or
// connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewMServiceClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) MServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &mServiceClient{
		listVideosForChannel: connect_go.NewClient[v1.YoutubeChanneListRequest, v1.YoutubeChannelListResponse](
			httpClient,
			baseURL+"/m.v1.MService/ListVideosForChannel",
			opts...,
		),
	}
}

// mServiceClient implements MServiceClient.
type mServiceClient struct {
	listVideosForChannel *connect_go.Client[v1.YoutubeChanneListRequest, v1.YoutubeChannelListResponse]
}

// ListVideosForChannel calls m.v1.MService.ListVideosForChannel.
func (c *mServiceClient) ListVideosForChannel(ctx context.Context, req *connect_go.Request[v1.YoutubeChanneListRequest]) (*connect_go.Response[v1.YoutubeChannelListResponse], error) {
	return c.listVideosForChannel.CallUnary(ctx, req)
}

// MServiceHandler is an implementation of the m.v1.MService service.
type MServiceHandler interface {
	ListVideosForChannel(context.Context, *connect_go.Request[v1.YoutubeChanneListRequest]) (*connect_go.Response[v1.YoutubeChannelListResponse], error)
}

// NewMServiceHandler builds an HTTP handler from the service implementation. It returns the path on
// which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewMServiceHandler(svc MServiceHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	mux := http.NewServeMux()
	mux.Handle("/m.v1.MService/ListVideosForChannel", connect_go.NewUnaryHandler(
		"/m.v1.MService/ListVideosForChannel",
		svc.ListVideosForChannel,
		opts...,
	))
	return "/m.v1.MService/", mux
}

// UnimplementedMServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedMServiceHandler struct{}

func (UnimplementedMServiceHandler) ListVideosForChannel(context.Context, *connect_go.Request[v1.YoutubeChanneListRequest]) (*connect_go.Response[v1.YoutubeChannelListResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("m.v1.MService.ListVideosForChannel is not implemented"))
}
