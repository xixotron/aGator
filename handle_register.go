package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/xixotron/aGator/internal/database"
)

func handleRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Usage: %s <username>", cmd.Name)
	}

	ctx := context.Background()

	user, err := s.db.CreateUser(
		ctx,
		database.CreateUserParams{
			ID:        uuid.New(),
			Name:      cmd.Args[0],
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	)
	if err != nil {
		return fmt.Errorf("Could not create user: %w", err)
	}

	if err := s.cfg.SetUser(user.Name); err != nil {
		return fmt.Errorf("Could not set current user: %w", err)
	}

	return nil
}
