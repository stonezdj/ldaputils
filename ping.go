package main

import (
	"net"
	"strings"
	"time"
)

// PingURL a host and port
func PingURL(url string) (bool, error) {
	splitURL := strings.Split(url, "://")
	var host, port, protocol string
	port = "389"
	if len(splitURL) == 2 {
		protocol = splitURL[0]
		parts := strings.Split(splitURL[1], ":")
		if len(parts) == 2 {
			host, port = parts[0], parts[1]
		} else {
			if strings.EqualFold("ldap", protocol) {
				port = "389"
			}
			if strings.EqualFold("ldaps", protocol) {
				port = "636"
			}
			host = parts[0]
		}

	} else {
		parts := strings.Split(url, ":")
		if len(parts) == 2 {
			host, port = parts[0], parts[1]
		} else {
			host = parts[0]
		}

	}

	seconds := 5
	timeOut := time.Duration(seconds) * time.Second
	_, err := net.DialTimeout("tcp", host+":"+port, timeOut)
	if err != nil {
		return false, err
	}
	return true, nil
}
