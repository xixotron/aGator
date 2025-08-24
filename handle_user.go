package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/xixotron/aGator/internal/database"
)

func handleDeleteAllUsers(s *state, cmd command) error {
	err := s.db.DeleteAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't delete users: %w", err)
	}

	fmt.Println("Database reset successfully!")
	return s.cfg.SetUser("")
}

func handleListUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't list users: %w", err)
	}

	if len(users) == 0 {
		fmt.Println("No users found")
		return nil
	}

	fmt.Println("Registered users:")
	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf(" * %v (current)\n", user.Name)
			continue
		}
		fmt.Printf(" * %v\n", user.Name)
	}

	return nil
}

func handleRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Usage: %s <username>", cmd.Name)
	}

	ctx := context.Background()

	name := cmd.Args[0]
	user, err := s.db.CreateUser(ctx,
		database.CreateUserParams{
			ID:        uuid.New(),
			Name:      name,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	if err != nil {
		return fmt.Errorf("couldn't create user: %w", err)
	}

	if err := s.cfg.SetUser(user.Name); err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Println("User created successfully:")
	printUser(user)
	return nil
}

func handleLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Usage: %s <username>", cmd.Name)
	}

	name := cmd.Args[0]

	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("couldn't find user: %w", err)
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Println("User switched successfully!")
	return nil
}

func printUser(user database.User) {
	fmt.Printf(" * ID:      %v\n", user.ID)
	fmt.Printf(" * Name:    %v\n", user.Name)
}
