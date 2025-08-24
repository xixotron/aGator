package main

import (
	"context"
	"fmt"

	"github.com/xixotron/aGator/internal/rss"
)

func handleAgg(s *state, cmd command) error {
	feed, err := rss.FetchFeed(
		context.Background(),
		"https://www.wagslane.dev/index.xml")
	//		"http://127.0.0.1:8000/index.xml")
	if err != nil {
		return fmt.Errorf("couldn't fetch feed: %w", err)
	}

	fmt.Printf("Feed: %+v\n", feed)
	/*	fmt.Println("Read from feed successfully:")
		fmt.Printf("Title:        %v\n", feed.Channel.Title)
		fmt.Printf("Description:  %v\n", feed.Channel.Description)
		fmt.Printf("Link:         %v\n", feed.Channel.Link)
		for i, item := range feed.Channel.Item[:4] {
			fmt.Printf("\nEntry %v:\n", i)
			fmt.Printf("Title:        %v\n", item.Title)
			fmt.Printf("Description:  %v\n", item.Description)
			fmt.Printf("Link:         %v\n", item.Link)
			fmt.Printf("Publish Date: %v\n", item.PubDate)
		} */
	return nil
}
