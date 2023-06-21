package xraycfg

import (
	"context"
	"fmt"
	"os"
	"raycast/internal/service"
	"raycast/utility"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
)

type sXrayCfg struct {
	xrayConfigFile string
	xrayApiAddr    string
	inboundsCfg    *gjson.Json
	inboundGroup   []*gjson.Json
	outboundsCfg   *gjson.Json
	outboundGroup  []*gjson.Json
	routeCfg       []*gjson.Json
	balancerCfg    []*gjson.Json
}

func init() {
	service.RegisterXrayCfg(New())
}

func New() *sXrayCfg {
	return &sXrayCfg{}
}

func firstKey(m *gvar.Var) string {
	for k := range m.MapStrAny() {
		return k
	}
	return "" // return empty string if map is empty
}

func (x *sXrayCfg) parseInbound(ctx context.Context) {
	cfgLen := x.inboundsCfg.Len(".")
	x.inboundGroup = make([]*gjson.Json, 0, cfgLen)
	for i := 0; i < cfgLen; i++ {
		t := x.inboundsCfg.GetJson(fmt.Sprintf("%d", i))
		tag := fmt.Sprintf("in-user-%d", i)
		k := firstKey(t.Var())
		switch k {
		case "http":
			n := utility.HttpProxyInbound{}
			x.inboundGroup = append(x.inboundGroup,
				n.FromCfg(t.GetJson(k), tag).Json())
		case "socks":
			n := utility.SocksProxyInbound{}
			x.inboundGroup = append(x.inboundGroup,
				n.FromCfg(t.GetJson(k), tag).Json())
		case "vmess":
			n := utility.VmessInbound{}
			x.inboundGroup = append(x.inboundGroup,
				n.FromCfg(t.GetJson(k), tag).Json())
		case "trojan":
			n := utility.TrojanInbound{}
			x.inboundGroup = append(x.inboundGroup,
				n.FromCfg(t.GetJson(k), tag).Json())
		}
	}
}

func (x *sXrayCfg) parseOutbound(ctx context.Context, outTag string, clear bool) {
	cfgLen := x.outboundsCfg.Len(".")
	if clear {
		x.outboundGroup = make([]*gjson.Json, 0, cfgLen)
	}
	for i := 0; i < cfgLen; i++ {
		t := x.outboundsCfg.GetJson(fmt.Sprintf("%d", i))
		tag := fmt.Sprintf("out-%s-%d", outTag, i)
		k := firstKey(t.Var())
		switch k {
		case "vmess":
			n := utility.VmessOutbound{}
			x.outboundGroup = append(x.outboundGroup,
				n.FromCfg(t.GetJson(k), tag).Json())
		case "trojan":
			n := utility.TrojanOutbound{}
			x.outboundGroup = append(x.outboundGroup,
				n.FromCfg(t.GetJson(k), tag).Json())
		case "direct":
			n := utility.DirectOutbound{}
			x.outboundGroup = append(x.outboundGroup,
				n.FromCfg(t.GetJson(k), tag).Json())
		case "block":
			n := utility.BlockOutbound{}
			x.outboundGroup = append(x.outboundGroup,
				n.FromCfg(t.GetJson(k), tag).Json())
		}
	}
}

func (x *sXrayCfg) parseRoutes(ctx context.Context) {
	inboundLen := x.inboundsCfg.Len(".")
	outboundLen := x.outboundsCfg.Len(".")
	inboundList := make([]string, inboundLen)
	systemOutboundList := make([]string, outboundLen)
	systemInboundList := make([]string, outboundLen)
	for i := 0; i < inboundLen; i++ {
		inboundList[i] = fmt.Sprintf("in-user-%d", i)
	}
	for i := 0; i < outboundLen; i++ {
		systemOutboundList[i] = fmt.Sprintf("out-system-%d", i)
	}
	for i := 0; i < outboundLen; i++ {
		systemInboundList[i] = fmt.Sprintf("in-system-%d", i)
	}
	// direct private addresses
	rt := utility.Route{
		Network:  "tcp,udp",
		TargetIp: []string{"geoip:private"},
		Outbound: []string{"direct"},
	}
	rtCfg, _, _ := rt.Json()
	x.routeCfg = append(x.routeCfg, rtCfg)
	// in-user-* > out-user-*
	rt = utility.Route{
		Network: "tcp,udp",
		Inbound: inboundList,
		Outbound: []string{
			"$out-user-",
		},
	}
	rtCfg, balCfg, _ := rt.Json()
	x.routeCfg = append(x.routeCfg, rtCfg)
	if balCfg != nil {
		x.balancerCfg = append(x.balancerCfg, balCfg)
	}
	// in-system-{0} > out-system-{0}
	for k, v := range systemInboundList {
		rt := utility.Route{
			Network:  "tcp,udp",
			Inbound:  []string{v},
			Outbound: []string{systemOutboundList[k]},
		}
		rtCfg, _, _ := rt.Json()
		x.routeCfg = append(x.routeCfg, rtCfg)
	}
	// block all
	rt = utility.Route{
		Network:  "tcp,udp",
		Outbound: []string{"block"},
	}
	rtCfg, _, _ = rt.Json()
	x.routeCfg = append(x.routeCfg, rtCfg)
	// cfgLen := x.outboundsCfg.Len(".")
	// x.outboundGroup = make([]*gjson.Json, 0, cfgLen)
	// for i := 0; i < cfgLen; i++ {
	// 	t := x.outboundsCfg.GetJson(fmt.Sprintf("%d", i))
	// 	k := firstKey(t.Var())
	// 	switch k {
	// 	case "vmess":
	// 		n := utility.VmessOutbound{}
	// 		x.outboundGroup = append(x.outboundGroup,
	// 			n.FromCfg(t.GetJson(k), fmt.Sprintf("out-user-%d", i)).Json())
	// 	case "direct":
	// 		n := utility.DirectOutbound{}
	// 		x.outboundGroup = append(x.outboundGroup,
	// 			n.FromCfg(t.GetJson(k), fmt.Sprintf("out-user-%d", i)).Json())
	// 	case "block":
	// 		n := utility.BlockOutbound{}
	// 		x.outboundGroup = append(x.outboundGroup,
	// 			n.FromCfg(t.GetJson(k), fmt.Sprintf("out-user-%d", i)).Json())
	// 	}
	// }
}

func (x *sXrayCfg) Generate(ctx context.Context) {
	g.Log().Infof(ctx, "[XrayCfg] Generating config file to %s", x.xrayConfigFile)
	s := utility.CfgFramework{}
	s.Init()
	s.Api(true)
	n := utility.ApiInbound{}
	s.Inbounds(n.FromCfg(x.xrayApiAddr).Json())
	s.Inbounds(x.inboundGroup...)
	s.Outbounds(x.outboundGroup...)
	direct := utility.DirectOutbound{}
	s.Outbounds(direct.FromCfg(nil, "direct").Json())
	block := utility.BlockOutbound{}
	s.Outbounds(block.FromCfg(nil, "block").Json())
	s.Routes(x.routeCfg...)
	s.Balancers(x.balancerCfg...)
	f, err := gfile.OpenWithFlagPerm(x.xrayConfigFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0640)
	if err != nil {
		g.Log().Errorf(ctx, "[XrayCfg] Failed to write config file: %s", err.Error())
		return
	}
	f.WriteString(s.Json().String())
	f.Close()
}

func (x *sXrayCfg) Start(ctx context.Context) {
	g.Log().Warning(ctx, "[service] Starting XrayCfg...")
	x.xrayConfigFile = g.Config().MustGet(ctx, "raycast.xrayConfig", "").String()
	x.xrayApiAddr = g.Config().MustGet(ctx, "raycast.xrayApiAddr", "").String()
	inbounds := gjson.New(g.Config().MustGet(ctx, "inbound", ""))
	outbounds := gjson.New(g.Config().MustGet(ctx, "outbound", ""))
	x.inboundsCfg = inbounds
	x.outboundsCfg = outbounds
	x.parseInbound(ctx)
	x.parseOutbound(ctx, "user", true)
	x.parseOutbound(ctx, "system", false)
	x.parseRoutes(ctx)
	x.Generate(ctx)
	// g.Log().Warningf(ctx, "%+v", firstKey(inbounds.Get("0")))
}

func (x *sXrayCfg) Stop(ctx context.Context) {
	g.Log().Warning(ctx, "[service] Stopping XrayCfg...")
}
