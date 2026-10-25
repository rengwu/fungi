package main

import (
	"io"
	"os"
)

func limits() map[string]int {
	return map[string]int{"jsonBytes": maxJSON, "artifactBytes": maxArtifact, "totalArtifactBytes": maxArtifacts}
}
func describe(args []string, w io.Writer) error {
	executable, err := os.Executable()
	if err != nil {
		return err
	}
	items := examples(executable)
	paths := []string{}
	for _, item := range items {
		if item.Name != "form" && item.Available {
			paths = append(paths, item.Source)
		}
	}
	topics := map[string]any{
		"examples": map[string]any{"items": items, "guide": "fungi help examples", "note": "Adapt these starting points; they are not a list of allowed workflows. Unavailable examples require a complete package."},
		"html": map[string]any{
			"launch": "fungi open ./index.html --delivery wait",
			"bridge": map[string]string{
				"ready()":                       "Promise<{input,draft,attachments,limits,delivery}>",
				"saveDraft(data, attachments?)": "Promise<unknown>; ignore the return value; attachments: [{key,name,blob}]. Acknowledged only after saving.",
				"submit(data, attachments?)":    "Promise<result>; attachments: [{key,name,blob}]",
				"cancel()":                      "Promise<result>",
			},
			"example":     `const {input,draft}=await fungi.ready(); await fungi.submit({notes:"Looks good"});`,
			"assets":      "Local files beneath the page directory; relative URLs, ES modules and dynamic import. No build or registration required. Bundle any libraries with the skill.",
			"draftImages": "ready() exposes restored attachments as {key,name,blob}; save bytes with saveDraft(data,attachments).",
			"ordering":    "saveDraft/submit/cancel execute in call order, including Blob encoding. Each failure rejects its own Promise; later calls may retry. Close/crash restores only acknowledged saves.",
			"rendering":   []string{"HTML/CSS and custom elements", "Canvas 2D/SVG/WebGL", "local and Blob workers", "WebAssembly", "audio/video subject to system codec and URL support", "sandboxed inline srcdoc previews"},
			"boundaries":  "Skill pages are trusted local code and can access Wails native APIs. Fungi confines served assets to the source root; browser storage isolation is best effort. Use saveDraft for durable state and sessionStorage for window-local state. Sandbox untrusted inline previews.",
			"media":       "Linux WebKitGTK cannot load direct wails:// media URLs. Fetch-to-Blob can work for small files but loads them fully; no transparent workaround or universal codec support is promised.",
			"examples":    paths,
			"authoring":   "fungi help html",
			"platforms":   "fungi help installation: runtime requirements and startup diagnostics. Feature-detect WebGL and codecs; their availability depends on the host.",
		},
		"form": map[string]any{
			"experimental": true, "launch": "fungi open ./form.json --input ./values.json --delivery wait",
			"example": map[string]any{"title": "Feedback", "fields": []any{
				map[string]any{"type": "range", "name": "duration", "label": "Duration", "min": 100, "max": 1000, "value": 300},
				map[string]any{"type": "textarea", "name": "notes", "label": "Notes"},
				map[string]any{"type": "image", "name": "reference", "label": "Paste or choose an image"},
			}},
			"types":      []string{"text", "textarea", "range", "image"},
			"precedence": "saved draft > input by field name > field value > empty",
		},
		"helper": map[string]any{
			"launch":    "fungi open --exec ./helper --delivery wait -- arguments",
			"transport": "NDJSON stdin/stdout; diagnostics only stderr. Launch protocol=1 includes session,input,draft,draftArtifacts,outputDir,limits.",
			"messages":  []any{map[string]any{"type": "ready"}, map[string]any{"id": "1", "type": "draft", "data": map[string]any{"notes": "draft"}}, map[string]any{"id": "2", "type": "result", "outcome": "submitted", "data": map[string]any{"notes": "done"}}},
			"ack":       "{id,type:'ack',ok:true,result?} or {id,type:'ack',ok:false,error:{code,message}}. Wait for ack before exit.",
			"artifacts": "[{key,name,path,mediaType}]; path relative to outputDir. Native helpers are trusted executables.",
		},
		"delivery": map[string]any{
			"auto":     "Exact inherited CODEX_THREAD_ID + codex queue. Otherwise explicit manual fallback. No discovery from latest conversations.",
			"wait":     "One pending shell tool; no Fungi human timeout. Your harness must keep it pending or provide automatic background completion. Do not poll.",
			"manual":   "Return a launch receipt. User says done, then run its collect argv once with collectEnv.",
			"codex":    "auto callback preferred; accepted queue message does not prove it was consumed.",
			"claude":   "wait via native background Bash in local interactive sessions; headless/cloud modes have time limits. Use manual if the host cannot sustain a human wait.",
			"opencode": "Installed v1 bounded Bash is unsuitable for AFK: use manual. Use wait only when the actual host exposes suitable native background completion.",
			"grok":     "wait using native background terminal completion; verify host duration/settings.",
			"kimi":     "wait via Bash run_in_background=true, disable_timeout=true; end turn and receive native completion. Do not poll.",
			"pi":       "wait using one pending Bash call with timeout omitted.",
			"recovery": "collect is non-destructive. Explicit reopen creates a new session from the saved draft. Never poll collect.",
			"limits":   "No harness setup required. Use manual when native capabilities are insufficient. See verification records; short probes do not establish indefinite AFK reliability.",
		},
		"installation": map[string]any{
			"install":  "fungi install --archive FILE --sha256 HASH [--dir DIRECTORY]; one platform package with Go CLI and Wails helper",
			"runtime":  "fungi doctor checks the installed system webview without a window. Exit 3 means a missing runtime. Bootstrap install.sh --setup-runtime or install.ps1 -SetupRuntime opts into official runtime installation; normal launches never install software.",
			"update":   "fungi update [--version vX.Y.Z]; verified GitHub release, or --archive FILE --sha256 HASH for an explicit local package",
			"rollback": "fungi rollback; switches to the previous retained release",
			"prune":    "fungi prune; retains current, previous, active processes and package assets referenced by recoverable sessions",
			"scope":    "Managed per-user installs only. Portable folders remain usable. No background updates or host policy changes. Public release/signing acceptance remains pending.",
		},
	}
	if len(args) > 1 {
		return fail("invalid_arguments", "describe accepts at most one topic")
	}
	if len(args) == 1 {
		v, ok := topics[args[0]]
		if !ok {
			return fail("invalid_arguments", "unknown topic")
		}
		return emit(w, map[string]any{"protocol": 1, "topic": args[0], "capability": v, "limits": limits(), "reference": []string{executable, "help", args[0]}})
	}
	return emit(w, map[string]any{"protocol": 1, "version": buildVersion, "topics": []string{"html", "form", "helper", "delivery", "installation", "examples"}, "commands": []string{"help", "doctor", "describe", "open", "collect", "reopen", "version", "install", "update", "rollback", "prune"}, "outcomes": []string{"submitted", "cancelled", "closed", "failed"}, "limits": limits(), "documentation": map[string]any{"entry": []string{executable, "help", "start"}, "topics": manualTopics()}, "next": "fungi help start"})
}
