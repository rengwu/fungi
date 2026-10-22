package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"fungi/internal/bundle"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const releaseBase = "https://github.com/rengwu/fungi/releases"

type packageManifest struct {
	Format       int    `json:"format"`
	Version      string `json:"version"`
	Platform     string `json:"platform"`
	Arch         string `json:"arch"`
	BinarySHA256 string `json:"binarySHA256"`
	Electron     string `json:"electron"`
}

func manageInstallation(args []string, out io.Writer) error {
	action := args[0]
	root, _, _ := executableInstallation()
	archive, digest, version := "", "", ""
	var reusedRuntimeBytes int64
	seen := map[string]bool{}
	for i := 1; i < len(args); i++ {
		key := args[i]
		if seen[key] || i+1 == len(args) {
			return fail("invalid_arguments", "missing or repeated installation option "+key)
		}
		seen[key] = true
		i++
		switch key {
		case "--dir":
			if action != "install" {
				return fail("invalid_arguments", "--dir is only accepted by install")
			}
			root = args[i]
		case "--archive":
			archive = args[i]
		case "--sha256":
			digest = strings.ToLower(args[i])
		case "--version":
			version = args[i]
		default:
			return fail("invalid_arguments", "unknown installation option "+key)
		}
	}
	if action == "rollback" || action == "prune" {
		if len(args) != 1 {
			return fail("invalid_arguments", action+" takes no arguments")
		}
	} else {
		if (archive == "") != (digest == "") {
			return fail("invalid_arguments", "--archive and --sha256 must be supplied together")
		}
		if digest != "" && !validDigest(digest) {
			return fail("invalid_arguments", "SHA-256 must be 64 hexadecimal characters")
		}
		if version != "" && (!versionPattern.MatchString(version) || archive != "" || action != "update") {
			return fail("invalid_arguments", "--version requires update and a vX.Y.Z release; it cannot be combined with --archive")
		}
		if action == "install" && archive == "" {
			return fail("invalid_arguments", "install needs --archive FILE --sha256 HASH; use the platform installer to download first")
		}
	}
	if root == "" {
		if action != "install" {
			return fail("not_installed", "this is a portable or development copy; use install to create a managed installation")
		}
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		root = filepath.Join(home, ".fungi")
	}
	root, err := filepath.Abs(root)
	if err != nil {
		return err
	}
	if err = os.MkdirAll(root, 0700); err != nil {
		return err
	}
	unlockUpdate, err := lockSession(filepath.Join(root, ".update.lock"))
	if err != nil {
		if asProblem(err).Code == "session_owned" {
			return fail("update_busy", "another install/update is already running; the current version remains available")
		}
		return err
	}
	defer unlockUpdate()
	state, err := readInstallation(root)
	if err != nil && !(action == "install" && os.IsNotExist(err)) {
		return err
	}
	if state.Format == 0 {
		// A small marker makes an interrupted first install retryable without
		// ever adopting or overwriting an unrelated directory.
		var marker struct {
			Format int `json:"format"`
		}
		markerPath := filepath.Join(root, "layout.json")
		if e := readJSON(markerPath, &marker); os.IsNotExist(e) {
			entries, e := os.ReadDir(root)
			if e != nil {
				return e
			}
			for _, entry := range entries {
				if entry.Name() != ".update.lock" {
					return fail("invalid_installation", "installation directory is not empty: "+root)
				}
			}
			if e = atomicJSON(markerPath, map[string]int{"format": 1}); e != nil {
				return e
			}
		} else if e != nil || marker.Format != 1 {
			return fail("invalid_installation", "unknown installation layout")
		}
	}
	// Only the process holding .update.lock owns these disposable staging files.
	stale, err := filepath.Glob(filepath.Join(root, ".stage-*"))
	if err != nil {
		return err
	}
	for _, name := range stale {
		if err = os.RemoveAll(name); err != nil {
			return err
		}
	}
	if action == "rollback" {
		if state.Previous == "" {
			return fail("no_previous_release", "no previous release is retained")
		}
		if _, err = inspectPackage(versionDir(root, state.Previous), state.Previous); err != nil {
			return err
		}
		state.Current, state.Previous = state.Previous, state.Current
	} else if action != "prune" {
		stage, e := os.MkdirTemp(root, ".stage-")
		if e != nil {
			return e
		}
		defer os.RemoveAll(stage)
		if archive == "" {
			archive, digest, err = downloadRelease(releaseBase, version, stage)
			if err != nil {
				return err
			}
		} else {
			// Snapshot caller-owned files before verifying, so extraction reads
			// the exact bytes hashed even if the original file is changed.
			name := "release.tar.gz"
			if strings.HasSuffix(archive, ".zip") {
				name = "release.zip"
			} else if !strings.HasSuffix(archive, ".tar.gz") {
				return fail("invalid_archive", "expected .zip or .tar.gz")
			}
			snapshot := filepath.Join(stage, name)
			if err = snapshotArchive(archive, snapshot); err != nil {
				return err
			}
			archive = snapshot
		}
		got, e := fileDigest(archive)
		if e != nil {
			return e
		}
		if got != digest {
			return fail("checksum_mismatch", "release archive SHA-256 differs; existing installation is unchanged")
		}
		contents := filepath.Join(stage, "package")
		if err = os.Mkdir(contents, 0700); err != nil {
			return err
		}
		if err = bundle.Extract(archive, contents); err != nil {
			return fail("invalid_archive", err.Error())
		}
		manifest, e := inspectPackage(contents, version)
		if e != nil {
			return e
		}
		if err = probePackage(contents, manifest.Version); err != nil {
			return err
		}
		if err = os.MkdirAll(filepath.Join(root, "versions"), 0700); err != nil {
			return err
		}
		destination := versionDir(root, manifest.Version)
		if _, e := os.Stat(destination); e == nil {
			var installed struct {
				SHA256 string `json:"sha256"`
			}
			if e = readJSON(filepath.Join(destination, ".archive.json"), &installed); e != nil {
				return e
			}
			if installed.SHA256 != digest {
				return fail("version_conflict", "an installed version has different bytes; release versions are immutable")
			}
			if _, err = inspectPackage(destination, manifest.Version); err != nil {
				return err
			}
		} else if !os.IsNotExist(e) {
			return e
		} else {
			if state.Current != "" {
				oldDir := versionDir(root, state.Current)
				var old packageManifest
				if readJSON(filepath.Join(oldDir, "RELEASE.json"), &old) == nil && old.Electron == manifest.Electron {
					reusedRuntimeBytes = reuseRuntimeFiles(oldDir, contents)
				}
			}
			if err = atomicJSON(filepath.Join(contents, ".archive.json"), map[string]string{"sha256": digest}); err != nil {
				return err
			}
			if err = os.Rename(contents, destination); err != nil {
				return err
			}
		}
		if state.Format == 0 {
			if err = os.MkdirAll(filepath.Join(root, "bin"), 0700); err != nil {
				return err
			}
			candidate := filepath.Join(stage, binaryName())
			if err = copyExecutable(filepath.Join(destination, binaryName()), candidate); err != nil {
				return err
			}
			if err = os.Rename(candidate, installedBinary(root)); err != nil {
				return err
			}
			state.Format = 1
		}
		if state.Current != manifest.Version {
			state.Previous, state.Current = state.Current, manifest.Version
		}
	}
	unlock, err := installationLock(root)
	if err != nil {
		return err
	}
	defer unlock()
	if action != "prune" {
		if err = atomicJSON(filepath.Join(root, installFile), state); err != nil {
			return err
		}
	}
	removed, warnings := pruneVersions(root, state)
	return emit(out, map[string]any{"protocol": 1, "kind": "installation", "action": action, "current": state.Current, "previous": state.Previous, "command": installedBinary(root), "removed": removed, "warnings": warnings, "reusedRuntimeBytes": reusedRuntimeBytes})
}

// Package files are immutable. Reuse identical regular renderer files without
// a runtime registry or privileged symlinks on Windows. Failure is an ordinary
// copy fallback; the new, verified bytes remain untouched until a link exists.
func reuseRuntimeFiles(previous, staged string) int64 {
	oldRoot, newRoot := filepath.Join(previous, "renderer", "runtime"), filepath.Join(staged, "renderer", "runtime")
	var reused int64
	_ = filepath.WalkDir(newRoot, func(name string, entry os.DirEntry, err error) error {
		if err != nil || !entry.Type().IsRegular() {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return nil
		}
		rel, err := filepath.Rel(newRoot, name)
		if err != nil {
			return nil
		}
		old := filepath.Join(oldRoot, rel)
		oldInfo, err := os.Lstat(old)
		if err != nil || !oldInfo.Mode().IsRegular() || oldInfo.Mode().Perm() != info.Mode().Perm() || oldInfo.Size() != info.Size() {
			return nil
		}
		newHash, err := fileDigest(name)
		if err != nil {
			return nil
		}
		oldHash, err := fileDigest(old)
		if err != nil || newHash != oldHash {
			return nil
		}
		link := name + ".reuse-" + newID()
		if os.Link(old, link) != nil {
			return nil
		}
		if os.Rename(link, name) != nil {
			_ = os.Remove(link)
			return nil
		}
		reused += info.Size()
		return nil
	})
	return reused
}

func snapshotArchive(source, destination string) error {
	in, err := os.Open(source)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(destination, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	n, err := io.Copy(out, io.LimitReader(in, bundle.MaxBytes+1))
	if err == nil && n > bundle.MaxBytes {
		err = fail("invalid_archive", "archive exceeds size limit")
	}
	if e := out.Close(); err == nil {
		err = e
	}
	return err
}

func validDigest(value string) bool {
	b, err := hex.DecodeString(value)
	return err == nil && len(b) == sha256.Size
}
func fileDigest(name string) (string, error) {
	f, err := os.Open(name)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err = io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func inspectPackage(dir, expected string) (packageManifest, error) {
	var m packageManifest
	err := readJSON(filepath.Join(dir, "RELEASE.json"), &m)
	if err != nil {
		return m, err
	}
	if m.Format != 1 || !versionPattern.MatchString(m.Version) || (expected != "" && m.Version != expected) || m.Platform != runtime.GOOS || m.Arch != runtime.GOARCH || !validDigest(m.BinarySHA256) || m.Electron == "" {
		return m, fail("invalid_package", "release identity, protocol or OS/CPU does not match this machine")
	}
	digest, err := fileDigest(filepath.Join(dir, binaryName()))
	if err != nil {
		return m, err
	}
	if digest != m.BinarySHA256 {
		return m, fail("invalid_package", "release executable SHA-256 differs")
	}
	runtimeDir := filepath.Join(dir, "renderer", "runtime")
	v, err := os.ReadFile(filepath.Join(runtimeDir, "version"))
	if err != nil {
		return m, err
	}
	if strings.TrimPrefix(strings.TrimSpace(string(v)), "v") != m.Electron {
		return m, fail("invalid_package", "bundled Electron version differs")
	}
	engine := "electron"
	if runtime.GOOS == "windows" {
		engine = "electron.exe"
	}
	if runtime.GOOS == "darwin" {
		engine = filepath.Join("Electron.app", "Contents", "MacOS", "Electron")
	}
	for _, name := range []string{filepath.Join(runtimeDir, engine), filepath.Join(dir, "renderer", "main.cjs"), filepath.Join(dir, "renderer", "preload.cjs")} {
		info, e := os.Stat(name)
		if e != nil {
			return m, e
		}
		if !info.Mode().IsRegular() {
			return m, fail("invalid_package", "required package entry is not a regular file: "+name)
		}
	}
	return m, nil
}

func probePackage(dir, version string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, filepath.Join(dir, binaryName()), "--version")
	b, err := cmd.Output()
	if err != nil {
		return fail("release_probe_failed", "new executable could not report its version: "+err.Error())
	}
	var reply struct {
		Protocol int    `json:"protocol"`
		Kind     string `json:"kind"`
		Version  string `json:"version"`
	}
	if err = json.Unmarshal(b, &reply); err != nil || reply.Protocol != 1 || reply.Kind != "version" || reply.Version != version {
		return fail("release_probe_failed", "new executable does not match release metadata")
	}
	return nil
}

func copyExecutable(source, destination string) error {
	in, err := os.Open(source)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(destination, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0755)
	if err != nil {
		return err
	}
	_, err = io.Copy(out, in)
	if err == nil {
		err = out.Sync()
	}
	if e := out.Close(); err == nil {
		err = e
	}
	return err
}

func releaseAssetName() string {
	ext := ".tar.gz"
	if runtime.GOOS == "windows" {
		ext = ".zip"
	}
	return "fungi-" + runtime.GOOS + "-" + runtime.GOARCH + ext
}

func downloadRelease(base, version, destination string) (string, string, error) {
	location := base + "/latest/download/"
	if version != "" {
		location = base + "/download/" + version + "/"
	}
	name := releaseAssetName()
	checksum := filepath.Join(destination, name+".sha256")
	if err := downloadFile(location+name+".sha256", checksum, 1024); err != nil {
		return "", "", err
	}
	b, err := os.ReadFile(checksum)
	if err != nil {
		return "", "", err
	}
	fields := strings.Fields(string(b))
	if len(fields) != 2 || !validDigest(fields[0]) || strings.TrimPrefix(fields[1], "*") != name {
		return "", "", fail("invalid_release", "invalid published checksum")
	}
	archive := filepath.Join(destination, name)
	if err = downloadFile(location+name, archive, bundle.MaxBytes); err != nil {
		return "", "", err
	}
	return archive, strings.ToLower(fields[0]), nil
}

func downloadFile(url, destination string, limit int64) error {
	client := http.Client{Timeout: 15 * time.Minute}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "fungi/"+buildVersion)
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fail("download_failed", fmt.Sprintf("release download returned HTTP %d; no release installed", res.StatusCode))
	}
	if res.ContentLength > limit {
		return fail("invalid_release", "release download exceeds size limit")
	}
	f, err := os.OpenFile(destination, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	n, err := io.Copy(f, io.LimitReader(res.Body, limit+1))
	if err == nil && n > limit {
		err = fail("invalid_release", "release download exceeds size limit")
	}
	if e := f.Close(); err == nil {
		err = e
	}
	return err
}
