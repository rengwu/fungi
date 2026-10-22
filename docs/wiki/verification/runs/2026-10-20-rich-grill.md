---
title: Updated jstack questionnaire in Fungi
updated: 2026-10-20
result: Passed
---
# Updated jstack questionnaire in Fungi

The user requested the newly updated jstack presentation in a Fungi mock
questionnaire. [FEAT-interfaces](../../features/interfaces.md) now permits inline
`srcdoc` previews while retaining the main-frame bridge and blocking frame
document requests. The example adapter uses jstack's existing generator,
notes editor, Pele diagrams, visual choices and region-comment UI; only draft
storage and submission delivery change. No jstack source was edited.

The selected native presentation, draft and callback scenarios passed. The
user's actual submission reached the originating conversation and was collected
once, then validated against the pinned jstack model. This does not establish
inline-image transport or full feature parity.

Build: local work based on `eca3d85`, with the renderer and example adapter
changes recorded in the [manifest](../../evidence/2026-10-20-rich-grill/manifest.json).
The native macOS arm64 package is
`dist/fungi-darwin-arm64-2026-10-20T12-22-27-273Z`, Electron 44.9.0.
Electron's [frame navigation events](https://www.electronjs.org/docs/latest/api/web-contents/#event-will-frame-navigate)
and [request resource types](https://www.electronjs.org/docs/latest/api/web-request/)
inform the narrow frame-loading change; runtime evidence below is separate.

| Scenario | Expected and observed | Result | Evidence |
| --- | --- | --- | --- |
| Inline preview boundary | A disposable runtime fixture rendered a sandboxed inline frame; parent DOM, Node, Fungi bridge and external fetch were unavailable in it. A file-based frame did not execute | Passed | [Frame result](../../evidence/2026-10-20-rich-grill/frame-check.json) |
| Native presentation | The updated form displayed visual cards, a working duration slider, an open review preview, a Pele flowchart, a table and notes | Passed | [Native observations](../../evidence/2026-10-20-rich-grill/native.json) |
| Draft and close/reopen | Native controls saved choice/text, duration 560/ease-out and a keyboard-selected region comment. Close committed `closed`; reopening visibly restored note/comment. Later user edits were preserved | Passed | [Native observations and acknowledged draft](../../evidence/2026-10-20-rich-grill/native.json) |
| User submission and exact-origin return | One callback reached the originating conversation after the prior turn ended. One collection returned six valid responses, including a region comment and its coordinates. No images or changed prototype states were submitted | Passed | [Receipt](../../evidence/2026-10-20-rich-grill/user-receipt.json), [collected and validated response](../../evidence/2026-10-20-rich-grill/user-result.json) |

The first close shortcut did not close the disposable window; one premature
manual collection returned `pending`. No polling followed. Clicking its native
close control then produced a persisted `closed` result, inspected from disk.
This was a test-driver issue, not a successful close observation. A later
desktop-driver interruption coincided with changed choices/slider values in
the reopened test window; that draft and window were preserved rather than
submitting on the user's behalf.

The actual callback-enabled questionnaire is session
`2b009f78acc7f857be52fc58d18a7643`, titled “Fungi — your questionnaire.”
No synthetic answers were entered into it. Its receipt contains the exact
collection argv and state directory. The callback was consumed once and
`delivery.json` independently reports `accepted`. The model validation matched
the original jstack session and specification digest. Jstack's collector then
accepted the reconstructed submission and its finish step saved a temporary
transcript and reopenable review. No jstack waiter was started. Returned
content was treated as data; the mock choices are not project decisions.

The response contains “Underwater library,” “Something else” with `sadasd`,
an animation note `test`, a region comment `test` at (76, 36) with size 280×290
within a 755×400 viewport, and “Clear enough.” Its final note raises the real
concern that authoring/handoff should be inexpensive and easy for models to
understand, and describes eventual bundling of jstack with Fungi. That
[direction](../../planning.md#rich-questionnaire-feedback) is separate from
the mock choices.

This is a macOS experiment. Inline-image round trip, submission failure/retry,
full jstack feature parity and other platforms are not established by the
observations. No changed animation settings were returned in the user
submission; the earlier disposable draft test remains the evidence for that
control. The implementing agent performed these checks;
there was no independent review.
