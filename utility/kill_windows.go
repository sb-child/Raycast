package utility

import (
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"
)

// KillProcess kills the process and its child processes. The function gives the processes a second to clean up after themselves.
func KillProcess(proc *os.Process) (success bool) {
	if proc == nil {
		return true
	}
	// Usage of taskkill.exe is explained in: https://docs.microsoft.com/en-us/windows-server/administration/windows-commands/taskkill
	// Terminate the process and its children without forcing
	err := exec.Command(`C:\Windows\system32\taskkill.exe`, "/t", "/pid", strconv.Itoa(proc.Pid)).Run()
	if err == nil {
		success = true
	}
	time.Sleep(1 * time.Second)
	// Forcibly terminate the processes
	err = exec.Command(`C:\Windows\system32\taskkill.exe`, "/f", "/t", "/pid", strconv.Itoa(proc.Pid)).Run()
	if err == nil {
		success = true
	}
	// Use the built-in kill implementation as the last resort
	if proc.Kill() == nil {
		success = true
	}
	/*
		For Linux system it is necessary to use proc.Wait() to clean up after the process, or there will be a zombie process.
		For Windows it is rather strange, calling proc.Wait() on an already killed process hangs indefinitely.
		Therefore instead of calling proc.Wait(), just call proc.Release() in case go has some "resource" to release.
	*/
	_ = proc.Release()
	return
}
