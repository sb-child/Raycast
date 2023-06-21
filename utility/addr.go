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

func ResolverSelect(x string) string {
	switch x {
	case "system":
		return "AsIs"
	case "xray":
		return "UseIP"
	case "xray4":
		return "UseIPv4"
	case "xray6":
		return "UseIPv6"
	}
	return "AsIs"
}
