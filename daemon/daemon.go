package daemon

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/weex/slpp/db"
	"github.com/weex/slpp/lnd"
	"github.com/weex/slpp/model"
	"github.com/weex/slpp/ui"
)

type Config struct {
	ListenSpec string

	Db  db.Config
	UI  ui.Config
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
	log.Printf("Connected to database.")

	// Lightning network
	lnd, err := lnd.InitLnd(cfg.Lnd)
	if err != nil {
		log.Printf("Error initializing Lightning: %v\n", err)
		return err
	}
	log.Printf("Connected to LND.")

	listener, err := net.Listen("tcp", cfg.ListenSpec)
	if err != nil {
		log.Printf("Error creating listener: %v\n", err)
		return err
	}
	log.Printf("Started HTTP server.")

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
