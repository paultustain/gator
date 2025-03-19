package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/paultustain/gator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("no arguments")
	}
	_, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	err = s.cfg.SetUser(cmd.args[0])

	if err != nil {
		return errors.New("failed to set user")
	}
	fmt.Printf("User has been set\n")
	return nil
}

func handlerRegister(s *state, cmd command) error {

	if len(cmd.args) == 0 {
		return errors.New("no arguments")
	}
	username := cmd.args[0]

	_, err := s.db.GetUser(context.Background(), username)
	if err == nil {
		return errors.New("duplicated user")
	}

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	})

	if err != nil {

		return err
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return err
	}
	fmt.Println("User registered")
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.Reset(context.Background())
	if err != nil {
		return err
	}

	fmt.Println("successfully reset db")
	return nil
}
