package command

import (
	"github.com/codegangsta/cli"
)

func NewUpCommand() cli.Command {
	return cli.Command{
		Name:  "up",
		Usage: "start cluster containers based on vanguard.yml config",
		Action: func(c *cli.Context) {
			upCommandFunc()
		},
	}
}

func upCommandFunc() {

}
