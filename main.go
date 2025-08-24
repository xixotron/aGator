package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/xixotron/aGator/internal/config"
	"github.com/xixotron/aGator/internal/database"

	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	programState, err := initState()
	if err != nil {
		log.Fatal(err)
	}

	cmds := prepareCommands()
	cmd, err := parseArgs()
	if err != nil {
		log.Fatal(err)
	}

	err = cmds.run(programState, cmd)
	if err != nil {
		log.Fatal(err)
	}
}

func initState() (*state, error) {
	cfg, err := config.Read()
	if err != nil {
		return &state{}, fmt.Errorf("error reading config: %w", err)
	}

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		return &state{}, fmt.Errorf("error connecting to db: %w", err)
	}

	return &state{
		cfg: &cfg,
		db:  database.New(db),
	}, nil
}

func prepareCommands() commands {
	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handleLogin)
	cmds.register("register", handleRegister)
	cmds.register("reset", handleDeleteAllUsers)
	cmds.register("users", handleListUsers)
	cmds.register("agg", handleAgg)
	cmds.register("addfeed", middlewareLoggedIn(handleAddFeed))
	cmds.register("feeds", handleListFeeds)
	cmds.register("follow", middlewareLoggedIn(handleFeedFollow))
	cmds.register("following", middlewareLoggedIn(handleFeedsFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handleUnfollow))

	return cmds
}

func parseArgs() (command, error) {
	if len(os.Args) < 2 {
		return command{},
			fmt.Errorf("Usage: %s <command> [args ...]\n", path.Base(os.Args[0]))
	}

	return command{
		Name: os.Args[1],
		Args: os.Args[2:],
	}, nil
}
