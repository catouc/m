package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ardanlabs/conf/v2"
	"github.com/ardanlabs/conf/v2/yaml"
	"github.com/catouc/m/internal/m/v1/mv1connect"
	"github.com/catouc/m/internal/server"
	"github.com/rs/zerolog"
	"golang.org/x/net/context"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

var (
	listenAddr = flag.String("listenAddr", ":8080", "Listen address for the server")
	cfgFile    = flag.String("cfg", "", "Specify a yaml config file")

	logger zerolog.Logger
)

func init() {
	flag.Parse()
}

func main() {
	cfg := server.Config{}
	help, err := parseConfig(&cfg, *cfgFile)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return
		}
	}
	s := server.New(context.Background(), cfg)
	err = s.Init()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initialise server")
	}

	path, handler := mv1connect.NewMServiceHandler(s)
	mux := http.NewServeMux()
	mux.Handle(path, handler)

	log.Fatal(http.ListenAndServe(*listenAddr, h2c.NewHandler(mux, &http2.Server{})))
}

func parseConfig(cfg *server.Config, cfgFile string) (string, error) {
	if cfgFile != "" {
		f, err := os.Open(cfgFile)
		if err != nil {
			return "", fmt.Errorf("failed to open config file: %w", err)
		}
		help, err := conf.Parse("", &cfg, yaml.WithReader(f))
		if err != nil {
			return help, err
		}
	} else {
		help, err := conf.Parse("", &cfg)
		if err != nil {
			return help, err
		}
	}

	return "", nil
}
