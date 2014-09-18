package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v1"
	"io/ioutil"
)

type Config struct {
	Application string
	Servers     []*Server
}

func (self *Config) GetTiers() ([]*Tier, error) {
	containersCount := 0
	for _, server := range self.Servers {
		containersCount += len(server.Containers)
	}
	launchedContainers := []string{}
	tiersContainersCount := 0
	tiers := []*Tier{}
	for {
		tier := &Tier{}
		tierLaunchedContainers := []string{}
		for _, server := range self.Servers {
			for _, container := range server.Containers {
				if NotIn(launchedContainers, container.Name()) &&
					container.CanBeLaunchedWith(launchedContainers) {
					tiersContainersCount++
					tier.AddContainer(server, container)
					if NotIn(tierLaunchedContainers, container.Name()) {
						tierLaunchedContainers = append(tierLaunchedContainers, container.Name())
					}
				}
			}
		}
		tiers = append(tiers, tier)
		launchedContainers = append(launchedContainers, tierLaunchedContainers...)
		if tiersContainersCount == containersCount {
			break
		}
		if len(tier.Servers) == 0 {
			return nil, errors.New("circular links dependency in vanguard.yml")
		}
	}
	return tiers, nil
}

type Server struct {
	Hostname   string
	Containers []*Container
}

type Container struct {
	Image string   `yml:"image,omitempty"`
	Tag   string   `yml:"tag,omitempty"`
	Count int      `yml:"count,omitempty"`
	Links []string `yml:"links,omitempty"`
}

func (self *Container) String() string {
	count := self.Count
	if count == 0 {
		count = 1
	}
	return fmt.Sprintf("%v - %v", self.Name(), count)
}

func (self *Container) Name() string {
	name := self.Image
	if self.Tag != "" {
		name += ":" + self.Tag
	}
	return name
}

func (self *Container) CanBeLaunchedWith(names []string) bool {
	var dependencyMissing bool
	for _, link := range self.Links {
		dependencyMissing = true
		for _, name := range names {
			if name == link {
				dependencyMissing = false
			}
		}
	}
	return !dependencyMissing
}

type Tier struct {
	Servers []*Server
}

func (self *Tier) AddContainer(server *Server, container *Container) {
	var usedServer *Server
	for _, s := range self.Servers {
		if s.Hostname == server.Hostname {
			usedServer = s
		}
	}
	if usedServer == nil {
		usedServer = &Server{Hostname: server.Hostname}
		self.Servers = append(self.Servers, usedServer)
	}
	usedServer.Containers = append(usedServer.Containers, container)
}

func (self *Tier) String() string {
	content := ""
	for _, server := range self.Servers {
		content += server.Hostname + ":\n"
		for _, container := range server.Containers {
			content += container.String() + "\n"
		}
	}
	return content
}

func ParseConfig(path string) (*Config, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	err = yaml.Unmarshal(content, &config)
	return config, err
}

func NotIn(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return false
		}
	}
	return true
}