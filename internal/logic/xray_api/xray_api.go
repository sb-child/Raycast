package xrayapi

import (
	"context"
	"raycast/internal/service"
	"raycast/utility"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/xtls/xray-core/app/proxyman"
	"github.com/xtls/xray-core/app/proxyman/command"
	"github.com/xtls/xray-core/common/net"
	"github.com/xtls/xray-core/common/serial"
	"github.com/xtls/xray-core/core"
	httpInbound "github.com/xtls/xray-core/proxy/http"
)

type sXrayApi struct {
	xrayApiAddr string
	xray        utility.XrayController
}

func init() {
	service.RegisterXrayApi(New())
}

func New() *sXrayApi {
	return &sXrayApi{}
}

func (x *sXrayApi) AddSystemInbound(ctx context.Context, addr string, tag string) (err error) {
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

func (x *sXrayApi) DelSystemInbound(ctx context.Context, tag string) (err error) {
	_, err = x.xray.HsClient.RemoveInbound(context.Background(), &command.RemoveInboundRequest{
		Tag: tag,
	})
	return
}

func (x *sXrayApi) Start(ctx context.Context) {
	g.Log().Warning(ctx, "[service] Starting XrayApi...")
	x.xrayApiAddr = g.Config().MustGet(ctx, "raycast.xrayApiAddr", "").String()
	err := x.xray.Init(x.xrayApiAddr)
	g.Log().Info(ctx, err)
}

func (x *sXrayApi) Stop(ctx context.Context) {
	g.Log().Warning(ctx, "[service] Stopping XrayApi...")
}
