package main

import (
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

	dbQueries := database.New(db)

	cmds := commands{
		handlers: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handleLogin)

	inputArgs := os.Args
	if len(inputArgs) < 2 {
		log.Fatal("usage: cli <command> [args...]")
	}
	cmdName := inputArgs[1]
	cmdArgs := inputArgs[2:]

	progState := &state{
		cfg: &c,
		db:  dbQueries,
	}

	err = cmds.run(progState, command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		log.Fatal(err)
	}
	c, err = config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	fmt.Println(c)
}
