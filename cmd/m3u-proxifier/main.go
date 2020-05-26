package main

import (
	"flag"
	"fmt"
	"github.com/gren236/m3u-proxifier/cmd/m3u-proxifier/cmd"
	"os"
)

func main() {
	// Define configuration path flag
	config := flag.String("config", "", "Absolute path to configuration file.")

	flag.Parse()

	if *config == "" {
		fmt.Fprintln(os.Stderr, "No configuration provided!")
		return
	}

	if err := cmd.Handle(*config); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
