package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type delivery struct {
	Mode    string `json:"mode"`
	Adapter string `json:"adapter,omitempty"`
	State   string `json:"state,omitempty"`
	Reason  string `json:"reason,omitempty"`
	Thread  string `json:"thread,omitempty"`
	Command string `json:"command,omitempty"`
}
type opened struct {
	Protocol   int               `json:"protocol"`
	Kind       string            `json:"kind"`
	Session    string            `json:"session"`
	Delivery   delivery          `json:"delivery"`
	Collect    []string          `json:"collect"`
	CollectEnv map[string]string `json:"collectEnv"`
}

func main()                         { os.Exit(run(os.Args[1:], os.Stdout)) }
func emit(w io.Writer, v any) error { return json.NewEncoder(w).Encode(v) }
func run(args []string, out io.Writer) int {
	err := command(args, out)
	if err == nil {
		return 0
	}
	if _, ok := err.(*alreadyPrinted); ok {
		return 1
	}
	p := asProblem(err)
	_ = emit(out, map[string]any{"protocol": 1, "kind": "error", "error": p})
	if strings.HasPrefix(p.Code, "invalid_") {
		return 2
	}
	return 1
}
func command(args []string, out io.Writer) error {
	if len(args) == 0 || args[0] == "--help" || args[0] == "help" {
		fmt.Fprintln(out, `Fungi — local interfaces for agent workflows
  fungi describe [html|form|helper|delivery]
  fungi open SOURCE [--input FILE] [--delivery auto|wait|manual]
  fungi open --exec HELPER [--delivery MODE] -- [ARGS...]
  fungi collect SESSION
  fungi reopen SESSION [--delivery MODE]

Use one pending --delivery wait call when your harness supports it.
Never poll collect. Auto uses an exact Codex callback or reports manual fallback.
FUNGI_STATE_DIR overrides private storage for isolated tests.`)
		return nil
	}
	switch args[0] {
	case "describe":
		return describe(args[1:], out)
	case "_run":
		if len(args) != 2 {
			return fail("invalid_arguments", "internal worker needs a session")
		}
		return serveSession(args[1], out)
	case "collect":
		if len(args) != 2 {
			return fail("invalid_arguments", "collect needs an exact session")
		}
		r, err := collect(args[1])
		if err != nil {
			return err
		}
		if err = emit(out, r); err != nil {
			return err
		}
		if v, ok := r.(*result); ok && v.Outcome == "failed" {
			return &alreadyPrinted{}
		}
		return nil
	case "open", "reopen":
		l, mode, err := parseLaunch(args)
		if err != nil {
			return err
		}
		s, err := sessionStore(l.Session)
		if err != nil {
			return err
		}
		if err = s.init(l); err != nil {
			return err
		}
		if l.ReopenedFrom != "" {
			old, _ := sessionStore(l.ReopenedFrom)
			d, err := old.readDraft()
			if err != nil {
				return err
			}
			for i := range d.Artifacts {
				a := &d.Artifacts[i]
				b, err := os.ReadFile(a.Path)
				if err != nil {
					return err
				}
				a.Path = s.path(filepath.Join("artifacts", a.ID))
				if err = atomicBytes(a.Path, b); err != nil {
					return err
				}
			}
			if err = s.publish(s.path("draft.json"), d); err != nil {
				return err
			}
		}
		receipt, connection, err := startWorker(l.Session, s)
		if connection != nil {
			defer connection.close()
		}
		if err != nil {
			if r, e := s.readResult(); e == nil {
				_ = emit(out, r)
				return &alreadyPrinted{}
			}
			return err
		}
		if mode != "wait" {
			return emit(out, receipt)
		}
		r, err := waitResult(s, connection)
		if err != nil {
			return err
		}
		if err = emit(out, r); err != nil {
			return err
		}
		if r.Outcome == "failed" {
			return &alreadyPrinted{}
		}
		return nil
	default:
		return fail("invalid_arguments", "unknown command; use fungi --help")
	}
}

// Failed interactions already have a structured result. Don't print two envelopes.
type alreadyPrinted struct{}

func (*alreadyPrinted) Error() string { return "failed interaction" }

func parseLaunch(args []string) (launch, string, error) {
	l := launch{Protocol: 1, Session: newID(), Input: json.RawMessage("null")}
	mode := "auto"
	input := ""
	positionals := []string{}
	for i := 1; i < len(args); i++ {
		a := args[i]
		if a == "--" {
			l.Args = append(l.Args, args[i+1:]...)
			break
		}
		if a == "--input" || a == "--exec" || a == "--delivery" {
			if i+1 >= len(args) {
				return l, mode, fail("invalid_arguments", "missing value for "+a)
			}
			i++
			switch a {
			case "--input":
				input = args[i]
			case "--exec":
				l.Exec = args[i]
			case "--delivery":
				mode = args[i]
			}
		} else if strings.HasPrefix(a, "-") {
			return l, mode, fail("invalid_arguments", "unknown option "+a)
		} else {
			positionals = append(positionals, a)
		}
	}
	if mode != "auto" && mode != "wait" && mode != "manual" {
		return l, mode, fail("invalid_arguments", "delivery must be auto, wait or manual")
	}
	if args[0] == "reopen" {
		if len(positionals) != 1 || l.Exec != "" || input != "" || len(l.Args) > 0 {
			return l, mode, fail("invalid_arguments", "reopen needs one session and optional delivery")
		}
		old, err := sessionStore(positionals[0])
		if err != nil {
			return l, mode, err
		}
		state, err := collect(positionals[0])
		if err != nil && asProblem(err).Code != "session_interrupted" {
			return l, mode, err
		}
		if r, ok := state.(*result); ok {
			if r.Outcome == "submitted" {
				return l, mode, fail("cannot_reopen", "submitted sessions cannot be reopened")
			}
		} else if err == nil {
			return l, mode, fail("cannot_reopen", "active sessions cannot be reopened")
		}
		var previous launch
		if err = readJSON(old.path("launch.json"), &previous); err != nil {
			return l, mode, err
		}
		l.Source = previous.Source
		l.Exec = previous.Exec
		l.Args = previous.Args
		l.Input = previous.Input
		l.ReopenedFrom = previous.Session
	} else {
		if (l.Exec == "" && len(positionals) != 1) || (l.Exec != "" && len(positionals) != 0) {
			return l, mode, fail("invalid_arguments", "provide one source OR --exec helper")
		}
		if l.Exec == "" {
			l.Source = positionals[0]
			if len(l.Args) > 0 {
				return l, mode, fail("invalid_arguments", "arguments after -- require --exec")
			}
		}
		if input != "" {
			if err := readJSON(input, &l.Input); err != nil {
				return l, mode, err
			}
		}
	}
	if l.Exec != "" {
		p, err := exec.LookPath(l.Exec)
		if err != nil {
			return l, mode, err
		}
		l.Exec, err = filepath.Abs(p)
		if err != nil {
			return l, mode, err
		}
	} else {
		p, err := filepath.Abs(l.Source)
		if err != nil {
			return l, mode, err
		}
		l.Source = p
		info, err := os.Stat(p)
		if err != nil {
			return l, mode, err
		}
		if !info.Mode().IsRegular() {
			return l, mode, fail("invalid_source", "source must be a file")
		}
		if strings.EqualFold(filepath.Ext(p), ".json") {
			var b json.RawMessage
			if err = readJSON(p, &b); err != nil {
				return l, mode, err
			}
			if err = validateForm(b); err != nil {
				return l, mode, err
			}
		}
		if _, _, err = rendererCommand(); err != nil {
			return l, mode, err
		}
	}
	l.Delivery = chooseDelivery(mode)
	return l, mode, nil
}
func chooseDelivery(mode string) delivery {
	if mode == "wait" {
		return delivery{Mode: "wait", Reason: "one pending command; caller must supply a suitable host wait/background mechanism"}
	}
	if mode == "manual" {
		return delivery{Mode: "manual", Reason: "requested manual delivery; say done, then collect once"}
	}
	// Never discover a caller from process names, executable presence or the
	// most recent conversation. Bind only an explicitly inherited exact ID.
	thread := os.Getenv("CODEX_THREAD_ID")
	if uuidPattern.MatchString(thread) {
		if cmd, err := exec.LookPath("codex"); err == nil {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			err := exec.CommandContext(ctx, cmd, "queue", "--help").Run()
			cancel()
			if err == nil {
				return delivery{Mode: "callback", Adapter: "codex", State: "armed", Thread: thread, Command: cmd}
			}
		}
	}
	return delivery{Mode: "manual", Reason: "no supported exact-caller callback detected; use a host-managed --delivery wait call or say done and collect once"}
}
func publicDelivery(d delivery) delivery { d.Thread = ""; d.Command = ""; return d }

type workerPipe struct {
	scanner *bufio.Scanner
	close   func()
}

func (w *workerPipe) next(value any) error {
	if !w.scanner.Scan() {
		if err := w.scanner.Err(); err != nil {
			return err
		}
		return io.EOF
	}
	return json.Unmarshal(w.scanner.Bytes(), value)
}
func startWorker(id string, s *store) (opened, *workerPipe, error) {
	var receipt opened
	exe, err := os.Executable()
	if err != nil {
		return receipt, nil, err
	}
	cmd := exec.Command(exe, "_run", id)
	detach(cmd)
	stderr, err := os.OpenFile(s.path("worker.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return receipt, nil, err
	}
	defer stderr.Close()
	cmd.Stderr = stderr
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		return receipt, nil, err
	}
	if err = cmd.Start(); err != nil {
		return receipt, nil, err
	}
	scanner := bufio.NewScanner(pipe)
	scanner.Buffer(make([]byte, 4096), maxJSON)
	connection := &workerPipe{scanner: scanner, close: func() { _ = pipe.Close(); go cmd.Wait() }}
	var reply struct {
		opened
		Error *problem `json:"error"`
	}
	done := make(chan error, 1)
	go func() { done <- connection.next(&reply) }()
	select {
	case err = <-done:
	case <-time.After(30 * time.Second):
		err = fail("startup_timeout", "interface did not become ready within 30 seconds")
		_ = cmd.Process.Kill()
	}
	if err != nil {
		return receipt, connection, err
	}
	receipt = reply.opened
	if reply.Error != nil {
		return receipt, connection, reply.Error
	}
	if receipt.Kind != "opened" {
		return receipt, connection, fail("startup_failed", "interface ended before readiness; session "+id)
	}
	return receipt, connection, nil
}
func collect(id string) (any, error) {
	s, err := sessionStore(id)
	if err != nil {
		return nil, err
	}
	if _, err = os.Stat(s.path("launch.json")); err != nil {
		return nil, err
	}
	if r, err := s.readResult(); err == nil {
		return r, nil
	} else if !os.IsNotExist(err) {
		return nil, err
	}
	release, err := lockSession(s.path("lock"))
	if err == nil {
		release()
	} else if asProblem(err).Code != "session_owned" {
		return nil, err
	}
	// The result can be published between the first read and lock observation.
	if r, e := s.readResult(); e == nil {
		return r, nil
	}
	if err != nil {
		return map[string]any{"protocol": 1, "kind": "pending", "session": id, "instruction": "Do not poll. Await the original tool/callback, or collect once after the user says done."}, nil
	}
	return nil, fail("session_interrupted", "session owner is unavailable; saved draft can be reopened explicitly")
}
func waitResult(s *store, connection *workerPipe) (*result, error) {
	var r result
	if err := connection.next(&r); err == nil && r.Kind == "result" {
		return &r, nil
	}
	if saved, err := s.readResult(); err == nil {
		return saved, nil
	}
	return nil, fail("wait_disconnected", "waiting connection ended before a committed result; collect once to diagnose or reopen an interrupted session")
}
func rendererCommand() (string, string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", "", err
	}
	dir := os.Getenv("FUNGI_RENDERER_DIR")
	if dir == "" {
		dir = filepath.Join(filepath.Dir(exe), "renderer")
	}
	dir, err = filepath.Abs(dir)
	if err != nil {
		return "", "", err
	}
	paths := []string{
		filepath.Join(dir, "runtime", "Electron.app", "Contents", "MacOS", "Electron"),
		filepath.Join(dir, "runtime", "electron.exe"), filepath.Join(dir, "runtime", "electron"),
		filepath.Join(dir, "node_modules", "electron", "dist", "Electron.app", "Contents", "MacOS", "Electron"),
		filepath.Join(dir, "node_modules", "electron", "dist", "electron.exe"), filepath.Join(dir, "node_modules", "electron", "dist", "electron"),
	}
	for _, p := range paths {
		if info, e := os.Stat(p); e == nil && info.Mode().IsRegular() {
			return p, dir, nil
		}
	}
	return "", "", fail("renderer_unavailable", "bundled Electron runtime missing beside fungi; development setup: npm ci --prefix renderer")
}

// Scanner is shared by helper protocol tests and the owner.
func scanMessages(r io.Reader, consume func(message) error) error {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 4096), maxJSON)
	for scanner.Scan() {
		var m message
		if err := json.Unmarshal(scanner.Bytes(), &m); err != nil {
			return fail("invalid_protocol", err.Error())
		}
		if err := consume(m); err != nil {
			return err
		}
	}
	return scanner.Err()
}
