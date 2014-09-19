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
	envVariables := make(map[string]string)
	for _, tier := range tiers {
		for _, server := range tier.Servers {
			vClient := client.NewClient(server.Hostname)
			for _, container := range server.Containers {
				err = checkContainerRunning(vClient, envVariables, server, container)
				if err != nil {
					fmt.Println(err)
					return
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

func checkContainerRunning(vClient *client.Client, envVariables map[string]string, server *config.Server, container *config.Container) error {
	for i := 0; i < container.GetCount(); i++ {
		variables := make(map[string]string)
		for _, link := range container.Links {
			if _, ok := envVariables[link]; ok {
				variables[link] = envVariables[link]
			}
		}
		serverContainer, err := vClient.Run(container.Name(), container.Image, container.Tag, container.ImageId, variables)
		return err
		if _, ok := envVariables[container.Name()]; !ok {
			envVariables[container.Name()] = serverContainer.Ip
		}
		fmt.Println(serverContainer)
	}
	return nil
}
