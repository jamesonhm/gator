package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jamesonhm/gator/internal/config"
)

type state struct {
	cfg *config.Config
}

func main() {
	c, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	progState := &state{
		cfg: &c,
	}

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
