package client

import (
	"encoding/json"
	"errors"
	"github.com/romanoff/vanguard/container"
	"io/ioutil"
	"net/http"
	"net/url"
)

func NewClient(hostname string) *Client {
	return &Client{Hostname: hostname}
}

type Client struct {
	Hostname string
}

func (self *Client) Run(name string, tag string, imageId string, variables map[string]string) (*container.Container, error) {
	resp, err := http.PostForm("http://"+self.Hostname+":2728/containers",
		url.Values{"name": {name}, "tag": {tag}, "image_id": {imageId}})
	if err != nil {
		return nil, errors.New("vanguard agent is not running on host " + self.Hostname)
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var c *container.Container
	err = json.Unmarshal(content, &c)
	if err != nil {
		return nil, err
	}
	return c, nil
}
