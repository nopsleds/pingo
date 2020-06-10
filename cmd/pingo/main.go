package main

import (
	"flag"
	"log"

	"../../impl/config"
	"../../impl/core"
	"../../impl/web"
)

const (
	defaultConfigPath = "./config.toml"
)

var (
	configFn = flag.String("f", defaultConfigPath, "path to toml config file")
)

func panicIf(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}

func main() {

	log.Println("=== rtower ===")
	flag.Parse()

	config, err := config.LoadOrInitFile(*configFn)
	panicIf(err)

	log.Printf("config = %+v", config)

	instance, err := core.New(*config)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	log.Println("starting instance...")
	panicIf(instance.Run())

	log.Println("starting web...")
	go panicIf(web.RunWeb(config.Web, instance))

}
