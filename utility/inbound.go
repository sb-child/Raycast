package utility

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
)

type HttpProxyInbound struct {
	Listen string
	Users  []string
	Tag    string
}

func (x *HttpProxyInbound) FromCfg(c *gjson.Json, tag string) *HttpProxyInbound {
	x.Listen = c.Get("listen", "").String()
	x.Users = c.Get("users", g.ArrayStr{}).Strings()
	x.Tag = tag
	return x
}

func (x *HttpProxyInbound) Json() *gjson.Json {
	host, port := Addr(x.Listen)
	r := gjson.New(g.Map{
		"protocol": "http",
		"listen":   host,
		"port":     port,
		"settings": g.Map{},
		"tag":      x.Tag,
	})
	users := g.List{}
	for _, v := range x.Users {
		u, p := SplitUserPwd(v)
		users = append(users, g.Map{
			"user": u,
			"pass": p,
		})
	}
	if len(users) != 0 {
		r.Set("settings.accounts", users)
	}
	return r
}

type SocksProxyInbound struct {
	Listen string
	Users  []string
	Udp    bool
	Tag    string
}

func (x *SocksProxyInbound) FromCfg(c *gjson.Json, tag string) *SocksProxyInbound {
	x.Listen = c.Get("listen", "").String()
	x.Udp = c.Get("udp", "").Bool()
	x.Users = c.Get("users", g.ArrayStr{}).Strings()
	x.Tag = tag
	return x
}

func (x *SocksProxyInbound) Json() *gjson.Json {
	host, port := Addr(x.Listen)
	r := gjson.New(g.Map{
		"protocol": "socks",
		"listen":   host,
		"port":     port,
		"settings": g.Map{
			"auth": "noauth",
			"udp":  x.Udp,
			"ip":   host,
		},
		"tag": x.Tag,
	})
	users := g.List{}
	for _, v := range x.Users {
		u, p := SplitUserPwd(v)
		users = append(users, g.Map{
			"user": u,
			"pass": p,
		})
	}
	if len(users) != 0 {
		r.Set("settings.auth", "password")
		r.Set("settings.accounts", users)
	}
	return r
}

type ApiInbound struct {
}
