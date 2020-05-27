package cmd

import (
	"github.com/gren236/m3u-proxifier/cmd/m3u-proxifier/config"
	"github.com/gren236/m3u-proxifier/pkg/playlist"
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

	pl, err := playlist.Retrieve(jconf.Location)
	if err != nil {
		return err
	}
	defer pl.Close()

	proxPl, err := pl.Proxify(jconf.Proxy, "test/gen_testProxifiedFile.m3u")
	if err != nil {
		return err
	}
	defer proxPl.Close()

	// TODO Compare and merge files line by line
	// TODO Write new file (write and create flags?)

	return nil
}
