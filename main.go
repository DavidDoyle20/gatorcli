package main

import (
	"database/sql"
	"fmt"
	"gatorcli/internal/database"
	"internal/config"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	//Read config
	cfg, _ := config.Read()

	// Open connection with database
	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		fmt.Println(err)
	}
	dbQueries := database.New(db)

	// Create initial state
	currentState := state{
		db:  dbQueries,
		cfg: &cfg,
	}

	// Initialize command map
	var cmds commands
	cmds.cmdToFunction = make(map[string]func(*state, command) error)

	// Register new commands here
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)

	args := os.Args
	if len(args) < 3 {
		fmt.Println("Must have more than 2 arguments")
		os.Exit(1)
	}
	var cmd command
	cmd.name = args[1]
	cmd.args = args[2:]
	cmds.run(&currentState, cmd)
}
