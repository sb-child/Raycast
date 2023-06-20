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
		switch firstKey(t.Var()) {
		case "http":
			n := utility.HttpProxyInbound{}
			x.inboundGroup = append(x.inboundGroup,
				n.FromCfg(t.GetJson("http"), fmt.Sprintf("in-user-%d", i)).Json())
		case "socks":
			n := utility.SocksProxyInbound{}
			x.inboundGroup = append(x.inboundGroup,
				n.FromCfg(t.GetJson("socks"), fmt.Sprintf("in-user-%d", i)).Json())
		case "vmess":
			n := utility.VmessInbound{}
			x.inboundGroup = append(x.inboundGroup,
				n.FromCfg(t.GetJson("vmess"), fmt.Sprintf("in-user-%d", i)).Json())
		}
	}
}

func (x *sXrayCfg) Generate(ctx context.Context) {
	g.Log().Infof(ctx, "[XrayCfg] Generating config file to %s", x.xrayConfigFile)
	s := utility.CfgFramework{}
	s.Init()
	s.Api(true)
	n := utility.ApiInbound{}
	n.FromCfg(x.xrayApiAddr)
	s.Inbounds(n.Json())
	s.Inbounds(x.inboundGroup...)
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
	x.Generate(ctx)
	// g.Log().Warningf(ctx, "%+v", firstKey(inbounds.Get("0")))
}

func (x *sXrayCfg) Stop(ctx context.Context) {
	g.Log().Warning(ctx, "[service] Stopping XrayCfg...")
}
