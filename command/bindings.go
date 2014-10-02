package command

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/romanoff/vanguard/client"
	"github.com/romanoff/vanguard/config"
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
	hosts := []string{}
	if host != "" {
		hosts = append(hosts, host)
	} else {
		host = "127.0.0.1"
		cfg, _ := config.ParseConfig("vanguard.yml")
		if cfg != nil && len(cfg.Servers) > 0 {
			host = cfg.Servers[0].Hostname
		}
		vClient := client.NewClient(host)
		remoteHosts, err := vClient.Hosts()
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, remoteHost := range remoteHosts {
			hosts = append(hosts, remoteHost.ExternalIp)
		}
	}
	for _, host := range hosts {
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
