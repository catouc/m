package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/Deichindianer/m/internal/server"
)

var (
	listenAddr = flag.String("listenAddr", ":8080", "Listen address for the server")
)

func init() {
	flag.Parse()
}

func main() {
	s := server.New()

	log.Fatal(http.ListenAndServe(*listenAddr, s))
}
