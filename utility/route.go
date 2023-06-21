package utility

import (
	"strings"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/grand"
)

type Route struct {
	Inbound    []string
	Outbound   []string
	Network    string
	Protocol   string
	Domain     []string
	SourceIp   []string
	TargetIp   []string
	User       []string
	HttpAttrs  g.Map
	SourcePort string
	TargetPort string
}

func (x *Route) Json() (routeCfg *gjson.Json, balancerCfg *gjson.Json, balancerTag string) {
	balancerTag = "bal-" + grand.S(16)
	bals := []string{}
	for _, v := range x.Outbound {
		if strings.HasPrefix(v, "$") {
			bals = append(bals, strings.TrimPrefix(v, "$"))
		}
	}
	balancerCfg = gjson.New(g.Map{
		"tag": balancerTag,
		"selector": bals,
	})
	routeCfg = gjson.New(g.Map{
		"domainMatcher": "hybrid",
		"type":          "field",
		"domain":        x.Domain,
		"ip":            x.TargetIp,
		"port":          x.TargetPort,
		"sourcePort":    x.SourcePort,
		"network":       x.Network,
		"source":        x.SourceIp,
		"user":          x.User,
		"inboundTag":    x.Inbound,
		"protocol":      SplitItems(x.Protocol),
		"attrs":         x.HttpAttrs,
	})
	if len(x.Network) == 0 {
		routeCfg.Remove("network")
	}
	if len(x.Protocol) == 0 {
		routeCfg.Remove("protocol")
	}
	if len(x.SourcePort) == 0 {
		routeCfg.Remove("sourcePort")
	}
	if len(x.TargetPort) == 0 {
		routeCfg.Remove("port")
	}
	if len(bals) == 0 {
		balancerTag = ""
		balancerCfg = nil
		if len(x.Outbound) == 0 {
			routeCfg.Set("outboundTag", "")
		} else {
			routeCfg.Set("outboundTag", x.Outbound[0])
		}
	} else {
		routeCfg.Set("balancerTag", balancerTag)
	}
	return
}
