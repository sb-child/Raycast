package xrayexec

import (
	"context"
	"raycast/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gproc"
)

type sXrayExec struct {
	xrayProgram    string
	xrayConfigFile string
	ray            *gproc.Process
}

func init() {
	service.RegisterXrayExec(New())
}

func New() *sXrayExec {
	return &sXrayExec{}
}

func (x *sXrayExec) xray(ctx context.Context, args []string) *gproc.Process {
	return gproc.NewProcess(x.xrayProgram, args)
}

func (x *sXrayExec) testXray(ctx context.Context) {
	x.ray = x.xray(ctx, []string{"help"})
	pid, err := x.ray.Start(ctx)
	if err != nil {
		g.Log().Warning(ctx, "Failed to start Xray:", err.Error())
		return
	}
	g.Log().Warningf(ctx, "Xray is started, PID is [%d]", pid)
}

func (x *sXrayExec) runXray(ctx context.Context) {

}

func (x *sXrayExec) configure(ctx context.Context) {
	x.xrayProgram = g.Config().MustGet(ctx, "raycast.xray", "").String()
	x.xrayConfigFile = g.Config().MustGet(ctx, "raycast.xrayConfig", "").String()
}

func (x *sXrayExec) Start(ctx context.Context) {
	g.Log().Warning(ctx, "[service] Starting XrayExec...")
	x.configure(ctx)
	x.testXray(ctx)
	x.runXray(ctx)
}

func (x *sXrayExec) Stop(ctx context.Context) {
	g.Log().Warning(ctx, "[service] Stopping XrayExec...")
	x.ray.Kill()
}
