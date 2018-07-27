package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"github.com/smallfish/simpleyaml"
	"github.com/weex/slpp/daemon"
)

var assetsPath string


func processFlags() *daemon.Config {
	// load yaml file
	data, err := ioutil.ReadFile("config/base.yaml")
	if err != nil {
		panic("Can't read base.yaml.")
	}

	y, err := simpleyaml.NewYaml(data)
	if err != nil {
		panic("Can't parse base.yaml.")
	}

	name, err := y.Get("mysql").Get("user").String()
	if err != nil {
		panic("Can't find config yaml node.")
	}
	fmt.Println("name:", name)

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
