package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"../common"
	"../structs"
	"github.com/urfave/cli"
)

type AddCommand struct {
}

func getAddFlags() []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:     "owner",
			Aliases:  []string{"o"},
			Usage:    "Immediately claims the task for the reporter.",
			Required: false,
		},
	}
}

func AddIt(c *cli.Context) error {

	var conf common.TaskerConfig
	var task structs.Task
	task.Reporter = conf.GetUsername()
	if c.NArg() > 0 && c.NArg() < 2 {
		task.Description = c.Args().Get(0)

		owner := c.IsSet("owner") || c.Bool("owner")
		if owner {
			task.Owner = conf.GetUsername()
		}

		bytesRepresentation, err := json.Marshal(task)
		if err != nil {
			log.Fatalln(err)
		}

		response, err := http.Post("http://localhost:8081/api/tasker/v1/save", "application/json", bytes.NewBuffer(bytesRepresentation))
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
			var task *structs.Task
			json.Unmarshal([]byte(contents), &task)
			layout := "2006-01-02T15:04:05"
			t, _ := time.Parse(layout, task.Time)
			fmt.Print("\n")
			fmt.Println("--------")
			fmt.Println("A task was added to the backlog.")
			fmt.Println("--------")
			fmt.Println("Key: " + task.TaskKey)
			fmt.Println("Description: " + task.Description)
			fmt.Println("Submitted: " + t.Format("January 2, 2006"))
			fmt.Println("--------")
			if !owner {
				fmt.Println("If you would like to claim this task use:")
				fmt.Println("tasker claim " + task.TaskKey)
				fmt.Println("--------")
			}
			fmt.Print("\n")
		}
	} else {
		fmt.Print("\n")
		fmt.Println("--------")
		fmt.Println("tasker add command was incorrectly formatted.")
		fmt.Println("Use the following syntax:")
		fmt.Println("-  tasker add \"description of the task\"")
		fmt.Println("-  tasker add -o \"description of the task\"")
		fmt.Println("--------")
		fmt.Print("\n")
	}

	return nil
}

func init() {
	flags := getAddFlags()
	common.RegisterCommand(&cli.Command{
		Name:   "add",
		Usage:  "Adds an entry for a new task.",
		Action: AddIt,
		Flags:  flags,
	})
}
