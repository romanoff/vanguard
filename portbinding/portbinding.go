package portbinding

import (
	"io"
	"net"
)

func New(port, host, hostPort string) (*PortBinding, error) {
	return &PortBinding{Port: port, Host: host, HostPort: hostPort}, nil
}

type PortBinding struct {
	Port     string
	Host     string
	HostPort string
	stop     chan bool
	Running bool
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
			break Loop;
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


func handleConnection(conn net.Conn, pb *PortBinding) {
	remote, err := net.Dial("tcp", pb.Host+":"+pb.HostPort)
	if err == nil {
		go io.Copy(conn, remote)
		go io.Copy(remote, conn)
	}
}
