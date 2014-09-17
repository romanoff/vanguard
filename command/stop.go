package command

import (
	"fmt"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"net/http"
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
	if containerId == "" {
		fmt.Println("Container id not specified")
		return
	}
	client := &http.Client{}
	req, _ := http.NewRequest("DELETE", "http://127.0.0.1:2728/containers/"+containerId, nil)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("vanguard agent is not running")
		return
	}
	content, _ := ioutil.ReadAll(resp.Body)
	fmt.Print(string(content))
}
