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
	return nil
}

func handleAddFeed(s *state, cmd command) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("Usage: %s <name> <url>", cmd.Name)
	}
	name := cmd.Args[0]
	url := cmd.Args[1]

	ctx := context.Background()
	user, err := s.db.GetUser(ctx, s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("couldn't find active user: %w", err)
	}

	feed, err := s.db.CreateFeed(ctx, database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("coudn't create feed: %w", err)
	}

	fmt.Println("feed added succesfully:")
	printFeed(feed)
	return nil
}

func handleListFeeds(s *state, cmd command) error {
	ctx := context.Background()
	feeds, err := s.db.GetFeeds(ctx)
	if err != nil {
		return fmt.Errorf("coudn't list feeds: %w", err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found")
		return nil
	}

	users := make(map[uuid.UUID]string)

	for _, feed := range feeds {
		_, ok := users[feed.UserID]
		if ok {
			continue
		}
		user, err := s.db.GetUserName(ctx, feed.UserID)
		if err != nil {
			return fmt.Errorf("couldn't find username: %w", err)
		}
		users[feed.UserID] = user
	}

	fmt.Println("Registered feeds:")
	for _, feed := range feeds {
		printFeed(feed)
		fmt.Printf(" * UserName: %v\n", users[feed.UserID])
		fmt.Println()
		fmt.Println("=====================================")
	}
	return nil
}
func printFeed(f database.Feed) {
	fmt.Printf(" * ID:       %v\n", f.ID)
	fmt.Printf(" * Name:     %v\n", f.Name)
	fmt.Printf(" * Created:  %v\n", f.CreatedAt)
	fmt.Printf(" * Updated:  %v\n", f.UpdatedAt)
	fmt.Printf(" * URL:      %v\n", f.Url)
	fmt.Printf(" * UserID:   %v\n", f.UserID)
}
