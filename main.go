package main

import (
	"fmt"
	"log"

	"github.com/paultustain/gator/internal/config"
)

func main() {

	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config one; %v", err)
	}
	fmt.Printf("Read once: %+v\n", cfg)

	err = cfg.SetUser("Paul")
	if err != nil {
		log.Fatalf("error setting user")
	}
	cfg, err = config.Read()
	if err != nil {
		log.Fatalf("error reading config")
	}
	fmt.Printf("Read twice: %+v\n", cfg)

}
