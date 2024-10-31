package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jamesonhm/gator/internal/database"
)

func handleLogin(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	user := cmd.Args[0]
	u, err := s.db.GetUser(context.Background(), user)
	if err != nil {
		return err
	}

	err = s.cfg.SetUser(u.Name)
	if err != nil {
		return err
	}
	fmt.Printf("User set to %s\n", user)
	return nil
}

func handleRegister(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	user := cmd.Args[0]
	u, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      user})
	if err != nil {
		return err
	}

	err = s.cfg.SetUser(user)
	if err != nil {
		return err
	}
	fmt.Println("User created and returned:")
	printUser(u)
	return nil
}

func printUser(user database.User) {
	fmt.Printf(" * ID:		%v\n", user.ID)
	fmt.Printf(" * Name:	%v\n", user.Name)
}
