---
title: HTML and flat JSON interfaces
updated: 2026-10-20
---
# HTML and flat JSON interfaces

**ID: FEAT-interfaces**

Skills can bring a local HTML directory or a small JSON form. The shared page bridge returns to [the session engine](sessions.md). This is the [accepted primitive UI scope](../planning.md#accepted-first-milestone-scope), not the future component architecture.

## Entry points and expected behavior

`fungi open examples/review.html` shows a local animation preview, slider, notes and an optional image file. `fungi open examples/feedback.json` renders flat text, textarea, range and image fields. `fungi describe html` and `fungi describe form` expose the actual API and examples.

The portable package also includes `skills/grill-form/`. Run `fungi open skills/grill-form/assets/form.html --input skills/grill-form/example.json` for a complete skill-owned questionnaire. It accepts the older jstack question schema, renders Markdown and Mermaid offline, saves plain notes and one optional image, and returns readable question/choice labels in `data.answers`. This does not extend the engine's built-in JSON form schema. [Migration decision](../planning.md#skill-authoring-and-longer-waits)

The bridge is `ready`, `saveDraft`, `submit` and `cancel`. Draft saving optionally takes attachments, so image bytes survive explicit reopen. HTML authors need no build system or registration manifest. The JSON schema remains experimental.

The Electron page has context isolation and sandboxing. Local assets are confined beneath the source directory; remote requests, navigation and OS permission requests are denied. Native helpers remain trusted executables with ordinary OS access. Reference image fields have a file picker and an explicit editable paste target.

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

File selection, image paste, image draft recovery, submission, cancellation, local HTML controls and actual write-failure recovery were observed through the native UI. Killing the isolated page renderer produced `failed` with `renderer_crashed`; explicit reopen restored `DRAFT-CRASH-RECOVER`, and closing that new window produced `closed`. The app uses a minimal native menu without development reload. No screenshot, canvas, nested JSON components or external live-update mechanism is implemented.

## Verification history

- [2026-10-20 skill proof](../verification/runs/2026-10-20-skill-proof.md) — historical grill-form port, native draft/image recovery, owner loss, retry and cancellation.
- [2026-10-20 first proof](../verification/runs/2026-10-20-first-proof.md) — includes asset-loader and clipboard-fixture corrections, native artifact submission, and crash recovery.
- [2026-10-20 no-setup compatibility](../verification/runs/2026-10-20-no-setup.md) — native form/image submission and originating-harness consumption.
