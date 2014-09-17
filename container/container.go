package container

import (
	"errors"
	"github.com/fsouza/go-dockerclient"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Container struct {
	Name        string            `json:"name,omitempty"`
	Tag         string            `json:"tag,omitempty"`
	ImageId     string            `json:"image_id"`
	ContainerId string            `json:"container_id"`
	Variables   map[string]string `json:"variables,omitempty"`
	Ip          string            `json:"ip"`
	Hostname    string            `json:"hostname"`
	CreatedAt   time.Time         `json:"created_at,omitempty"`
}

func (self *Container) String() string {
	name := ""
	if self.Name != "" {
		name += self.Name
		if self.Tag != "" {
			name += ":" + self.Tag
		}
	}
	if name != "" {
		name += ", "
	}
	name += "imageid: " + self.ImageId + ", ip: " + self.Ip
	return name
}

func (self *Container) Run() error {
	if self.ContainerId != "" {
		return errors.New("Container is already running - " + self.ContainerId)
	}
	var err error
	self.Hostname, err = os.Hostname()
	if err != nil {
		return err
	}
	self.Ip, err = ReserveIp()
	if err != nil {
		return err
	}
	if self.ImageId == "" {
		self.ImageId, err = GetImageId(self.Name, self.Tag)
		if err != nil {
			FreeIp(self.Ip)
			return err
		}
	}
	err = self.runWithWeave()
	if err != nil {
		FreeIp(self.Ip)
	}
	return err
}

func (self *Container) runWithWeave() error {
	args := []string{"run", self.Ip, "-i", "-t"}
	if self.Variables != nil {
		for key, value := range self.Variables {
			args = append(args, "-e", key+"="+value)
		}
	}
	args = append(args, self.ImageId)
	containerId, err := exec.Command("weave", args...).Output()
	if err == nil {
		self.ContainerId = strings.TrimSpace(string(containerId))
		self.CreatedAt = time.Now()
		Persist(self)
		log.Println("Started container " + self.String())
	}
	return err
}

func (self *Container) Stop() error {
	client, err := GetDockerClient()
	if err != nil {
		return err
	}
	if self.ContainerId == "" {
		return errors.New("Trying to stop container without container id")
	}
	err = client.StopContainer(self.ContainerId, 5)
	if err == nil {
		FreeIp(self.Ip)
		log.Println("Stopped container " + self.String())
		err = client.RemoveContainer(docker.RemoveContainerOptions{ID: self.ContainerId})
		if err != nil {
			return err
		}
	}
	return err
}

func GetImageId(name string, tag string) (string, error) {
	client, err := GetDockerClient()
	if err != nil {
		return "", err
	}
	imageName := name
	if tag != "" {
		name += ":" + tag
	}
	image, err := client.InspectImage(imageName)
	if err != nil {
		return "", err
	}
	return image.ID, nil
}
