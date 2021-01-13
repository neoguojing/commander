package commander

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

var (
	cmdServ *IPCServer
)

const (
	domain = "http://localhost:28888/stop"
)

var (
	rootCmd = &cobra.Command{
		Use:   "",
		Short: "start",
		Run: func(cmd *cobra.Command, args []string) {
			catchSignal()
		},
		PreRun: func(cmd *cobra.Command, args []string) {
			for _, server := range Servers.Dump() {
				go server.Start()
			}
		},
	}

	stopCmd = &cobra.Command{
		Use:   "stop",
		Short: "stop",
		Run: func(cmd *cobra.Command, args []string) {
			sendStopCmd()
		},
		PostRun: func(cmd *cobra.Command, args []string) {
		},
	}
)

func init() {
	initCmdServer()
	rootCmd.Use = os.Args[0]
	rootCmd.AddCommand(stopCmd)
}

func initCmdServer() {
	cmdServ = NewIPCServer()
	cmdServ.AddRoute("/stop", func(w http.ResponseWriter, r *http.Request) {
		for _, server := range Servers.Dump() {
			server.Stop()
		}

		fmt.Fprintf(w, "server exit")
	})
	Servers.Register(cmdServ)
}

func sendStopCmd() {
	client, err := GetIPCClient(cmdServ.UnixSockPath)
	if err == nil {
		resp, err := client.Get(domain)
		if err != nil {
			log.Fatal(err)
		} else {
			content, err := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(string(content))
		}
	} else {
		log.Fatal(err)
	}
}

func catchSignal() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case sig := <-sigs:
			for _, server := range Servers.Dump() {
				server.Stop()
			}
			_ = sig
			return
		}
	}
}
