package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"

	"../../impl/config"
	"../../impl/probe"
)

const (
	defaultConfigPath = "./config.toml"
)

var (
	configFn = flag.String("f", defaultConfigPath, "path to toml config file")
)

func main() {

	flag.Parse()

	f, err := os.Open(*configFn)
	if err != nil {
		panic(err)
	}

	var config config.Config
	_, err = toml.DecodeReader(f, &config)
	if err != nil {
		panic(err)
	}

	log.Printf("config = %+v", config)

	for targetId, target := range config.Targets {
		p, err := MakeProbe(target)
		if err != nil {
			log.Printf("error for probe '%s': %v", targetId, err)
		} else {
			log.Printf("probe %s  => %+v", targetId, p.Test())
		}
	}
}

func MakeProbe(target config.Target) (probe.Probe, error) {
	switch target.Type {
	case config.TypeHttp:
		return &probe.HttpProbe{URL: target.HttpUrl}, nil
	default:
		return nil, fmt.Errorf("unsupported probe type '%s'", target.Type)
	}
}
