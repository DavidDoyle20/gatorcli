package main

import (
	"context"
	"fmt"
	"gatorcli/internal/database"
	"github.com/google/uuid"
	"internal/config"
	"os"
	"time"
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

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("Need to provide a username for login")
	}
	username := cmd.args[0]
	_, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		fmt.Printf("User '%s' does not exist!\n", username)
		os.Exit(1)
	}
	err = config.SetUser(*s.cfg, username)
	if err != nil {
		return err
	}
	fmt.Println("The user has been successfuly set")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("Need to provide a username to register")
	}
	username := cmd.args[0]

	_, err := s.db.GetUser(context.Background(), username)
	if err == nil {
		fmt.Printf("A user with the name '%s' already exists", username)
		os.Exit(1)
	}

	tempUser := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	}
	usr, err := s.db.CreateUser(context.Background(), tempUser)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	config.SetUser(*s.cfg, username)
	fmt.Println("The user has been successfuly registered")
	fmt.Printf("%s: %s\n", usr.CreatedAt.String(), usr.Name)
	return nil
}
