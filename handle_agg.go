package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/xixotron/aGator/internal/database"
	"github.com/xixotron/aGator/internal/rss"
)

func handleAgg(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Usage: %s <time_between_reqests>", cmd.Name)
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("invalid duration: %w", err)
	}

	fmt.Printf("Collection feeds every %v...\n", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)
	defer ticker.Stop()
	for ; ; <-ticker.C {
		err := scrapeFeeds(s, user, timeBetweenRequests)
		if err != nil {
			return err
		}
	}
}

func scrapeFeeds(s *state, user database.User, timeBetweenRequests time.Duration) error {
	ctx := context.Background()
	feed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		return fmt.Errorf("couldn't get the next feed to fetch: %w", err)
	}

	_, err = s.db.MarkFeedFetched(ctx, database.MarkFeedFetchedParams{
		UpdatedAt: time.Now().UTC(),
		ID:        feed.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't mark the feed as read: %w", err)
	}

	content, err := rss.FetchFeed(
		ctx,
		feed.Url,
	)
	if err != nil {
		log.Printf("couldn't fetch feed: %v\n", err)
		return nil
	}
	printFeedContent(content, feed.Name)
	return nil
}

func printFeedContent(feed *rss.RSSFeed, name string) {
	fmt.Println()
	fmt.Printf("Feed:  %v\n", name)
	fmt.Printf("Title: %v\n", feed.Channel.Title)
	fmt.Printf("Link: %v\n", feed.Channel.Link)
	fmt.Println("Description:")
	fmt.Println(feed.Channel.Description)
	// fmt.Println()
	fmt.Println("Posts:")
	for _, post := range feed.Channel.Item {
		fmt.Printf(" * Title: %v\n", post.Title)
		//fmt.Printf("Link: %v\n", post.Link)
		//fmt.Println("Description:")
		//fmt.Println(post.Description)
		// fmt.Println("-------------------------------------")
	}
	fmt.Println("=====================================")
	// fmt.Println()
}
