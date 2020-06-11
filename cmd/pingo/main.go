package main

import (
	"flag"
	"log"

	"github.com/nopsleds/pingo/impl"
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

	config, err := impl.LoadOrInitFile(*configFn)
	panicIf(err)

	instance, err := impl.NewPingoInstance(*config)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	log.Println("starting instance...")
	panicIf(instance.Run())

	log.Println("starting web...")
	panicIf(impl.RunWeb(config.Web, instance))
}
