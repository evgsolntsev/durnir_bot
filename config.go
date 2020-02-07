package main

import (
	"encoding/json"
	"math/rand"
	"os"
	"time"
)

type Configuration struct {
	Token       string
	Debug       bool
	Spreadsheet string
}

func (c *Configuration) init(filename string) error {
	rand.Seed(time.Now().Unix())

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
