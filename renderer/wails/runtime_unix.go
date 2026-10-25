//go:build !windows

package main

import "runtime"

// The platform loader resolves the linked runtime before main starts. This
// check needs neither a display nor a running GTK/Cocoa event loop.
func availableRuntime() (string, error) {
	if runtime.GOOS == "darwin" {
		return "WKWebView", nil
	}
	return "WebKitGTK 4.1", nil
}
