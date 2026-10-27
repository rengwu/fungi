---
title: HTML and flat JSON interfaces
updated: 2026-10-27
---
# HTML and flat JSON interfaces

**ID: FEAT-interfaces**

Skills can bring a local HTML directory or a small JSON form. The shared page bridge returns to [the session engine](sessions.md). This is the [accepted primitive UI scope](../planning.md#accepted-first-milestone-scope), not the future component architecture.

## Entry points and expected behavior

`fungi open examples/review.html` shows a local animation preview, slider, notes and an optional image file. `fungi open examples/feedback.json` renders flat text, textarea, range and image fields. `fungi describe html` and `fungi describe form` expose the actual API and examples.

`open` accepts optional `--width` and `--height` for HTML/JSON windows, in
logical pixels (defaults 820 × 780, minimum 420 per dimension). `reopen`
inherits the original launch size unless overridden; user resizing is not
saved. The OS may adjust oversized windows. Native helpers own their windows
and reject these Fungi flags. See `fungi help cli` for validation rules.

Generic reply, comparison and form examples live in `examples/` and ship with
the portable package. Workflow skills and jstack integration live in the
separate `jstack-fungi-test` project; they use the same public page API.

The bridge is `ready`, `saveDraft`, `submit` and `cancel`. Draft saving optionally takes attachments, so image bytes survive explicit reopen. HTML authors need no build system or registration manifest. The JSON schema remains experimental.

Protocol 1 preserves the documented authoring API; optional metadata may be
added. Renderer internals and Wails native APIs are outside that compatibility
promise. `fungi help html` is the maintained contract. Bridge tests cover repeated
readiness, additive metadata, empty/binary attachments, JSON values including
false/zero/null, call ordering and retry after rejection.

Mutating calls preserve invocation order before Blob conversion begins. A failed save rejects that call without preventing later submission. Read `fungi help html` for complete replacement semantics and limits, and `fungi help delivery` for completion and recovery. `examples/compare/index.html` demonstrates a custom comparison canvas and PNG result through the same API.

The local renderer supports ES modules, custom elements, Canvas/SVG/WebGL, local and Blob workers, WebAssembly and media subject to the system engine. The handler supports byte ranges, but
Linux WebKitGTK cannot load direct Wails-scheme media URLs; this accepted limit
is still recorded as a failed media observation. Skills provide their libraries; no component catalog or framework is installed in the engine. The asset handler confines real paths to the source root, rejects directories/unsupported methods and escapes, and streams bytes instead of buffering entire media files.

The Wails helper treats skill pages as trusted local application code and
exposes native Wails APIs. The asset handler confines local files to the source
root and applies a local-content policy. The main page owns the Fungi bridge;
skills should sandbox untrusted inline previews. File pickers and an editable
paste target provide reference image input.

Browser storage isolation is best effort. Fungi's per-session drafts/results
stay separate; macOS/Linux `localStorage` can be shared. Prefer `saveDraft` for
durable state and `sessionStorage` for window-local state. Native APIs, manual
crash recovery and Linux direct-media limitations are
[accepted requirements](../planning.md#accepted-wails-limits-and-next-milestone-on-2026-10-25).
No framework fork or automatic recovery layer is required.

Submission closes the window after acknowledgement. Persistence failure retains the input and permits retry. Manual delivery is stated in the form. Window close saves a `closed` result; it preserves acknowledged drafts. A custom page must opt into saving its state.

## Verification recipe

Follow [setup](../verification/setup.md). Launch a new isolated form:

1. Enter text, move the slider and select the committed `docs/wiki/evidence/2026-10-20-first-proof/reference.png`.
2. Observe “Draft saved” and inspect `draft.json` plus the saved artifact hash.
3. Close the native window. Collect once: expect `closed`, without answer data.
4. Reopen that exact session. Expect a new session, restored values and the image preview.
5. Submit. Compare the returned artifact's SHA-256 to the fixture; collect again and compare envelopes.

For the HTML failure path, save notes, then make only that disposable session directory non-writable. Submit: expect an error and retained input, with no `result.json`. Restore its permissions and retry; expect submission. Never change permissions outside the test-created directory.

For paste, copy image pixels into the OS clipboard, focus the visible paste target and paste. Expect a preview, a saved draft and a returned image artifact. Clipboard re-encoding means it need not match the original file's hash. The [first Electron proof](../verification/runs/2026-10-20-first-proof.md) records the original paste observation; it does not qualify the Wails clipboard path on every desktop.

For full-helper recovery, save a distinctive draft in a disposable session,
terminate only that session's `fungi-renderer` helper, and collect its failed
result (`helper_exited`). Reopen the exact session: expect a new linked session
with the saved value. Closing that restored window must leave the original
failed result unchanged. A system web-content crash can leave the helper and
a blank window alive; close that window before collection and recovery.
Fungi does not automatically detect or restart that web-content process.

## Coverage and limitations

The [platform record](../platforms.md) separates current Wails expectations
from historical Electron evidence. Runtime setup is explicit and optional. A renderer that exits before readiness produces a saved
`renderer_startup_failed` result with its diagnostic log location. This does
not infer a cause or retry with a different security policy.

The [Wails integration run](../verification/runs/2026-10-25-wails-integration.md)
records native input, draft restoration, submission and artifact return.
The [Wails follow-up](../verification/runs/2026-10-25-wails-followup.md) records
full-helper exit and explicit recovery. The older `renderer_crashed` result
and isolated-page crash recipe belong to the Electron proof, not the current
Wails host. A drawing editor, screenshot capture, nested JSON components and
external live updates remain deferred. Feature-detect GPU and codec
capabilities; hosted rendering results vary by platform.

## Verification history

- [2026-10-27 repository cleanup](../verification/runs/2026-10-27-repo-cleanup.md) — removed workflow-specific examples/tests, corrected historical crash guidance, and verified the explicit-path smoke driver against a fresh Wails package.

- [2026-10-27 starting window size](../verification/runs/2026-10-27-window-size.md) — native Mac defaults, custom and minimum sizes, single-axis flags, reopen inheritance/override, and JSON form submission.

- [2026-10-27 contract and AFK checks](../verification/runs/2026-10-27-afk-contract.md) — current Wails bridge checks, older skill compatibility and completed two-hour idle/sleep-wake tests.

- [2026-10-26 offline discovery](../verification/runs/2026-10-26-discovery.md) — installed help, examples, update/rollback and native throwaway reply.

- [2026-10-25 Wails integration](../verification/runs/2026-10-25-wails-integration.md) — normal packaged launches, native input and recovery, artifact return, and updates with an active UI.

- [2026-10-25 Wails follow-up](../verification/runs/2026-10-25-wails-followup.md) — accepted simpler requirements, confirmed Linux custom-scheme media loader failure, and local full-host-exit draft recovery.

- [2026-10-25 Wails assessment](../verification/runs/2026-10-25-wails.md) — historical replacement-equivalence proof. Later follow-up and integration runs resolved its distribution work and recorded the accepted native-access, storage and Linux-media limits.

- [2026-10-22/23 Windows desktop handoff](../verification/runs/windows-2026-10-22-203907.md) — agent-operated native input, visible draft restoration and comparison PNG/hash passed with an empty child PATH; the full renderer fixture remained unexecuted because symlink creation lacked privilege. Aggregate Untested.

- [2026-10-22 verified baselines](../verification/runs/2026-10-22-baselines.md) — unchanged authoring API, explicit WebGL availability and diagnosable unsupported startup.

- [2026-10-22 hosted checks](../verification/runs/2026-10-22-hosted.md) — real Electron execution, Windows pipe repair and Linux sandbox/graphics limits.

- [2026-10-22 platform verification](../verification/runs/2026-10-22-platform.md) — call ordering, rendering primitives, saved/exported images, rich-form retry and a separate comparison example.

- [2026-10-20 rich jstack demo](../verification/runs/2026-10-20-rich-grill.md) — inline previews, native draft checks and a completed user questionnaire with region feedback and exact-origin return.

- [2026-10-20 skill proof](../verification/runs/2026-10-20-skill-proof.md) — historical grill-form port, native draft/image recovery, owner loss, retry and cancellation.
- [2026-10-20 first proof](../verification/runs/2026-10-20-first-proof.md) — includes asset-loader and clipboard-fixture corrections, native artifact submission, and crash recovery.
- [2026-10-20 no-setup compatibility](../verification/runs/2026-10-20-no-setup.md) — native form/image submission and originating-harness consumption.
