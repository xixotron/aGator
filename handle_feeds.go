package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/xixotron/aGator/internal/database"
)

func handleAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("Usage: %s <name> <url>", cmd.Name)
	}
	name := cmd.Args[0]
	url := cmd.Args[1]

	ctx := context.Background()
	feed, err := s.db.CreateFeed(ctx, database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed: %w", err)
	}

	feedFollow, err := s.db.CreateFollowFeed(ctx, database.CreateFollowFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		FeedID:    feed.ID,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("couldnt' create feed follow: %w", err)
	}

	fmt.Println("feed added succesfully:")
	printFeed(feed, user)
	fmt.Println()
	fmt.Println("feed followed successfuly:")
	printFeedFollow(feedFollow)
	fmt.Println("=====================================")
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

	fmt.Println("Registered feeds:")
	users := make(map[uuid.UUID]database.User)
	for _, feed := range feeds {
		user, ok := users[feed.UserID]
		if !ok {
			user, err = s.db.GetUserByID(ctx, feed.UserID)
			if err != nil {
				return fmt.Errorf("couldn't find username: %w", err)
			}
			users[feed.UserID] = user
		}

		printFeed(feed, user)
		fmt.Println("=====================================")
	}

	return nil
}

func printFeed(f database.Feed, user database.User) {
	fmt.Printf(" * ID:       %v\n", f.ID)
	fmt.Printf(" * Name:     %v\n", f.Name)
	fmt.Printf(" * Created:  %v\n", f.CreatedAt)
	fmt.Printf(" * Updated:  %v\n", f.UpdatedAt)
	fmt.Printf(" * URL:      %v\n", f.Url)
	fmt.Printf(" * UserName: %v\n", user.Name)
}
