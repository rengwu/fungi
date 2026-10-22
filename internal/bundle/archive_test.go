package bundle

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestArchiveRoundTrip(t *testing.T) {
	for _, extension := range []string{".zip", ".tar.gz"} {
		t.Run(extension, func(t *testing.T) {
			source := t.TempDir()
			if e := os.MkdirAll(filepath.Join(source, "folder with spaces"), 0755); e != nil {
				t.Fatal(e)
			}
			if e := os.WriteFile(filepath.Join(source, "folder with spaces", "tool"), []byte("payload"), 0755); e != nil {
				t.Fatal(e)
			}
			if runtime.GOOS != "windows" {
				if e := os.Symlink("folder with spaces", filepath.Join(source, "current")); e != nil {
					t.Fatal(e)
				}
			}
			archive := filepath.Join(t.TempDir(), "release"+extension)
			if e := Pack(source, archive); e != nil {
				t.Fatal(e)
			}
			dest := t.TempDir()
			if e := Extract(archive, dest); e != nil {
				t.Fatal(e)
			}
			if runtime.GOOS != "windows" {
				source = filepath.Join(dest, "current", "tool")
			} else {
				source = filepath.Join(dest, "folder with spaces", "tool")
			}
			b, e := os.ReadFile(source)
			if e != nil || string(b) != "payload" {
				t.Fatal("lost content/link", e)
			}
			info, e := os.Stat(source)
			if e != nil || (runtime.GOOS != "windows" && info.Mode()&0111 == 0) {
				t.Fatal("lost executable permissions", e)
			}
		})
	}
}

func TestRejectUnsafeArchives(t *testing.T) {
	for _, name := range []string{"../outside", "/absolute", "C:/drive", "dir\\escape", "a/../escape", "duplicate", "link"} {
		t.Run(name, func(t *testing.T) {
			archive := filepath.Join(t.TempDir(), "bad.zip")
			f, e := os.Create(archive)
			if e != nil {
				t.Fatal(e)
			}
			w := zip.NewWriter(f)
			h := &zip.FileHeader{Name: name}
			h.SetMode(0644)
			if name == "link" {
				h.SetMode(os.ModeSymlink | 0777)
			}
			writer, e := w.CreateHeader(h)
			if e != nil {
				t.Fatal(e)
			}
			if name == "link" {
				writer.Write([]byte("../outside"))
			} else {
				writer.Write([]byte("data"))
			}
			if name == "duplicate" {
				duplicate := &zip.FileHeader{Name: name}
				duplicate.SetMode(0644)
				writer, e = w.CreateHeader(duplicate)
				if e != nil {
					t.Fatal(e)
				}
				writer.Write([]byte("overwrite"))
			}
			w.Close()
			f.Close()
			if e = Extract(archive, t.TempDir()); e == nil {
				t.Fatal("unsafe archive accepted")
			}
		})
	}
	archive := filepath.Join(t.TempDir(), "hardlink.tar.gz")
	f, _ := os.Create(archive)
	gz := gzip.NewWriter(f)
	tw := tar.NewWriter(gz)
	tw.WriteHeader(&tar.Header{Name: "link", Typeflag: tar.TypeLink, Linkname: "../outside", Mode: 0644})
	tw.Close()
	gz.Close()
	f.Close()
	if e := Extract(archive, t.TempDir()); e == nil {
		t.Fatal("hardlink accepted")
	}
}
