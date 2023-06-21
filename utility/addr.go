package utility

import (
	"net/netip"
	"strings"

	"github.com/gogf/gf/v2/util/gconv"
)

func Addr(x string) (host string, port uint16) {
	a, err := netip.ParseAddrPort(x)
	if err != nil {
		host, port = DomainAddr(x)
		return
	}
	host = a.Addr().String()
	port = a.Port()
	return
}

func DomainAddr(x string) (host string, port uint16) {
	a := strings.SplitN(x, ":", 2)
	if len(a) != 2 {
		return
	}
	host = a[0]
	port = gconv.Uint16(a[1])
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
