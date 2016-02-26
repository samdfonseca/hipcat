package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"log"
)

type Config struct {
	authToken     string
	defaultRoomId string
	defaultRoomName	string
	test          bool
}

func getConfigPath() string {
	homedir := os.Getenv("HOME")
	if homedir == "" {
		exitErr(fmt.Errorf("$HOME not set"))
	}
	return homedir + "/.hipcat"
}

func readConfig() *Config {
	config := &Config{
		authToken:     "",
		defaultRoomId: "",
		test:          false,
	}
	lines := readLines(getConfigPath())

	for _, line := range lines {
		s := strings.Split(line, "=")
		if len(s) != 2 {
			exitErr(fmt.Errorf("failed to parse config at: %s\n", line))
		}
		key := strip(s[0])
		switch key {
		case "auth_token":
			config.authToken = strip(s[1])
		case "default_room_name":
			config.defaultRoomName = strip(s[1])
		case "default_room_id":
			config.defaultRoomId = strip(s[1])
		case "test":
			if strip(s[1]) == "true" {
				config.test = true
			} else if strip(s[1]) == "false" {
				config.test = false
			} else {
				output(fmt.Sprintf("unrecognized value for 'test' in config: %s\n", strip(s[1])))
				log.Println("")

			}
		default:
			output(fmt.Sprintf("unrecognized config parameter: %s\n", line))
		}
	}

	return config
}

func strip(s string) string {
	return strings.Replace(s, " ", "", -1)
}

func readLines(path string) []string {
	var lines []string

	file, err := os.Open(path)
	failOnError(err, "unable to read config", true)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() != "" {
			lines = append(lines, scanner.Text())
		}
	}
	return lines
}
