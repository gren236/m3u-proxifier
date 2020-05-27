package playlist

import (
	"bufio"
	"fmt"
	"io"
	"net/url"
	"os"
	"regexp"
)

const M3U_SIGN string = "#EXTM3U"

var AddrRegexp *regexp.Regexp = regexp.MustCompile(`^.*://.*@(.*)$`)

type Error string

func (e Error) Error() string {
	return string(e)
}

type Playlist struct {
	LocalPath      string
	RemoteLocation url.URL
	*os.File
}

func (p *Playlist) Proxify(addr string, resFileName string) (*Playlist, error) {
	// Setting the file pointer back to the beginning once reading is complete
	defer p.Seek(0, io.SeekStart)

	resPl, err := New(resFileName)
	if err != nil {
		return nil, err
	}
	// Setting the file pointer back to the beginning once writing is complete
	defer resPl.Seek(0, io.SeekStart)

	resW := bufio.NewWriter(resPl)

	scnr := bufio.NewScanner(p)
	for scnr.Scan() {
		line := scnr.Text()

		// If the line does not contain address, write it as-is
		if line[:1] == "#" {
			_, err := resW.Write([]byte(line + "\n"))
			if err != nil {
				return nil, err
			}
			continue
		}

		var resLine string
		// If match found the address, extract it and form a resulting line, otherwise write the line as-is
		if match := AddrRegexp.FindSubmatch([]byte(line)); len(match) > 1 {
			resLine = fmt.Sprintf("http://%s/udp/%s\n", addr, match[1])
		} else {
			resLine = line
		}

		_, err := resW.Write([]byte(resLine))
		if err != nil {
			return nil, err
		}
	}

	// Check scanner errors
	if err := scnr.Err(); err != nil {
		return nil, err
	}

	// Check writing errors
	if err := resW.Flush(); err != nil {
		return nil, err
	}

	return resPl, nil
}

// Merge adds entries from given playlist to receiver playlist
// Entries with the same address are preserved
func (p *Playlist) Merge(old *Playlist) error {
	// TODO
	return nil
}

func New(path string) (*Playlist, error) {
	tmp, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0775)
	if err != nil {
		return nil, err
	}

	res := Playlist{
		LocalPath:      path,
		RemoteLocation: url.URL{},
	}
	res.File = tmp

	return &res, nil
}

func LoadLocal(path string) (*Playlist, error) {
	res, err := New(path)
	if err != nil {
		return nil, err
	}

	if !isM3U(res) {
		return nil, Error("file is not an M3U playlist")
	}

	return res, nil
}

func LoadRemote(rawLoc string) (*Playlist, error) {
	// TODO Send HTTP request
	// TODO Download the file
	// TODO Check if file is M3U playlist
	// TODO Return playlist

	return &Playlist{}, nil
}

func isM3U(r io.ReadSeeker) bool {
	scnr := bufio.NewScanner(r)
	// Setting the file pointer back to the beginning once reading is complete
	defer r.Seek(0, io.SeekStart)
	scnr.Scan()

	if scnr.Text() == M3U_SIGN {
		return true
	}

	return false
}