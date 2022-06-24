package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/bufbuild/connect-go"
	"github.com/catouc/m/internal/hntop"
	"log"
	"net/http"

	v1 "github.com/catouc/m/internal/m/v1"
	"github.com/catouc/m/internal/m/v1/mv1connect"
	"github.com/catouc/m/internal/youtube"
	"golang.org/x/net/context"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

var (
	listenAddr = flag.String("listenAddr", ":8080", "Listen address for the server")
)

func init() {
	flag.Parse()
}

type MServer struct {
	ytClient *youtube.Client
	hnClient *hntop.Client
}

func New() *MServer {
	return &MServer{hnClient: hntop.New(context.Background())}
}

func (ms *MServer) ListVideosForChannel(_ context.Context, req *connect.Request[v1.YoutubeChanneListRequest]) (*connect.Response[v1.YoutubeVideoListResponse], error) {
	v, err := ms.ytClient.GetLatestVideosFromChannel(req.Msg.ChannelName)
	if err != nil {
		return nil, err
	}

	response := connect.NewResponse(&v1.YoutubeVideoListResponse{
		Videos: v,
	})
	return response, nil
}

func (ms *MServer) ListNewBlogPosts(_ context.Context, req *connect.Request[v1.ListNewBlogPostRequest]) (*connect.Response[v1.ListNewBlogPostResponse], error) {
	return nil, errors.New("not implemented")
}

func (ms *MServer) ListVideosForCategory(_ context.Context, req *connect.Request[v1.YoutubeCategoryListRequest]) (*connect.Response[v1.YoutubeVideoListResponse], error) {
	return nil, errors.New("not implemented")
}

func (ms *MServer) GetHNFrontpage(ctx context.Context, _ *connect.Request[v1.HNFrontpageRequest]) (*connect.Response[v1.HNFrontpageResponse], error) {
	frontPage, err := ms.hnClient.GetHNTop30Stories(ctx)
	if err != nil {
		return nil, err
	}

	fmt.Println(frontPage[0].URL)

	response := connect.NewResponse(&v1.HNFrontpageResponse{Stories: frontPage})
	return response, nil
}

func main() {
	path, handler := mv1connect.NewMServiceHandler(New())
	mux := http.NewServeMux()
	mux.Handle(path, handler)

	log.Fatal(http.ListenAndServe(*listenAddr, h2c.NewHandler(mux, &http2.Server{})))
}
