package command

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/romanoff/vanguard/client"
	"strings"
)

func NewExposeCommand() cli.Command {
	return cli.Command{
		Name:  "expose",
		Usage: "host:port internal_host:port",
		Action: func(c *cli.Context) {
			exposeCommandFunc(c)
		},
	}
}

func exposeCommandFunc(c *cli.Context) {
	if len(c.Args()) != 2 {
		showExposeErrorMessage()
		return
	}
	slice := strings.Split(c.Args()[0], ":")
	if len(slice) == 0 || len(slice) > 2 {
		showExposeErrorMessage()
		return
	}
	host := "127.0.0.1"
	port := slice[0]
	if len(slice) == 2 {
		host = slice[0]
		port = slice[1]
	}
	slice = strings.Split(c.Args()[1], ":")
	if len(slice) != 2 {
		showExposeErrorMessage()
		return
	}
	internalHost := slice[0]
	internalPort := slice[1]
	vClient := client.NewClient(host)
	portbinding, err := vClient.Expose(port, internalHost, internalPort)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(portbinding)
}

func showExposeErrorMessage() {
	fmt.Println("Expected to get 2 arguments - host:port internal_host:port")
}
