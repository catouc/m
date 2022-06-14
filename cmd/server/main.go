package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/bufbuild/connect-go"
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
}

func (ms *MServer) ListVideosForChannel(_ context.Context, req *connect.Request[v1.YoutubeChanneListRequest]) (*connect.Response[v1.YoutubeChannelListResponse], error) {
	v, err := ms.ytClient.GetLatestVideosFromChannel(req.Msg.ChannelName)
	if err != nil {
		return nil, err
	}

	response := connect.NewResponse(&v1.YoutubeChannelListResponse{
		Videos: v,
	})
	return response, nil
}

func main() {
	ms := MServer{}

	path, handler := mv1connect.NewMServiceHandler(&ms)
	mux := http.NewServeMux()
	mux.Handle(path, handler)

	log.Fatal(http.ListenAndServe(*listenAddr, h2c.NewHandler(mux, &http2.Server{})))
}
