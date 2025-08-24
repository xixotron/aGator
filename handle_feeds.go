package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/xixotron/aGator/internal/database"
)

func handleFeedFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Usage: %s <url>", cmd.Name)
	}
	url := cmd.Args[0]

	ctx := context.Background()
	feed, err := s.db.GetFeedByURL(ctx, url)
	if err != nil {
		return fmt.Errorf("couldn't find feed with provided url: %w", err)
	}

	feedFollow, err := s.db.CreateFollowFeed(ctx, database.CreateFollowFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID:    feed.ID,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("coudn't create feed follow: %w", err)
	}

	fmt.Println("Created feed follow:")
	printFeedFollow(feedFollow)

	return nil
}

func handleFeedsFollowing(s *state, cmd command, user database.User) error {
	feedsFollows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("couldn't find feeds followed by user: %w", err)
	}

	if len(feedsFollows) == 0 {
		fmt.Println("User follows no feeds")
		return nil
	}

	fmt.Println("Feeds being followed:")
	for _, feedFollow := range feedsFollows {
		printUserFeedFollow(feedFollow)
		fmt.Println()
		fmt.Println("=====================================")
	}
	return nil
}

func handleAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("Usage: %s <name> <url>", cmd.Name)
	}
	name := cmd.Args[0]
	url := cmd.Args[1]

	ctx := context.Background()
	feed, err := s.db.CreateFeed(ctx, database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed: %w", err)
	}

	feedFollow, err := s.db.CreateFollowFeed(ctx, database.CreateFollowFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
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

func printFeedFollow(ff database.CreateFollowFeedRow) {
	fmt.Printf(" * ID:       %v\n", ff.ID)
	fmt.Printf(" * Created:  %v\n", ff.CreatedAt)
	fmt.Printf(" * Updated:  %v\n", ff.UpdatedAt)
	fmt.Printf(" * FeedName: %v\n", ff.FeedName)
	fmt.Printf(" * UserName: %v\n", ff.UserName)
}

func printUserFeedFollow(ff database.GetFeedFollowsForUserRow) {
	fmt.Printf(" * ID:       %v\n", ff.ID)
	fmt.Printf(" * Created:  %v\n", ff.CreatedAt)
	fmt.Printf(" * Updated:  %v\n", ff.UpdatedAt)
	fmt.Printf(" * FeedName: %v\n", ff.FeedName)
	fmt.Printf(" * UserName: %v\n", ff.UserName)
}
