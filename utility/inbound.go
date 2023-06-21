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
	x.Udp = c.Get("udp", false).Bool()
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

type VmessInbound struct {
	Listen   string
	Users    []string
	Secure   bool
	Tag      string
	Security *gjson.Json
}

func (x *VmessInbound) FromCfg(c *gjson.Json, tag string) *VmessInbound {
	x.Listen = c.Get("listen", "").String()
	x.Users = c.Get("users", g.ArrayStr{}).Strings()
	x.Secure = c.Get("secure", false).Bool()
	x.Tag = tag
	x.Security = AutoSecurityJson(c.GetJson("security"), true)
	return x
}

func (x *VmessInbound) Json() *gjson.Json {
	host, port := Addr(x.Listen)
	r := gjson.New(g.Map{
		"protocol": "vmess",
		"listen":   host,
		"port":     port,
		"settings": g.Map{
			"disableInsecureEncryption": x.Secure,
		},
		"streamSettings": x.Security,
		"tag":            x.Tag,
	})
	users := g.List{}
	for _, v := range x.Users {
		u, p := SplitUserPwd(v)
		users = append(users, g.Map{
			"level":   0,
			"alterId": 0,
			"id":      u,
			"email":   p + "@cast.ray",
		})
	}
	if len(users) != 0 {
		r.Set("settings.clients", users)
	}
	return r
}

type TrojanInbound struct {
	Listen   string
	Fallback string
	Users    []string
	Tag      string
	Security *gjson.Json
}

func (x *TrojanInbound) FromCfg(c *gjson.Json, tag string) *TrojanInbound {
	x.Listen = c.Get("listen", "").String()
	x.Fallback = c.Get("fallback", "").String()
	x.Users = c.Get("users", g.ArrayStr{}).Strings()
	x.Tag = tag
	x.Security = AutoSecurityJson(c.GetJson("security"), true)
	return x
}

func (x *TrojanInbound) Json() *gjson.Json {
	host, port := Addr(x.Listen)
	fbDest, fbVer := SplitFallback(x.Fallback)
	r := gjson.New(g.Map{
		"protocol": "trojan",
		"listen":   host,
		"port":     port,
		"settings": g.Map{
			"clients": g.List{},
			"fallbacks": g.List{
				g.Map{
					"dest": fbDest,
					"xver": fbVer,
				},
			},
		},
		"streamSettings": x.Security,
		"tag":            x.Tag,
	})
	users := g.List{}
	for _, v := range x.Users {
		u, p := SplitUserPwd(v)
		users = append(users, g.Map{
			"level":    0,
			"password": u,
			"email":    p + "@cast.ray",
		})
	}
	if len(users) != 0 {
		r.Set("settings.clients", users)
	}
	return r
}

type ApiInbound struct {
	Listen string
}

func (x *ApiInbound) FromCfg(addr string) *ApiInbound {
	x.Listen = addr
	return x
}

func (x *ApiInbound) Json() *gjson.Json {
	host, port := Addr(x.Listen)
	r := gjson.New(g.Map{
		"protocol": "dokodemo-door",
		"listen":   host,
		"port":     port,
		"settings": g.Map{
			"address": host,
		},
		"tag": "api",
	})
	return r
}
