//go:build windows

package main

import (
	"fmt"
	"os/exec"
	"syscall"
	"unsafe"
)

var (
	modKernel32                  = syscall.NewLazyDLL("kernel32.dll")
	procCreateJobObject          = modKernel32.NewProc("CreateJobObjectW")
	procSetInformationJobObject  = modKernel32.NewProc("SetInformationJobObject")
	procAssignProcessToJobObject = modKernel32.NewProc("AssignProcessToJobObject")
	jobHandle                    uintptr
)

type ioCounters struct{ a, b, c, d, e, f uint64 }
type basicLimit struct {
	PerProcessUserTimeLimit int64
	PerJobUserTimeLimit     int64
	LimitFlags              uint32
	MinimumWorkingSetSize   uintptr
	MaximumWorkingSetSize   uintptr
	ActiveProcessLimit      uint32
	Affinity                uintptr
	PriorityClass           uint32
	SchedulingClass         uint32
}
type extendedLimit struct {
	Basic                 basicLimit
	Io                    ioCounters
	ProcessMemoryLimit    uintptr
	JobMemoryLimit        uintptr
	PeakProcessMemoryUsed uintptr
	PeakJobMemoryUsed     uintptr
}

// All child processes are placed in one job that is killed when the bridge exits,
// so a crash or closed window never leaves orphaned MCP servers or tunnels behind.
func init() {
	h, _, _ := procCreateJobObject.Call(0, 0)
	if h == 0 { return }
	var info extendedLimit
	info.Basic.LimitFlags = 0x2000 // JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE
	r, _, _ := procSetInformationJobObject.Call(h, 9, uintptr(unsafe.Pointer(&info)), unsafe.Sizeof(info))
	if r != 0 { jobHandle = h }
}

func prepareCmd(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000}
}

func afterStart(cmd *exec.Cmd) {
	if jobHandle == 0 || cmd.Process == nil { return }
	ph, err := syscall.OpenProcess(0x0100|0x0001, false, uint32(cmd.Process.Pid)) // SET_QUOTA | TERMINATE
	if err != nil { return }
	procAssignProcessToJobObject.Call(jobHandle, uintptr(ph))
	syscall.CloseHandle(ph)
}

func killTree(cmd *exec.Cmd) {
	if cmd.Process == nil { return }
	k := exec.Command("taskkill", "/T", "/F", "/PID", fmt.Sprint(cmd.Process.Pid))
	prepareCmd(k)
	if k.Run() != nil { cmd.Process.Kill() }
}
