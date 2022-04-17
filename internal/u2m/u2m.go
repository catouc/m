package u2m

import (
	"fmt"
	"net/url"

	md "github.com/JohannesKaufmann/html-to-markdown"
)

func Convert(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("input is not a valid URL: %w", err)
	}

	converter := md.NewConverter(u.Host, true, nil)

	markdown, err := converter.ConvertURL(u.String())
	if err != nil {
		return "", fmt.Errorf("failed to convert to markdown: %w", err)
	}

	return markdown, nil
}
