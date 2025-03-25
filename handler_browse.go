package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/paultustain/gator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	var limit int
	var err error

	if len(cmd.args) != 0 {
		limit, err = strconv.Atoi(cmd.args[0])
		if err != nil {
			return errors.New("failed to convert to integer")
		}
	} else {
		limit = 2
	}

	posts, err := s.db.GetPostsForUser(
		context.Background(),
		database.GetPostsForUserParams{
			UserID: uuid.NullUUID{
				UUID:  user.ID,
				Valid: true,
			},
			Limit: int32(limit),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to get posts: %w", err)
	}
	fmt.Println(posts)
	for _, post := range posts {
		fmt.Printf("%s from %s\n", post.PublishedAt, post.FeedName)
		fmt.Printf("--- %s ---\n", post.Title)
		fmt.Printf("    %v\n", post.Description)
		fmt.Printf("Link: %s\n", post.Url)
		fmt.Println("=====================================")
	}
	return nil
}
