package youtube

import (
	"fmt"
	"net/url"
	"os"

	"golang.org/x/net/context"
	"google.golang.org/api/option"
	yt "google.golang.org/api/youtube/v3"
)

type Client struct {
	APIKey   string
	Channels map[string][]Video
}

type Video struct {
	ID          string
	URL         string
	ParsedTitle string
	yt.SearchResultSnippet
}

func (c *Client) GetLatestVideosFromChannel(channelName string) ([]Video, error) {
	s, err := yt.NewService(context.Background(), option.WithAPIKey(os.Getenv("M_YT_API_KEY")))
	if err != nil {
		return nil, err
	}

	chanReq := s.Channels.List([]string{"id"})
	chanReq.ForUsername(channelName)
	channel, err := chanReq.Do()
	if err != nil {
		return nil, err
	}

	req := s.Search.List([]string{"snippet"})
	req.Order("date")
	req.MaxResults(20)
	req.ChannelId(channel.Items[0].Id)

	vs, err := req.Do()
	if err != nil {
		return nil, err
	}

	videos := make([]Video, 0, len(vs.Items))
	for _, v := range vs.Items {
		videoTitle, err := url.QueryUnescape(v.Snippet.Title)
		if err != nil {
			fmt.Printf("failed to query escape video: %s\n", err)
		}

		videos = append(videos, Video{
			ID:                  v.Id.VideoId,
			URL:                 fmt.Sprintf("https://youtube.com?v=%s", v.Id.VideoId),
			ParsedTitle:         videoTitle,
			SearchResultSnippet: *v.Snippet,
		})
	}

	return videos, nil
}
