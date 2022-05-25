package util

import (
	"net"
)

func GetLocalIpV4() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:8")
	if err != nil {
		return "", nil
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr).IP.To4().String()
	return localAddr, nil
}
