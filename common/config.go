package common

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type TaskerConfig struct {
	User        string `json:"user"`
	CurrentTask string `json:"currentTask"`
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
	err := ioutil.WriteFile("/etc/tasker/tasker.json", []byte(str), 0644)
	if err != nil {
		log.Fatal(err)
	}

	return err
}

func (c *TaskerConfig) SetUsername(username string) error {
	c.User = username
	err := writeConfig(c)
	return err
}

func (c *TaskerConfig) GetUsername() string {
	data := readConfig()
	json.Unmarshal([]byte(data), &c)
	return c.User
}
