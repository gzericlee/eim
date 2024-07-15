package client

import rpcxclient "github.com/smallnest/rpcx/client"

type FileClient struct {
	*rpcxclient.XClientPool
}
