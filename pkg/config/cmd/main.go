package main

// Command-line tool for converting a global TOML service configuration to JSON.

import (
	"log"

	"github.com/luigi-riefolo/eGO/pkg/config"
)

func main() {

	if err := config.PrintJSONConfig(); err != nil {
		log.Fatalf("Cannot print services JSON configuration: %v", err)

	}
}
