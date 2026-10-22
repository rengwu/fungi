// Development packaging tool; it is not shipped or needed on the user's machine.
package main

import (
	"crypto/sha256"
	"fmt"
	"fungi/internal/bundle"
	"io"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: archive DIRECTORY OUTPUT.zip|OUTPUT.tar.gz")
		os.Exit(2)
	}
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
func run() error {
	if err := bundle.Pack(os.Args[1], os.Args[2]); err != nil {
		return err
	}
	f, err := os.Open(os.Args[2])
	if err != nil {
		return err
	}
	defer f.Close()
	h := sha256.New()
	if _, err = io.Copy(h, f); err != nil {
		return err
	}
	return os.WriteFile(os.Args[2]+".sha256", []byte(fmt.Sprintf("%x  %s\n", h.Sum(nil), filepath.Base(os.Args[2]))), 0644)
}
