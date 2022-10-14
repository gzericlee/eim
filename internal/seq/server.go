package seq

import (
	"eim/internal/seq/rpc"
)

func InitSeqServer(ip string, port int, etcdEndpoints []string) error {
	return rpc.StartServer(ip, port, etcdEndpoints)
}
