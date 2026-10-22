package main

import (
	"os"
	"syscall"
)

func lockSession(path string) (func(), error) {
	p, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return nil, err
	}
	h, err := syscall.CreateFile(p, syscall.GENERIC_READ|syscall.GENERIC_WRITE, 0, nil, syscall.OPEN_ALWAYS, syscall.FILE_ATTRIBUTE_NORMAL, 0)
	if err != nil {
		// ERROR_SHARING_VIOLATION is the exclusive-open contention signal.
		// Access, path and storage errors do not establish a living owner.
		if err == syscall.Errno(32) {
			return nil, fail("session_owned", err.Error())
		}
		return nil, err
	}
	f := os.NewFile(uintptr(h), path)
	return func() { f.Close() }, nil
}
