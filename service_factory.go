package commander

import (
	"sync"
)

var (
	Servers *ServerFactory
)

type ServerFactory struct {
	servers []IServer
	sync.Mutex
}

func NewServerFactory() *ServerFactory {
	return &ServerFactory{
		servers: make([]IServer, 0),
	}
}

func (this *ServerFactory) Register(server IServer) {
	this.Lock()
	defer this.Unlock()

	this.servers = append(this.servers, server)
}

func (this *ServerFactory) Dump() []IServer {
	return this.servers
}

func init() {
	Servers = NewServerFactory()
}
