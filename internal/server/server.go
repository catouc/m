package server

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/catouc/m/internal/feed"
	"github.com/catouc/m/internal/hntop"
	v1 "github.com/catouc/m/internal/m/v1"
	"github.com/catouc/m/internal/youtube"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"golang.org/x/net/context"
)

type Server struct {
	logger   zerolog.Logger
	museum   *feed.Museum
	hnClient *hntop.Client
	ytClient *youtube.Client
}

func New(ctx context.Context) *Server {
	return &Server{
		logger:   zerolog.New(os.Stdout),
		museum:   feed.NewMuseum(time.Hour),
		hnClient: hntop.New(ctx),
	}
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

func (s *Server) ListVideosForChannel(_ context.Context, req *connect.Request[v1.YoutubeChanneListRequest]) (*connect.Response[v1.YoutubeVideoListResponse], error) {
	v, err := s.ytClient.GetLatestVideosFromChannel(req.Msg.ChannelName)
	if err != nil {
		return nil, err
	}

	response := connect.NewResponse(&v1.YoutubeVideoListResponse{
		Videos: v,
	})
	return response, nil
}

func (s *Server) ListNewBlogPosts(_ context.Context, req *connect.Request[v1.ListNewBlogPostRequest]) (*connect.Response[v1.ListNewBlogPostResponse], error) {
	return nil, errors.New("not implemented")
}

func (s *Server) ListVideosForCategory(_ context.Context, req *connect.Request[v1.YoutubeCategoryListRequest]) (*connect.Response[v1.YoutubeVideoListResponse], error) {
	return nil, errors.New("not implemented")
}

func (s *Server) GetHNFrontpage(ctx context.Context, _ *connect.Request[v1.HNFrontpageRequest]) (*connect.Response[v1.HNFrontpageResponse], error) {
	frontPage, err := s.hnClient.GetHNTop30Stories(ctx)
	if err != nil {
		return nil, err
	}

	fmt.Println(frontPage[0].URL)

	response := connect.NewResponse(&v1.HNFrontpageResponse{Stories: frontPage})
	return response, nil
}
