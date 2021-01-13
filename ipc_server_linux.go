// +build linux

package commander

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

// NewIPCServer ...
func NewIPCServer() *IPCServer {
	uinxServ := &IPCServer{
		UnixSockPath: fmt.Sprintf("/tmp/%s.sock", os.Args[0]),
	}

	return uinxServ
}

type IPCServer struct {
	UnixSockPath string
	listener     net.Listener
}

func (this *IPCServer) Start() {
	os.Remove(this.UnixSockPath)

	addr, err := net.ResolveUnixAddr("unix", this.UnixSockPath)
	if err != nil {
		log.Fatal(err)
	}

	this.listener, err = net.ListenUnix("unix", addr)
	if err != nil {
		log.Fatal(err)
	}

	http.Serve(this.listener, nil)
}

func (this *IPCServer) AddRoute(path string, handler http.HandlerFunc) {
	http.HandleFunc(path, handler)
}

func (this *IPCServer) Stop() {
	this.listener.Close()
	err := os.Remove(this.UnixSockPath)
	if err != nil {
		log.Println(err)
	}
}
