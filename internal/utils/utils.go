package utils

import (
	"fmt"
	"net"
)

func GetIP(remoteAddr string) (net.IP, error) {
	ip, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return nil, fmt.Errorf("utils: %q is not IP:port", remoteAddr)
	}

	userIP := net.ParseIP(ip)
	if userIP == nil {
		return nil, fmt.Errorf("utils: %q is not IP:port", remoteAddr)
	}

	return userIP, nil
}
