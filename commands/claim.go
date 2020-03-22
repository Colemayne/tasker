package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"../common"
	"../structs"
	"github.com/urfave/cli"
)

type ClaimCommand struct {
}

func ClaimIt(c *cli.Context) error {

	var claim structs.Claim
	var conf common.TaskerConfig
	if c.NArg() > 0 || c.NArg() < 2 {
		claim.TaskKey = c.Args().Get(0)
		claim.Owner = conf.GetUsername()
		bytesRepresentation, err := json.Marshal(claim)
		if err != nil {
			log.Fatalln(err)
		}
		response, err := http.Post("http://localhost:8081/api/tasker/v1/claim", "application/json", bytes.NewBuffer(bytesRepresentation))

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
			var result bool
			json.Unmarshal([]byte(contents), &result)

			fmt.Print("\n")
			fmt.Println("--------")
			if result {
				fmt.Println("you have claimed: " + c.Args().Get(0))
			} else {
				fmt.Println(c.Args().Get(0) + " was not claimed.")
			}
			fmt.Println("--------")
			fmt.Print("\n")
		}

	} else {
		fmt.Println("This command requires one argument")
	}
	return nil
}

func init() {
	common.RegisterCommand(&cli.Command{
		Name:   "claim",
		Usage:  "Claim a task",
		Action: ClaimIt,
	})
}
