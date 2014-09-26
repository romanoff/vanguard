package command

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/romanoff/vanguard/client"
	"strings"
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
	host := "127.0.0.1"
	port := c.Args().First()
	slice := strings.Split(c.Args()[0], ":")
	if len(slice) == 2 {
		host = slice[0]
		port = slice[1]
	}
	bindingHost := ""
	bindingPort := ""
	if len(c.Args()) == 2 {
		slice := strings.Split(c.Args()[1], ":")
		if len(slice) == 2 {
			bindingHost = slice[0]
			bindingPort = slice[1]
		}
	}
	vClient := client.NewClient(host)
	err := vClient.Hide(port, bindingHost, bindingPort)
	if err != nil {
		fmt.Println(err)
		return
	}
	bindingsCommandFunc(host)
}
