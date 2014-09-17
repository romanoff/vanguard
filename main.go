package main

import (
	"github.com/codegangsta/cli"
	"github.com/romanoff/vanguard/command"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "vanguard"
	app.Usage = "multihost docker orchestration"
	app.Commands = []cli.Command{
		command.NewAgentCommand(),
		command.NewRunCommand(),
		command.NewStopCommand(),
		command.NewPsCommand(),
	}
	app.Run(os.Args)
}
