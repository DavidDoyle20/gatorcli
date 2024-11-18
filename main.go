package main

import (
	"database/sql"
	"fmt"
	"github.com/DavidDoyle20/gatorcli/internal/config"
	"github.com/DavidDoyle20/gatorcli/internal/database"
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
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))

	args := os.Args
	var cmd command

	if len(args) == 1 {
		// Do nothing, no command given
		os.Exit(0)
	}
	if len(args) == 2 {
		// No arguments other than the name given
		cmd.name = args[1]
	} else {
		cmd.name = args[1]
		cmd.args = args[2:]
	}

	err = cmds.run(&currentState, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
