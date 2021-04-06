package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	UserName string `json:"username"`
	PassWord string `json:"password"`
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Load(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	datajson := []byte(data)
	err = json.Unmarshal(datajson, c)
	if err != nil {
		return err
	}
	return nil
}
