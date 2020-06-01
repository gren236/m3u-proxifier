package playlist

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"time"
)

const M3U_SIGN string = `#EXTM3U`

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

// Merge adds entries from given playlist to receiving playlist and returns a new one.
// Entries with the same address are preserved
func (p *Playlist) Merge(new *Playlist, resFileName string) (*Playlist, error) {
	// Setting the file pointer back to the beginning once reading is complete
	defer p.Seek(0, io.SeekStart)
	defer new.Seek(0, io.SeekStart)

	// Create a new file and copy all current addresses to map for a further checking
	resPl, err := New(resFileName)
	if err != nil {
		return nil, err
	}
	// Setting the file pointer back to the beginning once writing is complete
	defer resPl.Seek(0, io.SeekStart)

	resW := bufio.NewWriter(resPl)

	// Map to check for new entries
	chk := make(map[string]bool)

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

		// If line is address, add it to map
		chk[line] = true

		_, err := resW.Write([]byte(line + "\n"))
		if err != nil {
			return nil, err
		}
	}

	// Iterate over given entries and check if there are new ones
	nScnr := bufio.NewScanner(new)
	nScnr.Scan()       // Pass the first line
	var descBuf string // Hold the last description line for a further writing
	for nScnr.Scan() {
		line := nScnr.Text()

		// If the line does not contain address, write it to description buffer
		if line[:1] == "#" {
			descBuf = line
			continue
		}

		// If line is not present in map, write it to the file
		if _, ok := chk[line]; !ok {
			_, err := resW.Write([]byte(descBuf + "\n" + line + "\n"))
			if err != nil {
				return nil, err
			}

			// Write it to the map to avoid duplications
			chk[line] = true
		}
	}

	// Check scanner errors
	if err := scnr.Err(); err != nil {
		return nil, err
	}
	if err := nScnr.Err(); err != nil {
		return nil, err
	}
	// Check writing errors
	if err := resW.Flush(); err != nil {
		return nil, err
	}

	return resPl, nil
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
	// Send HTTP request
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", rawLoc, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Download the file
	tmp, err := ioutil.TempFile("", "m3u-proxifier_")
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(tmp, resp.Body)
	if err != nil {
		return nil, err
	}

	// Check if file is M3U playlist
	tmp.Seek(0, io.SeekStart) // Set file pointer offset to zero after writing
	if !isM3U(tmp) {
		return nil, Error("file is not an M3U playlist")
	}

	// Return playlist
	res, err := New(tmp.Name())
	if err != nil {
		return nil, err
	}

	return res, nil
}

func isM3U(r io.ReadSeeker) bool {
	scnr := bufio.NewScanner(r)
	// Setting the file pointer back to the beginning once reading is complete
	defer r.Seek(0, io.SeekStart)
	scnr.Scan()

	if ok, _ := regexp.Match(M3U_SIGN, []byte(scnr.Text())); ok {
		return true
	}

	return false
}
