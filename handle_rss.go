package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/xixotron/aGator/internal/database"
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

func handleAddFeed(s *state, cmd command) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("Usage: %s <feedname> <url>", cmd.Name)
	}
	ctx := context.Background()
	user, err := s.db.GetUser(ctx, s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("couldn't find active user: %w", err)
	}

	feed, err := s.db.CreateFeed(ctx,
		database.CreateFeedParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      cmd.Args[0],
			Url:       cmd.Args[1],
			UserID:    user.ID,
		})
	if err != nil {
		return fmt.Errorf("coudn't add feed: %w", err)
	}

	fmt.Println("feed added succesfully:")
	printFeed(feed)
	return nil
}

func printFeed(f database.Feed) {
	fmt.Printf(" * ID:      %v\n", f.ID)
	fmt.Printf(" * Name:    %v\n", f.Name)
	fmt.Printf(" * URL:     %v\n", f.Url)
	fmt.Printf(" * UserID:  %v\n", f.UserID)
}
