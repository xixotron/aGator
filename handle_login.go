package main

import (
	"context"
	"fmt"
)

func handleLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Usage: %s <username>", cmd.Name)
	}

	ctx := context.Background()
	userName := cmd.Args[0]
	user, err := s.db.GetUser(ctx, userName)
	if err != nil {
		return fmt.Errorf("Could not authenticate user: %w", err)
	}

	if user.Name != userName {
		return fmt.Errorf("Retrieved Username doesn't match provided username")
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("Could not set user: %w", err)
	}

	fmt.Println("User switched successfuly!")
	return nil
}
