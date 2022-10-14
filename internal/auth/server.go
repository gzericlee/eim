package auth

import (
	"eim/internal/auth/rpc"
)

func InitAuthServer(ip string, port int, etcdEndpoints []string) error {
	return rpc.StartServer(ip, port, etcdEndpoints)
}
