package config

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/url"
)

type ConfigJSON struct {
	Location        *url.URL
	Proxy, Old, New string
}

func ParseJSON(r io.Reader) (*ConfigJSON, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var raw map[string]interface{}

	err = json.Unmarshal(data, &raw)
	if err != nil {
		return nil, err
	}

	// Parse basic types
	res := ConfigJSON{
		Old:   raw["old"].(string),
		New:   raw["new"].(string),
		Proxy: raw["proxy"].(string),
	}

	// Parse each string containing URL to URL type
	if tmp, err := url.Parse(raw["location"].(string)); err != nil {
		return nil, err
	} else {
		res.Location = tmp
	}

	return &res, nil
}
