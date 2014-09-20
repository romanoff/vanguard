package command

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/romanoff/vanguard/client"
	"strings"
)

func NewRunCommand() cli.Command {
	return cli.Command{
		Name:  "run",
		Usage: "run container",
		Flags: []cli.Flag{
			cli.StringSliceFlag{
				Name:  "e",
				Value: &cli.StringSlice{},
				Usage: "environment variables",
			},
			cli.StringSliceFlag{
				Name:  "dns",
				Value: &cli.StringSlice{},
				Usage: "dns servers for container",
			},
		},
		Action: func(c *cli.Context) {
			runCommandFunc(c)
		},
	}
}

func runCommandFunc(c *cli.Context) {
	name := c.Args().First()
	if name == "" {
		fmt.Println("No image has been specified")
		return
	}
	hostname := "127.0.0.1"
	if len(c.Args()) > 1 {
		hostname = getHostname(c.Args()[1])
	}
	vClient := client.NewClient(hostname)
	variables := make(map[string]string)
	for _, envVariable := range c.StringSlice("e") {
		envVar := strings.Split(envVariable, "=")
		if len(envVar) == 2 {
			variables[envVar[0]] = envVar[1]
		}
	}
	dnsServers := []string{}
	for _, dns := range c.StringSlice("dns") {
		dnsServers = append(dnsServers, dns)
	}
	container, err := vClient.Run(name, name, "", "", variables, dnsServers)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(container)
}
