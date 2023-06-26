package xrayapi

import (
	"bytes"
	"context"
	"raycast/internal/service"
	"raycast/utility"
	"sync"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/xtls/xray-core/app/proxyman"
	"github.com/xtls/xray-core/app/proxyman/command"
	statsService "github.com/xtls/xray-core/app/stats/command"
	"github.com/xtls/xray-core/common/net"
	"github.com/xtls/xray-core/common/serial"
	"github.com/xtls/xray-core/core"
	confSerial "github.com/xtls/xray-core/infra/conf/serial"
	httpInbound "github.com/xtls/xray-core/proxy/http"
)

type sXrayApi struct {
	xrayApiAddr string
	xray        utility.XrayController
	lock        sync.Mutex
}

func init() {
	service.RegisterXrayApi(New())
}

func New() *sXrayApi {
	return &sXrayApi{}
}

func (x *sXrayApi) AddOutbound(ctx context.Context, json *gjson.Json) error {
	x.lock.Lock()
	defer x.lock.Unlock()
	framework := gjson.New(g.Map{
		"outbounds": g.Slice{json},
	})
	cfg, err := confSerial.DecodeJSONConfig(bytes.NewBufferString(framework.String()))
	if err != nil {
		return err
	}
	if len(cfg.OutboundConfigs) == 0 {
		return gerror.New("cannot find outbound config")
	}
	outbound := cfg.OutboundConfigs[0]
	o, err := outbound.Build()
	if err != nil {
		return err
	}
	_, err = x.xray.HsClient.AddOutbound(ctx, &command.AddOutboundRequest{
		Outbound: o,
	})
	return err
}

func (x *sXrayApi) DelOutbound(ctx context.Context, tag string) (err error) {
	x.lock.Lock()
	defer x.lock.Unlock()
	_, err = x.xray.HsClient.RemoveOutbound(ctx, &command.RemoveOutboundRequest{
		Tag: tag,
	})
	return
}

func (x *sXrayApi) AddSystemInbound(ctx context.Context, addr string, tag string) (err error) {
	x.lock.Lock()
	defer x.lock.Unlock()
	host, port := utility.Addr(addr)
	_, err = x.xray.HsClient.AddInbound(ctx, &command.AddInboundRequest{
		Inbound: &core.InboundHandlerConfig{
			Tag: tag,
			ReceiverSettings: serial.ToTypedMessage(
				&proxyman.ReceiverConfig{
					Listen:   net.NewIPOrDomain(net.ParseAddress(host)),
					PortList: &net.PortList{Range: []*net.PortRange{net.SinglePortRange(net.Port(port))}},
				},
			),
			ProxySettings: serial.ToTypedMessage(&httpInbound.ServerConfig{}),
		},
	})
	return
}

func (x *sXrayApi) DelInbound(ctx context.Context, tag string) (err error) {
	x.lock.Lock()
	defer x.lock.Unlock()
	_, err = x.xray.HsClient.RemoveInbound(ctx, &command.RemoveInboundRequest{
		Tag: tag,
	})
	return
}

func (x *sXrayApi) Stat(ctx context.Context, inbound bool, tag string, down bool) (val int64, err error) {
	x.lock.Lock()
	defer x.lock.Unlock()
	p := ""
	if inbound {
		p = "inbound>>>"
	} else {
		p = "outbound>>>"
	}
	p += tag + ">>>traffic>>>"
	if down {
		p += "downlink"
	} else {
		p += "uplink"
	}
	r, err := x.xray.SsClient.QueryStats(ctx, &statsService.QueryStatsRequest{
		Pattern: p,
		Reset_:  false,
	})
	if err != nil {
		return
	}
	s := r.GetStat()
	if len(s) == 0 {
		return 0, gerror.New("no stat found")
	}
	val = s[0].Value
	return
}

func (x *sXrayApi) Start(ctx context.Context) {
	g.Log().Warning(ctx, "[service] Starting XrayApi...")
	x.xrayApiAddr = g.Config().MustGet(ctx, "raycast.xrayApiAddr", "").String()
	err := x.xray.Init(x.xrayApiAddr)
	if err != nil {
		g.Log().Errorf(ctx, "[XrayApi] Api service is connecting to %s with err %s", x.xrayApiAddr, err.Error())
	} else {
		g.Log().Infof(ctx, "[XrayApi] Api service is connecting to %s", x.xrayApiAddr)
	}
}

func (x *sXrayApi) Stop(ctx context.Context) {
	g.Log().Warning(ctx, "[service] Stopping XrayApi...")
}
