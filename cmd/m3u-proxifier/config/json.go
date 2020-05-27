package config

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

type Json struct {
	Location, Proxy, Old, New string
}

func ParseJSON(r io.Reader) (*Json, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	res := new(Json)

	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
