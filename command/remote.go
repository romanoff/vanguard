package command

import (
	"github.com/codegangsta/cli"
)

func NewRemoteCommand() cli.Command {
	return cli.Command{
		Name:  "remote",
		Usage: "options for remote storage",
		Subcommands: []cli.Command{
			{
				Name:  "push",
				Usage: "push file to remote storage",
				Action: func(c *cli.Context) {
					PushFileToRemote(c)
				},
			},
			{
				Name:  "pull",
				Usage: "pull file from remote storage",
				Action: func(c *cli.Context) {
					PullFileFromRemote(c)
				},
			},
		},
	}
}

func PushFileToRemote(c *cli.Context) {

}

func PullFileFromRemote(c *cli.Context) {

}
