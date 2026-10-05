//go:build !windows

package main

import "os/exec"

func prepareCmd(cmd *exec.Cmd) {}
func afterStart(cmd *exec.Cmd) {}
func killTree(cmd *exec.Cmd) { if cmd.Process != nil { cmd.Process.Kill() } }
