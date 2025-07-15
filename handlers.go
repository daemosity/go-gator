package main

import (
	"context"
	"fmt"
	"time"

	"github.com/daemosity/go-gator/internal/database"
	"github.com/google/uuid"
)

func getCommands() commands {
	commands := initCommands()
	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)
	commands.register("reset", handlerReset)
	commands.register("users", handlerUsers)

	return commands
}

func handlerLogin(s *state, cmd command) error {
	if !cmd.hasArgs() {
		return fmt.Errorf("error: %s requires one (1) argument: [username]", cmd.name)
	}

	ctx := context.Background()
	givenUser := cmd.args[0]

	if _, err := s.db.GetUser(ctx, givenUser); err != nil {
		return fmt.Errorf("error: username %s is not registered, use 'register [username]' command", givenUser)
	}

	if err := s.config.SetUser(givenUser); err != nil {
		return err
	}

	fmt.Printf("INFO: %s has been set as current user.\n", givenUser)

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if !cmd.hasArgs() {
		return fmt.Errorf("error: %s requires one (1) argument: [username]", cmd.name)
	}

	ctx := context.Background()
	userName := cmd.args[0]

	_, err := s.db.GetUser(ctx, userName)
	if err == nil {
		return fmt.Errorf("error: username %s already registered, use 'login [username]' command", userName)
	}

	entries := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      userName,
	}

	user, err := s.db.CreateUser(ctx, entries)
	if err != nil {
		return err
	}

	if err := s.config.SetUser(userName); err != nil {
		return err
	}

	fmt.Printf("INFO: User %s registered in system\n", userName)
	fmt.Printf("%v\n", user)
	return nil
}

func handlerReset(s *state, cmd command) error {
	ctx := context.Background()

	if err := s.db.DeleteAllUsers(ctx); err != nil {
		return err
	}

	fmt.Printf("INFO: users table has been reset\n")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	ctx := context.Background()

	users, err := s.db.ListAllUsers(ctx)
	if err != nil {
		return err
	}
	current_user := s.config.Current_user_name

	for _, user := range users {
		if user.Name == current_user {
			fmt.Printf("%s (current)\n", user.Name)
		} else {
			fmt.Printf("%s\n", user.Name)
		}
	}
	return nil
}
