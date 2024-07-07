package httputil

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"time"

	"go.uber.org/zap"

	"eim/pkg/log"
)

type Proxy struct {
	Address []string `json:"address"`
	Cert    string   `json:"cert"`
	Key     string   `json:"key"`
}

func NewProxyHttpClient(proxy Proxy) (*http.Client, error) {
	var transport *http.Transport
	var httpClient *http.Client

	transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       120 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
		ResponseHeaderTimeout: 120 * time.Second,
	}

	if proxy.Cert != "" && proxy.Key != "" {
		cert, err := tls.X509KeyPair([]byte(proxy.Cert), []byte(proxy.Key))
		if err != nil {
			log.Error("Error load cert and key", zap.Error(err))
			return httpClient, err
		}
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true, Certificates: []tls.Certificate{cert}}
	} else {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	if len(proxy.Address) > 0 {
		tunnel := ""
		for _, addr := range proxy.Address {
			if addr == "" {
				continue
			}
			tunnel += addr + ","
		}
		if tunnel != "" {
			tunnel = strings.TrimRight(tunnel, ",")
			proxyAddr := proxy.Address[len(proxy.Address)-1]
			transport.Proxy = func(req *http.Request) (*url.URL, error) {
				req.Header.Add("X-Tunnel", tunnel)
				log.Debug("Configure proxy address:", zap.String("address", proxyAddr), zap.String("x-tunnel", tunnel))
				return url.Parse(proxyAddr)
			}
		}
	}

	httpClient = &http.Client{Transport: transport}

	return httpClient, nil
}
