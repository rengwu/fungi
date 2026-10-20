---
title: Skill authoring and longer-idle verification
updated: 2026-10-20
result: Passed
---
# Skill authoring and longer-idle verification

The simpler committed grill-form now runs as a Fungi skill without an engine API change. Its exercised native UI scenarios passed, the Linux arm64 core suite passed, and a real interview returned to its originating conversation. One user-controlled macOS sleep/wake cycle also passed, including retained notes/image and automatic return. The two-hour helper completed and its exact-origin callback was collected once. The aggregate is **Passed** for these selected scenarios; additional desktop and arbitrary-duration claims remain outside this proof.

## Scope and build

Selected features: [FEAT-interfaces](../../features/interfaces.md), [FEAT-sessions](../../features/sessions.md), [FEAT-delivery](../../features/delivery.md). This implements the [authorized next step](../../planning.md#skill-authoring-and-longer-waits), using historical jstack assets rather than its ongoing uncommitted rework.

The baseline is jstack `7a7c919`; `git diff` confirms its relevant assets/model match the repository tip at the end of 2026-10-19, `32324e47aae7d3ad3f92e6a4aad8966e9451d691`. No jstack files were changed. The Fungi core and Electron bridge are unchanged from the earlier compatibility proof. New code is skill-owned UI, a development idle helper, focused question-model tests and inclusion of skills in packaging.

Native UI environment: macOS arm64, Go 1.26.4, Electron 44.9.0, Node 22.19.0 for development only. The tested package is `dist/fungi-darwin-arm64-2026-10-20T09-41-25-905Z`. It launched the skill with `PATH=/usr/bin:/bin`, without a separately available Node or Go executable. [Build and source manifest](../../evidence/2026-10-20-skill-proof/manifest.json)

Linux tests ran through installed OrbStack in the official `golang:1.26.4` image, pinned by digest. The test container had networking disabled, UID/GID 65534, a read-only copy of the required source and no host credentials or account configuration. It was removed after completion. This exercises Linux arm64 processes and filesystems, not a desktop window or public distribution. [Exact command](../../evidence/2026-10-20-skill-proof/linux-command.json)

## Results

| Scenario | Expected and observed | Result | Evidence |
| --- | --- | --- | --- |
| Historical layout in the native package | Markdown table, flowchart, recommendation, ordinary radio choices and plain textarea rendered; manual completion guidance visible | Passed | [Native observations](../../evidence/2026-10-20-skill-proof/native-ui.json) |
| Answer validation | Other with empty notes kept submission disabled; nonempty notes enabled it. Model checks also cover notes-only answers, invalid IDs and changed-question recovery rejection | Passed | [Native observations](../../evidence/2026-10-20-skill-proof/native-ui.json), [model tests](../../evidence/2026-10-20-skill-proof/model-tests.log) |
| Image and close/reopen | Native file picker selected the fixture; acknowledged draft saved notes, choice and image. Close returned `closed`; new linked session visibly restored notes and preview | Passed | [Session evidence](../../evidence/2026-10-20-skill-proof/native-ui.json) |
| Abrupt owner termination | Killed only the disposable session owner with SIGKILL; collection returned `session_interrupted`. Reopening created another linked session with the acknowledged notes and image | Passed | [Session evidence](../../evidence/2026-10-20-skill-proof/native-ui.json) |
| Storage failure and retry | Made only the test session directory non-writable; native Submit showed failure and retained inputs, with no result published. Restored permissions; retry returned submitted answers plus image with the expected SHA-256 | Passed | [Session evidence](../../evidence/2026-10-20-skill-proof/native-ui.json) |
| Unsupported diagram and literal HTML | Invalid Mermaid showed source/error while controls remained usable; raw script markup displayed as text. Cancel returned `cancelled` with no answers/artifacts | Passed | [Native observations](../../evidence/2026-10-20-skill-proof/native-ui.json) |
| Skill and JavaScript checks | Skill metadata valid; four focused model tests and JavaScript syntax checks passed | Passed | [Model tests](../../evidence/2026-10-20-skill-proof/model-tests.log), [manifest](../../evidence/2026-10-20-skill-proof/manifest.json) |
| Real skill interview and automatic return | User answered questions 18–19 in the new skill; exact-origin notice arrived after the previous agent turn; one collection returned the actual selected choices and notes | Passed | [Collected interview](../../evidence/2026-10-20-skill-proof/environment-result.json) |
| Linux arm64 core execution | Race-enabled suite passed, including concurrent isolation, lost waiter, manual receipt, interrupted recovery, stdout drain, locks, artifact confinement and uncertain callback behavior | Passed | [Linux run](../../evidence/2026-10-20-skill-proof/linux-core.log) |
| Idle helper smoke check | A single held call returned a submitted timing artifact after at least five seconds; the draft was acknowledged before readiness | Passed | [Short result](../../evidence/2026-10-20-skill-proof/idle-short.json) |
| Two-hour idle and exact-origin return | One local timer with no polling or heartbeat returned after 7,200.0148 helper-measured seconds. One Codex notice reached the originating conversation and one collection returned the result; timing artifact matched data and SHA-256, owner/helper exited | Passed | [Launch receipt and draft](../../evidence/2026-10-20-skill-proof/idle-long.json), [collected result and validation](../../evidence/2026-10-20-skill-proof/idle-long-result.json) |
| Manual sleep/wake with native UI | macOS recorded clamshell sleep at 18:09:06 +0800 and final lid/HID wake at 18:13:38, with an intervening maintenance dark wake/sleep. The same session returned its original note, the user's `AFTER-WAKE` edit and unchanged image; exact-origin callback collected once | Passed | [Pre-sleep baseline](../../evidence/2026-10-20-skill-proof/sleep-wake-baseline.json), [collected result and validation](../../evidence/2026-10-20-skill-proof/sleep-wake-result.json), [power events](../../evidence/2026-10-20-skill-proof/sleep-wake-power.log) |

The image fixture is 180 bytes and hashes to `1bd7e798fbfd3bd20de881435c7bfa7811db8466e73714a8ab7c1af6b7460e9a`. Native UI actions used the connected desktop driver; no direct result writes or page-handler invocations substituted for clicks, typing or file selection. Output files and artifact hashes were inspected independently.

The long probe began at `2026-10-20T09:41:51.107501Z` and finished at `2026-10-20T11:46:08.980722Z`. Its UTC interval is 7,457.873221 seconds (2 h 4 min 17.873 s); its separately reported `time.Since(start)` measurement is 7,200.01478275 seconds. The verified manual sleep/wake cycle occurred inside that interval. Both measurements are retained rather than describing the wall-clock interval as exactly two hours. The 271-byte timing artifact hashes to `2a8c5a13db37e81205254dcbd69a067fdebb27d963e8471ca31e34bd91dc2138`, and its parsed JSON equals the collected result data.

## Tooling observations

The native app selector also started a default Electron launcher. It was distinguished from the test session by its exact executable and absence of renderer/session arguments, then only that launcher was stopped. Existing user-edited windows and drafts were preserved.

Python's HTTPS client lacked a usable certificate chain for upstream research. Normal `curl` certificate validation succeeded; certificate checks were not disabled. The pinned Pele npm tarball passed its published SHA-512 check, and a separate three-diagram probe passed. Those are research observations, not a change to the questionnaire renderer. [Pele assessment](../../research/pele.md)

## Limits and continuation

The two-hour test exercised the native helper, session owner and Codex callback. It does not establish two hours of Electron UI interaction or arbitrary-duration reliability across every harness. Its saved response and originating-conversation consumption are now both observed. No model polling or additional waiter was used.

Sleep/wake evidence covers one short user-controlled cycle on this Mac, not an overnight wait, reboot or guarantee across harnesses. A clean-Mac downloaded launch, Windows execution and Linux desktop UI remain untested. The [user's answers](../../journal/2026-10-20-grill-desktop-arrangements.md) defer additional desktop checks until last. No public signing/notarization, new OS installation or permanent harness integration was performed. OrbStack was started to run the Linux check; unrelated containers were not changed.

Verification was performed by the implementing agent. There was no independent review. The historical user-edited demo remains outside these fixtures. All callbacks in this verification round have been consumed; no probe remains awaiting collection.
