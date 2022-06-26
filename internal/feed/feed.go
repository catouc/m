package feed

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/rs/zerolog"
)

type Config struct {
	UpdateInterval time.Duration `conf:"required,env:M_MUSEUM_UPDATE_INTERVAL" yaml:"updateInterval"`
	SaveLocation   string        `conf:"env:M_MUSEUM_SAVELOCATION" yaml:"saveLocation"`
}

type Museum struct {
	Logger zerolog.Logger
	Config Config

	// where the key is the feed URL
	Feeds       map[string]*gofeed.Feed
	FeedURLFile *os.File
}

type Option func(m *Museum)

func WithSaveLocation(location string) Option {
	return func(m *Museum) {
		m.Config.SaveLocation = location
	}
}

func NewMuseum(updateInterval time.Duration, opts ...Option) *Museum {
	m := Museum{
		Config: Config{
			UpdateInterval: updateInterval,
		},
		Feeds: make(map[string]*gofeed.Feed),
	}

	for _, o := range opts {
		o(&m)
	}

	if m.Config.SaveLocation == "" {
		m.Config.SaveLocation = path.Join(os.TempDir(), "feed-museum")
	}

	return &m
}

func (m *Museum) Init() error {
	err := os.MkdirAll(m.Config.SaveLocation, 0755)
	if err != nil {
		return fmt.Errorf("failed to create temporary dir for museum: %w", err)
	}

	err = os.Chown(m.Config.SaveLocation, os.Getuid(), os.Getgid())
	if err != nil {
		return fmt.Errorf("failed to chown temporary dir for museum: %w", err)
	}

	m.FeedURLFile, err = os.OpenFile(path.Join(m.Config.SaveLocation, "feeds.txt"), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return fmt.Errorf("failed to open save loc writer: %w", err)
	}

	feedURLScanner := bufio.NewScanner(m.FeedURLFile)
	for feedURLScanner.Scan() {
		err = m.Register(feedURLScanner.Text())
		if err != nil {
			m.Logger.Error().
				Err(err).
				Str("FeedURL", feedURLScanner.Text()).
				Msg("failed to register feed")
			continue
		}
	}

	return nil
}

func (m *Museum) Register(feedURL string) error {
	_, exists := m.Feeds[feedURL]
	if exists {
		return nil
	}

	_, err := m.FeedURLFile.WriteString(feedURL + "\n")
	if err != nil {
		return fmt.Errorf("failed to write URL to disk: %w", err)
	}

	fp := gofeed.NewParser()

	feed, err := fp.ParseURL(feedURL)
	if err != nil {
		return fmt.Errorf("failed to parse feed %s: %w", feedURL, err)
	}

	m.Feeds[feedURL] = feed
	return nil
}
