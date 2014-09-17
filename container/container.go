package container

import (
	"errors"
	"os/exec"
	"time"
)

type Container struct {
	Name        string
	Tag         string
	ImageId     string
	ContainerId string
	Variables   map[string]string
	Ip          string
	CreatedAt   time.Time
}

func (self *Container) Run() error {
	if self.ContainerId != "" {
		return errors.New("Container is already running - " + self.ContainerId)
	}
	var err error
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
		self.ContainerId = string(containerId)
		self.CreatedAt = time.Now()
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
