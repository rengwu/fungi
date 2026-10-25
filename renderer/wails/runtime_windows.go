package main

import (
	"errors"
	"github.com/wailsapp/go-webview2/webviewloader"
)

func availableRuntime() (string, error) {
	version, err := webviewloader.GetAvailableCoreWebView2BrowserVersionString("")
	if err != nil {
		return "", err
	}
	if version == "" {
		return "", errors.New("Microsoft WebView2 Runtime is not installed")
	}
	// Match the minimum required by the pinned Wails v2.17 SDK.
	comparison, err := webviewloader.CompareBrowserVersions(version, "94.0.992.31")
	if err != nil || comparison < 0 {
		return version, errors.New("Microsoft WebView2 Runtime needs updating")
	}
	return "WebView2 " + version, nil
}
