package main

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/paultustain/gator/internal/database"
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

	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		feedURL,
		nil,
	)

	if err != nil {
		return &RSSFeed{}, err
	}

	req.Header.Set("User-Agent", "gator")

	client := http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return &RSSFeed{}, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var rssFeed RSSFeed
	err = xml.Unmarshal(
		body,
		&rssFeed,
	)
	if err != nil {
		return &RSSFeed{}, err
	}

	return &rssFeed, nil
}

func scrapeFeeds(s *state) error {
	next_feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("Found a feed to fetch!")

	err = s.db.MarkFetchedFeed(context.Background(), next_feed.ID)

	if err != nil {
		return err
	}

	rss, err := fetchFeed(context.Background(), next_feed.Url)

	if err != nil {
		return err
	}

	for _, rssItem := range rss.Channel.Item {

		publish_date, err := time.Parse(time.RFC3339, rssItem.PubDate)

		if err != nil {
			fmt.Printf("Warning: couldn't parse date %s: %v\n", rssItem.PubDate, err)
			publish_date = time.Now() // fallback
		}

		_, err = s.db.CreatePost(
			context.Background(),
			database.CreatePostParams{
				ID:          uuid.New(),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				Title:       html.UnescapeString(rssItem.Title),
				Url:         rssItem.Link,
				Description: html.UnescapeString(rssItem.Description),
				PublishedAt: publish_date,
				FeedID: uuid.NullUUID{
					UUID:  next_feed.ID,
					Valid: true,
				},
			},
		)

		if err != nil {
			if strings.Contains(err.Error(), "duplicate") {
				// Ignore duplicate errors as specified
				fmt.Printf("Skipping duplicate post: %s\n", rssItem.Title)
				continue
			}

			fmt.Println("couldn't add post to posts - skipping: %w", err)
		}
	}

	return nil
}

func handlerAgg(s *state, cmd command) error {

	if len(cmd.args) < 1 {
		return errors.New("not enough arguements")
	}
	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])

	if err != nil {
		return err
	}

	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}

}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return errors.New("not enough arguments")
	}

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID: uuid.NullUUID{
			UUID:  user.ID,
			Valid: true,
		},
	})

	if err != nil {
		return fmt.Errorf("couldn't fetch feed: %w", err)
	}

	s.db.CreateFeedFollow(context.Background(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID: uuid.NullUUID{
				UUID:  user.ID,
				Valid: true,
			},
			FeedID: uuid.NullUUID{
				UUID:  feed.ID,
				Valid: true,
			},
		})
	fmt.Println(feed)

	return nil

}

func handlerFeeds(s *state, cmd command) error {
	feed, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	fmt.Println(feed)

	return nil
}
