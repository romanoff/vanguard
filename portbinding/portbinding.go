package portbinding

import (
	"io"
	"net"
)

func New(port, host, hostPort string) (*PortBinding, error) {
	return &PortBinding{Port: port, Backends: []*Backend{{Host: host, Port: hostPort}}}, nil
}

type PortBinding struct {
	Port     string
	stop     chan bool
	Backends []*Backend
	Running  bool
}

type Backend struct {
	Host string
	Port string
}

func (self *Backend) String() string {
	return self.Host + ":" + self.Port
}

func (self *PortBinding) Start() error {
	conn, err := net.Listen("tcp", ":"+self.Port)
	if err != nil {
		return err
	}
	self.stop = make(chan bool)
	self.Running = true
	connectionChannel := getConnectionChannel(conn)
Loop:
	for {
		select {
		case <-self.stop:
			break Loop
		case hostConn := <-connectionChannel:
			handleConnection(hostConn, self)
		}
	}
	conn.Close()
	return nil
}

func (self *PortBinding) Stop() error {
	self.stop <- true
	self.Running = false
	return nil
}

func (self *PortBinding) String() string {
	content := self.Port + " binded to"
	for i, backend := range self.Backends {
		if i != 0 {
			content += ","
		}
		content += " " + backend.String()
	}
	return content
}

func getConnectionChannel(conn net.Listener) <-chan net.Conn {
	out := make(chan net.Conn)
	go func() {
		for {
			hostConn, _ := conn.Accept()
			out <- hostConn
		}
	}()
	return out
}

func copy(wc io.WriteCloser, r io.Reader) {
	defer wc.Close()
	io.Copy(wc, r)
}

func handleConnection(conn net.Conn, pb *PortBinding) {
	if len(pb.Backends) == 0 {
		return
	}
	remote, err := net.Dial("tcp", pb.Backends[0].String())
	if err == nil {
		go copy(conn, remote)
		go copy(remote, conn)
	}
}
