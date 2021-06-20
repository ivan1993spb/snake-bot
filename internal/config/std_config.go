package config

import (
	"flag"
	"os"
)

func StdConfig() (Config, error) {
	f := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	cfg, err := ParseFlags(f, os.Args[1:], DefaultConfig())
	return cfg, err
}
