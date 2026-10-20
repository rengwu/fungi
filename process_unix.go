//go:build !windows

package main

import (
	"os/exec"
	"os/signal"
	"syscall"
)

func detach(c *exec.Cmd) { c.SysProcAttr = &syscall.SysProcAttr{Setsid: true} }

// A disconnected launch/wait pipe must never cancel a human interaction.
func ignoreBrokenPipe() { signal.Ignore(syscall.SIGPIPE) }
