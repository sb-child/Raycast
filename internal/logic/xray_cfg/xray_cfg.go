package xraycfg

import (
	"context"
	"raycast/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

type sXrayCfg struct {
}

func init() {
	service.RegisterXrayCfg(New())
}

func New() *sXrayCfg {
	return &sXrayCfg{}
}

func (x *sXrayCfg) Start(ctx context.Context) {
	g.Log().Warning(ctx, "[service] Starting XrayCfg...")
}

func (x *sXrayCfg) Stop(ctx context.Context) {
	g.Log().Warning(ctx, "[service] Stopping XrayCfg...")
}
