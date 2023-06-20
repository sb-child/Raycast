package utility

import (
	"net/netip"
)

func Addr(x string) (host string, port uint16) {
	a, err := netip.ParseAddrPort(x)
	if err != nil {
		host = ""
		port = 0
		return
	}
	host = a.Addr().String()
	port = a.Port()
	return
}
