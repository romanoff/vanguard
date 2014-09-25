package command

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/romanoff/vanguard/client"
)

func NewHideCommand() cli.Command {
	return cli.Command{
		Name:  "hide",
		Usage: "hide port",
		Action: func(c *cli.Context) {
			hideCommandFunc(c)
		},
	}
}

func hideCommandFunc(c *cli.Context) {
	port := c.Args().First()
	host := "127.0.0.1"
	if len(c.Args()) == 2 {
		host = c.Args()[1]
	}
	vClient := client.NewClient(host)
	binding, err := vClient.Hide(port)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("hidden " + host + " port " + binding.String())
}
