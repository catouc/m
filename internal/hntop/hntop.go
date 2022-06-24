package hntop

import (
	"encoding/json"
	"fmt"
	v1 "github.com/catouc/m/internal/m/v1"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	"golang.org/x/net/context"
)

const (
	apiEndpoint = "https://hacker-news.firebaseio.com/v0/"
)

type Client struct {
	HTTPClient *http.Client
	GCTick     time.Duration

	mut   sync.RWMutex
	Cache map[int]*HNStory
}

func New(ctx context.Context) *Client {
	c := http.DefaultClient
	c.Timeout = 5 * time.Second

	client := Client{
		HTTPClient: c,
		GCTick:     10 * time.Second,
		Cache:      make(map[int]*HNStory),
	}

	go client.CollectGarbage(ctx)
	return &client
}

type HNStory struct {
	By           string    `json:"by"`
	Descendants  int       `json:"descendants"`
	ID           int       `json:"id"`
	Kids         []int     `json:"kids"`
	Score        int       `json:"score"`
	Time         int       `json:"time"`
	Title        string    `json:"title"`
	Type         string    `json:"type"`
	URL          string    `json:"url"`
	LastAccessed time.Time `json:"last_accessed"`
}

func (c *Client) GetHNTop30Stories(ctx context.Context) ([]*v1.HNStory, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiEndpoint+"topstories.json", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to construct http request: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call api: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read bodyBytes: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got non 200 response %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var body []int
	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal bodyBytes into []int: %w", err)
	}

	stories := make([]*v1.HNStory, 0)
	for _, id := range body[:30] {
		story, err := c.GetHNStory(ctx, id)
		if err != nil {
			// TODO: logging
			continue
		}

		fmt.Println(story.URL)

		tmp := v1.HNStory{
			Author: story.By,
			ID:     int32(story.ID),
			Title:  story.Title,
			URL:    story.URL,
		}

		stories = append(stories, &tmp)
	}

	return stories, nil
}

func (c *Client) GetHNStory(ctx context.Context, id int) (*HNStory, error) {
	c.mut.RLock()
	existingStory, exists := c.Cache[id]
	c.mut.RUnlock()
	if exists {
		return existingStory, nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiEndpoint+"item/"+strconv.Itoa(id)+".json", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to construct http request: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call api: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got non 200 response %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var hnStory HNStory
	err = json.Unmarshal(bodyBytes, &hnStory)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	hnStory.LastAccessed = time.Now()
	c.mut.Lock()
	c.Cache[id] = &hnStory
	c.mut.Unlock()

	return &hnStory, nil
}

func (c *Client) CollectGarbage(ctx context.Context) {
	ticker := time.NewTicker(c.GCTick)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.mut.Lock()
			for id, story := range c.Cache {
				if time.Now().Sub(story.LastAccessed) > 10*time.Second {
					delete(c.Cache, id)
				}
			}
			c.mut.Unlock()
		}
	}
}
