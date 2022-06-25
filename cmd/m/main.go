package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bufbuild/connect-go"
	v1 "github.com/catouc/m/internal/m/v1"
	"github.com/catouc/m/internal/m/v1/mv1connect"
)

func main() {
	c := mv1connect.NewMServiceClient(http.DefaultClient, "http://127.0.0.1:8080")

	posts, err := c.ListNewBlogPosts(context.Background(), connect.NewRequest(&v1.ListNewBlogPostRequest{}))
	if err != nil {
		panic(err)
	}

	for _, p := range posts.Msg.Posts {
		fmt.Println(p.Title)
		fmt.Println(p.Content)
		fmt.Println("----------")
	}
}
