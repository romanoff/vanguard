package command

import (
	"errors"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/romanoff/vanguard/client"
	"github.com/romanoff/vanguard/config"
	"github.com/romanoff/vanguard/container"
	"regexp"
	"strings"
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
	for _, server := range config.Servers {
		c := client.NewClient(server.Hostname)
		bindings, err := c.Bindings()
		for _, binding := range bindings {
			err = c.Hide(binding.Port, "", "")
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
	manager := &ContainerManager{
		EnvVariables:       make(map[string]string),
		Clients:            make(map[string]*client.Client),
		RunningContainers:  make(map[string][]*container.Container),
		UsedContainerNames: make(map[string][]string),
	}
	for _, tier := range tiers {
		for _, server := range tier.Servers {
			for _, cont := range server.Containers {
				manager.Manage(server.Hostname, cont)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
	manager.StopUnused()
}

type ContainerManager struct {
	EnvVariables       map[string]string
	Clients            map[string]*client.Client
	RunningContainers  map[string][]*container.Container
	UsedContainerNames map[string][]string
}

func (self *ContainerManager) GetRunningContainers(host string) ([]*container.Container, error) {
	if self.RunningContainers[host] != nil {
		return self.RunningContainers[host], nil
	}
	c := self.Clients[host]
	if c == nil {
		self.Clients[host] = client.NewClient(host)
		c = self.Clients[host]
	}
	containers, err := c.Index(true)
	if err != nil {
		return nil, err
	}
	self.RunningContainers[host] = containers
	return containers, nil
}

func (self *ContainerManager) GetRunningContainersByName(host string, name string) ([]*container.Container, error) {
	runningContainers, err := self.GetRunningContainers(host)
	if err != nil {
		return nil, err
	}
	containers := []*container.Container{}
	for _, c := range runningContainers {
		if c.LabelName() == name {
			containers = append(containers, c)
		}
	}
	return containers, nil
}

func (self *ContainerManager) Manage(host string, container *config.Container) error {
	if config.NotIn(self.UsedContainerNames[host], container.Name()) {
		self.UsedContainerNames[host] = append(self.UsedContainerNames[host], container.Name())
	}
	err := self.Launch(host, container)
	if err != nil {
		return err
	}
	return self.StopExtra(host, container)
}

func (self *ContainerManager) Expose(host string, cont *config.Container, serverContainer *container.Container) error {
	vClient := self.Clients[host]
	if cont.Expose != nil && len(cont.Expose) > 0 {
		for _, expose := range cont.Expose {
			ports := strings.Split(expose, ":")
			if len(ports) != 2 {
				return errors.New(fmt.Sprintf("Invalid expose syntax: %v", expose))
			}
			hostPort := ports[0]
			containerPort := ports[1]
			binding, err := vClient.Expose(hostPort, serverContainer.Ip, containerPort)
			if err != nil {
				return err
			}
			fmt.Println(host + ":" + binding.String())
		}
	}
	return nil
}

func (self *ContainerManager) Launch(host string, cont *config.Container) error {
	runningContainers, err := self.GetRunningContainersByName(host, cont.Name())
	if err != nil {
		return err
	}
	containersToLaunch := cont.GetCount() - len(runningContainers)
	vClient := self.Clients[host]
	if containersToLaunch > 0 {
		for i := 0; i < containersToLaunch; i++ {
			variables := make(map[string]string)
			for _, link := range cont.Links {
				if _, ok := self.EnvVariables[link]; ok {
					variables[link] = self.EnvVariables[link]
				}
			}
			if cont.Variables != nil {
				for _, variable := range cont.Variables {
					values := strings.Split(variable, ":")
					if len(values) == 2 {
						variables[values[0]] = self.VariableValue(values[1])
					}
				}
			}
			serverContainer, err := vClient.Run(cont.Name(), cont.Image, cont.Tag, cont.ImageId, variables, cont.DNS)
			if err != nil {
				return err
			}
			if _, ok := self.EnvVariables[cont.Name()]; !ok {
				self.EnvVariables[cont.Name()] = serverContainer.Ip
			}
			fmt.Println(serverContainer)
			err = self.Expose(host, cont, serverContainer)
			if err != nil {
				return err
			}
		}
	} else if len(runningContainers) > 0 {
		for _, rc := range runningContainers {
			err = self.Expose(host, cont, rc)
			if err != nil {
				return err
			}
		}
		if _, ok := self.EnvVariables[cont.Name()]; !ok {
			self.EnvVariables[cont.Name()] = runningContainers[0].Ip
		}
	}
	return nil
}

var ipVariableRegexp *regexp.Regexp = regexp.MustCompile("ip\\((.*)\\)")

func (self *ContainerManager) VariableValue(value string) string {
	label := ipVariableRegexp.FindAllStringSubmatch(value, -1)
	if label == nil {
		return value
	}
	labelValue := label[0][1]
	if self.EnvVariables[labelValue] != "" {
		return self.EnvVariables[labelValue]
	}
	return value
}

func (self *ContainerManager) StopExtra(host string, cont *config.Container) error {
	runningContainers, err := self.GetRunningContainersByName(host, cont.Name())
	if err != nil {
		return err
	}
	containersToStop := len(runningContainers) - cont.GetCount()
	vClient := self.Clients[host]
	if containersToStop > 0 {
		for i := 0; i < containersToStop; i++ {
			c := runningContainers[len(runningContainers)-1-i]
			err = vClient.Stop(c.ContainerId)
			if err != nil {
				return err
			}
			fmt.Println("stopped " + c.String())
		}
	}
	return nil
}

func (self *ContainerManager) StopUnused() error {
	for host, containers := range self.RunningContainers {
		for _, container := range containers {
			if config.NotIn(self.UsedContainerNames[host], container.LabelName()) {
				vClient := self.Clients[host]
				err := vClient.Stop(container.ContainerId)
				if err != nil {
					return err
				}
				fmt.Println("stopped " + container.String())
			}
		}
	}
	return nil
}

func ShowTiers(tiers []*config.Tier) {
	for i, tier := range tiers {
		fmt.Printf("Tier %v:\n", i+1)
		fmt.Println("------------------")
		fmt.Println(tier)
	}
}
