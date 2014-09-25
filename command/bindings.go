package command

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/romanoff/vanguard/client"
)

func NewBindingsCommand() cli.Command {
	return cli.Command{
		Name:  "bindings",
		Usage: "shows port bindings",
		Action: func(c *cli.Context) {
			bindingsCommandFunc(c)
		},
	}
}

func bindingsCommandFunc(c *cli.Context) {
	host := c.Args().First()
	if host == "" {
		host = "127.0.0.1"
	}
	vClient := client.NewClient(host)
	bindings, err := vClient.Bindings()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, binding := range bindings {
		fmt.Println(host + " port " + binding.String())
	}
}
