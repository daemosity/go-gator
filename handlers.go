package main

import (
	"context"
	"fmt"
	"net/url"
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
	commands.register("agg", handlerAgg)
	commands.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	commands.register("feeds", handlerFeeds)
	commands.register("follow", middlewareLoggedIn(handlerFollow))
	commands.register("following", middlewareLoggedIn(handlerFollowing))

	return commands
}

func handlerLogin(s *state, cmd command) error {
	if !cmd.hasArgs() {
		return fmt.Errorf("error: %s requires one (1) argument:\n\nUsage: %s [username]", cmd.name, cmd.name)
	}

	ctx := context.Background()
	givenUser := cmd.args[0]

	if _, err := s.db.GetUser(ctx, givenUser); err != nil {
		return fmt.Errorf("error: username %s is not registered, use:\n register [username]", givenUser)
	}

	if err := s.config.SetUser(givenUser); err != nil {
		return err
	}

	fmt.Printf("INFO: %s has been set as current user.\n", givenUser)

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if !cmd.hasArgs() {
		return fmt.Errorf("error: %s requires one (1) argument:\n\nUsage: %s [username]", cmd.name, cmd.name)
	}

	ctx := context.Background()
	userName := cmd.args[0]

	_, err := s.db.GetUser(ctx, userName)
	if err == nil {
		return fmt.Errorf("error: username %s already registered, use\n\nUsage: login [username]", userName)
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
	if cmd.hasArgs() {
		return fmt.Errorf("error: %s takes no arguments", cmd.name)
	}
	ctx := context.Background()

	if err := s.db.DeleteAllUsers(ctx); err != nil {
		return err
	}

	fmt.Printf("INFO: users table has been reset\n")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	if cmd.hasArgs() {
		return fmt.Errorf("error: %s takes no arguments", cmd.name)
	}

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

func handlerAgg(s *state, cmd command) error {
	// if !cmd.hasArgs() {
	// 	return fmt.Errorf("error: %s requires one (1) argument: [rss-url]", cmd.name)
	// }
	ctx := context.Background()

	// url_to_fetch := cmd.args[0]

	_, err := fetchFeed(ctx, "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("error: %s requires (2) arguments:\n\nUsage: %s [feed-name] [feed-url]", cmd.name, cmd.name)
	}
	ctx := context.Background()

	feedName, feedURL := cmd.args[0], cmd.args[1]
	parsedURL, err := url.Parse(feedURL)
	if err != nil {
		return fmt.Errorf("error parsing url: %s", feedURL)
	}

	feedEntry := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       parsedURL.String(),
		UserID:    user.ID,
	}

	feedInfo, err := s.db.CreateFeed(ctx, feedEntry)
	if err != nil {
		return fmt.Errorf("error: problem creating new feed: %w", err)
	}

	fmt.Printf("%v\n\n", feedInfo)

	cmd.name = "follow"
	cmd.args = []string{feedInfo.Url}
	err = handlerFollow(s, cmd, user)
	if err != nil {
		return err
	}

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	if cmd.hasArgs() {
		return fmt.Errorf("error: %s takes no arguments", cmd.name)
	}

	ctx := context.Background()

	feedResults, err := s.db.ListAllFeeds(ctx)
	if err != nil {
		return err
	}

	for _, feed := range feedResults {
		fmt.Printf("feed: %s\nurl: %s\nsubmitted by: %s\n\n", feed.FeedName, feed.FeedURL, feed.CreatedBy.String)
	}
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if !cmd.hasArgs() || len(cmd.args) > 1 {
		return fmt.Errorf("error: %s takes only one argument.\nUsage: %s [feed-url]", cmd.name, cmd.name)
	}

	ctx := context.Background()

	urlToFollow := cmd.args[0]

	parsedURL, err := url.Parse(urlToFollow)
	if err != nil {
		return fmt.Errorf("error parsing url: %s", urlToFollow)
	}

	feed, err := s.db.GetFeedByURL(ctx, parsedURL.String())
	if err != nil {
		return err
	}

	feedFollow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	feedFollowInfo, err := s.db.CreateFeedFollow(ctx, feedFollow)
	if err != nil {
		return err
	}

	fmt.Printf("%s is now following %s\n", feedFollowInfo.UserName, feedFollowInfo.FeedName)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	if cmd.hasArgs() {
		return fmt.Errorf("error: %s takes no arguments", cmd.name)
	}

	ctx := context.Background()

	follows, err := s.db.GetFeedFollowsForUser(ctx, user.Name)
	if err != nil {
		return err
	}

	fmt.Printf("%s, you are currently following:\n", user.Name)
	for _, follow := range follows {
		fmt.Printf("- %s\n", follow.FeedName)
	}

	return nil
}
