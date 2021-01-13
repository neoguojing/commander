// +build !linux

package commander

import (
	"errors"
	"net/http"
)

//GetIPCClient ...
func GetIPCClient(addr string) (*http.Client, error) {
	return nil, errors.New("not avaliable")
}
