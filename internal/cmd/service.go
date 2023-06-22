package cmd

import (
	"context"
	"raycast/internal/consts"
	"raycast/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func StartPanelServer(ctx context.Context) {
	if !g.Config().MustGet(ctx, "panel.enabled", false).Bool() {
		return
	}
	s := g.Server(consts.RAYCAST_PANEL_NAME)
	// s.EnablePProf()
	s.Group(g.Config().MustGet(ctx, "panel.rootDir", "/").String(), func(group *ghttp.RouterGroup) {
		// group.GET("/_", duo.ParamsHandler)
		// group.GET("/_/:prompt", duo.PromptHandler)
		// group.GET("/b/:b64", duo.Base64Handler)
	})
	s.Start()
}

func StopPanelServer(ctx context.Context) {
	if !g.Config().MustGet(ctx, "panel.enabled", false).Bool() {
		return
	}
	s := g.Server(consts.RAYCAST_PANEL_NAME)
	s.Shutdown()
}

func StartAllServices(ctx context.Context) {
	service.XrayCfg().Start(ctx)
	service.XrayExec().Start(ctx)
	service.XrayApi().Start(ctx)
	service.Ctrl().Start(ctx)
	g.Log().Info(ctx, "All services started")
}

func StopAllServices(ctx context.Context) {
	service.Ctrl().Stop(ctx)
	service.XrayApi().Stop(ctx)
	service.XrayExec().Stop(ctx)
	service.XrayCfg().Stop(ctx)
	g.Log().Info(ctx, "All services stopped")
}
