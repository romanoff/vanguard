package client

import (
	"encoding/json"
	"errors"
	"fmt"
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

func (self *Client) Run(label string, name string, tag string, imageId string, variables map[string]string) (*container.Container, error) {
	values := url.Values{"label": {label}, "name": {name}, "tag": {tag}, "image_id": {imageId}}
	if variables != nil {
		for key, value := range variables {
			if key != "label" && key != "name" && key != "tag" && key != "image_id" {
				values.Add(key, value)
			}
		}
	}
	resp, err := http.PostForm("http://"+self.Hostname+":2728/containers", values)
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

func (self *Client) Index(check bool) ([]*container.Container, error) {
	url := "http://" + self.Hostname + ":2728/containers"
	if check {
		url += "?check=true"
	}
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.New("vanguard agent is not running on host " + self.Hostname)
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var containers []*container.Container
	err = json.Unmarshal(content, &containers)
	if err != nil {
		return nil, err
	}
	return containers, nil
}

func (self *Client) Stop(containerId string) error {
	if containerId == "" {
		return errors.New("Container id to stop not specified")
	}
	client := &http.Client{}
	req, _ := http.NewRequest("DELETE", "http://"+self.Hostname+":2728/containers/"+containerId, nil)
	resp, err := client.Do(req)
	if err != nil {
		return errors.New("vanguard agent is not running on host " + self.Hostname)
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var data map[string]interface{}
	err = json.Unmarshal(content, &data)
	if err != nil {
		return err
	}
	if data["error"] != nil {
		return errors.New(fmt.Sprintf("%v", data["error"]))
	}
	return nil
}
