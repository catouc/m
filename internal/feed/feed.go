package feed

import (
	"fmt"
	"time"

	"github.com/mmcdole/gofeed"
)

type Museum struct {
	Tick time.Duration

	// where the key is the feed URL
	Feeds map[string]*gofeed.Feed
}

func NewMuseum(tick time.Duration) *Museum {
	return &Museum{
		Tick:  tick,
		Feeds: make(map[string]*gofeed.Feed),
	}
}

func (m *Museum) Register(feedURL string) error {
	_, exists := m.Feeds[feedURL]
	if exists {
		return nil
	}

	fp := gofeed.NewParser()

	feed, err := fp.ParseURL(feedURL)
	if err != nil {
		return fmt.Errorf("failed to parse feed %s: %w", feedURL, err)
	}

	m.Feeds[feedURL] = feed
	return nil
}
