package command

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/romanoff/vanguard/client"
)

func NewHostsCommand() cli.Command {
	return cli.Command{
		Name:  "hosts",
		Usage: "available hosts",
		Action: func(c *cli.Context) {
			hostsCommandFunc(c)
		},
	}
}

func hostsCommandFunc(c *cli.Context) {
	host := c.Args().First()
	if host == "" {
		host = "127.0.0.1"
	}
	vClient := client.NewClient(host)
	hosts, err := vClient.Hosts()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, host := range hosts {
		fmt.Println(host)
	}
}
