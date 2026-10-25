package main

import (
	"encoding/json"
	"fungi/internal/bundle"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func releaseFixture(t *testing.T, version string) (string, string) {
	t.Helper()
	dir := t.TempDir()
	cmd := exec.Command("go", "build", "-trimpath", "-ldflags", "-X main.buildVersion="+version, "-o", filepath.Join(dir, binaryName()), ".")
	if b, e := cmd.CombinedOutput(); e != nil {
		t.Fatalf("build fixture: %v: %s", e, b)
	}
	for name, source := range map[string]string{"helper" + strings.TrimPrefix(binaryName(), "fungi"): testFixture} {
		b, e := os.ReadFile(source)
		if e != nil {
			t.Fatal(e)
		}
		if e = os.WriteFile(filepath.Join(dir, name), b, 0755); e != nil {
			t.Fatal(e)
		}
	}
	helperPath := filepath.Join(dir, rendererRelativePath())
	if err := os.MkdirAll(filepath.Dir(helperPath), 0700); err != nil {
		t.Fatal(err)
	}
	if err := copyExecutable(testFixture, helperPath); err != nil {
		t.Fatal(err)
	}
	rendererDigest, err := fileDigest(helperPath)
	if err != nil {
		t.Fatal(err)
	}
	digest, e := fileDigest(filepath.Join(dir, binaryName()))
	if e != nil {
		t.Fatal(e)
	}
	if e = atomicJSON(filepath.Join(dir, "RELEASE.json"), packageManifest{Format: 2, Version: version, Platform: runtime.GOOS, Arch: runtime.GOARCH, BinarySHA256: digest, Wails: "test", RendererSHA256: rendererDigest}); e != nil {
		t.Fatal(e)
	}
	archive := filepath.Join(t.TempDir(), releaseAssetName())
	if e = bundle.Pack(dir, archive); e != nil {
		t.Fatal(e)
	}
	digest, e = fileDigest(archive)
	if e != nil {
		t.Fatal(e)
	}
	return archive, digest
}

func installedCLI(t *testing.T, binary, state string, wantOK bool, args ...string) []byte {
	t.Helper()
	cmd := exec.Command(binary, args...)
	cmd.Env = append(os.Environ(), "FUNGI_STATE_DIR="+state, "CODEX_THREAD_ID=")
	b, e := cmd.CombinedOutput()
	if (e == nil) != wantOK {
		t.Fatalf("%s %v: %v\n%s", binary, args, e, b)
	}
	return b
}

func eventually(t *testing.T, check func() bool) {
	t.Helper()
	// Match the product's startup bound while allowing each fixture to finish
	// immediately on observable readiness. Never used to wait for a person.
	deadline := time.Now().Add(30 * time.Second)
	for !check() {
		if time.Now().After(deadline) {
			t.Fatal("bounded fixture condition did not complete")
		}
		time.Sleep(20 * time.Millisecond)
	}
}

func TestManagedInstallUpdateRollbackAndSessions(t *testing.T) {
	a, shaA := releaseFixture(t, "v0.1.0-test.1")
	b, shaB := releaseFixture(t, "v0.1.0-test.2")
	c, shaC := releaseFixture(t, "v0.1.0-test.3")
	root, state := filepath.Join(t.TempDir(), "install with spaces"), t.TempDir()
	entry := installedBinary(root)
	// A failed first attempt can be retried without hand-cleaning its directory.
	installedCLI(t, testBinary, state, false, "install", "--dir", root, "--archive", a, "--sha256", strings.Repeat("0", 64))
	installedCLI(t, testBinary, state, true, "install", "--dir", root, "--archive", a, "--sha256", shaA)
	installedCLI(t, testBinary, state, true, "install", "--dir", root, "--archive", a, "--sha256", shaA)
	root, _ = filepath.EvalSymlinks(root)
	entry = installedBinary(root)
	if got := installedCLI(t, entry, state, true, "--version"); !strings.Contains(string(got), "v0.1.0-test.1") {
		t.Fatal(string(got))
	}
	// A runtime failure must not activate a new release or break the CLI.
	t.Setenv("FUNGI_TEST_RUNTIME_MISSING", "1")
	installedCLI(t, entry, state, false, "update", "--archive", b, "--sha256", shaB)
	if got := installedCLI(t, entry, state, true, "--version"); !strings.Contains(string(got), "v0.1.0-test.1") {
		t.Fatal(string(got))
	}
	t.Setenv("FUNGI_TEST_RUNTIME_MISSING", "")
	gate := filepath.Join(t.TempDir(), "gate")
	helper := filepath.Join(versionDir(root, "v0.1.0-test.1"), "helper"+strings.TrimPrefix(binaryName(), "fungi"))
	var receipt opened
	if e := json.Unmarshal(installedCLI(t, entry, state, true, "open", "--exec", helper, "--delivery", "manual", "--", "update-session", gate), &receipt); e != nil {
		t.Fatal(e)
	}
	if receipt.Collect[0] != entry {
		t.Fatal("callback path is not stable", receipt)
	}
	installedCLI(t, entry, state, false, "update", "--archive", b, "--sha256", strings.Repeat("0", 64))
	if current, e := readInstallation(root); e != nil || current.Current != "v0.1.0-test.1" {
		t.Fatal("failed update changed current", current, e)
	}
	// Updaters serialize without preventing use of the selected release.
	unlock, e := lockSession(filepath.Join(root, ".update.lock"))
	if e != nil {
		t.Fatal(e)
	}
	busy := installedCLI(t, entry, state, false, "update", "--archive", b, "--sha256", shaB)
	installedCLI(t, entry, state, true, "--version")
	unlock()
	if !strings.Contains(string(busy), "update_busy") {
		t.Fatal(string(busy))
	}
	installedCLI(t, entry, state, true, "update", "--archive", b, "--sha256", shaB)
	installedCLI(t, entry, state, true, "update", "--archive", c, "--sha256", shaC)
	if _, e = os.Stat(helper); e != nil {
		t.Fatal("active session package pruned", e)
	}
	// Close on A, then reopen through C. The old helper and acknowledged draft
	// remain usable even though A is no longer current or the rollback target.
	if e = os.WriteFile(gate, []byte("closed"), 0600); e != nil {
		t.Fatal(e)
	}
	eventually(t, func() bool { _, e := os.Stat(filepath.Join(state, receipt.Session, "result.json")); return e == nil })
	closed := installedCLI(t, receipt.Collect[0], state, true, receipt.Collect[1:]...)
	installedCLI(t, entry, state, true, "prune")
	if _, e = os.Stat(helper); e != nil {
		t.Fatal("recoverable draft assets pruned", e)
	}
	if e = os.WriteFile(gate, []byte("submitted"), 0600); e != nil {
		t.Fatal(e)
	}
	var result result
	if e = json.Unmarshal(installedCLI(t, entry, state, true, "reopen", receipt.Session, "--delivery", "wait"), &result); e != nil {
		t.Fatal(e)
	}
	if result.Outcome != "submitted" || !strings.Contains(string(result.Data), "saved before update") || len(result.Artifacts) != 1 {
		t.Fatal("recovery failed", result)
	}
	data, e := os.ReadFile(result.Artifacts[0].Path)
	if e != nil || string(data) != "artifact survives update" {
		t.Fatal("artifact lost", e)
	}
	if again := installedCLI(t, entry, state, true, "collect", receipt.Session); string(again) != string(closed) {
		t.Fatal("old result changed")
	}
	installedCLI(t, entry, state, true, "rollback")
	if got := installedCLI(t, entry, state, true, "--version"); !strings.Contains(string(got), "v0.1.0-test.2") {
		t.Fatal("rollback failed", string(got))
	}
	// Once the user has discarded the original recoverable session, its
	// reference no longer pins A. The new submitted artifact remains in state.
	if e = os.RemoveAll(filepath.Join(state, receipt.Session)); e != nil {
		t.Fatal(e)
	}
	eventually(t, func() bool {
		entries, _ := os.ReadDir(filepath.Join(root, "leases", "v0.1.0-test.1"))
		return len(entries) == 0
	})
	installedCLI(t, entry, state, true, "prune")
	if _, e = os.Stat(versionDir(root, "v0.1.0-test.1")); !os.IsNotExist(e) {
		t.Fatal("obsolete version not pruned", e)
	}
	if _, e = os.Stat(result.Artifacts[0].Path); e != nil {
		t.Fatal("pruning removed artifact", e)
	}
	installedCLI(t, entry, state, true, "collect", result.Session)
}

func TestReleaseDownloadFailures(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/large":
			w.Write([]byte(strings.Repeat("x", 100)))
		case "/partial":
			w.Header().Set("Content-Length", "100")
			w.Write([]byte("partial"))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	for _, path := range []string{"/missing", "/large", "/partial"} {
		if e := downloadFile(server.URL+path, filepath.Join(t.TempDir(), "download"), 10); e == nil {
			t.Fatalf("accepted %s", path)
		}
	}
}

func TestUnsupportedPackageFormatsAreRejected(t *testing.T) {
	for _, format := range []int{0, 1, 3} {
		dir := t.TempDir()
		manifest := packageManifest{Format: format, Version: "v0.1.0-unsupported", Platform: runtime.GOOS, Arch: runtime.GOARCH, BinarySHA256: strings.Repeat("0", 64)}
		if err := atomicJSON(filepath.Join(dir, "RELEASE.json"), manifest); err != nil {
			t.Fatal(err)
		}
		// Reject before reading or executing any legacy/future package binary.
		_, err := inspectPackage(dir, manifest.Version)
		if err == nil || asProblem(err).Code != "invalid_package" || !strings.Contains(err.Error(), "format 2") {
			t.Fatalf("format %d: %v", format, err)
		}
	}
}
