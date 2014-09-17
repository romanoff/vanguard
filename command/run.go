package command

import (
	"fmt"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"net/http"
	"net/url"
)

func NewRunCommand() cli.Command {
	return cli.Command{
		Name:  "run",
		Usage: "run container",
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
	resp, err := http.PostForm("http://127.0.0.1:2728/containers",
		url.Values{"name": {name}})
	if err != nil {
		fmt.Println("vanguard agent is not running")
		return
	}
	content, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(content))
}
