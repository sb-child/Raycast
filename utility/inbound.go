package utility

type HttpProxyInbound struct {
	Listen string
	Users  []string
}

type SocksProxyInbound struct {
	Listen string
	Users  []string
	Udp    bool
}

type ApiInbound struct {
}
