package utility

import (
	"os"
	"syscall"
	"time"
)

// KillProcess kills the process and its child processes. The function gives the processes a second to clean up after themselves.
func KillProcess(proc *os.Process) (success bool) {
	if proc == nil {
		return true
	}
	// Send SIGTERM to the process group (if any) and the process itself
	if killErr := syscall.Kill(-proc.Pid, syscall.SIGTERM); killErr == nil {
		success = true
	}
	if killErr := syscall.Kill(proc.Pid, syscall.SIGTERM); killErr == nil {
		success = true
	}
	// Wait a second for the process to clean up after itself, then force their termination.
	time.Sleep(1 * time.Second)
	if killErr := syscall.Kill(-proc.Pid, syscall.SIGKILL); killErr == nil {
		success = true
	}
	if killErr := syscall.Kill(proc.Pid, syscall.SIGKILL); killErr == nil {
		success = true
	}
	// Use the built-in kill implementation as the last resort
	if proc.Kill() == nil {
		success = true
	}
	/*
		A killed process remains in process table, laitos as the parent process must retrieve
		the exit status, or the killed process will become a zombie.
	*/
	_, _ = proc.Wait()
	_ = proc.Release()
	return
}
