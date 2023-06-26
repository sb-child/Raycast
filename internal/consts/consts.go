package consts

import "time"

const RAYCAST_PANEL_NAME = "raycast-panel"

type ProcessState uint8

const (
	PROC_RUNNING ProcessState = iota
	PROC_EXITED  ProcessState = iota
	PROC_FAILED  ProcessState = iota
)

type Proxy struct {
	Tag   string
	Name  string
	Index int
}

type ProxyWithLatency struct {
	Proxy
	Latency time.Duration
}

// func (x *ProxyWithLatency) LatencyFloat64() float64 {
// 	return x.Latency.Seconds()
// }
