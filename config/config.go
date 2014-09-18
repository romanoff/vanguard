package config

import (
	"gopkg.in/yaml.v1"
	"io/ioutil"
)

type Config struct {
	Application string
	Servers     []*Server
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

func ParseConfig(path string) (*Config, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	err = yaml.Unmarshal(content, &config)
	return config, err
}
