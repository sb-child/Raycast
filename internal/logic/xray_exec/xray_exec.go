package xrayexec

import (
	"context"
	"raycast/internal/consts"
	"raycast/internal/service"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gproc"
)

type sXrayExec struct {
	xrayProgram    string
	xrayConfigFile string
	ray            *gproc.Process
	rayLock        sync.Mutex
	rayDone        chan error
}

func init() {
	service.RegisterXrayExec(New())
}

func New() *sXrayExec {
	return &sXrayExec{}
}

func (x *sXrayExec) xrayWait(ctx context.Context, p *gproc.Process) {
	g.Log().Warning(ctx, "[XrayExec] Xray is running")
	err := p.Wait()
	if err == nil {
		g.Log().Warning(ctx, "[XrayExec] Xray is exited")
	} else {
		g.Log().Warning(ctx, "[XrayExec] Xray is crashed:", err.Error())
	}
	x.rayDone <- err
}

func (x *sXrayExec) xray(ctx context.Context, args []string) *gproc.Process {
	return gproc.NewProcess(x.xrayProgram, args)
}

func (x *sXrayExec) startXray(ctx context.Context) error {
	x.rayLock.Lock()
	defer x.rayLock.Unlock()
	if !(x.ray == nil || (x.ray.Pid() == 0) || (x.ray.Pid() == -1)) {
		x.ray.Kill()
	}
	x.ray = x.xray(ctx, []string{"run", "-c", x.xrayConfigFile})
	for len(x.rayDone) > 0 {
		<-x.rayDone
	}
	x.rayDone = make(chan error, 1)
	pid, err := x.ray.Start(ctx)
	go x.xrayWait(ctx, x.ray)
	if err != nil {
		g.Log().Warning(ctx, "Failed to start Xray:", err.Error())
		return err
	}
	g.Log().Warningf(ctx, "Xray is started, PID is [%d]", pid)
	return nil
}

func (x *sXrayExec) stopXray(ctx context.Context) error {
	x.rayLock.Lock()
	defer x.rayLock.Unlock()
	if x.ray == nil {
		return nil
	}
	err := x.ray.Kill()
	select {
	case <-x.rayDone:
		break
	case <-time.After(time.Second):
		g.Log().Warningf(ctx, "Xray process timeout while stopping")
	}
	x.ray = nil
	return err
}

func (x *sXrayExec) xrayStatus(ctx context.Context) consts.ProcessState {
	x.rayLock.Lock()
	defer x.rayLock.Unlock()
	if x.ray == nil || x.ray.Pid() == 0 {
		return consts.PROC_EXITED
	}
	if x.ray.ProcessState != nil && x.ray.ProcessState.Exited() {
		if x.ray.ProcessState.ExitCode() == 0 {
			return consts.PROC_EXITED
		}
		return consts.PROC_FAILED
	}
	return consts.PROC_RUNNING
}

// func (x *sXrayExec) testXrayWithCfg(ctx context.Context) {
// 	x.ray = x.xray(ctx, []string{"version"})
// 	pid, err := x.ray.Start(ctx)
// 	if err != nil {
// 		g.Log().Warning(ctx, "Failed to start Xray:", err.Error())
// 		return
// 	}
// 	g.Log().Warningf(ctx, "Xray is started, PID is [%d]", pid)
// }

func (x *sXrayExec) runXrayWithCfg(ctx context.Context) {

}

func (x *sXrayExec) configure(ctx context.Context) {
	x.xrayProgram = g.Config().MustGet(ctx, "raycast.xray", "").String()
	x.xrayConfigFile = g.Config().MustGet(ctx, "raycast.xrayConfig", "").String()
}

func (x *sXrayExec) Start(ctx context.Context) {
	g.Log().Warning(ctx, "[service] Starting XrayExec...")
	x.rayDone = make(chan error, 1)
	x.configure(ctx)
	// x.testXrayWithCfg(ctx)
	// x.runXrayWithCfg(ctx)
	g.Log().Infof(ctx, "status %d", x.xrayStatus(ctx))
	x.startXray(ctx)
	g.Log().Infof(ctx, "status %d", x.xrayStatus(ctx))
	time.Sleep(time.Second)
	g.Log().Infof(ctx, "status %d", x.xrayStatus(ctx))
	// x.stopXray(ctx)
	// g.Log().Infof(ctx, "status %d", x.xrayStatus(ctx))
	// x.startXray(ctx)
	// g.Log().Infof(ctx, "status %d", x.xrayStatus(ctx))
	// x.stopXray(ctx)
	// g.Log().Infof(ctx, "status %d", x.xrayStatus(ctx))
}

func (x *sXrayExec) Stop(ctx context.Context) {
	g.Log().Warning(ctx, "[service] Stopping XrayExec...")
	x.stopXray(ctx)
}
