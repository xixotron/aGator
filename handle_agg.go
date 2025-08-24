package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
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
		err := scrapeFeeds(s.db, timeBetweenRequests)
		if err != nil {
			return err
		}
	}
}

func scrapeFeeds(db *database.Queries, timeBetweenRequests time.Duration) error {
	minDate := time.Now().Add(-timeBetweenRequests)
	ctx := context.Background()

	for {
		feed, err := db.GetNextFeedToFetch(ctx)
		if err != nil {
			return fmt.Errorf("couldn't get the next feed to fetch: %w", err)
		}
		if feed.LastFetchedAt.Valid && feed.LastFetchedAt.Time.After(minDate) {
			log.Println("All feeds scrapped!")
			return nil
		}

		_, err = db.MarkFeedFetched(ctx, database.MarkFeedFetchedParams{
			UpdatedAt: time.Now().UTC(),
			ID:        feed.ID,
		})
		if err != nil {
			return fmt.Errorf("couldn't mark the feed %s as fetched: %w", feed.Name, err)
		}

		scrapeFeed(db, feed)
	}
}

func scrapeFeed(db *database.Queries, feed database.Feed) {
	content, err := rss.FetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Printf("couldn't fetch feed %s: %v\n", feed.Name, err)
		return
	}

	log.Printf("Scraping feed: %v", feed.Name)
	log.Printf("Found: %v posts", len(content.Channel.Item))

	for _, post := range content.Channel.Item {
		err = savePost(db, post, feed.ID)
		if err != nil {
			log.Printf("Error saving post: %v", err)
		}
	}
}

func savePost(db *database.Queries, post rss.RSSItem, feedID uuid.UUID) error {
	if len(post.Title) == 0 {
		if len(post.Description) == 0 {
			log.Println("Post Title and Description empty")
			return nil
		}
		post.Title = post.Description[:min(20, len(post.Description))] + "..."
	}

	if len(post.Link) == 0 {
		log.Println("Post Link empty")
		return nil
	}

	if len(post.PubDate) == 0 {
		log.Println("Post without date")
		return nil
	}

	date, err := parsePostDate(post.PubDate)
	if err != nil {
		log.Println(err)
		return nil
	}

	_, err = db.CreatePost(context.Background(), database.CreatePostParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Title: sql.NullString{
			String: post.Title,
			Valid:  len(post.Title) > 0,
		},
		Url: post.Link,
		Description: sql.NullString{
			String: post.Description,
			Valid:  len(post.Description) > 0,
		},
		PublishedAt: date,
		FeedID:      feedID,
	})

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint \"posts_url_key\"") {
			return nil
		}
	}

	return err
}

func parsePostDate(s string) (time.Time, error) {
	dateFormats := [...]string{
		"Mon, _2 Jan 2006 15:04:05 -0700",
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822Z,
		time.RFC822,
		time.RFC850,
		time.RFC3339,
	}
	for _, format := range dateFormats {
		date, err := time.Parse(format, s)
		if err == nil {
			return date, nil
		}
		log.Printf("date format error with: %+v\n", format)
	}
	return time.Time{}, fmt.Errorf("Unknown date format for pubDate %q", s)
}
