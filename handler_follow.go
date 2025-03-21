package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/paultustain/gator/internal/database"
)

func handlerFollow(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("not enough arguements")
	}

	feed, err := s.db.GetFeed(context.Background(), "https://hnrss.org/newest")
	if err != nil {
		return err
	}

	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
	}

	s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
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

	return nil
}

func handlerFollowing(s *state, cmd command) error {
	current_user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
	}
	follows, err := s.db.GetFeedFollowsUser(
		context.Background(),
		uuid.NullUUID{
			UUID:  current_user.ID,
			Valid: true,
		},
	)

	if err != nil {
		return err
	}

	fmt.Println(follows)

	return nil
}
