package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	connect_go "github.com/bufbuild/connect-go"
	mv1 "github.com/catouc/m/internal/m/v1"
	"github.com/catouc/m/internal/m/v1/mv1connect"
)

func main() {
	c := mv1connect.NewMServiceClient(http.DefaultClient, "http://127.0.0.1:8080")

	if len(os.Args) < 2 {
		fmt.Println("need at least one arg")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "read":
		posts, err := c.ListNewBlogPosts(context.TODO(), &connect_go.Request[mv1.ListNewBlogPostRequest]{})
		if err != nil {
			fmt.Printf("failed to list new blog posts: %s\n", err)
			os.Exit(1)
		}

		for _, p := range posts.Msg.Posts {
			fmt.Printf("%s by %s\n", p.Title, p.Author)
		}
	case "register":
		_, err := c.RegisterBlog(context.TODO(), connect_go.NewRequest(&mv1.RegisterBlogRequest{FeedURL: os.Args[2]}))
		if err != nil {
			fmt.Printf("failed to add new blog feed: %s\n", err)
			os.Exit(1)
		}
	}
}
