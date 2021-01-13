package commander

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

// RUNNER ...
type RUNNER func(cmd *cobra.Command, args []string)

// Option ...
type Option func(*cobra.Command)

// WithDescribe ...
func WithDescribe(desc string) Option {
	return func(cmd *cobra.Command) {
		cmd.Long = desc
	}
}

// Commander ...
type Commander struct {
	servers []IServer
	root    *cobra.Command
	cmds    map[string]*cobra.Command
	cmdServ *IPCServer
}

// NewCommander ...
func NewCommander() *Commander {

	c := &Commander{}
	c.root = rootCmd
	c.cmds = map[string]*cobra.Command{"stop": stopCmd}
	c.servers = make([]IServer, 0)
	c.cmdServ = NewIPCServer()
	return c
}

// AddCommand ...
func (c *Commander) AddCommand(parent, name string, runner RUNNER, opts ...Option) *Commander {
	cmd := &cobra.Command{
		Use:   name,
		Short: name,
		Run:   runner,
	}

	for _, opt := range opts {
		opt(cmd)
	}

	if parent == "" {
		c.root.AddCommand(cmd)
	} else {
		if p, ok := c.cmds[parent]; ok {
			p.AddCommand(cmd)
			c.cmds[name] = cmd
		} else {
			log.Fatal("invalid parent")
		}
	}

	return c
}

// Run block function
func (c *Commander) Run() error {
	c.root.PreRun = func(cmd *cobra.Command, args []string) {
		for _, server := range c.servers {
			go server.Start()
		}
	}
	c.root.Run = func(cmd *cobra.Command, args []string) {
		c.catchSignal()
	}

	c.initIPCServer()

	err := c.root.Execute()
	return err
}

// Register ...
func (c *Commander) Register(servers ...IServer) {
	c.servers = append(c.servers, servers...)
}

// Dump ...
func (c *Commander) Dump() []IServer {
	return c.servers
}

func (c *Commander) catchSignal() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case sig := <-sigs:
			for _, server := range c.servers {
				server.Stop()
			}
			_ = sig
			return
		}
	}
}

func (c *Commander) initIPCServer() {
	c.cmdServ.AddRoute("/stop", func(w http.ResponseWriter, r *http.Request) {
		for _, server := range c.servers {
			server.Stop()
		}

		fmt.Fprintf(w, "server exit")
	})
	c.servers = append(c.servers, c.cmdServ)
}
