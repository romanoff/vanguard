package command

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/romanoff/vanguard/client"
	"github.com/romanoff/vanguard/config"
)

func NewUpCommand() cli.Command {
	return cli.Command{
		Name:  "up",
		Usage: "start cluster containers based on vanguard.yml config",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "dry",
				Usage: "show containers lunching process by tiers",
			},
		},
		Action: func(c *cli.Context) {
			upCommandFunc(c)
		},
	}
}

func upCommandFunc(c *cli.Context) {
	config, err := config.ParseConfig("vanguard.yml")
	if err != nil {
		fmt.Println(err)
		return
	}
	tiers, err := config.GetTiers()
	if err != nil {
		fmt.Println(err)
		return
	}
	if c.Bool("dry") {
		ShowTiers(tiers)
		return
	}
	for _, tier := range tiers {
		for _, server := range tier.Servers {
			vClient := client.NewClient(server.Hostname)
			for _, container := range server.Containers {
				for i := 0; i < container.GetCount(); i++ {
					serverContainer, err := vClient.Run(container.Image, container.Tag, "", nil)
					if err != nil {
						fmt.Println(err)
						return
					}
					fmt.Println(serverContainer)
				}
			}
		}
	}
}

func ShowTiers(tiers []*config.Tier) {
	for i, tier := range tiers {
		fmt.Printf("Tier %v:\n", i+1)
		fmt.Println("------------------")
		fmt.Println(tier)
	}
}
