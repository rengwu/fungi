package main

// Wails v2.17's file dialog references UTType on current macOS SDKs.

// #cgo LDFLAGS: -framework UniformTypeIdentifiers
import "C"
