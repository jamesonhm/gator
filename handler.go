package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jamesonhm/gator/internal/database"
)

func handleAddFeed(s *state, cmd command) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: addfeed <feed-name> <url>")
	}

	ctx := context.Background()
	user, err := s.db.GetUser(ctx, s.cfg.CurrUser)
	if err != nil {
		return fmt.Errorf("error getting user from db: %v", err)
	}

	feed, err := s.db.CreateFeed(ctx, database.CreateFeedParams{
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("error adding feed to database: %v", err)
	}

	fmt.Println(" * Name:", feed.Name)
	fmt.Println(" * Url:", feed.Url)
	fmt.Println(" * User:", feed.UserID)
	fmt.Println(" * CreatedAt:", feed.CreatedAt)
	return nil
}

func handleAgg(s *state, cmd command) error {
	url := "https://www.wagslane.dev/index.xml"
	feed, err := fetchFeed(context.Background(), url)
	if err != nil {
		return err
	}
	fmt.Println(feed)
	return nil
}

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

func handleReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return err
	}
	fmt.Printf("Users deleted from db\n")
	return nil
}

func handleUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, u := range users {
		fmt.Printf("* %s", u.Name)
		if u.Name == s.cfg.CurrUser {
			fmt.Printf(" (current)")
		}
		fmt.Printf("\n")
	}
	return nil
}

func printUser(user database.User) {
	fmt.Printf(" * ID:		%v\n", user.ID)
	fmt.Printf(" * Name:	%v\n", user.Name)
}
