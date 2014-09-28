package host

import (
	"github.com/romanoff/vanguard/portbinding"
	"os"
)

func New(hostname, ip string) (*Host, error) {
	var err error
	if hostname == "" {
		hostname, err = os.Hostname()
		if err != nil {
			return nil, err
		}
	}
	if ip == "" {
		ip = "10.0.1.1"
	}
	return &Host{
		Hostname:     hostname,
		Ip:           ip,
		PortBindings: []*portbinding.PortBinding{},
	}, nil
}

type Host struct {
	Hostname     string
	Ip           string
	PortBindings []*portbinding.PortBinding
}

func (self *Host) Persist() error {
	return nil
}
