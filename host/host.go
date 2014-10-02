package host

import (
	"encoding/json"
	"errors"
	"github.com/romanoff/vanguard/portbinding"
	"github.com/romanoff/vanguard/storage"
	"os"
)

func New(hostname, externalInterface, internalInterface, externalIp, internalIp string) (*Host, error) {
	var err error
	if hostname == "" {
		hostname, err = os.Hostname()
		if err != nil {
			return nil, err
		}
	}
	if externalIp == "" {
		externalIp, err = GetIpAddress(externalInterface)
		if err != nil {
			return nil, err
		}
	}
	if internalIp == "" {
		internalIp, err = GetIpAddress(internalInterface)
		if err != nil {
			return nil, err
		}
	}
	host, err := GetHost(internalIp)
	if err == nil {
		for _, binding := range host.PortBindings {
			go binding.Start()
		}
	}
	if host == nil {
		host = &Host{
			Hostname:     hostname,
			ExternalIp:   externalIp,
			InternalIp:   internalIp,
			PortBindings: []*portbinding.PortBinding{},
		}
	}
	currentHost = host
	return host, nil
}

type Host struct {
	Hostname     string
	ExternalIp   string
	InternalIp   string
	PortBindings []*portbinding.PortBinding
}

func (self *Host) Persist() error {
	db := storage.GetStorage()
	jsonBytes, err := json.Marshal(self)
	if err != nil {
		return err
	}
	err = db.Set(self.InternalIp, string(jsonBytes))
	if err != nil {
		return err
	}
	return nil
}

func (self *Host) String() string {
	return "hostname: " + self.Hostname + " ip: " + self.ExternalIp + ", weave ip: " + self.InternalIp
}

func GetHost(ip string) (*Host, error) {
	db := storage.GetStorage()
	jsonString, err := db.Get(ip)
	if err != nil {
		return nil, err
	}
	var host *Host
	err = json.Unmarshal([]byte(jsonString), &host)
	if err != nil {
		return nil, err
	}
	return host, nil
}

var currentHost *Host

func GetCurrentHost() (*Host, error) {
	if currentHost == nil {
		return nil, errors.New("current host is not set")
	}
	return currentHost, nil
}

func GetByIp(ip string) (*Host, error) {
	db := storage.GetStorage()
	jsonString, err := db.Get(ip)
	if err != nil {
		return nil, err
	}
	var host *Host
	err = json.Unmarshal([]byte(jsonString), &host)
	if err != nil {
		return nil, err
	}
	return host, nil
}

func GetHosts() ([]*Host, error) {
	db := storage.GetStorage()
	ips, err := db.Keys("10.0.1.")
	if err != nil {
		return nil, err
	}
	hosts := []*Host{}
	for _, ip := range ips {
		host, err := GetByIp(ip)
		if err != nil {
			return nil, err
		}
		if host.ExternalIp != "" && host.InternalIp != "" {
			hosts = append(hosts, host)
		}
	}
	return hosts, nil
}
