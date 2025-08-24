package rss

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
	"time"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "gator")
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	rssFeed := RSSFeed{}
	err = xml.Unmarshal(dat, &rssFeed)
	if err != nil {
		return nil, err
	}
	rssFeed.Unescape()

	return &rssFeed, nil
}

func (f *RSSFeed) Unescape() {
	f.Channel.Title = html.UnescapeString(f.Channel.Title)
	f.Channel.Link = html.UnescapeString(f.Channel.Link)
	f.Channel.Description = html.UnescapeString(f.Channel.Description)

	for _, item := range f.Channel.Item {
		item.Unescape()
	}
}

func (i *RSSItem) Unescape() {
	i.Title = html.UnescapeString(i.Title)
	i.Link = html.UnescapeString(i.Link)
	i.Description = html.UnescapeString(i.Description)
	i.PubDate = html.UnescapeString(i.PubDate)
}
