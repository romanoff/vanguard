package command

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/romanoff/vanguard/client"
)

func NewHostCommand() cli.Command {
	return cli.Command{
		Name:  "host",
		Usage: "current host",
		Action: func(c *cli.Context) {
			hostCommandFunc()
		},
	}
}

func hostCommandFunc() {
	vClient := client.NewClient("127.0.0.1")
	host, err := vClient.Host()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(host)
}
