package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/paultustain/gator/internal/config"
	"github.com/paultustain/gator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {

	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config one; %v", err)
	}

	dbURL := cfg.DBUrl

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("error opening the dbURL")
	}
	dbQueries := database.New(db)

	s := state{
		db:  dbQueries,
		cfg: &cfg,
	}

	cmds := commands{
		handlers: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)

	args := os.Args
	if len(args) < 2 {
		fmt.Println("Error: Not enough arguments")
		os.Exit(1)
	}

	cmd := command{
		name: os.Args[1],
		args: os.Args[2:],
	}

	if err := cmds.run(&s, cmd); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

}
