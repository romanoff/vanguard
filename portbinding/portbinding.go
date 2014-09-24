package portbinding

func New(port, host, hostPort string) (*PortBinding, error) {
	return &PortBinding{Port: port, Host: host, HostPort: hostPort}, nil
}

type PortBinding struct {
	Port string
	Host string
	HostPort string
}

func (self *PortBinding) Start() error {
	return nil
}

func (self *PortBinding) Stop() error {
	return nil
}
