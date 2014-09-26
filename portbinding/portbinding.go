package portbinding

import (
	"io"
	"math/rand"
	"net"
	"time"
)

func New(port string) (*PortBinding, error) {
	return &PortBinding{Port: port, Backends: []*Backend{}}, nil
}

type Backend struct {
	Host string
	Port string
}

func (self *Backend) String() string {
	return self.Host + ":" + self.Port
}

type PortBinding struct {
	Port     string
	stop     chan bool
	Backends []*Backend
	Running  bool
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

func (self *PortBinding) GetBackend(host, port string) (int, *Backend) {
	for i, backend := range self.Backends {
		if backend.Host == host && backend.Port == port {
			return i, backend
		}
	}
	return -1, nil
}

func (self *PortBinding) AddBackend(host, port string) {
	_, backend := self.GetBackend(host, port)
	if backend == nil {
		self.Backends = append(self.Backends, &Backend{Host: host, Port: port})
	}
}

func (self *PortBinding) RemoveBackend(host, port string) {
	i, _ := self.GetBackend(host, port)
	if i != -1 {
		copy(self.Backends[i:], self.Backends[i+1:])
		self.Backends[len(self.Backends)-1] = nil
		self.Backends = self.Backends[:len(self.Backends)-1]
	}
}

func (self *PortBinding) GetRandomBackend() *Backend {
	if len(self.Backends) == 0 {
		return nil
	}
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(len(self.Backends))
	return self.Backends[i]
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

func copyandclose(wc io.WriteCloser, r io.Reader) {
	defer wc.Close()
	io.Copy(wc, r)
}

func handleConnection(conn net.Conn, pb *PortBinding) {
	backend := pb.GetRandomBackend()
	if backend == nil {
		return
	}
	remote, err := net.Dial("tcp", backend.String())
	if err == nil {
		go copyandclose(conn, remote)
		go copyandclose(remote, conn)
	}
}
