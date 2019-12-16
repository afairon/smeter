package main

import (
	"flag"
	"log"

	"github.com/afairon/smeter/config"
	"github.com/afairon/smeter/server"
	_ "github.com/lib/pq"
)

func main() {

	var configFile string

	flag.StringVar(&configFile, "config", config.ConfigFile, "path to configuration file")

	flag.Parse()

	cfg := config.NewConfig()

	err := cfg.Load(configFile)
	if err != nil {
		log.Printf("Error while opening %s", configFile)
		err = cfg.Write(configFile)
		log.Printf("Creating configuration file %s", configFile)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	server := server.NewServer(cfg)

	server.Run()
}
