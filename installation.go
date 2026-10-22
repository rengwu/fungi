package main

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
)

var buildVersion = "0.1.0-dev"
var versionPattern = regexp.MustCompile(`^v[0-9]+\.[0-9]+\.[0-9]+(?:-[A-Za-z0-9]+(?:[.-][A-Za-z0-9]+)*)?$`)

const installFile = "installation.json"

type installation struct {
	Format   int    `json:"format"`
	Current  string `json:"current"`
	Previous string `json:"previous,omitempty"`
}

func binaryName() string {
	if runtime.GOOS == "windows" {
		return "fungi.exe"
	}
	return "fungi"
}
func installedBinary(root string) string     { return filepath.Join(root, "bin", binaryName()) }
func versionDir(root, version string) string { return filepath.Join(root, "versions", version) }

func readInstallation(root string) (installation, error) {
	var state installation
	err := readJSON(filepath.Join(root, installFile), &state)
	if err == nil && (state.Format != 1 || !versionPattern.MatchString(state.Current) || (state.Previous != "" && !versionPattern.MatchString(state.Previous))) {
		err = fail("invalid_installation", "unsupported or invalid installation record")
	}
	return state, err
}

// Only our own filesystem layout identifies a managed installation. Portable
// archives and development checkouts continue to work without installation.
func executableInstallation() (root, version string, launcher bool) {
	exe, err := os.Executable()
	if err != nil {
		return
	}
	exe, err = filepath.EvalSymlinks(exe)
	if err != nil {
		return
	}
	dir := filepath.Dir(exe)
	if filepath.Base(dir) == "bin" {
		candidate := filepath.Dir(dir)
		if _, err := os.Stat(filepath.Join(candidate, installFile)); err == nil {
			return candidate, "", true
		}
	}
	if versionPattern.MatchString(filepath.Base(dir)) && filepath.Base(filepath.Dir(dir)) == "versions" {
		candidate := filepath.Dir(filepath.Dir(dir))
		if _, err := os.Stat(filepath.Join(candidate, installFile)); err == nil {
			return candidate, filepath.Base(dir), false
		}
	}
	return
}

// The small stable entry point delegates to the selected version. It is never
// overwritten during an update, including while Windows has it open.
func installationCommand(args []string, out io.Writer) (bool, func(), error) {
	root, version, launcher := executableInstallation()
	if root == "" {
		return false, func() {}, nil
	}
	unlock, err := installationLock(root)
	if err != nil {
		return false, nil, err
	}
	if launcher {
		state, e := readInstallation(root)
		if e != nil {
			unlock()
			return false, nil, e
		}
		version = state.Current
	}
	lease, err := takeVersionLease(root, version)
	unlock()
	if err != nil {
		return false, nil, err
	}
	if !launcher {
		return false, lease, nil
	}
	defer lease()
	cmd := exec.Command(filepath.Join(versionDir(root, version), binaryName()), args...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, out, os.Stderr
	if err = cmd.Run(); err != nil {
		var exited *exec.ExitError
		if errors.As(err, &exited) {
			return true, func() {}, &commandExit{exited.ExitCode()}
		}
	}
	return true, func() {}, err
}

type commandExit struct{ code int }

func (e *commandExit) Error() string { return "versioned command exited" }

func installationLock(root string) (func(), error) {
	deadline := time.Now().Add(5 * time.Second)
	for {
		unlock, err := lockSession(filepath.Join(root, ".activation.lock"))
		if err == nil {
			return unlock, nil
		}
		if asProblem(err).Code != "session_owned" {
			return nil, err
		}
		if time.Now().After(deadline) {
			return nil, fail("installation_busy", "another command is activating or inspecting this installation; retry later")
		}
		time.Sleep(20 * time.Millisecond)
	}
}

// One OS-held lease per process protects immutable package files while in use.
// A killed process cannot leave a permanently held lock. No PID polling is used.
func takeVersionLease(root, version string) (func(), error) {
	if _, err := os.Stat(filepath.Join(versionDir(root, version), binaryName())); err != nil {
		return nil, err
	}
	dir := filepath.Join(root, "leases", version)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, err
	}
	name := filepath.Join(dir, newID())
	unlock, err := lockSession(name)
	if err != nil {
		return nil, err
	}
	return func() { unlock(); _ = os.Remove(name) }, nil
}

func collectionExecutable() string {
	root, _, _ := executableInstallation()
	if root != "" {
		return installedBinary(root)
	}
	exe, _ := os.Executable()
	return exe
}

// Retain package-owned skill/helper assets for explicit draft reopening. These
// references cover alternate FUNGI_STATE_DIR roots without scanning the machine.
func pinSessionPackages(l launch, s *store) error {
	root, _, _ := executableInstallation()
	if root == "" {
		return nil
	}
	statePath, err := filepath.EvalSymlinks(s.dir)
	if err != nil {
		return err
	}
	if rel, e := filepath.Rel(root, statePath); e == nil && filepath.IsLocal(rel) {
		return fail("invalid_state_directory", "session storage must be outside the managed installation")
	}
	unlock, err := installationLock(root)
	if err != nil {
		return err
	}
	defer unlock()
	entries, err := os.ReadDir(filepath.Join(root, "versions"))
	if err != nil {
		return err
	}
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	paths := append([]string{l.Source, l.Exec, cwd}, l.Args...)
	for _, entry := range entries {
		if !entry.IsDir() || !versionPattern.MatchString(entry.Name()) {
			continue
		}
		prefix := versionDir(root, entry.Name()) + string(os.PathSeparator)
		for _, p := range paths {
			if _, value, ok := strings.Cut(p, "="); ok {
				p = value
			}
			if canonical, e := filepath.EvalSymlinks(p); e == nil {
				p = canonical
			}
			if absolute, e := filepath.Abs(p); e == nil {
				p = absolute
			}
			if runtime.GOOS == "windows" {
				p = strings.ToLower(p)
				prefix = strings.ToLower(prefix)
			}
			if strings.Contains(p+string(os.PathSeparator), prefix) {
				dir := filepath.Join(root, "references", entry.Name())
				if err = os.MkdirAll(dir, 0700); err != nil {
					return err
				}
				if err = atomicJSON(filepath.Join(dir, l.Session+".json"), map[string]string{"sessionDir": s.dir}); err != nil {
					return err
				}
				break
			}
		}
	}
	return nil
}

// Called with the activation lock held. Failure to inspect a reference retains
// it: uncertainty must not destroy a potentially recoverable draft.
func versionNeeded(root, version string) (bool, error) {
	leases, err := os.ReadDir(filepath.Join(root, "leases", version))
	if err != nil && !os.IsNotExist(err) {
		return true, err
	}
	for _, entry := range leases {
		name := filepath.Join(root, "leases", version, entry.Name())
		unlock, err := lockSession(name)
		if err != nil {
			if asProblem(err).Code == "session_owned" {
				return true, nil
			}
			return true, err
		}
		unlock()
		if err = os.Remove(name); err != nil {
			return true, err
		}
	}
	refs, err := os.ReadDir(filepath.Join(root, "references", version))
	if err != nil && !os.IsNotExist(err) {
		return true, err
	}
	for _, entry := range refs {
		name := filepath.Join(root, "references", version, entry.Name())
		var ref struct {
			SessionDir string `json:"sessionDir"`
		}
		if err = readJSON(name, &ref); err != nil || ref.SessionDir == "" {
			return true, fail("invalid_installation", "cannot inspect retained session reference "+name)
		}
		var r result
		err = readJSON(filepath.Join(ref.SessionDir, "result.json"), &r)
		if err != nil || r.Outcome != "submitted" {
			if _, e := os.Stat(ref.SessionDir); !os.IsNotExist(e) {
				return true, nil
			}
		}
		if err = os.Remove(name); err != nil {
			return true, err
		}
	}
	return false, nil
}

func pruneVersions(root string, state installation) ([]string, []string) {
	var removed, warnings []string
	entries, err := os.ReadDir(filepath.Join(root, "versions"))
	if err != nil {
		return nil, []string{err.Error()}
	}
	for _, entry := range entries {
		version := entry.Name()
		if !entry.IsDir() || !versionPattern.MatchString(version) || version == state.Current || version == state.Previous {
			continue
		}
		needed, err := versionNeeded(root, version)
		if err != nil {
			warnings = append(warnings, err.Error())
			continue
		}
		if needed {
			continue
		}
		if err = os.RemoveAll(versionDir(root, version)); err != nil {
			warnings = append(warnings, err.Error())
			continue
		}
		removed = append(removed, version)
		for _, kind := range []string{"leases", "references"} {
			if err = os.RemoveAll(filepath.Join(root, kind, version)); err != nil {
				warnings = append(warnings, err.Error())
			}
		}
	}
	return removed, warnings
}
