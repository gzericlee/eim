package storage

import (
	"eim/internal/storage/rpc"
)

func InitStorageServer(ip string, port int, etcdEndpoints []string) error {
	return rpc.StartServer(ip, port, etcdEndpoints)
}
