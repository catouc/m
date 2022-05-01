package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/Deichindianer/m/internal/server"
	"golang.org/x/net/context"
)

var (
	listenAddr = flag.String("listenAddr", ":8080", "Listen address for the server")
)

func init() {
	flag.Parse()
}

func main() {
	rootCtx := context.TODO()
	s := server.New(rootCtx)

	log.Fatal(http.ListenAndServe(*listenAddr, s))
}
