package daemon

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"../db"
    "../lnd"
	"../model"
	"../ui"
)

type Config struct {
	ListenSpec string

	Db db.Config
	UI ui.Config
    Lnd lnd.Config
}

func Run(cfg *Config) error {
	log.Printf("Starting, HTTP on: %s\n", cfg.ListenSpec)

    // Database
    db, err := db.InitDb(cfg.Db)
	if err != nil {
		log.Printf("Error initializing database: %v\n", err)
		return err
	}

	m := model.New(db)

    // Lightning network
    lnd, err := lnd.InitLnd(cfg.Lnd)
	if err != nil {
		log.Printf("Error initializing Lightning: %v\n", err)
		return err
	}
    log.Printf("got lnd %v", lnd)

	listener, err := net.Listen("tcp", cfg.ListenSpec)
	if err != nil {
		log.Printf("Error creating listener: %v\n", err)
		return err
	}

	ui.Start(cfg.UI, m, lnd, listener)

	waitForSignal()

	return nil
}

func waitForSignal() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	s := <-ch
	log.Printf("Got signal: %v, exiting.", s)
}
