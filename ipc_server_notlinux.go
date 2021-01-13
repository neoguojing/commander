// +build !linux

package commander

import (
	"net"
	"net/http"
)

// NewIPCServer ...
func NewIPCServer() *IPCServer {
	uinxServ := &IPCServer{}

	return uinxServ
}

type IPCServer struct {
	listener net.Listener
}

func (this *IPCServer) Start() {
	http.Serve(this.listener, nil)
}

func (this *IPCServer) AddRoute(path string, handler http.HandlerFunc) {
	http.HandleFunc(path, handler)
}

func (this *IPCServer) Stop() {

}
