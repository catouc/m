package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/catouc/m/internal/m/v1/mv1connect"
	"github.com/catouc/m/internal/server"
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

func main() {
	s := server.New(context.Background())

	path, handler := mv1connect.NewMServiceHandler(s)
	mux := http.NewServeMux()
	mux.Handle(path, handler)

	log.Fatal(http.ListenAndServe(*listenAddr, h2c.NewHandler(mux, &http2.Server{})))
}
