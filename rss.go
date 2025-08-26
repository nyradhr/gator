package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
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

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch feed: status %d", res.StatusCode)
	}
	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	feed := RSSFeed{}
	err = xml.Unmarshal(resBody, &feed)
	if err != nil {
		return nil, err
	}
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
	}
	return &feed, nil
}

func scrapeFeeds(s *state) {
	ctx := context.Background()
	dbFeed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		log.Println("Feed to fetch not found:", err)
		return
	}

	err = s.db.MarkFeedFetched(ctx, dbFeed.ID)
	if err != nil {
		log.Printf("Error while marking feed %s as fetched: %v", dbFeed.Name, err)
		return
	}
	rssFeed, err := fetchFeed(ctx, dbFeed.Url)
	if err != nil {
		log.Printf("Error while collecting feed %s: %v", dbFeed.Name, err)
		return
	}
	log.Println("Showing posts from", dbFeed.Name)
	for _, item := range rssFeed.Channel.Item {
		fmt.Println(item.Title)
	}
	log.Printf("Feed %s collected, %v posts found", dbFeed.Name, len(rssFeed.Channel.Item))
}
