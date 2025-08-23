package main

import (
	"github.com/xixotron/aleyGator/internal/config"
	"log"
	"os"
	"path"
)

type state struct {
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error config: %v", err)
	}

	programState := &state{
		cfg: &cfg,
	}

	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handleLogin)

	if len(os.Args) < 2 {
		log.Fatalf("Usage %s <command> [args ...]\n", path.Base(os.Args[0]))
	}

	cmd := command{
		Name: os.Args[1],
		Args: os.Args[2:],
	}

	err = cmds.run(programState, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
