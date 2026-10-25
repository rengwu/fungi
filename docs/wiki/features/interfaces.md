---
title: HTML and flat JSON interfaces
updated: 2026-10-25
---
# HTML and flat JSON interfaces

**ID: FEAT-interfaces**

Skills can bring a local HTML directory or a small JSON form. The shared page bridge returns to [the session engine](sessions.md). This is the [accepted primitive UI scope](../planning.md#accepted-first-milestone-scope), not the future component architecture.

## Entry points and expected behavior

`fungi open examples/review.html` shows a local animation preview, slider, notes and an optional image file. `fungi open examples/feedback.json` renders flat text, textarea, range and image fields. `fungi describe html` and `fungi describe form` expose the actual API and examples.

The portable package also includes `skills/grill-form/`. Run `fungi open skills/grill-form/assets/form.html --input skills/grill-form/example.json` for a complete skill-owned questionnaire. It accepts the older jstack question schema, renders Markdown and Mermaid offline, saves plain notes and one optional image, and returns readable question/choice labels in `data.answers`. This does not extend the engine's built-in JSON form schema. [Migration decision](../planning.md#skill-authoring-and-longer-waits)

The bridge is `ready`, `saveDraft`, `submit` and `cancel`. Draft saving optionally takes attachments, so image bytes survive explicit reopen. HTML authors need no build system or registration manifest. The JSON schema remains experimental.

Mutating calls preserve invocation order before Blob conversion begins. A failed save rejects that call without preventing later submission. The [authoring guide](../authoring.md) describes the complete replacement semantics, delivery procedure and limits. `examples/compare/index.html` demonstrates a custom comparison canvas and PNG result through the same API.

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

1. Enter text, move the slider and select the generated `verification-output/reference.png`.
2. Observe “Draft saved” and inspect `draft.json` plus the saved artifact hash.
3. Close the native window. Collect once: expect `closed`, without answer data.
4. Reopen that exact session. Expect a new session, restored values and the image preview.
5. Submit. Compare the returned artifact's SHA-256 to the fixture; collect again and compare envelopes.

For the HTML failure path, save notes, then make only that disposable session directory non-writable. Submit: expect an error and retained input, with no `result.json`. Restore its permissions and retry; expect submission. Never change permissions outside the test-created directory.

For paste, copy image pixels into the OS clipboard, focus the visible paste target and paste. Expect a preview, a saved draft and a returned image artifact. This passed using a separate, one-image form that kept the paste target within the viewport. Plain text produced a useful error; pasting pixels from Preview produced a 401-byte PNG, a saved draft and a submitted artifact whose recorded hash matched its bytes. Clipboard re-encoding means it need not match the original file's hash.

For renderer recovery, save a distinctive draft in a disposable session, terminate only that session's page-renderer process, and collect its failed result. Reopen the exact session: expect a new linked session with the saved value. Closing that restored window must leave the original failed result unchanged.

## Coverage and limitations

The [platform record](../platforms.md) separates current Wails expectations
from historical Electron evidence. Runtime setup is explicit and optional. A renderer that exits before readiness produces a saved
`renderer_startup_failed` result with its diagnostic log location. This does
not infer a cause or retry with a different security policy.

File selection, image paste, image draft recovery, submission, cancellation, local HTML controls and actual write-failure recovery were observed through the native UI. Killing the isolated page renderer produced `failed` with `renderer_crashed`; explicit reopen restored `DRAFT-CRASH-RECOVER`, and closing that new window produced `closed`. The app uses a minimal native menu without development reload. The comparison example exports a custom canvas; a drawing editor, screenshot capture, nested JSON components and external live updates remain deferred. Feature-detect GPU and codec capabilities; hosted rendering results vary by platform.

## Verification history

- [2026-10-25 Wails follow-up](../verification/runs/2026-10-25-wails-followup.md) — accepted simpler requirements, confirmed Linux custom-scheme media loader failure, and local full-host-exit draft recovery.

- [2026-10-25 Wails assessment](../verification/runs/2026-10-25-wails.md) — historical replacement-equivalence proof. Later user decisions accept native access and best-effort browser storage; Linux media and installation integration still require follow-up.

- [2026-10-22/23 Windows desktop handoff](../verification/runs/windows-2026-10-22-203907.md) — agent-operated native input, visible draft restoration and comparison PNG/hash passed with an empty child PATH; the full renderer fixture remained unexecuted because symlink creation lacked privilege. Aggregate Untested.

- [2026-10-22 verified baselines](../verification/runs/2026-10-22-baselines.md) — unchanged authoring API, explicit WebGL availability and diagnosable unsupported startup.

- [2026-10-22 hosted checks](../verification/runs/2026-10-22-hosted.md) — real Electron execution, Windows pipe repair and Linux sandbox/graphics limits.

- [2026-10-22 platform verification](../verification/runs/2026-10-22-platform.md) — call ordering, rendering primitives, saved/exported images, rich-form retry and a separate comparison example.

- [2026-10-20 rich jstack demo](../verification/runs/2026-10-20-rich-grill.md) — inline previews, native draft checks and a completed user questionnaire with region feedback and exact-origin return.

- [2026-10-20 skill proof](../verification/runs/2026-10-20-skill-proof.md) — historical grill-form port, native draft/image recovery, owner loss, retry and cancellation.
- [2026-10-20 first proof](../verification/runs/2026-10-20-first-proof.md) — includes asset-loader and clipboard-fixture corrections, native artifact submission, and crash recovery.
- [2026-10-20 no-setup compatibility](../verification/runs/2026-10-20-no-setup.md) — native form/image submission and originating-harness consumption.
