package command

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/romanoff/vanguard/client"
)

func NewRunCommand() cli.Command {
	return cli.Command{
		Name:  "run",
		Usage: "run container",
		Action: func(c *cli.Context) {
			runCommandFunc(c)
		},
	}
}

func runCommandFunc(c *cli.Context) {
	name := c.Args().First()
	if name == "" {
		fmt.Println("No image has been specified")
		return
	}
	vClient := client.NewClient("127.0.0.1")
	container, err := vClient.Run(name, "", "", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(container)
}
