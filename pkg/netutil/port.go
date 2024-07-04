package netutil

import (
	"errors"
	"fmt"
	"net"
)

func RandomPort() (int, error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port, nil
}

func IncrAvailablePort(defaultPort int) (int, error) {
	for port := defaultPort; port <= 65535; port++ {
		address := fmt.Sprintf(":%d", port)
		listener, err := net.Listen("tcp", address)
		if err == nil {
			listener.Close()
			return port, nil
		}
	}
	return 0, fmt.Errorf("no available port found")
}

func CheckPort(port int) bool {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		var err *net.OpError
		if errors.As(err, &err) && err.Op == "listen" {
			return true
		}
	}
	if listener != nil {
		listener.Close()
	}
	return false
}
