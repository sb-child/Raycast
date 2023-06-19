package consts

const RAYCAST_PANEL_NAME = "raycast-panel"

type ProcessState uint8

const (
	PROC_RUNNING ProcessState = iota
	PROC_EXITED ProcessState = iota
	PROC_FAILED  ProcessState = iota
)
