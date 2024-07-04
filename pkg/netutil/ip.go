package netutil

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func InternalIP() (string, error) {
	inters, err := net.Interfaces()
	if err != nil {
		return "", fmt.Errorf("get net interfaces -> %w", err)
	}
	for _, inter := range inters {
		if inter.Flags&net.FlagUp != 0 && !strings.HasPrefix(inter.Name, "lo") {
			addrs, err := inter.Addrs()
			if err != nil {
				continue
			}
			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						return ipnet.IP.String(), nil
					}
				}
			}
		}
	}
	return "", fmt.Errorf("no local ipv4 address found")
}

func PodIP() (string, error) {
	podIp := os.Getenv("POD_IP")
	if podIp == "" {
		return "", fmt.Errorf("env POD_IP is empty")
	}
	return podIp, nil
}
