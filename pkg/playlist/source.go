package playlist

import (
	"net/url"
	"os"
)

type RetrieveError string

func (r RetrieveError) Error() string {
	return string(r)
}

func Retrieve(loc string) (*Playlist, error) {
	// Check if playlist is local or remote
	if _, err := os.Stat(loc); err == nil {
		return LoadLocal(loc)
	}

	if _, err := url.Parse(loc); err == nil {
		return LoadRemote(loc)
	}

	return nil, RetrieveError("source location not supported")
}
