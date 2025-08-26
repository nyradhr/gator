package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/nyradhr/gator/internal/database"
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
	for _, item := range rssFeed.Channel.Item {
		pubDate, err := parseDate(item.PubDate)
		if err != nil {
			log.Printf("Error parsing publishing date for item in feed %s: %v", dbFeed.Name, err)
			return
		}
		params := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: toNullString(item.Description),
			PublishedAt: toNullTime(pubDate),
			FeedID:      dbFeed.ID,
		}
		_, err = s.db.CreatePost(ctx, params)
		if err != nil {
			if isUniqueViolation(err) {
				continue
			}
			log.Printf("Error while creating post: %v", err)
		}
	}
}

var fmts = []string{
	time.RFC1123Z, time.RFC1123,
	time.RFC822Z, time.RFC822,
	time.RFC3339, "Mon, 02 Jan 2006 15:04:05 MST",
}

func parseDate(s string) (time.Time, error) {
	for _, f := range fmts {
		if t, err := time.Parse(f, s); err == nil {
			return t, nil
		}
	}
	return time.Now(), fmt.Errorf("unrecognized date: %q", s)
}

func toNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

func isUniqueViolation(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return string(pqErr.Code) == "23505"
	}
	return false
}

func toNullTime(t time.Time) sql.NullTime {
	if t.IsZero() {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: t, Valid: true}
}
