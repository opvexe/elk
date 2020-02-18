package etcd

import (
	"fmt"
	"net"
)

var localIPSouce []string

func init() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(fmt.Sprintf("local ip failed:", err))
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				localIPSouce = append(localIPSouce, ipnet.IP.String())
			}
		}
	}
}
