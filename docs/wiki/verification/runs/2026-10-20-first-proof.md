---
title: First prototype verification
updated: 2026-10-20
result: Untested
---
# First prototype verification

This report preserves the initial run and its original scope. The [no-setup follow-up](2026-10-20-no-setup.md) records the later user decision, successful Kimi authentication, native image journeys and manual-fallback verification.

The prototype implements the Go owner, Electron bridge, flat forms, HTML authoring, artifacts, recovery and delivery modes. **The accepted six-harness AFK milestone is incomplete.** The aggregate is Untested because required paths remain unavailable or unverified.

## Scope and build

Selected features: [FEAT-sessions](../../features/sessions.md), [FEAT-interfaces](../../features/interfaces.md), [FEAT-delivery](../../features/delivery.md). The run emphasizes result correctness and return paths before UI breadth.

There is no Git revision. [Source/build manifest](../../evidence/2026-10-20-first-proof/manifest.json) records exact content hashes and tested binary identities. The final harness batch used binary SHA-256 `b52f2fd7cb5e28f981d2653495b2998c0a7aa882dfda1a2e8889e051350512ec`. Subsequent edits added form-default validation and adjusted the reference paste/menu UI; they did not change the pipe transport. The final core suite was rerun. Separate native tests exercised the final package's paste target, cancellation and crash recovery.

Environment: macOS arm64; Go 1.26.4; Electron 44.9.0; Node 22.19.0 for development. See [setup](../setup.md) for fixture isolation, account state and actual commands.

## Results

| Scenario | Expected and observed | Result | Evidence |
| --- | --- | --- | --- |
| Core lifecycle and storage | Real CLI/helper tests cover immutable results, repeated collection, concurrent sessions, image drafts, recovery, confinement, duplicate ownership and publication failure | Passed | [Go tests](../../evidence/2026-10-20-first-proof/go-tests.log) |
| Lost waiting CLI | Kill the launching waiter; detached owner still saves the later submitted result | Passed | [Go tests](../../evidence/2026-10-20-first-proof/go-tests.log) |
| Helper exit and stdout drain | Exit without result is failed; a fully emitted submitted result survives a later nonzero exit | Passed | [Go tests](../../evidence/2026-10-20-first-proof/go-tests.log) |
| Uncertain notification | Saved result survives; an ambiguous failed acknowledgement is recorded and not resent | Passed | [Go tests](../../evidence/2026-10-20-first-proof/go-tests.log) |
| JSON form with selected image | Native file picker → saved image draft → close → new linked session → restored preview → submission; returned file hash matches fixture | Passed | [Saved UI results](../../evidence/2026-10-20-first-proof/ui-results.json), [fixture image](../../evidence/2026-10-20-first-proof/reference.png) |
| Custom HTML controls and storage retry | Local slider/preview and notes work; non-writable session directory rejects submit without losing input; restoring access permits submission | Passed | [Saved UI results](../../evidence/2026-10-20-first-proof/ui-results.json) |
| Native image paste | Plain text rejected; pixels copied from Preview produce preview, saved draft and submitted 401-byte PNG with matching stored hash | Passed | [Saved UI results](../../evidence/2026-10-20-first-proof/ui-results.json), [pasted image](../../evidence/2026-10-20-first-proof/pasted-image.png) |
| Native Cancel button | Explicit Cancel closes the window and saves `cancelled` without answer data | Passed | [Saved UI results](../../evidence/2026-10-20-first-proof/ui-results.json) |
| Renderer crash and reopen | Kill only the fixture's page renderer: `failed`/`renderer_crashed`; reopen restores acknowledged draft in a new linked window; closing it saves `closed` | Passed | [Saved UI results](../../evidence/2026-10-20-first-proof/ui-results.json) |
| Codex pending command | Actual Fungi helper response returned in the same disposable conversation after a ten-second silent wait | Passed | [Harness probes](../../evidence/2026-10-20-first-proof/harness-probes.json) |
| OpenCode pending command | Canonical working directory; one Bash invocation; actual result consumed after a ten-second wait | Passed | [Harness probes](../../evidence/2026-10-20-first-proof/harness-probes.json) |
| Pi pending command | One Bash call, timeout null; actual response consumed after 180 seconds | Passed | [Harness probes](../../evidence/2026-10-20-first-proof/harness-probes.json) |
| Claude background completion | Background launch → WAITING/end turn → completion event → one output-file read → actual result after 180 seconds | Passed | [Harness probes](../../evidence/2026-10-20-first-proof/harness-probes.json) |
| Grok background completion | Background launch → WAITING → completion → one completed-output retrieval → actual result after 180 seconds | Passed | [Harness probes](../../evidence/2026-10-20-first-proof/harness-probes.json) |
| Codex exact-conversation callback consumption | Native-helper notice queued during the active turn reached the originating conversation after that turn ended; one collect returned submitted data and its artifact | Passed | [Callback observation](../../evidence/2026-10-20-first-proof/codex-callback.json) |
| Corrected self-contained Codex notice | Native-form notice reached the originating conversation with full executable, arguments and state-directory environment; one collect using those values returned `closed` | Passed | [Corrected callback observation](../../evidence/2026-10-20-first-proof/codex-callback-corrected.json) |
| Kimi result delivery | Installed CLI rejected the test before tool execution: selected provider has no credential | Untested | [Harness probes](../../evidence/2026-10-20-first-proof/harness-probes.json) |
| Arbitrary AFK duration and all-six native UI loop | Three-minute synthetic waits do not establish unlimited waits, sleep/wake reliability, or a native human interaction in every harness | Untested | [Delivery limitations](../../features/delivery.md#coverage-and-limitations) |
| Portable local launch | Corrected directory's internal symlinks remain within the package; HTML UI launches with PATH restricted to `/usr/bin:/bin`, so no external Node/npm is used | Passed | [Package metadata](../../evidence/2026-10-20-first-proof/manifest.json), [UI result](../../evidence/2026-10-20-first-proof/ui-results.json) |
| Downloaded clean-machine release | No notarized distribution, quarantine/Gatekeeper or fully disconnected clean-machine test | Untested | [Setup limits](../setup.md#observation-and-cleanup) |

The three-minute batches took approximately 190 seconds (Pi), 193 seconds (Claude) and 194 seconds (Grok), including model startup/final response time. Claude supplied a 600,000 ms tool timeout. That successful bounded test must not be described as an unlimited AFK guarantee. No repeated model status calls occurred during these waits; Claude and Grok each retrieved output once after completion.

The portable folder contains approximately 315 MiB before the final documentation refresh. One idle `ps` sample of the open demo found 5.8 MiB RSS for its Go owner and 262.6 MiB across four Electron processes, each showing 0% sampled CPU. This is one observation, not a peak-memory or cross-machine benchmark; the demo could also receive user edits. [Resource sample](../../evidence/2026-10-20-first-proof/resources.json)

## Earlier failures and corrections

1. Pi's default `gpt-5.5` was rejected by the existing provider. A supported model was selected only for the test invocation.
2. Kimi failed for missing credentials. No attempt was made to obtain or change credentials.
3. Codex read-only mode could not create session storage. Workspace-write exposed the initial loopback-listener restriction. Removing HTTP in favor of pipes and OS file locks made the subsequent Codex probe pass.
4. OpenCode rejected a non-canonical `/var/...` working directory as external. Using the canonical `/private/var/...` path resolved the test's permission mismatch; permissions were not broadened.
5. Electron's first asset loader called a network loader that its own request filter blocked. A confined direct file handler corrected it; native forms and HTML then rendered.
6. The first package copy altered framework symlinks and failed to find ICU data. Preserving literal symlinks corrected the next build. The failed build is not a distributable.
7. Early image-paste attempts did not establish a working clipboard/focus fixture. An explicit editable target was added; a separate one-image form kept it visible and passed native paste, draft save and submission. Other demo windows received user edits, which were preserved. The default development reload menu was removed; explicit close/reopen is the supported restoration path. A final one-line adjustment clears a previous paste-error banner when a valid image is accepted; it does not change artifact transport.
8. The first Codex callback notice omitted a non-default state directory. That notice was subsequently consumed successfully using the launch context. New notices include the executable, argument array and state-directory environment; a subsequent native-form notice was consumed using those supplied values and returned the saved `closed` outcome.

## Limits and next decision

The core is usable for further authoring experiments, and the macOS folder is a development package. The milestone is not complete. Resolve whether optional harness integrations may be loaded when ordinary shell tools have bounded lifetimes, and have the user authenticate Kimi. Then finish exact-caller/idle-agent delivery, longer AFK tests and the combined UI journey in all six harnesses.

Testing was performed by the implementing agent; there was no independent reviewer. Wiki structure validation checks documentation integrity only.
