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
		return &state{}, fmt.Errorf("Error config: %w", err)
	}

	programState := &state{cfg: &cfg}

	db, err := sql.Open("postgres", programState.cfg.DBURL)
	if err != nil {
		return &state{}, fmt.Errorf("Error connecting to DB: %w", err)
	}

	programState.db = database.New(db)
	return programState, nil
}

func prepareCommands() commands {
	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handleLogin)
	cmds.register("register", handleRegister)

	return cmds
}

func parseArgs() (command, error) {
	if len(os.Args) < 2 {
		return command{}, fmt.Errorf("Usage %s <command> [args ...]\n", path.Base(os.Args[0]))
	}

	cmd := command{}
	cmd.Name = os.Args[1]
	cmd.Args = os.Args[2:]

	return cmd, nil
}
