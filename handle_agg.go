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
	return nil
}
