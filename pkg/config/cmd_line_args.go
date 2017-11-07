package config

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	confFile string
	//logLevel    int
	opts          string
	showInfo      bool
	printJSONConf bool

	confSet    = flag.NewFlagSet("config", flag.PanicOnError)
	infoSet    = flag.NewFlagSet("info", flag.PanicOnError)
	serviceSet = flag.NewFlagSet("service", flag.PanicOnError)
)

func usage() {
	fmt.Printf("Usage:\n\n")
	fmt.Println("info:")
	infoSet.PrintDefaults()
	fmt.Println("\nservice:")
	serviceSet.PrintDefaults()
}

func parseArgs() bool {

	if len(os.Args) == 1 {
		flag.PrintDefaults()
		log.Fatal("Please supply a sub-command: conf, info, service")
	}

	// configuration flagset
	confSet.StringVar(&confFile, "config", "", "TOML configuration file")
	confSet.BoolVar(&printJSONConf, "json", false, "Show services configuration in JSON")

	// info flagset
	infoSet.BoolVar(&showInfo, "info", false, "Show service information")

	// service flagset
	serviceSet.StringVar(&confFile, "config", "", "TOML configuration file")
	serviceSet.StringVar(&opts, "opts", "", "Comma-separated list of configuration values to be overwritten")

	/*
		serviceSet.IntVar(&logLevel, "log-level", 1, "log level: 0 Info, 1 Debug [default]")
	*/

	ok := false
	switch os.Args[1] {

	case "help":
		usage()

	case "config":
		if err := confSet.Parse(os.Args[2:]); err != nil {
			log.Fatal("cannot parse command-line arguments")
		}
		ok = true

	case "info":

	case "service":
		if err := serviceSet.Parse(os.Args[2:]); err != nil {
			log.Fatal("cannot parse command-line arguments")
		}
		ok = true

	default:
		serviceSet.PrintDefaults()
		log.Fatal("sub-command not supported")
	}

	return ok
}
