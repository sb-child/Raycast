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
