package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"strconv"
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
	commands.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	commands.register("browse", middlewareLoggedIn(handlerBrowse))

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
	if !cmd.hasArgs() {
		return fmt.Errorf("error: %s requires one (1) argument: [time_between_reqs (e.g. 1s; 1m; 1h)]", cmd.name)
	}

	time_between_reqs := cmd.args[0]
	parsedInterval, err := time.ParseDuration(time_between_reqs)
	if err != nil {
		return err
	}

	fmt.Printf("\n\nCollecting feeds every %s\n\n", time_between_reqs)
	ticker := time.NewTicker(parsedInterval)
	for ; ; <-ticker.C {
		if err := scrapeFeeds(s); err != nil {
			return err
		}
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

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if !cmd.hasArgs() || len(cmd.args) > 1 {
		return fmt.Errorf("error: %s takes only one argument.\nUsage: %s [feed-url]", cmd.name, cmd.name)
	}

	ctx := context.Background()

	urlToRemove := cmd.args[0]

	parsedURL, err := url.Parse(urlToRemove)
	if err != nil {
		return fmt.Errorf("error parsing url: %s", urlToRemove)
	}

	feedToRemove := database.RemoveFollowForUserParams{
		UserID: user.ID,
		Url:    parsedURL.String(),
	}

	if _, err := s.db.RemoveFollowForUser(ctx, feedToRemove); err != nil {
		return err
	}

	fmt.Printf("INFO: RSS -- %s -- successfully removed from follows\n", parsedURL.String())
	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	if len(cmd.args) > 1 {
		return fmt.Errorf("error: %s takes only one optional argument.\nUsage: %s (limit)", cmd.name, cmd.name)
	}
	ctx := context.Background()

	var numPosts string
	if len(cmd.args) == 1 {
		numPosts = cmd.args[0]
	} else {
		numPosts = "2"
	}

	result, err := strconv.ParseInt(numPosts, 10, 32)
	if err != nil {
		return err
	}
	limit := int32(result)

	searchParams := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  limit,
	}

	posts, err := s.db.GetPostsForUser(ctx, searchParams)
	if err != nil {
		return err
	}

	dateTimeLayout := "2006-01-02 03:04:05 PM CST"
	fmt.Printf("%s, here are the %s most recent posts from your feed:\n\n", user.Name, numPosts)
	for _, post := range posts {
		fmt.Printf("\nPublished: %s\nTitle: %s\nURL: %s\nDescription: %s\n\n", post.PublishedAt.Format(dateTimeLayout), post.Title.String, post.Url, post.Description.String)
	}

	return nil
}

func scrapeFeeds(s *state) error {
	ctx := context.Background()
	feed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		return err
	}

	s.db.MarkFeedFetched(ctx, feed.ID)
	rssFeed, err := fetchFeed(ctx, feed.Url)
	if err != nil {
		return err
	}

	fmt.Printf("\nFetching feed: %v\n", rssFeed)
	for _, feedItem := range rssFeed.Channel.Item {
		title := buildSQLNullString(feedItem.Title)
		description := buildSQLNullString(feedItem.Description)
		published, err := parseDate(feedItem.PubDate)
		if err != nil {
			return err
		}

		newPost := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       title,
			Url:         feedItem.Link,
			Description: description,
			PublishedAt: published,
			FeedID:      feed.ID,
		}

		s.db.CreatePost(ctx, newPost)
	}

	return nil
}

func buildSQLNullString(text string) sql.NullString {
	var entry sql.NullString
	if len(text) > 0 {
		entry = sql.NullString{
			String: text,
			Valid:  true,
		}
	} else {
		entry = sql.NullString{
			String: "",
			Valid:  false,
		}
	}

	return entry
}

// A slice of common date layouts to try when parsing.
// The order matters; more common or specific formats should come first.
var dateLayouts = []string{
	time.RFC1123Z,
	time.RFC1123,
	time.RFC3339,
	time.RFC822,
	time.RFC822Z,
	"2006-01-02T15:04:05Z", // A variation of RFC3339
	"2006-01-02 15:04:05",  // Common SQL-like format
}

// parseDate attempts to parse a date string using a series of known layouts.
func parseDate(dateStr string) (time.Time, error) {
	for _, layout := range dateLayouts {
		t, err := time.Parse(layout, dateStr)
		if err == nil {
			return t, nil // Success!
		}
	}
	// If none of the layouts worked, return an error.
	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}
