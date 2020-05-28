package cmd

import (
	"github.com/gren236/m3u-proxifier/cmd/m3u-proxifier/config"
	"github.com/gren236/m3u-proxifier/pkg/playlist"
	"io/ioutil"
	"os"
	"path/filepath"
)

func Handle(conf string) error {
	file, err := os.OpenFile(conf, os.O_RDONLY, 0775)
	if err != nil {
		return err
	}
	defer file.Close()

	// Parse the configuration JSON.
	jconf, err := config.ParseJSON(file)
	if err != nil {
		return err
	}

	// Load updated playlist.
	pl, err := playlist.Retrieve(jconf.Location)
	if err != nil {
		return err
	}
	defer pl.Close()

	// Proxify the updated playlist and put it to the unique temp directory.
	tempDir, err := ioutil.TempDir("", "m3u-proxifier")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	proxPl, err := pl.Proxify(jconf.Proxy, filepath.Join(tempDir, "proxifiedFile.m3u"))
	if err != nil {
		return err
	}
	defer proxPl.Close()

	// Open the old playlist.
	oldPl, err := playlist.LoadLocal(jconf.Old)
	if err != nil {
		return err
	}
	defer oldPl.Close()

	// Create a new file with merged entries.
	newPl, err := oldPl.Merge(proxPl, jconf.New)
	if err != nil {
		return err
	}
	defer newPl.Close()

	return nil
}
