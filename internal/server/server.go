package server

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/catouc/m/internal/feed"
	"github.com/catouc/m/internal/hntop"
	v1 "github.com/catouc/m/internal/m/v1"
	"github.com/catouc/m/internal/youtube"
	"github.com/rs/zerolog"
	"golang.org/x/net/context"
)

type Config struct {
	MuseumConfig feed.Config `conf:"required" yaml:"museum"`
}

type Server struct {
	Config   Config
	Logger   zerolog.Logger
	Museum   *feed.Museum
	HNClient *hntop.Client
	YTClient *youtube.Client
}

func New(ctx context.Context, cfg Config) *Server {
	return &Server{
		Config:   cfg,
		Logger:   zerolog.New(os.Stdout),
		Museum:   feed.NewMuseum(time.Hour),
		HNClient: hntop.New(ctx),
	}
}

func (s *Server) Init() error {
	s.Logger.Info().Msg("Initialising server")
	s.Logger.Info().Msg("Initialising museum")
	err := s.Museum.Init()
	if err != nil {
		return fmt.Errorf("failed to init museum: %w", err)
	}
	s.Logger.Info().Msg("Initialising museum done")

	return nil
}

type RegisterFeedInMuseumBody struct {
	FeedURL string `json:"feed_url"`
}

func (s *Server) RegisterBlog(ctx context.Context, req *connect.Request[v1.RegisterBlogRequest]) (*connect.Response[v1.RegisterBlogResponse], error) {
	err := s.Museum.Register(req.Msg.FeedURL)
	if err != nil {
		s.Logger.Error().Err(err).Msg("failed to register feed")
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	return connect.NewResponse(&v1.RegisterBlogResponse{}), nil
}

func (s *Server) ListVideosForChannel(_ context.Context, req *connect.Request[v1.YoutubeChanneListRequest]) (*connect.Response[v1.YoutubeVideoListResponse], error) {
	v, err := s.YTClient.GetLatestVideosFromChannel(req.Msg.ChannelName)
	if err != nil {
		return nil, err
	}

	response := connect.NewResponse(&v1.YoutubeVideoListResponse{
		Videos: v,
	})
	return response, nil
}

func (s *Server) ListNewBlogPosts(_ context.Context, req *connect.Request[v1.ListNewBlogPostRequest]) (*connect.Response[v1.ListNewBlogPostResponse], error) {
	today := time.Now().Add(-24 * time.Hour)

	// This is all based on the assumption that the items in the feed struct are ordered in desc order by published date
	// to save us from iterating through the entire backlog of posts everytime.
	posts := make([]*v1.BlogPost, 0)
	for _, f := range s.Museum.Feeds {
		for _, i := range f.Items {
			if i.PublishedParsed == nil {
				s.Logger.Error().
					Str("BlogTitle", i.Title).
					Str("PublishedString", i.Published).
					Msg("PublishedParsed date of blog post is nil!")
				continue
			}
			if dateEqual(*i.PublishedParsed, today) {
				posts = append(posts, &v1.BlogPost{
					Title:   i.Title,
					Content: i.Content,
				})
			}
			break
		}
	}

	response := connect.NewResponse(&v1.ListNewBlogPostResponse{Posts: posts})
	return response, nil
}

func (s *Server) ListVideosForCategory(_ context.Context, req *connect.Request[v1.YoutubeCategoryListRequest]) (*connect.Response[v1.YoutubeVideoListResponse], error) {
	return nil, errors.New("not implemented")
}

func (s *Server) GetHNFrontpage(ctx context.Context, _ *connect.Request[v1.HNFrontpageRequest]) (*connect.Response[v1.HNFrontpageResponse], error) {
	frontPage, err := s.HNClient.GetHNTop30Stories(ctx)
	if err != nil {
		return nil, err
	}

	fmt.Println(frontPage[0].URL)

	response := connect.NewResponse(&v1.HNFrontpageResponse{Stories: frontPage})
	return response, nil
}

func dateEqual(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}
