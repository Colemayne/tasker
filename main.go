package main

import (
	"fmt"
	"os"

	"./common"
	cli_helpers "./helpers/cli"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	// Blank import to initialize commands in the commands directory.
	_ "./commands"
)

func main() {

	defer func() {
		if r := recover(); r != nil {
			// log panics forces exit
			if _, ok := r.(*logrus.Entry); ok {
				os.Exit(1)
			}
			panic(r)
		}
	}()

	app := cli.NewApp()
	app.Name = "Tasker"
	app.Usage = "Let's you keep track of who is doing what."
	app.Authors = []*cli.Author{
		{
			Name:  "Coleman Beiler",
			Email: "coleman@beilers.com",
		},
	}

	// Return the commands registered in commands directory.
	app.Commands = common.GetCommands()
	// Logic for what happens when a command is not recoginized.
	app.CommandNotFound = func(context *cli.Context, command string) {
		//logrus.Fatalln("Command", command, "not found.")
		fmt.Println("tasker: '" + command + "' is not a tasker command")
		fmt.Println("See 'tasker -h'")
	}
	// InitCli initializes the Windows console window by activating virtual terminal features.
	// Calling this function enables colored terminal output.
	cli_helpers.InitCli()

	// starts the application
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
