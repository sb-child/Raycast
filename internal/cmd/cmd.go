package cmd

import (
	"context"
	"os"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gproc"
)

func MainFunc(ctx context.Context, parser *gcmd.Parser) (err error) {
	StartAllServices(ctx)
	StartPanelServer(ctx)
	gproc.AddSigHandlerShutdown(func(sig os.Signal) {
		g.Log().Infof(ctx, "%s Signal received, stopping service...", sig.String())
		StopAllServices(ctx)
		StopPanelServer(ctx)
	})

	// go func() {
	// 	err := MainCmd(ctx, parser)
	// 	if err != nil {
	// 		g.Log().Warning(ctx, "main process exited with error:", err)
	// 		return
	// 	}
	// 	g.Log().Warning(ctx, "main process exited")
	// }()
	time.Sleep(1 * time.Second)
	// e := service.XrayApi().AddSystemInbound(ctx, "127.0.0.1:2018", "in-system-3")
	// g.Log().Info(ctx, e)
	// b := utility.BlockOutbound{}
	// e = service.XrayApi().AddOutbound(ctx, b.FromCfg(nil, "out-user-1").Json())
	// g.Log().Info(ctx, e)
	gproc.Listen()
	return nil
}

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start Raycast daemon",
		Func:  MainFunc,
	}
)
