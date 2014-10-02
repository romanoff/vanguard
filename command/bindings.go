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
			host := c.Args().First()
			if host != "clear" {
				bindingsCommandFunc(host)
			} else {
				host = ""
				if len(c.Args()) == 2 {
					host = c.Args()[1]
				}
				clearBindingsCommandFunc(host)
			}
		},
	}
}

func bindingsCommandFunc(host string) {
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
		fmt.Println(host + ":" + binding.String())
	}
}

func clearBindingsCommandFunc(host string) {
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
		err := vClient.Hide(binding.Port, "", "")
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
