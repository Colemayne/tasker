package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"../common"
	"../structs"
	"github.com/urfave/cli"
)

type ListCommand struct {
}

func getListFlags() []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:     "mine",
			Aliases:  []string{"m"},
			Usage:    "Show only my tasks",
			Required: false,
		},
		&cli.BoolFlag{
			Name:     "other",
			Aliases:  []string{"o"},
			Usage:    "Show only the tasks claimed by others",
			Required: false,
		},
		&cli.BoolFlag{
			Name:     "unclaimed",
			Aliases:  []string{"u"},
			Usage:    "Show only the tasks that are unclaimed",
			Required: false,
		},
	}
}

func ListIt(c *cli.Context) error {

	numberOfFlags := c.NumFlags()
	var mode int
	if numberOfFlags == 0 || numberOfFlags == 2 {
		if c.IsSet("mine") || c.Bool("mine") {
			mode = 1
		}
		if c.IsSet("other") || c.Bool("other") {
			mode = 2
		}
		if c.IsSet("unclaimed") || c.Bool("unclaimed") {
			mode = 3
		}
	} else {
		fmt.Print("\n")
		fmt.Println("--------")
		fmt.Println("The List command only allows on flag")
		fmt.Println("--------")
		fmt.Print("\n")
		return nil
	}

	var conf common.TaskerConfig
	err := conf.LoadConfig()
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	if conf.GetUsername() == "" {
		fmt.Print("\n")
		fmt.Println("--------")
		fmt.Println("Please set a username")
		fmt.Println("Run the register command")
		fmt.Println("--------")
		fmt.Print("\n")
		os.Exit(1)
	}
	response, err := http.Get(conf.GetHost() + "/api/tasker/v1/select/" + conf.GetUsername())
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}

		var tasks *structs.ListReponse
		json.Unmarshal([]byte(contents), &tasks)
		fmt.Print("\n")
		if mode == 0 || mode == 1 {
			fmt.Println("--------")
			fmt.Println("My Tasks:")
			fmt.Println("--------")
			myTasks := tasks.Mine
			if len(myTasks) > 0 {
				for _, elem := range myTasks {
					if elem.Description != "" {
						fmt.Println("(" + elem.TaskKey + ") - " + elem.Description)
					}
				}
			} else {
				fmt.Println("You have not claimed any tasks.")
			}
		}
		if mode == 0 || mode == 3 {
			unclaimed := tasks.Unclaimed
			fmt.Println("--------")
			fmt.Println("Available Tasks:")
			fmt.Println("--------")
			if len(unclaimed) > 0 {
				for _, elem := range unclaimed {
					if elem.Description != "" {
						fmt.Println("(" + elem.TaskKey + ") - " + elem.Description)
					}
				}
			} else {
				fmt.Println("No tasks are available in the backlog.")
			}
		}
		if mode == 0 || mode == 2 {
			others := tasks.Others
			fmt.Println("--------")
			fmt.Println("Claimed Tasks:")
			fmt.Println("--------")
			if len(others) > 0 {
				for _, elem := range others {
					if elem.Description != "" {
						fmt.Println("(" + elem.TaskKey + ")(" + elem.Owner + ") - " + elem.Description)
					}
				}
			} else {
				fmt.Println("Nobody else has claimed any tasks.")
			}
		}
		fmt.Println("--------")
		fmt.Print("\n")
	}
	return nil
}

func init() {
	flags := getListFlags()
	common.RegisterCommand(&cli.Command{
		Name:   "list",
		Usage:  "Lists the backlog tasks.",
		Action: ListIt,
		Flags:  flags,
	})
}
