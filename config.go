package main

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	Token string
	Debug bool
}

func (c *Configuration) init(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(c); err != nil {
		return err
	}

	return nil
}
