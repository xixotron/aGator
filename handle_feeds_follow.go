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
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
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

func handleUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Usage: %s <url>", cmd.Name)
	}
	url := cmd.Args[0]

	ctx := context.Background()

	feed, err := s.db.GetFeedByURL(ctx, url)
	if err != nil {
		return fmt.Errorf("couldn't find feed: %w", err)
	}

	err = s.db.DeleteFollowFeed(ctx, database.DeleteFollowFeedParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("coudn't delete feed follow: %w", err)
	}

	fmt.Println("Successfuly unfollowed:")
	fmt.Printf(" * User: %v\n", user.Name)
	fmt.Printf(" * Feed: %v\n", feed.Name)
	return nil
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
