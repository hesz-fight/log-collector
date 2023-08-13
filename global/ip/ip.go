package ip

import (
	"errors"
	"net"
	"strings"
)

var localIP string

func GetOutBoundIP() (string, error) {
	if localIP != "" {
		return localIP, nil
	}
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	addr := conn.LocalAddr().(*net.UDPAddr)
	if len(strings.Split(addr.IP.String(), ":")) == 0 {
		return "", errors.New("get local ip fail.")
	}

	localIP = strings.Split(addr.IP.String(), ":")[0]

	return localIP, nil
}
