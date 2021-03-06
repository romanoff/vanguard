package config

import (
	"errors"
	"fmt"
	"github.com/romanoff/vanguard/remote"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Application string
	Servers     []*Server
	Remote      *Remote
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

func (self *Config) GetRemote() (remote.Remote, error) {
	if self.Remote == nil {
		return nil, errors.New("No remote found")
	}
	r := self.Remote
	if r.Type != "s3" {
		return nil, errors.New("Right now only s3 remote is supported")
	}
	return &remote.S3{Bucket: r.Bucket, AccessKey: r.Access_Key, SecretKey: r.Secret_Key, Region: r.Region}, nil
}

func (self *Config) PullRemoteFiles(replaceExisting bool) error {
	if self.Remote == nil {
		return nil
	}
	if self.Remote.Files == nil || len(self.Remote.Files) == 0 {
		return nil
	}
	r, err := self.GetRemote()
	if err != nil {
		return err
	}
	files := self.Remote.Files
	for _, file := range files {
		destination := file.GetDestination()
		if !replaceExisting {
			if _, err := os.Stat(destination); err == nil {
				continue
			}
		}
		fmt.Printf("downloading remote file %v into %v\n", file.Name, destination)
		err = r.Pull(file.Name, destination)
		if err != nil {
			return errors.New(fmt.Sprintf("Error while downloading file from remote '%v' to '%v': %v", file.Name, destination, err))
		}
	}
	return nil
}

type Server struct {
	Label      string `yml:"label,omitempty"`
	Hostname   string
	Expose     []string `yml:"expose,omitempty"`
	Containers []*Container
}

type Remote struct {
	Type       string        `yml:"type,omitempty"`
	Bucket     string        `yml:"bucket,omitempty"`
	Access_Key string        `yml:"access_key,omitempty"`
	Secret_Key string        `yml:"secret_key,omitempty"`
	Region     string        `yml:"region,omitempty"`
	Files      []*RemoteFile `yml:"files,omitempty"`
}

type RemoteFile struct {
	Name string `yml:"name,omitempty"`
	Path string `yml:"path,omitempty"`
}

func (self *RemoteFile) GetDestination() string {
	destination := self.Path
	if strings.HasSuffix(destination, "/") || destination == ""{
		destination += filepath.Base(self.Name)
	}
	return destination
}

type Container struct {
	Label      string   `yml:"label,omitempty"`
	Image      string   `yml:"image,omitempty"`
	Tag        string   `yml:"tag,omitempty"`
	ImageId    string   `yml:"image_id,omitempty"`
	Count      int      `yml:"count,omitempty"`
	Links      []string `yml:"links,omitempty"`
	DNS        []string `yml:"dns,omitempty"`
	Expose     []string `yml:"expose,omitempty"`
	Variables  []string `yml:"variables,omitempty"`
	Volumes    []string `yml:"volumes,omitempty"`
	Command    string   `yml:"command,omitempty"`
	Dockerfile string   `yml:"dockerfile,omitempty"`
}

func (self *Container) String() string {
	return fmt.Sprintf("%v - %v", self.Name(), self.GetCount())
}

func (self *Container) GetCount() int {
	count := self.Count
	if count == 0 {
		count = 1
	}
	return count
}

func (self *Container) Name() string {
	if self.Label != "" {
		return self.Label
	}
	name := self.Image
	if self.Tag != "" {
		name += "_" + self.Tag
	}
	if name == "" {
		return self.ImageId
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
	content = []byte(os.ExpandEnv(string(content)))
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
