package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jamesonhm/gator/internal/database"
)

func handleAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: addfeed <feed-name> <url>")
	}

	ctx := context.Background()

	feed, err := s.db.CreateFeed(ctx, database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("error adding feed to database: %v", err)
	}

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("error creating new follow: %v", err)
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

func handleFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}
	url := cmd.Args[0]

	f, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("error getting feed by url: %v", err)
	}

	ff, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    f.ID,
	})
	if err != nil {
		return fmt.Errorf("error getting feed follow: %v", err)
	}
	fmt.Println(ff.UserName, "is now following", ff.FeedName)
	return nil
}

func handleFollowing(s *state, cmd command, user database.User) error {
	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("error getting feed follows for user: %v", err)
	}

	if len(follows) == 0 {
		fmt.Println("you are not following any feeds")
		return nil
	}
	for _, follow := range follows {
		fmt.Println("* ", follow.FeedName)
	}
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

func handleFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error listing feeds: %v", err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found")
		return nil
	}

	for i, f := range feeds {
		user, err := s.db.GetUserByID(context.Background(), f.UserID)
		if err != nil {
			return fmt.Errorf("error getting user: %s for feed: %s. Error: %v", f.UserID, f.Name, err)
		}
		fmt.Printf("%d - Feed: %s, URL: %s, User: %s\n", i+1, f.Name, f.Url, user.Name)
	}
	return nil
}

func handleUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <feed-url>", cmd.Name)
	}

	err := s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		Url:    cmd.Args[0],
	})
	if err != nil {
		return fmt.Errorf("error deleting follow: %v", err)
	}
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
