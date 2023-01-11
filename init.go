package commander

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	unixSockPath = fmt.Sprintf("/tmp/%s.sock", filepath.Base(os.Args[0]))
)

const (
	domain = "http://localhost:28888/stop"
)

var (
	rootCmd = &cobra.Command{
		Use:   "",
		Short: "start",
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
	rootCmd.Use = filepath.Base(os.Args[0])
	rootCmd.AddCommand(stopCmd)
}

func sendStopCmd() {
	client, err := GetIPCClient(unixSockPath)
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
