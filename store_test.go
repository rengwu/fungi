package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func testStore(t *testing.T) *store {
	t.Helper()
	s := &store{dir: filepath.Join(t.TempDir(), newID()), publish: atomicJSON}
	if err := s.init(launch{Session: filepath.Base(s.dir), Input: json.RawMessage("null")}); err != nil {
		t.Fatal(err)
	}
	return s
}
func TestDurableSubmissionAndDuplicate(t *testing.T) {
	s := testStore(t)
	id := filepath.Base(s.dir)
	if err := os.WriteFile(s.path("output/a.txt"), []byte("user bytes"), 0600); err != nil {
		t.Fatal(err)
	}
	m := message{Outcome: "submitted", Data: json.RawMessage(`{"answer":42}`), Artifacts: []attachment{{"a", "a.txt", "a.txt", "text/plain"}}}
	first, err := s.commit(id, m)
	if err != nil {
		t.Fatal(err)
	}
	repeat, err := s.commit(id, m)
	if err != nil {
		t.Fatal(err)
	}
	a, _ := json.Marshal(first)
	b, _ := json.Marshal(repeat)
	if !bytes.Equal(a, b) {
		t.Fatal("duplicate changed result")
	}
	if _, err = s.commit(id, message{Outcome: "closed"}); asProblem(err).Code != "already_completed" {
		t.Fatal("close replaced submit", err)
	}
	if err = s.saveDraft(message{Data: json.RawMessage(`{}`)}); err == nil {
		t.Fatal("draft changed completed session")
	}
	stored, err := s.readResult()
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(stored.Artifacts[0].Path)
	if err != nil || string(data) != "user bytes" {
		t.Fatal("artifact not durable", err)
	}
}
func TestFailedCommitPreservesPreviousDraft(t *testing.T) {
	s := testStore(t)
	if err := os.WriteFile(s.path("output/image"), []byte("image bytes"), 0600); err != nil {
		t.Fatal(err)
	}
	m := message{Data: json.RawMessage(`{"notes":"original","image":{"artifact":"img"}}`), Artifacts: []attachment{{"img", "image.png", "image", "image/png"}}}
	if err := s.saveDraft(m); err != nil {
		t.Fatal(err)
	}
	diskFull := errors.New("injected disk full at publication")
	s.publish = func(string, any) error { return diskFull }
	m.Data = json.RawMessage(`{"notes":"new"}`)
	if err := s.saveDraft(m); !errors.Is(err, diskFull) {
		t.Fatal(err)
	}
	if _, err := s.commit("x", message{Outcome: "submitted", Data: m.Data, Artifacts: m.Artifacts}); !errors.Is(err, diskFull) {
		t.Fatal(err)
	}
	d, err := s.readDraft()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(d.Data, []byte("original")) {
		t.Fatal("previous draft lost")
	}
	if _, err = os.Stat(d.Artifacts[0].Path); err != nil {
		t.Fatal(err)
	}
	if _, err = s.readResult(); !os.IsNotExist(err) {
		t.Fatal("failed commit published success")
	}
}
func TestArtifactConfinementAndLimits(t *testing.T) {
	s := testStore(t)
	outside := filepath.Join(t.TempDir(), "outside")
	os.WriteFile(outside, []byte("private"), 0600)
	if err := os.Symlink(outside, s.path("output/escape")); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"../launch.json", outside, "escape", "missing", "."} {
		if _, err := s.copyArtifacts([]attachment{{"k", "x", path, ""}}); err == nil {
			t.Errorf("accepted %q", path)
		}
	}
	os.WriteFile(s.path("output/ok"), []byte("ok"), 0600)
	if _, err := s.copyArtifacts([]attachment{{"k", "x", "ok", ""}, {"k", "y", "ok", ""}}); err == nil {
		t.Fatal("duplicate keys accepted")
	}
	f, err := os.Create(s.path("output/large"))
	if err != nil {
		t.Fatal(err)
	}
	f.Truncate(maxArtifact + 1)
	f.Close()
	if _, err := s.copyArtifacts([]attachment{{"k", "large", "large", ""}}); err == nil {
		t.Fatal("oversized artifact accepted")
	}
}
func TestDraftImageSurvivesClose(t *testing.T) {
	s := testStore(t)
	os.WriteFile(s.path("output/a"), []byte("selected image"), 0600)
	if err := s.saveDraft(message{Data: json.RawMessage(`{"image":{"artifact":"a"}}`), Artifacts: []attachment{{"a", "image.png", "a", "image/png"}}}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.commit("x", message{Outcome: "closed"}); err != nil {
		t.Fatal(err)
	}
	d, err := s.readDraft()
	if err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(d.Artifacts[0].Path)
	if err != nil || string(b) != "selected image" {
		t.Fatal("close lost draft image")
	}
}
func TestMalformedProtocolAndForms(t *testing.T) {
	for _, raw := range []string{`no json`, strings.Repeat("x", maxJSON+1)} {
		if err := scanMessages(strings.NewReader(raw), func(message) error { return nil }); err == nil {
			t.Fatal("accepted invalid protocol")
		}
	}
	for _, raw := range []string{
		`{"title":"a","fields":[{"name":"x","type":"invented"}]}`,
		`{"title":"a","fields":[{"name":"x","type":"range","min":1,"max":0}]}`,
		`{"title":"a","fields":[{"name":"x","type":"text"},{"name":"x","type":"text"}]}`,
		`{"title":"a","unknown":3,"fields":[]}`,
		`{"title":"a","fields":[{"name":"x","type":"text","value":12}]}`,
		`{"title":"a","fields":[{"name":"x","type":"range","min":0,"max":10,"value":11}]}`,
	} {
		if err := validateForm([]byte(raw)); err == nil {
			t.Fatal("accepted invalid form", raw)
		}
	}
}
func TestLockRejectsSecondOwner(t *testing.T) {
	path := filepath.Join(t.TempDir(), "lock")
	release, err := lockSession(path)
	if err != nil {
		t.Fatal(err)
	}
	if release2, err := lockSession(path); err == nil {
		release2()
		t.Fatal("duplicate owner")
	}
	release()
	release, err = lockSession(path)
	if err != nil {
		t.Fatal("lock not released", err)
	}
	release()
}
