package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	maxJSON      = 1 << 20
	maxArtifact  = 20 << 20
	maxArtifacts = 100 << 20
)

type problem struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (p *problem) Error() string      { return p.Message }
func fail(code, message string) error { return &problem{code, message} }
func asProblem(err error) *problem {
	var p *problem
	if errors.As(err, &p) {
		return p
	}
	return &problem{"operation_failed", err.Error()}
}

type attachment struct {
	Key       string `json:"key"`
	Name      string `json:"name"`
	Path      string `json:"path"`
	MediaType string `json:"mediaType"`
}
type artifact struct {
	ID        string `json:"id"`
	Key       string `json:"key"`
	Name      string `json:"name"`
	MediaType string `json:"mediaType"`
	Path      string `json:"path"`
	Bytes     int64  `json:"bytes"`
	SHA256    string `json:"sha256"`
}
type message struct {
	ID        string          `json:"id,omitempty"`
	Type      string          `json:"type"`
	Outcome   string          `json:"outcome,omitempty"`
	Data      json.RawMessage `json:"data,omitempty"`
	Artifacts []attachment    `json:"artifacts,omitempty"`
	Error     *problem        `json:"error,omitempty"`
}
type result struct {
	Protocol  int             `json:"protocol"`
	Kind      string          `json:"kind"`
	Session   string          `json:"session"`
	Outcome   string          `json:"outcome"`
	Data      json.RawMessage `json:"data,omitempty"`
	Artifacts []artifact      `json:"artifacts"`
	Error     *problem        `json:"error,omitempty"`
}
type draft struct {
	Data      json.RawMessage `json:"data"`
	Artifacts []artifact      `json:"artifacts"`
}
type launch struct {
	Protocol     int             `json:"protocol"`
	Session      string          `json:"session"`
	Source       string          `json:"source,omitempty"`
	Exec         string          `json:"exec,omitempty"`
	Args         []string        `json:"args,omitempty"`
	Input        json.RawMessage `json:"input"`
	ReopenedFrom string          `json:"reopenedFrom,omitempty"`
	Delivery     delivery        `json:"delivery"`
}
type store struct {
	dir string
	// A narrow fault-injection boundary for real filesystem commit tests.
	publish func(string, any) error
}

func newID() string { return hex.EncodeToString(randomBytes(16)) }
func randomBytes(n int) []byte {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return b
}

var idPattern = regexp.MustCompile(`^[a-f0-9]{32}$`)

func dataRoot() (string, error) {
	if p := os.Getenv("FUNGI_STATE_DIR"); p != "" {
		return filepath.Abs(p)
	}
	p, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(p, "fungi", "sessions"), nil
}
func sessionStore(id string) (*store, error) {
	if !idPattern.MatchString(id) {
		return nil, fail("invalid_session", "session must be its exact 32-character identifier")
	}
	root, err := dataRoot()
	if err != nil {
		return nil, err
	}
	return &store{dir: filepath.Join(root, id), publish: atomicJSON}, nil
}
func (s *store) path(name string) string { return filepath.Join(s.dir, name) }
func (s *store) init(l launch) error {
	if err := os.MkdirAll(filepath.Dir(s.dir), 0700); err != nil {
		return err
	}
	if err := os.Mkdir(s.dir, 0700); err != nil {
		return err
	}
	for _, d := range []string{"output", "artifacts"} {
		if err := os.Mkdir(s.path(d), 0700); err != nil {
			return err
		}
	}
	return s.publish(s.path("launch.json"), l)
}
func readJSON(path string, out any) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	b, err := io.ReadAll(io.LimitReader(f, maxJSON+1))
	if err != nil {
		return err
	}
	if len(b) > maxJSON {
		return fail("json_too_large", "JSON exceeds 1 MiB")
	}
	if err := json.Unmarshal(b, out); err != nil {
		return fail("invalid_json", err.Error())
	}
	return nil
}
func atomicJSON(path string, value any) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	if len(b) > maxJSON {
		return fail("json_too_large", "JSON exceeds 1 MiB")
	}
	return atomicBytes(path, append(b, '\n'))
}
func syncDir(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return f.Sync()
}
func atomicBytes(path string, b []byte) error {
	f, err := os.CreateTemp(filepath.Dir(path), ".pending-")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	if _, err = f.Write(b); err == nil {
		err = f.Sync()
	}
	closeErr := f.Close()
	if err != nil {
		return err
	}
	if closeErr != nil {
		return closeErr
	}
	if err = os.Rename(f.Name(), path); err != nil {
		return err
	}
	return syncDir(filepath.Dir(path))
}
func (s *store) readResult() (*result, error) {
	var r result
	if err := readJSON(s.path("result.json"), &r); err != nil {
		return nil, err
	}
	return &r, nil
}
func (s *store) readDraft() (*draft, error) {
	var d draft
	err := readJSON(s.path("draft.json"), &d)
	if os.IsNotExist(err) {
		return &draft{Data: json.RawMessage("null"), Artifacts: []artifact{}}, nil
	}
	return &d, err
}

// Output files are opened beneath a held root; traversal and escaping symlinks
// cannot make the engine copy a file outside the helper's private output folder.
func (s *store) copyArtifacts(items []attachment) ([]artifact, error) {
	out := []artifact{}
	seen := map[string]bool{}
	root, err := os.OpenRoot(s.path("output"))
	if err != nil {
		return nil, err
	}
	defer root.Close()
	var total int64
	for _, item := range items {
		if item.Key == "" || seen[item.Key] || item.Name == "" || len(item.Key) > 200 || len(item.Name) > 255 {
			return nil, fail("invalid_artifact", "artifacts need unique nonempty keys and names (200/255 byte limits)")
		}
		seen[item.Key] = true
		if !filepath.IsLocal(item.Path) {
			return nil, fail("invalid_artifact", "artifact path must be relative to outputDir")
		}
		// Stat first rejects FIFOs before Open can block. Helpers are trusted
		// native programs, but a raced file is still checked after opening.
		info, err := root.Stat(item.Path)
		if err != nil {
			return nil, err
		}
		if !info.Mode().IsRegular() || info.Size() > maxArtifact {
			return nil, fail("invalid_artifact", "artifact must be a regular file at most 20 MiB")
		}
		f, err := root.Open(item.Path)
		if err != nil {
			return nil, err
		}
		info, err = f.Stat()
		if err != nil || !info.Mode().IsRegular() {
			f.Close()
			return nil, fail("invalid_artifact", "artifact is not a regular file")
		}
		b, err := io.ReadAll(io.LimitReader(f, maxArtifact+1))
		f.Close()
		if err != nil {
			return nil, err
		}
		if len(b) > maxArtifact {
			return nil, fail("artifact_too_large", "artifact exceeds 20 MiB")
		}
		total += int64(len(b))
		if total > maxArtifacts {
			return nil, fail("artifacts_too_large", "attachments exceed 100 MiB")
		}
		h := sha256.Sum256(b)
		// Content-addressed files avoid accumulating copies during autosave.
		id := hex.EncodeToString(h[:])
		path := s.path(filepath.Join("artifacts", id))
		if err := atomicBytes(path, b); err != nil {
			return nil, err
		}
		out = append(out, artifact{id, item.Key, filepath.Base(item.Name), item.MediaType, path, int64(len(b)), id})
	}
	return out, nil
}
func (s *store) saveDraft(m message) error {
	defer s.pruneUnreferenced()
	if _, err := os.Stat(s.path("result.json")); err == nil {
		return fail("already_completed", "session is already complete")
	}
	arts, err := s.copyArtifacts(m.Artifacts)
	if err != nil {
		return err
	}
	if len(m.Data) == 0 {
		m.Data = json.RawMessage("null")
	}
	if err := s.publish(s.path("draft.json"), draft{m.Data, arts}); err != nil {
		return err
	}
	return s.pruneArtifacts(arts)
}
func (s *store) pruneArtifacts(keep []artifact) error {
	ids := map[string]bool{}
	for _, a := range keep {
		ids[a.ID] = true
	}
	entries, err := os.ReadDir(s.path("artifacts"))
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if !ids[entry.Name()] {
			_ = os.Remove(s.path(filepath.Join("artifacts", entry.Name())))
		}
	}
	return nil
}
func (s *store) commit(id string, m message) (*result, error) {
	defer s.pruneUnreferenced()
	switch m.Outcome {
	case "submitted", "cancelled", "closed", "failed":
	default:
		return nil, fail("invalid_outcome", "unknown terminal outcome")
	}
	if m.Outcome != "submitted" {
		m.Data = nil
		m.Artifacts = nil
	} else if len(m.Data) == 0 {
		m.Data = json.RawMessage("null")
	}
	if m.Outcome == "failed" && m.Error == nil {
		m.Error = &problem{"helper_failed", "interface reported a failure"}
	}
	if m.Outcome != "failed" {
		m.Error = nil
	}
	arts, err := s.copyArtifacts(m.Artifacts)
	if err != nil {
		return nil, err
	}
	r := &result{1, "result", id, m.Outcome, m.Data, arts, m.Error}
	// Draft images remain recoverable on close/cancel/failure. Bound the union,
	// so preserving them cannot silently exceed the advertised session budget.
	kept := append([]artifact{}, arts...)
	if d, err := s.readDraft(); err == nil {
		kept = append(kept, d.Artifacts...)
	}
	unique := map[string]bool{}
	var total int64
	for _, a := range kept {
		if !unique[a.ID] {
			unique[a.ID] = true
			total += a.Bytes
		}
	}
	if total > maxArtifacts {
		return nil, fail("artifacts_too_large", "saved draft and result exceed the 100 MiB session limit")
	}
	if old, err := s.readResult(); err == nil {
		a, _ := json.Marshal(old)
		b, _ := json.Marshal(r)
		if string(a) == string(b) {
			return old, nil
		}
		return nil, fail("already_completed", "a different terminal result is already committed")
	} else if !os.IsNotExist(err) {
		return nil, err
	}
	if err := s.publish(s.path("result.json"), r); err != nil {
		return nil, err
	}
	// Keep draft attachments for explicit recovery, as well as final artifacts.
	if d, err := s.readDraft(); err == nil {
		arts = append(arts, d.Artifacts...)
	}
	_ = s.pruneArtifacts(arts)
	return r, nil
}
func (s *store) pruneUnreferenced() {
	keep := []artifact{}
	if d, err := s.readDraft(); err == nil {
		keep = append(keep, d.Artifacts...)
	} else {
		return
	}
	if r, err := s.readResult(); err == nil {
		keep = append(keep, r.Artifacts...)
	} else if !os.IsNotExist(err) {
		return
	}
	_ = s.pruneArtifacts(keep)
}
func validateForm(b []byte) error {
	var f struct {
		Title   string `json:"title"`
		Context string `json:"context"`
		Fields  []struct {
			Type     string          `json:"type"`
			Name     string          `json:"name"`
			Label    string          `json:"label"`
			Value    json.RawMessage `json:"value"`
			Min      *float64        `json:"min"`
			Max      *float64        `json:"max"`
			Step     *float64        `json:"step"`
			Required bool            `json:"required"`
		} `json:"fields"`
	}
	d := json.NewDecoder(strings.NewReader(string(b)))
	d.DisallowUnknownFields()
	if err := d.Decode(&f); err != nil {
		return fail("invalid_form", err.Error())
	}
	if f.Title == "" || f.Fields == nil {
		return fail("invalid_form", "form needs a title and fields array")
	}
	seen := map[string]bool{}
	for _, field := range f.Fields {
		if field.Name == "" || seen[field.Name] || field.Name == "__proto__" || field.Name == "constructor" || field.Name == "prototype" {
			return fail("invalid_form", "field names must be unique, nonempty and not reserved")
		}
		seen[field.Name] = true
		switch field.Type {
		case "text", "textarea":
			if len(field.Value) > 0 {
				var value string
				if string(field.Value) == "null" || json.Unmarshal(field.Value, &value) != nil {
					return fail("invalid_form", "text defaults must be strings: "+field.Name)
				}
			}
		case "image":
			if len(field.Value) > 0 && string(field.Value) != "null" {
				return fail("invalid_form", "image defaults must be empty; supply image bytes through the page bridge")
			}
		case "range":
			if field.Min == nil || field.Max == nil || *field.Min >= *field.Max || (field.Step != nil && *field.Step <= 0) {
				return fail("invalid_form", fmt.Sprintf("invalid range: %s", field.Name))
			}
			if len(field.Value) > 0 {
				var value float64
				if string(field.Value) == "null" || json.Unmarshal(field.Value, &value) != nil || value < *field.Min || value > *field.Max {
					return fail("invalid_form", "range default must be a number within min/max: "+field.Name)
				}
			}
		default:
			return fail("invalid_form", "unsupported field type: "+field.Type)
		}
	}
	return nil
}
