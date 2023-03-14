package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Name             string
	DisplayName      string
	Description      string
	UserName         string
	Arguments        []string
	Executable       string
	LoggingType      []int
	Dependencies     []string
	WorkingDirectory string
	ChRoot           string
	Control          []string
	Transport        struct {
		Address  string
		Port     string
		Login    string
		Password string
		Http     bool
		Mqtt     bool
	}
	LogFileEnable bool
	LogFile       string
}

func (c *Config) Load() {
	file, _ := os.Open("config.json")
	byteValue, _ := ioutil.ReadAll(file)
	defer file.Close()
	json.Unmarshal(byteValue, &c)
	log.Printf("%v", c)
}
