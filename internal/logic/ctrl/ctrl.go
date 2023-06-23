package ctrl

import (
	"context"
	"fmt"
	"raycast/internal/service"
	"raycast/utility"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/grand"
	"golang.org/x/exp/slices"
)

type sCtrl struct {
	delayTimeout                float64
	outboundDelays              []float64
	userOutboundUploadTraffic   []int64
	userOutboundDownloadTraffic []int64
	userOutbounds               []string
	enabledOutbounds            []string
	outboundNextSpeedtest       []*gtime.Time
	ctxCancel                   context.CancelFunc
	userOutboundLock            sync.Mutex
	taskLock                    sync.WaitGroup
}

func init() {
	service.RegisterCtrl(New())
}

func New() *sCtrl {
	return &sCtrl{}
}

func (x *sCtrl) traffic(ctx context.Context, tag string) (up int64, dn int64) {
	dn, err := service.XrayApi().Stat(ctx, false, tag, true)
	if err != nil {
		g.Log().Warningf(ctx, "[Ctrl] Failed to get traffic for %s: %s", tag, err)
		return 0, 0
	}
	up, err = service.XrayApi().Stat(ctx, false, tag, false)
	if err != nil {
		g.Log().Warningf(ctx, "[Ctrl] Failed to get traffic for %s: %s", tag, err)
		return 0, 0
	}
	return
}

func (x *sCtrl) trafficLoop(ctx context.Context, tag string) {
	x.taskLock.Add(1)
	defer x.taskLock.Done()
	n, err := utility.ExtractNumber(tag)
	if err != nil {
		g.Log().Warningf(ctx, "[Ctrl/Traffic|%s] Failed to start", tag)
		return
	}
	tk := time.NewTicker(time.Second)
	defer tk.Stop()
	// g.Log().Infof(ctx, "[Ctrl/speedtest|%s] Next test at %s",
	// 	tag,
	// )
	lastUp := int64(0)
	lastDn := int64(0)
	f := func() {
		x.userOutboundLock.Lock()
		_, needTest := slices.BinarySearch(x.enabledOutbounds, tag)
		x.userOutboundLock.Unlock()
		if !needTest {
			x.userOutboundUploadTraffic[n] = 0
			x.userOutboundDownloadTraffic[n] = 0
			lastUp = 0
			lastDn = 0
			return
		}
		up, dn := x.traffic(ctx, tag)
		if up < lastUp || dn < lastDn {
			// maybe stats are reset
			// g.Log().Warningf(ctx, "[Ctrl/Traffic|%s] reset Up %dB/s Down %dB/s", tag, up, dn)
		} else {
			up -= lastUp
			dn -= lastDn
		}
		lastUp = up
		lastDn = dn
		x.userOutboundUploadTraffic[n] = up
		x.userOutboundDownloadTraffic[n] = dn
		g.Log().Warningf(ctx, "[Ctrl/Traffic|%s] Up %dB/s Down %dB/s", tag, up, dn)
	}
	for {
		if len(tk.C) >= 1 {
			<-tk.C
			f()
		}
		if utility.CheckCancel(ctx) {
			break
		}
	}
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
	resp, err := c.Get(ctx, "https://maps.gstatic.com/generate_204")
	if err != nil {
		g.Log().Infof(ctx, "[Ctrl] Failed to speedtest %s, deleting %s: %s", tag, sysInbound, err.Error())
		service.XrayApi().DelInbound(ctx, sysInbound)
		return timeout
	}
	stopTime := gtime.Now()
	resp.Close()
	delay := stopTime.Sub(startTime)
	// g.Log().Infof(ctx, "[Ctrl] Speedtest %s done: %s, deleting %s", tag, delay.String(), sysInbound)
	service.XrayApi().DelInbound(ctx, sysInbound)
	return delay.Seconds()
}

func (x *sCtrl) speedtestLoop(ctx context.Context, tag string) {
	x.taskLock.Add(1)
	defer x.taskLock.Done()
	n, _ := utility.ExtractNumber(tag)
	delay := time.Millisecond * time.Duration(x.delayTimeout*1000)
	x.outboundNextSpeedtest[n] = gtime.Now().Add(
		grand.D(time.Millisecond*500,
			time.Second+delay),
	)
	g.Log().Infof(ctx, "[Ctrl/speedtest|%s] Next test at %s",
		tag, x.outboundNextSpeedtest[n].String(),
	)
	for {
		if utility.CheckCancel(ctx) {
			break
		}
		if x.outboundNextSpeedtest[n].Before(gtime.Now()) {
			r := x.speedtest(ctx, tag, x.delayTimeout)
			rd := time.Millisecond * time.Duration(r*1000)
			x.outboundDelays[n] =
				x.outboundDelays[n]*(1-((r/x.delayTimeout)*0.8+0.1)) +
					r*((r/x.delayTimeout)*0.8+0.1)
			x.outboundNextSpeedtest[n] = x.outboundNextSpeedtest[n].Add(
				grand.D(rd+delay*4, rd+delay*8),
			)
			g.Log().Infof(ctx, "[Ctrl/speedtest|%s] Result is %f, Next test at %s",
				tag, r, x.outboundNextSpeedtest[n].String(),
			)
		} else {
			time.Sleep(time.Millisecond * 100)
		}
	}
}

func (x *sCtrl) EnableOutbound(ctx context.Context, tag string) {
	x.userOutboundLock.Lock()
	defer x.userOutboundLock.Unlock()
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
	x.userOutboundLock.Lock()
	defer x.userOutboundLock.Unlock()
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
	var cctx context.Context
	cctx, x.ctxCancel = context.WithCancel(ctx)
	x.userOutbounds = service.XrayCfg().GetUserOutboundList(cctx)
	x.outboundDelays = make([]float64, len(x.userOutbounds))
	x.outboundNextSpeedtest = make([]*gtime.Time, len(x.userOutbounds))
	x.userOutboundUploadTraffic = make([]int64, len(x.userOutbounds))
	x.userOutboundDownloadTraffic = make([]int64, len(x.userOutbounds))
	x.delayTimeout = g.Cfg().MustGet(cctx, "controller.delayTestTimeout", 5000.0).Float64() / 1000
	for i := 0; i < len(x.outboundDelays); i++ {
		x.outboundDelays[i] = x.delayTimeout
		go x.speedtestLoop(cctx, fmt.Sprintf("out-system-%d", i))
		go x.trafficLoop(cctx, fmt.Sprintf("out-user-%d", i))
	}
	go func() {
		for {
			time.Sleep(time.Second * 5)
			// for i := 0; i < len(x.userOutbounds); i++ {
			// 	d := x.speedtest(ctx, fmt.Sprintf("out-system-%d", i), x.delayTimeout)
			// 	x.outboundDelays[i] =
			// 		x.outboundDelays[i]*(1-((d/x.delayTimeout)*0.8+0.1)) +
			// 			d*((d/x.delayTimeout)*0.8+0.1)
			// }
			// up, dn := x.traffic(ctx, "out-user-22")
			// g.Log().Warningf(cctx, "%d %d", up, dn)
			for i := 0; i < len(x.userOutbounds); i++ {
				g.Log().Warningf(cctx, "delay %d: %f", i, x.outboundDelays[i])
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
	x.ctxCancel()
	x.taskLock.Wait()
}
