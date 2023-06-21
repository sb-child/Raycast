package utility

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
)

type VmessOutbound struct {
	Through  string
	Server   string
	User     string
	Tag      string
	Security *gjson.Json
}

func (x *VmessOutbound) FromCfg(c *gjson.Json, tag string) *VmessOutbound {
	x.Through = c.Get("through", "").String()
	x.Server = c.Get("server", "").String()
	x.User = c.Get("user", "").String()
	x.Tag = tag
	x.Security = AutoSecurityJson(c.GetJson("security"), false)
	return x
}

func (x *VmessOutbound) Json() *gjson.Json {
	host, port := Addr(x.Server)
	r := gjson.New(g.Map{
		"sendThrough": x.Through,
		"protocol":    "vmess",
		"settings": g.Map{
			"vnext": g.List{g.Map{
				"address": host,
				"port":    port,
				"users": g.List{g.Map{
					"id":       x.User,
					"alterId":  0,
					"security": "auto",
					"level":    0,
				}},
			}},
		},
		"streamSettings": x.Security,
		"tag":            x.Tag,
	})
	return r
}

type TrojanOutbound struct {
	Through  string
	Server   string
	User     string
	Tag      string
	Security *gjson.Json
}

func (x *TrojanOutbound) FromCfg(c *gjson.Json, tag string) *TrojanOutbound {
	x.Through = c.Get("through", "").String()
	x.Server = c.Get("server", "").String()
	x.User = c.Get("user", "").String()
	x.Tag = tag
	x.Security = AutoSecurityJson(c.GetJson("security"), false)
	return x
}

func (x *TrojanOutbound) Json() *gjson.Json {
	host, port := Addr(x.Server)
	r := gjson.New(g.Map{
		"sendThrough": x.Through,
		"protocol":    "trojan",
		"settings": g.Map{
			"servers": g.List{g.Map{
				"address":  host,
				"port":     port,
				"password": x.User,
			}},
		},
		"streamSettings": x.Security,
		"tag":            x.Tag,
	})
	return r
}

type DirectOutbound struct {
	Through  string
	Resolver string
	Tag      string
}

func (x *DirectOutbound) FromCfg(c *gjson.Json, tag string) *DirectOutbound {
	x.Through = c.Get("through", "").String()
	x.Resolver = c.Get("resolver", "").String()
	x.Tag = tag
	return x
}

func (x *DirectOutbound) Json() *gjson.Json {
	r := gjson.New(g.Map{
		"protocol": "freedom",
		"settings": g.Map{
			"domainStrategy": ResolverSelect(x.Resolver),
			"redirect":       x.Through,
			"userLevel":      0,
		},
		"tag": x.Tag,
	})
	if len(x.Through) == 0 {
		r.Remove("protocol.redirect")
	}
	return r
}

type BlockOutbound struct {
	Tag string
}

func (x *BlockOutbound) FromCfg(c *gjson.Json, tag string) *BlockOutbound {
	x.Tag = tag
	return x
}

func (x *BlockOutbound) Json() *gjson.Json {
	r := gjson.New(g.Map{
		"protocol": "blackhole",
		"settings": g.Map{
			"response": g.Map{
				"type": "none",
			},
		},
		"tag": x.Tag,
	})
	return r
}
