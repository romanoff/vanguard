package command

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/romanoff/vanguard/client"
)

func NewStopCommand() cli.Command {
	return cli.Command{
		Name:  "stop",
		Usage: "stop container",
		Action: func(c *cli.Context) {
			stopCommandFunc(c)
		},
	}
}

func stopCommandFunc(c *cli.Context) {
	containerId := c.Args().First()
	vClient := client.NewClient("127.0.0.1")
	err := vClient.Stop(containerId)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("success")
}
