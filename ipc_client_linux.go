package commander

import (
	"net"
	"net/http"
)

//GetIPCClient ...
func GetIPCClient(uinxSockPath string) (*http.Client, error) {
	unixDial := func(network, addr string) (net.Conn, error) {
		return net.Dial("unix", uinxSockPath)
	}

	tr := &http.Transport{
		Dial: unixDial,
	}
	client := &http.Client{Transport: tr}
	return client, nil
}
