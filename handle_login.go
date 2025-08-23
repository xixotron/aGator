package main

import (
	"fmt"
)

func handleLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Usage: %s <username>", cmd.Name)
	}

	err := s.cfg.SetUser(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("Could not set user: %w", err)
	}

	fmt.Println("User switched successfuly!")
	return nil
}
