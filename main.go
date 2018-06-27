package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/weex/gowebapp/daemon"
)

var assetsPath string

func processFlags() *daemon.Config {
	cfg := &daemon.Config{}

	flag.StringVar(&cfg.ListenSpec, "listen", ":3001", "HTTP listen spec")

	flag.StringVar(&cfg.Db.ConnectString, "db-connect", "postgres://postgres:postgres@db/gowebapp?sslmode=disable", "DB Connect String")

	flag.StringVar(&cfg.Lnd.DataDir, "datadir", "/home/dsterry/.lnd", "Path to lnd data dir")

    flag.StringVar(&assetsPath, "assets-path", "assets", "Path to assets dir")

	flag.Parse()
	return cfg
}

func setupHttpAssets(cfg *daemon.Config) {
	log.Printf("Assets served from %q.", assetsPath)
	cfg.UI.Assets = http.Dir(assetsPath)
}

func main() {
	cfg := processFlags()

	setupHttpAssets(cfg)

	if err := daemon.Run(cfg); err != nil {
		log.Printf("Error in main(): %v", err)
	}
}
