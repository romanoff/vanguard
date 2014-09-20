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
	Label       string            `json:"label,omitempty"`
	Name        string            `json:"name,omitempty"`
	Tag         string            `json:"tag,omitempty"`
	ImageId     string            `json:"image_id"`
	ContainerId string            `json:"container_id"`
	Variables   map[string]string `json:"variables,omitempty"`
	Ip          string            `json:"ip"`
	Hostname    string            `json:"hostname"`
	CreatedAt   time.Time         `json:"created_at,omitempty"`
	DNS         []string          `json:"dns,omitempty"`
}

func (self *Container) LabelName() string {
	if self.Label != "" {
		return self.Label
	}
	name := self.Name
	if self.Tag != "" {
		name += "_" + self.Tag
	}
	if name == "" {
		return self.ImageId
	}
	return name
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
	name += "containerid: " + self.ContainerId + " , ip: " + self.Ip + " , host: " + self.Hostname
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
	args := []string{"run", self.Ip + "/24", "-i", "-t"}
	if self.Variables != nil {
		for key, value := range self.Variables {
			args = append(args, "-e", strings.ToUpper(key)+"="+value)
		}
	}
	if self.DNS != nil {
		for _, dns := range self.DNS {
			args = append(args, "--dns", dns)
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

func (self *Container) Check() bool {
	hostname, err := os.Hostname()
	if err != nil {
		return false
	}
	if self.Hostname != hostname {
		return false
	}
	client, err := GetDockerClient()
	if err != nil {
		return false
	}
	if self.ContainerId == "" {
		return false
	}
	container, err := client.InspectContainer(self.ContainerId)
	if err != nil {
		FreeIp(self.Ip)
		return false
	}
	if container.State.Running != true {
		FreeIp(self.Ip)
		return false
	}
	return true
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
