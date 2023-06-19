package xrayexec

import (
	"context"
	"raycast/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

type sXrayExec struct {
}

func init() {
	service.RegisterXrayExec(New())
}

func New() *sXrayExec {
	return &sXrayExec{}
}

func (x *sXrayExec) Start(ctx context.Context) {
	g.Log().Warning(ctx, "[service] Starting XrayExec...")
}

func (x *sXrayExec) Stop(ctx context.Context) {
	g.Log().Warning(ctx, "[service] Stopping XrayExec...")
}
