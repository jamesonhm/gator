package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/jamesonhm/gator/internal/config"
	"github.com/jamesonhm/gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

func main() {
	c, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	db, err := sql.Open("postgres", c.DBurl)
	if err != nil {
		log.Fatalf("error opening connection to DB: %v", err)
	}
	defer db.Close()
	dbQueries := database.New(db)

	progState := &state{
		cfg: &c,
		db:  dbQueries,
	}

	cmds := commands{
		handlers: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handleLogin)
	cmds.register("register", handleRegister)
	cmds.register("reset", handleReset)
	cmds.register("users", handleUsers)
	cmds.register("agg", handleAgg)
	cmds.register("addfeed", middlewareLoggedIn(handleAddFeed))
	cmds.register("feeds", handleFeeds)
	cmds.register("follow", middlewareLoggedIn(handleFollow))
	cmds.register("following", middlewareLoggedIn(handleFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handleUnfollow))

	inputArgs := os.Args
	if len(inputArgs) < 2 {
		log.Fatal("usage: cli <command> [args...]")
	}
	cmdName := inputArgs[1]
	cmdArgs := inputArgs[2:]

	err = cmds.run(progState, command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		log.Fatal(err)
	}
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrUser)
		if err != nil {
			return fmt.Errorf("error getting user from db: %v", err)
		}
		err = handler(s, cmd, user)
		if err != nil {
			return err
		}
		return nil
	}
}
