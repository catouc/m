package main

import (
	"fmt"

	"github.com/catouc/m/internal/youtube"
)

func main() {
	yt := youtube.Client{}
	videos, err := yt.GetLatestVideosFromChannel("LinusTechTips")
	if err != nil {
		panic(err)
	}

	for _, v := range videos {
		fmt.Printf("%s => %s\n", v.ParsedTitle, v.URL)
	}
}
