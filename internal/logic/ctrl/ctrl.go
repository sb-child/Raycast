package ctrl

import (
	"context"
	"raycast/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

type sCtrl struct {
}

func init() {
	service.RegisterCtrl(New())
}

func New() *sCtrl {
	return &sCtrl{}
}

func (x *sCtrl) Start(ctx context.Context) {
	g.Log().Warning(ctx, "[service] Starting Ctrl...")
}

func (x *sCtrl) Stop(ctx context.Context) {
	g.Log().Warning(ctx, "[service] Stopping Ctrl...")
}
