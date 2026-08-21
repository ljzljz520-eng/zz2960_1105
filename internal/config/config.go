package config

import (
	"flag"
	"os"
)

type Config struct {
	DBPath string
	Listen string
}

func Default() Config { return Config{DBPath: "inventoryseal.db", Listen: ":8090"} }

func FromEnvironment() Config {
	value := Default()
	if path := os.Getenv("INVENTORYSEAL_DB"); path != "" {
		value.DBPath = path
	}
	if listen := os.Getenv("INVENTORYSEAL_LISTEN"); listen != "" {
		value.Listen = listen
	}
	return value
}

func Parse(args []string) Config {
	value := FromEnvironment()
	flags := flag.NewFlagSet("inventoryctl", flag.ContinueOnError)
	flags.StringVar(&value.DBPath, "db", value.DBPath, "database path")
	flags.StringVar(&value.Listen, "listen", value.Listen, "http listen address")
	_ = flags.Parse(args)
	return value
}
