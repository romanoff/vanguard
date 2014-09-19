package command

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/romanoff/vanguard/client"
	"github.com/romanoff/vanguard/config"
)

func NewPsCommand() cli.Command {
	return cli.Command{
		Name:  "ps",
		Usage: "shows running containers",
		Action: func(c *cli.Context) {
			psCommandFunc(c, false)
		},
	}
}

func NewPsckCommand() cli.Command {
	return cli.Command{
		Name:  "psck",
		Usage: "shows locally running containers and checks if they are running",
		Action: func(c *cli.Context) {
			psCommandFunc(c, true)
		},
	}
}

func psCommandFunc(c *cli.Context, check bool) {
	hostname := c.Args().First()
	if hostname == "" {
		hostname = "127.0.0.1"
	}
	cfg, _ := config.ParseConfig("vanguard.yml")
	if cfg != nil {
		for _, server := range cfg.Servers {
			if server.Label == hostname {
				hostname = server.Hostname
			}
		}
	}
	vClient := client.NewClient(hostname)
	containers, err := vClient.Index(check)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, c := range containers {
		fmt.Println(c)
	}
}
