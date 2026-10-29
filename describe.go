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
	topics := map[string]any{
		"helper": map[string]any{
			"launch":    "fungi open --exec ./helper --input ./input.json --delivery auto -- arguments",
			"transport": "Protocol 2 NDJSON over stdin/stdout; diagnostics on stderr. Helpers and their runtimes are supplied by other projects.",
			"messages":  []string{"ready", "draft", "result"},
			"ack":       "Draft/result messages receive {id,type:'ack',ok:true,result?} or {id,type:'ack',ok:false,error:{code,message}}. Wait for acknowledgment before exit.",
			"artifacts": "[{key,name,path,mediaType}]; paths relative to outputDir. Saved artifacts include absolute path, bytes and sha256.",
			"recovery":  "Each draft/result replaces all data and attachments. Reopen restores the last acknowledged draft in a new session.",
		},
		"delivery": map[string]any{
			"auto":         "Exact inherited CODEX_THREAD_ID + codex queue; otherwise explicit manual fallback.",
			"wait":         "One pending command with no Fungi human timeout; requires host-managed pending/background completion. Do not poll.",
			"manual":       "Return a launch receipt; collect once on the user's done using its collect argv and collectEnv.",
			"notification": "Saved results survive notification failure. Accepted delivery does not prove host consumption; uncertain callbacks are not resent automatically.",
		},
		"installation": map[string]any{
			"install":  "fungi install --archive FILE --sha256 HASH [--dir DIRECTORY]; engine-only package format 3",
			"doctor":   "Check atomic session storage and exclusive locks without launching a helper.",
			"update":   "fungi update [--version vX.Y.Z | --archive FILE --sha256 HASH]",
			"rollback": "Select the previous retained engine release; legacy UI packages are not accepted.",
			"prune":    "Retain current, previous, active and session-referenced versions. External helpers remain their provider's responsibility.",
		},
		"examples": map[string]any{
			"guide": "fungi help examples",
			"note":  "The manual includes a complete Python standard-library helper. The source checkout also has examples/helper/main.go. No example or UI executable is bundled.",
		},
		"migration": map[string]any{
			"from": 1, "to": protocolVersion,
			"changes":  "Use --exec with an external protocol-2 helper. Direct HTML/JSON-form launches, window flags, page APIs and FUNGI_RENDERER_DIR are removed.",
			"sessions": "Use the original executable and state directory for protocol-1 sessions. No automatic session conversion.",
		},
	}
	if len(args) > 1 {
		return fail("invalid_arguments", "describe accepts at most one topic")
	}
	if len(args) == 1 {
		v, ok := topics[args[0]]
		if !ok {
			return fail("invalid_arguments", "unknown topic; use fungi describe or fungi help migration")
		}
		return emit(w, map[string]any{"protocol": protocolVersion, "topic": args[0], "capability": v, "limits": limits(), "reference": []string{executable, "help", args[0]}})
	}
	return emit(w, map[string]any{
		"protocol": protocolVersion, "version": buildVersion, "engine": "headless",
		"topics":   []string{"helper", "delivery", "installation", "examples", "migration"},
		"commands": []string{"help", "doctor", "describe", "open", "collect", "reopen", "version", "install", "update", "rollback", "prune"},
		"outcomes": []string{"submitted", "cancelled", "closed", "failed"}, "limits": limits(),
		"documentation": map[string]any{"entry": []string{executable, "help", "start"}, "topics": manualTopics()},
		"next":          "fungi help start",
	})
}
