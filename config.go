package main

import (
	"encoding/json"
	"errors"
	"io"
	"os"
)

type Config struct {
	Port        uint64            `json:"port"`
	Nameservers map[string]string `json:"nameservers"`
}

func NewConfigFromPath(path string) (c *Config, err error) {
	file, err := os.Open(path)
	if err != nil {
		err = errors.New("Error when opening file from path " + path + ": " + err.Error())
		return
	}
	content, err := io.ReadAll(file)
	if err != nil {
		err = errors.New("Error when reading file from path " + path + ": " + err.Error())
		return
	}
	err = json.Unmarshal(content, &c)
	if err != nil {
		err = errors.New("Error when unmarshalling JSON c file from path " + path + ": " + err.Error())
		return
	}
	return
}
