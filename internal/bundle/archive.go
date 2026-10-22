// Package bundle reads and writes Fungi's portable release archives using only
// the standard library. Archives contain the package files directly at the root.
package bundle

import (
	"archive/tar"
	"archive/zip"
	"compress/flate"
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const MaxBytes int64 = 2 << 30

func Pack(dir, destination string) (err error) {
	out, err := os.OpenFile(destination, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer func() {
		if e := out.Close(); err == nil {
			err = e
		}
		if err != nil {
			os.Remove(destination)
		}
	}()
	var add func(string, fs.FileInfo, string, io.Reader) error
	var finish func() error
	if strings.HasSuffix(destination, ".zip") {
		w := zip.NewWriter(out)
		w.RegisterCompressor(zip.Deflate, func(w io.Writer) (io.WriteCloser, error) { return flate.NewWriter(w, flate.BestCompression) })
		finish = w.Close
		add = func(name string, info fs.FileInfo, link string, content io.Reader) error {
			h, e := zip.FileInfoHeader(info)
			if e != nil {
				return e
			}
			h.Name, h.Method = name, zip.Deflate
			if info.IsDir() {
				h.Name += "/"
			}
			f, e := w.CreateHeader(h)
			if e != nil {
				return e
			}
			if link != "" {
				_, e = io.WriteString(f, link)
			} else if content != nil {
				_, e = io.Copy(f, content)
			}
			return e
		}
	} else if strings.HasSuffix(destination, ".tar.gz") {
		gz, e := gzip.NewWriterLevel(out, gzip.BestCompression)
		if e != nil {
			return e
		}
		w := tar.NewWriter(gz)
		finish = func() error {
			if e := w.Close(); e != nil {
				return e
			}
			return gz.Close()
		}
		add = func(name string, info fs.FileInfo, link string, content io.Reader) error {
			h, e := tar.FileInfoHeader(info, link)
			if e != nil {
				return e
			}
			h.Name, h.Uid, h.Gid, h.Uname, h.Gname = name, 0, 0, "", ""
			if e = w.WriteHeader(h); e != nil {
				return e
			}
			if content != nil {
				_, e = io.Copy(w, content)
			}
			return e
		}
	} else {
		return fmt.Errorf("archive must end in .zip or .tar.gz")
	}
	err = filepath.WalkDir(dir, func(p string, entry fs.DirEntry, e error) error {
		if e != nil || p == dir {
			return e
		}
		info, e := entry.Info()
		if e != nil {
			return e
		}
		rel, e := filepath.Rel(dir, p)
		if e != nil {
			return e
		}
		name := filepath.ToSlash(rel)
		if info.Mode()&os.ModeSymlink != 0 {
			link, e := os.Readlink(p)
			if e != nil {
				return e
			}
			return add(name, info, link, nil)
		}
		if info.IsDir() {
			return add(name, info, "", nil)
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf("unsupported file %s", name)
		}
		f, e := os.Open(p)
		if e != nil {
			return e
		}
		defer f.Close()
		return add(name, info, "", f)
	})
	if e := finish(); err == nil {
		err = e
	}
	return err
}

// Extract requires an empty private destination. os.Root confines every write;
// links are created last, so an archive cannot write through its own symlinks.
func Extract(archive, destination string) error {
	r, err := os.OpenRoot(destination)
	if err != nil {
		return err
	}
	defer r.Close()
	seen := map[string]bool{}
	links := map[string]string{}
	var total int64
	add := func(name string, mode fs.FileMode, size int64, content io.Reader) error {
		name = strings.TrimSuffix(name, "/")
		if !fs.ValidPath(name) || name == "." || strings.ContainsAny(name, "\\:") {
			return fmt.Errorf("unsafe archive path %q", name)
		}
		if seen[name] {
			return fmt.Errorf("duplicate archive path %q", name)
		}
		seen[name] = true
		if size < 0 || size > MaxBytes || total+size > MaxBytes {
			return fmt.Errorf("archive exceeds extraction limit")
		}
		total += size
		if err := r.MkdirAll(path.Dir(name), 0755); err != nil {
			return err
		}
		if mode.IsDir() {
			return r.MkdirAll(name, 0755)
		}
		if mode&os.ModeSymlink != 0 {
			if size > 4096 {
				return fmt.Errorf("symlink too long")
			}
			b, e := io.ReadAll(io.LimitReader(content, 4097))
			if e != nil {
				return e
			}
			link := string(b)
			if link == "" || path.IsAbs(link) || strings.ContainsAny(link, "\\:") || !filepath.IsLocal(filepath.FromSlash(path.Join(path.Dir(name), link))) {
				return fmt.Errorf("unsafe symlink %q", name)
			}
			links[name] = link
			return nil
		}
		if !mode.IsRegular() {
			return fmt.Errorf("unsupported archive entry %q", name)
		}
		perm := fs.FileMode(0644)
		if mode&0111 != 0 {
			perm = 0755
		}
		f, e := r.OpenFile(name, os.O_CREATE|os.O_EXCL|os.O_WRONLY, perm)
		if e != nil {
			return e
		}
		n, e := io.Copy(f, io.LimitReader(content, size+1))
		if e == nil && n != size {
			e = fmt.Errorf("wrong size for %s", name)
		}
		if e == nil {
			e = f.Sync()
		}
		if closeErr := f.Close(); e == nil {
			e = closeErr
		}
		return e
	}
	if strings.HasSuffix(archive, ".zip") {
		z, e := zip.OpenReader(archive)
		if e != nil {
			return e
		}
		defer z.Close()
		for _, f := range z.File {
			if f.UncompressedSize64 > uint64(MaxBytes) {
				return fmt.Errorf("archive entry too large")
			}
			reader, e := f.Open()
			if e != nil {
				return e
			}
			e = add(f.Name, f.Mode(), int64(f.UncompressedSize64), reader)
			reader.Close()
			if e != nil {
				return e
			}
		}
	} else if strings.HasSuffix(archive, ".tar.gz") {
		f, e := os.Open(archive)
		if e != nil {
			return e
		}
		defer f.Close()
		gz, e := gzip.NewReader(f)
		if e != nil {
			return e
		}
		defer gz.Close()
		t := tar.NewReader(gz)
		for {
			h, e := t.Next()
			if e == io.EOF {
				break
			}
			if e != nil {
				return e
			}
			content, size := io.Reader(t), h.Size
			if h.Typeflag == tar.TypeSymlink {
				content, size = strings.NewReader(h.Linkname), int64(len(h.Linkname))
			}
			if h.Typeflag != tar.TypeReg && h.Typeflag != tar.TypeRegA && h.Typeflag != tar.TypeDir && h.Typeflag != tar.TypeSymlink {
				return fmt.Errorf("unsupported tar entry %q", h.Name)
			}
			if e = add(h.Name, h.FileInfo().Mode(), size, content); e != nil {
				return e
			}
		}
		// Consume the gzip footer as well, including its integrity check.
		if _, e = io.Copy(io.Discard, io.LimitReader(gz, MaxBytes+1)); e != nil {
			return e
		}
	} else {
		return fmt.Errorf("archive must end in .zip or .tar.gz")
	}
	for name, target := range links {
		if e := r.Symlink(target, name); e != nil {
			return e
		}
	}
	for name := range links {
		if _, e := r.Stat(name); e != nil {
			return fmt.Errorf("invalid archive symlink %s: %w", name, e)
		}
	}
	return nil
}
