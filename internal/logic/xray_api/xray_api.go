package xrayapi

import (
	"context"
	"raycast/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

type sXrayApi struct {
}

func init() {
	service.RegisterXrayApi(New())
}

func New() *sXrayApi {
	return &sXrayApi{}
}

func (x *sXrayApi) Start(ctx context.Context) {
	g.Log().Warning(ctx, "[service] Starting XrayApi...")
}

func (x *sXrayApi) Stop(ctx context.Context) {
	g.Log().Warning(ctx, "[service] Stopping XrayApi...")
}
