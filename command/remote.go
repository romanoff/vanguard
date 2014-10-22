package command

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/romanoff/vanguard/config"
	"github.com/romanoff/vanguard/remote"
	"path/filepath"
	"strings"
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
			{
				Name:  "ls",
				Usage: "list remote files",
				Action: func(c *cli.Context) {
					ShowFilesList(c)
				},
			},
			{
				Name:  "rm",
				Usage: "remove remote file",
				Action: func(c *cli.Context) {
					RemoveFileFromRemote(c)
				},
			},
		},
	}
}

func PushFileToRemote(c *cli.Context) {
	if len(c.Args()) == 0 {
		fmt.Println("Usage: vanguard remote push <filepath> <destination>")
		return
	}
	path := c.Args()[0]
	destination := filepath.Base(path)
	if len(c.Args()) > 1 {
		destination = c.Args()[1]
		if strings.HasSuffix(destination, "/") {
			destination += filepath.Base(path)
		}
	}
	r, err := getRemote()
	if err != nil {
		fmt.Println("No remote found:")
		fmt.Println(err)
		return
	}
	err = r.Push(path, destination)
	if err != nil {
		fmt.Println(err)
	}
}

func PullFileFromRemote(c *cli.Context) {
	if len(c.Args()) == 0 {
		fmt.Println("Usage: vanguard remote pull <filepath> <destination>")
		return
	}
	filepath := c.Args()[0]
	destination := "."
	if len(c.Args()) > 1 {
		destination = c.Args()[1]
	}
	r, err := getRemote()
	if err != nil {
		fmt.Println("No remote found:")
		fmt.Println(err)
		return
	}
	err = r.Pull(filepath, destination)
	if err != nil {
		fmt.Println(err)
	}
}

func getRemote() (remote.Remote, error) {
	config, err := config.ParseConfig("vanguard.yml")
	if err != nil {
		return nil, err
	}
	return config.GetRemote()
}

func ShowFilesList(c *cli.Context) {
	prefix := ""
	if len(c.Args()) > 0 {
		prefix = c.Args()[0]
	}
	r, err := getRemote()
	if err != nil {
		fmt.Println("No remote found:")
		fmt.Println(err)
		return
	}
	files, err := r.FilesList(prefix)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, filename := range files {
		fmt.Println(filename)
	}
}

func RemoveFileFromRemote(c *cli.Context) {
	if len(c.Args()) == 0 {
		fmt.Println("Usage: vanguard remote remove <filepath>")
		return
	}
	filepath := c.Args()[0]
	r, err := getRemote()
	if err != nil {
		fmt.Println("No remote found:")
		fmt.Println(err)
		return
	}
	err = r.Remove(filepath)
	if err != nil {
		fmt.Println(err)
	}
}
