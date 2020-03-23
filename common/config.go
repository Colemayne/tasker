package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const (
	Host       = "http://localhost:8081"
	ConfigFile = "/etc/tasker/tasker.json"
)

type TaskerConfig struct {
	User string `json:"user"`
	Host string `json:"host"`
}

func readConfig() []byte {
	data, err := ioutil.ReadFile("/etc/tasker/tasker.json")
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func writeConfig(c *TaskerConfig) error {
	str, _ := json.Marshal(c)
	err := ioutil.WriteFile("/etc/tasker/tasker.json", []byte(str), 0777)
	if err != nil {
		log.Fatal(err)
	}

	return err
}

// LoadConfig attempts to load configuration. Will create a default config file if one doesn't exist.
func (c *TaskerConfig) LoadConfig() error {
	// attempt to see if config file is present
	_, err := os.Stat(ConfigFile)

	if os.IsNotExist(err) {
		fmt.Print("\n")
		fmt.Println("--------")
		fmt.Println("No configuration file found. Creating one now.")
		fmt.Println("--------")
		fmt.Print("\n")
		c.Host = "http://localhost:8081"
		c.User = ""
		writeConfig(c)
		return nil
	} else if err != nil {
		return err
	}

	return nil
}

// SetUsername sets username to use for tasker.
func (c *TaskerConfig) SetUsername(username string) error {
	c.User = username
	err := writeConfig(c)
	return err
}

// GetUsername returns the username as found in the config.
func (c *TaskerConfig) GetUsername() string {
	data := readConfig()
	json.Unmarshal([]byte(data), &c)
	return c.User
}

// SetHost sets host adderess to use for tasker.
func (c *TaskerConfig) SetHost(host string) error {
	c.Host = host
	err := writeConfig(c)
	return err
}
