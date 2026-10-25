---
title: First implementation contract
updated: 2026-10-22
status: Prototype contract
---
# First implementation contract

This contract describes the initial implementation of the [accepted first milestone](planning.md#accepted-first-milestone-scope). The prototype exists; the [verification report](verification/runs/2026-10-20-first-proof.md) identifies passing scenarios and remaining gaps. The result boundary is intended to stabilize before broader UI composition.

Go owns sessions, persistence and delivery. Electron owns the native web window. Skills provide local HTML/assets, a small JSON form, or an executable helper. A skill sequences workflows; the engine runs one interaction at a time per session, with independent sessions allowed concurrently.

The [authoring guide](authoring.md) is the compact entry point for new skills. Protocol 1's existing meanings are the compatibility boundary: preserve them for additive extensions and explicitly version breaking changes. This does not make the experimental flat form schema or untested platforms stable.

## Four CLI operations

| Command | Contract |
| --- | --- |
| `fungi describe [topic]` | Return a compact machine-readable overview, or details for `html`, `form`, `helper`, or `delivery`, including examples and limits |
| `fungi open ./review.html --input ./input.json` | Create a session and launch the interface; `.json` selects the reference form renderer |
| `fungi collect SESSION_ID` | Return the saved result, a pending response, or an explicit error; never delete data |
| `fungi reopen SESSION_ID` | Start a new, linked interaction using saved launch input and the latest saved draft |

`open` also accepts `--exec /path/to/helper -- arg1 arg2` instead of a page/form path. Helper arguments are passed as an argument array, without shell interpolation. A source path is sufficient; no manifest, registration service or package installation is required.

`open` and `reopen` accept `--delivery auto|wait|manual`, defaulting to `auto`. Input is an optional JSON value read from a file and snapshotted at launch. Ordinary command output is JSON; diagnostics go to stderr. `--help` remains short human-readable text. Discoverable limits and format examples belong in the binary, so an agent can learn the interface without source checkout access.

The exact session identifier is the handle for all later operations. Never accept “latest session” as an implicit target.

## Waiting and return

| Mode | Launching command | Completion |
| --- | --- | --- |
| `auto` | Return a receipt after the session owner and interface are ready | Use an established exact-caller callback; otherwise report manual delivery in the receipt |
| `wait` | Keep one call pending, silent after startup diagnostics | Return one final result when the interaction ends; use a host-managed background tool when appropriate |
| `manual` | Return a receipt after readiness | Preserve the result for collection after the user says “done” |

`auto` cannot choose the caller's tool options. A harness that supplies background completion uses `wait` through that tool; a harness that can retain a pending call can use `wait` directly. The delivery description must explain the appropriate invocation. Finding an executable on PATH is not evidence of an active callback route.

Example launch receipt:

```json
{
  "protocol": 1,
  "kind": "opened",
  "session": "SESSION_ID",
  "delivery": {
    "mode": "callback",
    "adapter": "codex",
    "state": "armed"
  },
  "collect": ["fungi", "collect", "SESSION_ID"],
  "collectEnv": {"FUNGI_STATE_DIR": "/absolute/state/directory"}
}
```

`armed` means the route was prepared, not that a future message is guaranteed. A manual receipt instead reports `mode: "manual"` and a useful reason. To collect later, execute the `collect` argument array with `collectEnv`; a later tool invocation need not inherit the original storage override. There are no periodic progress messages, scheduled model checks or repeated wait calls. Ordinary human waiting has no Fungi deadline; startup and communication operations may have separate bounded timeouts.

Separate the command's connection from the session owner. Losing a waiting connection is not human cancellation. A harness may still terminate related processes; preserve committed state and support explicit recovery rather than promise survival of every host shutdown.

## Result and lifecycle

A session produces at most one authoritative terminal result:

| Outcome | Meaning |
| --- | --- |
| `submitted` | Human response and referenced artifact bytes have been committed |
| `cancelled` | Explicit cancellation without a submitted answer |
| `closed` | Window closed without submission; latest saved draft remains available |
| `failed` | Known renderer/helper failure, including exit without a valid outcome |

The owner serializes competing terminal events. Once a submission is committed, window exit cannot replace it. An identical repeated submission returns the existing receipt; a different submission after completion reports a conflict. Neither creates a second result or a fresh notification.

Example result, shared by a pending return and later collection:

```json
{
  "protocol": 1,
  "kind": "result",
  "session": "SESSION_ID",
  "outcome": "submitted",
  "data": {"duration": 300, "notes": "A little slower", "reference": {"artifact": "reference"}},
  "artifacts": [
    {
      "id": "ARTIFACT_ID",
      "key": "reference",
      "name": "reference.png",
      "mediaType": "image/png",
      "path": "/absolute/session/path/artifacts/ARTIFACT_ID.png",
      "bytes": 12345
    }
  ]
}
```

The session ID also identifies its sole result; a second independent result ID is unnecessary initially. An artifact's `key` is a unique name within the submission, chosen by the page/helper so its data can reference the attachment before the owner assigns a content-derived ID and storage path. The core preserves `data` as opaque JSON. Non-submitted outcomes have no answer data; `failed` includes a structured error. Drafts are not silently returned as submitted answers. Images travel as files referenced by the envelope, not base64 text in the model conversation.

`collect` is repeatable. While the owner is alive and no terminal record exists, it returns `kind: "pending"`; the skill must not call it repeatedly to wait. If the owner died without committing a result, report `session_interrupted`, preserve its files and offer explicit reopen. Do not manufacture a successful result from partial files or infer ownership from a reusable PID alone.

Reopen never rewrites the old result. It creates a new session with `reopenedFrom`, captures the current caller afresh, and restores saved input/draft. Initially permit reopening closed, cancelled, failed or interrupted sessions; reject reopening an active or submitted session. The original source must remain available, or reopen fails clearly. Restoration does not promise to reconstruct arbitrary unsaved page state or a historical copy of every source asset.

Exit code 0 means the command succeeded, including a closed/cancelled result or a pending collection response. Operational failures, including a failed interaction, use exit code 1; invalid requests use 2. Failure output remains structured. A tool timeout is a transport event, not a new human outcome.

## Page authoring

The page gets a deliberately small bridge:

| Method | Contract |
| --- | --- |
| `await fungi.ready()` | Return launch `input`, saved `draft` or null, restored Blob attachments, limits and delivery mode |
| `await fungi.saveDraft(data, attachments?)` | Persist a JSON draft and optional Blob attachments; resolve after saving |
| `await fungi.submit(data, attachments?)` | Commit JSON and optional `{key, name, blob}` attachments; resolve with a receipt after persistence |
| `await fungi.cancel()` | Commit a cancelled outcome |

The renderer handles Blob/File transfer internally. Page authors do not choose output paths, callback destinations or shell commands. Relative local assets remain usable; the default example needs no network or build step. Wails skill pages and native helpers are trusted local application code with native access. Browser storage isolation is best effort; Linux direct-media limitations are accepted.

Save, submit and cancel calls are ordered from invocation, before asynchronous Blob encoding. A failure rejects its own call and does not poison later retries. Each save/submit supplies a complete state and attachment list; omitted attachments mean none. Closing or crashing recovers only acknowledged saves, including with queued calls. Library bundles, local/Blob workers, WebAssembly, canvas and media (subject to the system engine) remain page capabilities rather than new bridge methods.

After successful submission, the reference renderer closes the window. If persistence fails, it stays open, shows the failure and retains the user's current input for retry. A committed response remains successful even if agent notification later fails; the interface can report that collection is still available.

The reference form saves named values as drafts when they change and on a normal close where possible. Custom HTML opts into draft saving. A crash can recover only the last acknowledged save; no universal DOM-state capture is promised.

## Flat JSON reference form

Keep this schema experimental and separate from the stable session/result boundary. Start with `title`, ordered `fields`, and optional text context. Initial field types are `text`, `textarea`, `range`, and `image`; submission/cancellation controls come from the renderer. Field names are unique. Range values are numbers, text values are strings, and image fields contribute artifact references by field name.

```json
{
  "title": "Animation feedback",
  "fields": [
    {"type": "range", "name": "duration", "label": "Duration", "min": 100, "max": 1000, "value": 300},
    {"type": "textarea", "name": "notes", "label": "Notes"},
    {"type": "image", "name": "reference", "label": "Paste or choose an image"}
  ]
}
```

For this renderer, a supplied input object overrides field defaults by name; a restored draft takes precedence over both. Image fields accept one image each initially. Their submitted value is `{"artifact":"FIELD_NAME"}`, referring to an artifact key, or null when empty. File metadata lives in `artifacts`. Image draft recovery requires saving the selected bytes as well as the reference; the reference renderer owns that transfer. Unknown fields/types or invalid ranges produce a useful validation error.

Nested layouts, custom JSON components, expressions and external live updates remain deferred. Local JavaScript can still drive an HTML preview continuously without involving the agent.

## Native helper contract

Use newline-delimited JSON on stdin/stdout and diagnostics on stderr. Empty lines are ignored; other malformed records fail the session. The owner sends one launch message containing `protocol`, `session`, `input`, `draft`, and a private `outputDir`. The helper emits `ready`, then optional draft messages and one result message:

```jsonl
{"type":"ready"}
{"type":"draft","data":{"notes":"In progress"}}
{"type":"result","outcome":"submitted","data":{"notes":"Done"},"artifacts":[{"key":"reference","name":"reference.png","path":"reference.png","mediaType":"image/png"}]}
```

Artifact paths resolve inside `outputDir`. The owner validates and commits the files, then replies on stdin with a committed receipt; the helper can close its interface and exit. On commit failure, reply with an error so the helper may retain its UI and retry. Protocol version is negotiated by the launch message; unsupported versions fail explicitly.

Drain complete output messages before treating process exit as failure. A ready message, zero exit code or disappearing window alone is never a submitted answer. Draft and final outcomes are sufficient initially; no general RPC framework or plugin SDK is needed.

## Persistence and ownership

Use one session-owning Go process and a private session directory. The detached worker returns readiness and completion through a pipe; closing that pipe does not cancel it. An OS file lock establishes ownership. There is no local HTTP listener or permanent daemon. Electron and native helpers feed the same owner-side commit function. No permanent system service is required. Keep input, draft, terminal result, artifact bytes and delivery bookkeeping distinct. Artifact IDs are SHA-256 values; unreferenced copies are removed after a save attempt to bound draft churn. Native helpers remain responsible for their temporary output directory.

For the macOS proof, use staged files and an explicitly tested commit sequence: save/flush referenced artifacts before publishing the terminal result; acknowledge only after the commit completes. Preserve the previous valid draft/result on failure. Cross-platform file replacement and durability need separate verification before Windows/Linux support is claimed. The contract specifies observable behavior, not a new storage framework.

Use owner-generated, content-addressed artifact paths and reject traversal, symlink escapes and incomplete files. Initial limits are 1 MiB per JSON control message, 20 MiB per artifact, and 100 MiB total artifact bytes per session. Advertise these through discovery and `ready()`; reject excess input before reporting success. Limits may change with evidence without expanding the protocol.

Collection is non-destructive. First-milestone sessions remain until explicitly removed; no cleanup is inferred from a successful callback or collection. Shutdown must release child processes and live connections without deleting saved results.

## Delivery boundary and acceptance checks

Bind the caller when launching. Page content cannot override it. Prefer one host-managed return route; do not simultaneously arm both a callback and a second model-facing completion watcher. Callback messages identify the saved session to collect rather than embedding answer content as instructions.

Track saved result and notification state separately. Use bounded local retries only when non-delivery is known or the receiver supports deduplication. An ambiguous acknowledgement must not trigger blind repeated model turns. An accepted callback does not prove consumption. Preserve the result regardless.

The [six-harness inventory](research/delivery.md#first-milestone-harness-inventory) supplies candidate routes, not passing tests. The implementation must verify the actual installed harness behavior and retain the accepted manual fallback where unavailable. The [fourth interview](planning.md#no-setup-delivery-decision-on-2026-10-20) explicitly forbids required harness setup: no plugin, hook, startup configuration or upgrade prerequisite. An unsupported version's manual result is recorded separately from automatic-return verification. No configuration edits or substitute conversations are part of this contract.

| Check | Required observation |
| --- | --- |
| HTML, flat JSON and native helper | Same lifecycle and result/artifact envelope |
| Human wait in each of six harnesses | No repeated model checks; native completion reaches the original conversation where supported, otherwise explicitly verify the accepted manual fallback |
| Two simultaneous conversations | Results and notices remain isolated |
| Submit twice; submit then window exit | One immutable result; no duplicate notification |
| Cancel, close, crash | Distinct outcomes; latest saved draft recoverable |
| Disk failure or interrupted artifact | No success acknowledgement; previous committed state preserved |
| Owner exit or uncertain callback | Explicit error/delivery state; saved results collectable |
| Reopen a closed draft | New linked session and current caller, old result unchanged |
| Portable package, offline launch | No separately installed Node/npm or browser runtime required on the tested baseline |

Record actual wait duration, version, environment and resource measurements. A short test does not prove indefinite reliability.

## Implementation order

1. Probe return mechanisms in the six harnesses with a minimal fixture, before UI breadth or adapter abstractions. Treat any required setup or unavailable mechanism as a concrete finding.
2. Implement the Go session/result path and tiny helper fixture, including interrupted commits, duplicate submission and independent sessions.
3. Add Electron, the page bridge, equivalent HTML/flat-JSON forms and image attachments. Exercise the real human-facing return path.
4. Verify all six harnesses end to end and package the macOS portable build. Record remaining Windows/Linux work explicitly.

The core and reference renderers are implemented. The [current compatibility report](verification/runs/2026-10-20-no-setup.md) applies the accepted no-setup/manual-fallback boundary and identifies exact tested routes. Broader JSON authoring stays open for the user’s later design discussion.
