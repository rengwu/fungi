package main

import (
	"io"
)

func limits() map[string]int {
	return map[string]int{"jsonBytes": maxJSON, "artifactBytes": maxArtifact, "totalArtifactBytes": maxArtifacts}
}
func describe(args []string, w io.Writer) error {
	topics := map[string]any{
		"html": map[string]any{
			"launch": "fungi open ./index.html --delivery wait",
			"bridge": map[string]string{
				"ready()":                       "Promise<{input,draft,attachments,limits,delivery}>",
				"saveDraft(data, attachments?)": "Promise<void>; attachments: [{key,name,blob}]. Acknowledged only after saving.",
				"submit(data, attachments?)":    "Promise<result>; attachments: [{key,name,blob}]",
				"cancel()":                      "Promise<result>",
			},
			"example":     `const {input,draft}=await fungi.ready(); await fungi.submit({notes:"Looks good"});`,
			"assets":      "Local files beneath the page directory. HTML/CSS/JS, no build required. Network and navigation denied by default.",
			"draftImages": "ready() exposes restored attachments as {key,name,blob}; save bytes with saveDraft(data,attachments).",
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
	}
	if len(args) > 1 {
		return fail("invalid_arguments", "describe accepts at most one topic")
	}
	if len(args) == 1 {
		v, ok := topics[args[0]]
		if !ok {
			return fail("invalid_arguments", "unknown topic")
		}
		return emit(w, map[string]any{"protocol": 1, "topic": args[0], "capability": v, "limits": limits()})
	}
	return emit(w, map[string]any{"protocol": 1, "version": "0.1.0-dev", "topics": []string{"html", "form", "helper", "delivery"}, "commands": []string{"describe", "open", "collect", "reopen"}, "outcomes": []string{"submitted", "cancelled", "closed", "failed"}, "limits": limits(), "next": "fungi describe html"})
}
