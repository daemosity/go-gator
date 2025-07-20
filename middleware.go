package main

import (
	"context"

	"github.com/daemosity/go-gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		ctx := context.Background()
		current_user := s.config.Current_user_name
		user, err := s.db.GetUser(ctx, current_user)
		if err != nil {
			return err
		}

		return handler(s, cmd, user)
	}
}
