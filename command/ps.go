package command

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/romanoff/vanguard/container"
	"io/ioutil"
	"net/http"
)

func NewPsCommand() cli.Command {
	return cli.Command{
		Name:  "ps",
		Usage: "shows running containers",
		Action: func(c *cli.Context) {
			psCommandFunc(c)
		},
	}
}

func psCommandFunc(c *cli.Context) {
	resp, err := http.Get("http://127.0.0.1:2728/containers")
	if err != nil {
		fmt.Println("vanguard agent is not running")
		return
	}
	content, _ := ioutil.ReadAll(resp.Body)
	var containers []*container.Container
	err = json.Unmarshal(content, &containers)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, c := range containers {
		fmt.Println(c)
	}
}
