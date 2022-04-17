package server

import (
	"net/http"
	"os"
	"time"

	"github.com/Deichindianer/m/internal/feed"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type Server struct {
	listenAddr string
	mux        *echo.Echo
	logger     zerolog.Logger
	museum     *feed.Museum
}

func New() *Server {
	s := Server{}

	s.logger = zerolog.New(os.Stdout)

	mux := echo.New()
	mux.GET("feeds/museum", s.ListFeedsInMuseum)
	mux.POST("/feeds/museum", s.RegisterFeedInMuseum)

	s.mux = mux
	s.museum = feed.NewMuseum(time.Hour)
	return &s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

type RegisterFeedInMuseumBody struct {
	FeedURL string `json:"feed_url"`
}

func (s *Server) RegisterFeedInMuseum(ctx echo.Context) error {
	var body RegisterFeedInMuseumBody
	err := ctx.Bind(&body)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to bind body")
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	err = s.museum.Register(body.FeedURL)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to register")
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return ctx.NoContent(http.StatusCreated)
}

type ListFeedsInMuseumResponse struct {
	Feeds []string `json:"feeds"`
}

func (s *Server) ListFeedsInMuseum(ctx echo.Context) error {
	resp := ListFeedsInMuseumResponse{
		Feeds: make([]string, 0),
	}

	for key, _ := range s.museum.Feeds {
		resp.Feeds = append(resp.Feeds, key)
	}

	return ctx.JSON(http.StatusOK, resp)
}
