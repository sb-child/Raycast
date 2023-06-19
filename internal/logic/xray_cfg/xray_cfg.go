package xraycfg

import (
	"context"
	"raycast/internal/service"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
)

type sXrayCfg struct {
	inboundsCfg     *gjson.Json
	outboundsCfg    *gjson.Json
	certificatesCfg *gjson.Json
}

func init() {
	service.RegisterXrayCfg(New())
}

func New() *sXrayCfg {
	return &sXrayCfg{}
}

func (x *sXrayCfg) Start(ctx context.Context) {
	g.Log().Warning(ctx, "[service] Starting XrayCfg...")
	inbounds := gjson.New(g.Config().MustGet(ctx, "inbound", ""))
	outbounds := gjson.New(g.Config().MustGet(ctx, "outbound", ""))
	certificates := gjson.New(g.Config().MustGet(ctx, "certificates", ""))
	x.inboundsCfg = inbounds
	x.outboundsCfg = outbounds
	x.certificatesCfg = certificates
	// g.Log().Warningf(ctx, "%+v", inbounds)
}

func (x *sXrayCfg) Stop(ctx context.Context) {
	g.Log().Warning(ctx, "[service] Stopping XrayCfg...")
}
