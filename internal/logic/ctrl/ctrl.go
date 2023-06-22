package ctrl

import (
	"context"
	"fmt"
	"raycast/internal/service"
	"raycast/utility"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type sCtrl struct {
	userOutbounds    []string
	enabledOutbounds []string
	outboundDelays   []float64
	delayTimeout     float64
}

func init() {
	service.RegisterCtrl(New())
}

func New() *sCtrl {
	return &sCtrl{}
}

func (x *sCtrl) speedtest(ctx context.Context, tag string, timeout float64) float64 {
	n, err := utility.ExtractNumber(tag)
	if err != nil {
		g.Log().Warningf(ctx, "[Ctrl] Failed to speedtest %s", tag)
		return timeout
	}
	port := n + 11000
	sysListen := fmt.Sprintf("127.0.0.1:%d", port)
	sysInbound := fmt.Sprintf("in-system-%d", n)
	err = service.XrayApi().AddSystemInbound(ctx, sysListen, sysInbound)
	if err != nil {
		g.Log().Warningf(ctx, "[Ctrl] Failed to speedtest %s: add inbound failed: %s", tag, err)
		return timeout
	}
	g.Log().Infof(ctx, "[Ctrl] Created inbound %s listen to %s", sysInbound, sysListen)
	c := g.Client().Timeout(time.Millisecond * time.Duration(timeout*1000)).
		Proxy("http://" + sysListen)
	startTime := gtime.Now()
	_, err = c.Get(ctx, "https://www.google.com/")
	if err != nil {
		g.Log().Infof(ctx, "[Ctrl] Failed to speedtest %s, deleting %s: %s", tag, sysInbound, err.Error())
		service.XrayApi().DelInbound(ctx, sysInbound)
		return timeout
	}
	stopTime := gtime.Now()
	delay := stopTime.Sub(startTime)
	g.Log().Infof(ctx, "[Ctrl] Speedtest %s done: %s, deleting %s", tag, delay.String(), sysInbound)
	service.XrayApi().DelInbound(ctx, sysInbound)
	return delay.Seconds()
}

func (x *sCtrl) EnableOutbound(ctx context.Context, tag string) {
	n, err := utility.ExtractNumber(tag)
	if err != nil {
		g.Log().Warningf(ctx, "[Ctrl] Failed to enable outbound %s", tag)
		return
	}
	json := service.XrayCfg().GetOutboundSetting(ctx, n, tag)
	if json == nil {
		g.Log().Warningf(ctx, "[Ctrl] Failed to enable outbound %s: config is empty", tag)
		return
	}
	err = service.XrayApi().AddOutbound(ctx, json)
	if err != nil {
		g.Log().Warningf(ctx, "[Ctrl] Failed to enable outbound %s: %s", tag, err.Error())
		return
	}
	x.enabledOutbounds = append(x.enabledOutbounds, tag)
	g.Log().Infof(ctx, "[Ctrl] Enabled outbound %s", tag)
}

func (x *sCtrl) DisableOutbound(ctx context.Context, tag string) {
	err := service.XrayApi().DelOutbound(ctx, tag)
	if err != nil {
		g.Log().Warningf(ctx, "[Ctrl] Failed to disable outbound %s: %s", tag, err.Error())
		return
	}
	x.enabledOutbounds = utility.RemoveElement[string](x.enabledOutbounds, tag)
	g.Log().Infof(ctx, "[Ctrl] Disabled outbound %s", tag)
}

func (x *sCtrl) loop(ctx context.Context) {

}

func (x *sCtrl) Start(ctx context.Context) {
	g.Log().Warning(ctx, "[service] Starting Ctrl...")
	x.userOutbounds = service.XrayCfg().GetUserOutboundList(ctx)
	x.outboundDelays = make([]float64, len(x.userOutbounds))
	x.delayTimeout = g.Cfg().MustGet(ctx, "controller.delayTestTimeout", 5000.0).Float64() / 1000
	for i := 0; i < len(x.outboundDelays); i++ {
		x.outboundDelays[i] = x.delayTimeout
	}
	go func() {
		for {
			time.Sleep(time.Second * 5)
			for i := 0; i < len(x.userOutbounds); i++ {
				d := x.speedtest(ctx, fmt.Sprintf("out-system-%d", i), x.delayTimeout)
				x.outboundDelays[i] =
					x.outboundDelays[i]*(1-((d/x.delayTimeout)*0.8+0.1)) +
						d*((d/x.delayTimeout)*0.8+0.1)
			}
			for i := 0; i < len(x.userOutbounds); i++ {
				g.Log().Warningf(ctx, "delay %d: %f", i, x.outboundDelays[i])
			}
			// x.EnableOutbound(ctx, "out-user-0")
			// time.Sleep(time.Second * 5)
			// x.DisableOutbound(ctx, "out-user-0")
		}
	}()
	// go func() {
	// 	x.loop(ctx)
	// }()
}

func (x *sCtrl) Stop(ctx context.Context) {
	g.Log().Warning(ctx, "[service] Stopping Ctrl...")
}
