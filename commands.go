package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gatorcli/internal/database"
	"internal/config"
	"internal/rss"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	cmdToFunction map[string]func(*state, command) error
}

func scrapeFeeds(s *state) error {
	//get next feed from db
	next_feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}
	//mark it as fetched
	err = s.db.MarkFeedFetched(context.Background(), next_feed.ID)
	if err != nil {
		return err
	}
	//fetch the feed using the url
	//feed, err := s.db.GetFeedByUrl(context.Background(), next_feed.Url)
	feed, err := rss.FetchFeed(context.Background(), next_feed.Url)
	if err != nil {
		return err
	}

	//iterate over the items in the feed and print their titles to the console
	for _, item := range feed.Channel.Item {
		post := database.CreatePostParams{
			ID:          uuid.New(),
			Title:       item.Title,
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: true},
			PublishedAt: time.Now(),
			FeedID:      next_feed.ID,
		}
		/*
			fmt.Println("*------------")
			fmt.Printf("| %s \n", post.Title)
			fmt.Printf("| %s \n", post.Url)
			fmt.Printf("| %s \n", post.Description.String)
			fmt.Println("*------------")
		*/
		stuff, err := s.db.CreatePost(context.Background(), post)
		if err != nil {
			var e *pq.Error
			if errors.As(err, &e) && e.Code.Name() != "unique_violation" {
				fmt.Printf("%s \n", e.Code.Name())
			}
		} else {
			fmt.Printf("* %s\n", stuff.Title)
		}
	}
	fmt.Println("Fetched All Feeds")
	return nil
}

// registers a new handler function to a command name
func (c *commands) register(name string, f func(*state, command) error) {
	c.cmdToFunction[name] = f
}

// runs the given command with the provided state if it exists
func (c *commands) run(s *state, cmd command) error {
	f, ok := c.cmdToFunction[cmd.name]
	if !ok {
		return fmt.Errorf("Command not found")
	}
	err := f(s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		current_user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return err
		}
		return handler(s, cmd, current_user)
	}
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("Need to provide a username for login")
	}
	username := cmd.args[0]
	_, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("User '%s' does not exist!\n", username)
	}
	err = config.SetUser(*s.cfg, username)
	if err != nil {
		return err
	}
	fmt.Println("The user has been successfuly set")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("Need to provide a username to register")
	}
	username := cmd.args[0]

	_, err := s.db.GetUser(context.Background(), username)
	if err == nil {
		return fmt.Errorf("A user with the name '%s' already exists", username)
	}

	tempUser := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	}
	usr, err := s.db.CreateUser(context.Background(), tempUser)
	if err != nil {
		return err
	}
	config.SetUser(*s.cfg, username)
	fmt.Println("The user has been successfuly registered")
	fmt.Printf("%s: %s\n", usr.CreatedAt.String(), usr.Name)
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.ResetUsers(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, u := range users {
		if u.Name == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", u.Name)
		} else {
			fmt.Printf("* %s\n", u.Name)
		}
	}
	return nil
}

// takes a single parameter
func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("Not enough args provided")
	}
	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Printf("Collecting feeds every %s\n", timeBetweenRequests.String())

	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
	return nil
}

// Takes a name and a url
func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("Not enough args provided")
	}

	name := cmd.args[0]
	url := cmd.args[1]

	feed := database.CreateFeedParams{
		ID:        uuid.New(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
		CreatedAt: time.Now(),
	}
	new_feed, err := s.db.CreateFeed(context.Background(), feed)
	if err != nil {
		return err
	}

	feed_follows := database.CreateFeedFollowParams{
		ID:     uuid.New(),
		UserID: user.ID,
		FeedID: new_feed.ID,
	}

	_, err = s.db.CreateFeedFollow(context.Background(), feed_follows)
	if err != nil {
		return err
	}

	fmt.Println(feed)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feedsUsers, err := s.db.GetFeedUsers(context.Background())
	if err != nil {
		return err
	}
	for _, f := range feedsUsers {
		fmt.Printf("* %s %s %s\n", f.Name, f.Url, f.UserName.String)
	}
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("Not enough args provided")
	}

	url := cmd.args[0]
	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return err
	}

	feed_follows := database.CreateFeedFollowParams{
		ID:     uuid.New(),
		UserID: user.ID,
		FeedID: feed.ID,
	}

	feed_follow_row, err := s.db.CreateFeedFollow(context.Background(), feed_follows)
	if err != nil {
		return err
	}

	fmt.Printf("Feed: %s User: %s\n", feed_follow_row.FeedName, feed_follow_row.UserName)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	following_feeds, err := s.db.GetFeedsByUser(context.Background(), user.ID)
	if err != nil {
		return err
	}
	for _, feed := range following_feeds {
		fmt.Printf("* %s %s\n", feed.Name, feed.Url)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("Not enough args provided")
	}

	url := cmd.args[0]
	err := s.db.DeleteFeedFollowRecord(context.Background(), database.DeleteFeedFollowRecordParams{user.ID, url})
	if err != nil {
		return err
	}
	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.args) >= 1 {
		var err error
		limit, err = strconv.Atoi(cmd.args[0])
		if err != nil {
			return err
		}
	}
	posts, err := s.db.GetPostsFromUser(context.Background(), database.GetPostsFromUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return err
	}
	for _, post := range posts {
		fmt.Printf("* %s\n", post.Title)
	}
	return nil
}
