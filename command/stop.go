package command

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/romanoff/vanguard/client"
	"github.com/romanoff/vanguard/config"
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
	hostname := "127.0.0.1"
	if len(c.Args()) > 1 {
		hostname = getHostname(c.Args()[1])
	}
	if containerId == "all" {
		cfg, _ := config.ParseConfig("vanguard.yml")
		if cfg != nil && hostname == "127.0.0.1" {
			for _, server := range cfg.Servers {
				err := stopAllOnHost(server.Hostname)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		} else {
			err := stopAllOnHost(hostname)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		return
	}
	vClient := client.NewClient(hostname)
	err := vClient.Stop(containerId)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("success")
}

func stopAllOnHost(hostname string) error {
	vClient := client.NewClient(hostname)
	containers, err := vClient.Index(true)
	if err != nil {
		return err
	}
	for _, container := range containers {
		err = vClient.Stop(container.ContainerId)
		if err != nil {
			return err
		}
		fmt.Println("stopped " + container.String())
	}
	return nil
}
