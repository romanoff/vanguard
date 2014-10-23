package command

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/romanoff/vanguard/config"
	"github.com/romanoff/vanguard/container"
)

func NewBuildCommand() cli.Command {
	return cli.Command{
		Name:  "build",
		Usage: "build containers that have dockerfile config if they don't exist",
		Action: func(c *cli.Context) {
			buildMissingImages(c, false)
		},
	}
}

func NewRebuildCommand() cli.Command {
	return cli.Command{
		Name:  "rebuild",
		Usage: "rebuild containers that have dockerfile config",
		Action: func(c *cli.Context) {
			buildMissingImages(c, true)
		},
	}
}

func buildMissingImages(c *cli.Context, force bool) {
	label := c.Args().First()
	config, err := config.ParseConfig("vanguard.yml")
	if err != nil {
		fmt.Println(err)
		return
	}
	messages := []string{}
	for _, server := range config.Servers {
		for _, cont := range server.Containers {
			_, err = container.GetImageId(cont.Image, cont.Tag)
			if force || err != nil {
				if cont.Dockerfile != "" && (label == "" || cont.Label == label) {
					_, err := container.BuildImage(cont.Dockerfile, cont.Image)
					if err != nil {
						messages = append(messages, fmt.Sprintf("%v image had errors: %v ", cont.Image, err))
					} else {
						messages = append(messages, fmt.Sprintf("%v image has been built successfully ", cont.Image))
					}
				}
			}
		}
	}
	if len(messages) == 0 {
		fmt.Println("No images to build")
	} else {
		for _, message := range messages {
			fmt.Println(message)
		}
	}
}
