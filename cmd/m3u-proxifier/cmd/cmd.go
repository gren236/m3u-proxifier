package cmd

import (
	"fmt"
	"github.com/gren236/m3u-proxifier/cmd/m3u-proxifier/config"
	"os"
)

func Handle(conf string) error {
	file, err := os.OpenFile(conf, os.O_RDONLY, 0775)
	if err != nil {
		return err
	}
	defer file.Close()

	jconf, err := config.ParseJSON(file)
	if err != nil {
		return err
	}

	fmt.Println(jconf)
	// TODO Check if playlist is remote
	// TODO Get file contents
	// TODO Parse file contents
	// TODO Compare and merge files line by line
	// TODO Write new file (write and create flags?)

	return nil
}
