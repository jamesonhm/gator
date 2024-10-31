package main

import (
	"fmt"
	"log"

	"github.com/jamesonhm/gator/internal/config"
)

func main() {
	c, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	fmt.Println(c)

	err = c.SetUser("jhm")
	c, err = config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	fmt.Println(c)
}
