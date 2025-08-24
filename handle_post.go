package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/xixotron/aGator/internal/database"
)

const defaultPostLimit int32 = 2

func handleBrowse(s *state, cmd command, user database.User) error {
	postsLimit := defaultPostLimit
	if len(cmd.Args) != 1 {
		fmt.Printf("Usage: %s [num_posts default=%v]\n", cmd.Name, postsLimit)
	} else {
		if tmp, err := strconv.Atoi(cmd.Args[0]); err == nil {
			postsLimit = int32(tmp)
		}
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  postsLimit,
	})
	if err != nil {
		return fmt.Errorf("couldn't get posts: %w", err)
	}

	if len(posts) == 0 {
		fmt.Println("No posts found!")
		return nil
	}

	feedsFollow, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("couldn't get user followd feeds: %w", err)
	}
	feedsNames := make(map[uuid.UUID]string)
	for _, feed := range feedsFollow {
		feedsNames[feed.FeedID] = feed.FeedName
	}

	lastFeedID := uuid.Nil
	for _, post := range posts {
		if lastFeedID != post.FeedID {
			lastFeedID = post.FeedID
			fmt.Println()
			fmt.Println("-------------------------------------")
			fmt.Printf("Post From: %v\n", feedsNames[post.FeedID])
		}
		printPost(post)
	}
	return nil
}

func printPost(post database.Post) {
	fmt.Printf("Title:       %v\n", post.Title.String)
	fmt.Printf("Published:   %v\n", post.PublishedAt)
	fmt.Printf("Link:        %v\n", post.Url)
	if post.Description.Valid {
		fmt.Printf("Description: \n")
		fmt.Println(post.Description.String[:min(100, len(post.Description.String))])
	}
}
