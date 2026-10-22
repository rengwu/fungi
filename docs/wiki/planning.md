---
title: Requirements and discussion
updated: 2026-10-23
---
# Requirements and discussion

## Requirements supplied by the user

- Provide a reliable, flexible, extensible engine for skill-authored human interaction workflows. The engine is reusable across agents and harnesses.
- Present information, collect responses and artifacts, and deliver them back to the exact calling agent conversation.
- Provide one core CLI executable with no installer or separately installed language runtime. The accepted portable distribution may include renderer runtime files beside it. Skills may supply assets and helper executables, and fungi may load optional capabilities from a nearby folder.
- Provide an AI-friendly CLI with discoverable capabilities and APIs; MCP is optional.
- Support desktop platforms rather than depending on macOS alone.
- Eventually enable highly composable interfaces: Markdown, diagrams, charts, forms, native capture, drawing, custom web content, and live controls. Skills own the workflow; fungi supplies execution and communication infrastructure.
- Support a dedicated native window with embedded web content as a possible presentation mechanism, avoiding dependence on an agent controlling the user's browser.
- Human interaction must not cause repeated model tool calls, status polling, continued token expenditure, or compactions while the user is away. Automatic response delivery is central.
- Use jstack's wiki, plan, grill-form, and rule-lean workflows for this project. Implementation was authorized on 2026-10-20 by the user’s `work-until-grill` instruction.

## Scope clarifications on 2026-10-19

The user first clarified that the default core should include no UI capabilities: prioritize a reliable and extensible base on which primitive and complex interfaces can be built.

The user then allowed primitive UI functionality and components for testing. This permits a small reference UI to validate the engine. It does not establish a requirement to build a broad widget library into the core.

The user requested unusually deep research, including niche implementations, widely available component libraries, frameworks, and rendering engines.

The user subsequently emphasized ease of extensibility and simplicity, explicitly invoking rule-lean. Skill authors must not have to learn a complicated plugin framework to add UI capabilities. Research breadth must not turn into implementation scope.

## Accepted decisions from the first interview

The [submitted answers](journal/2026-10-19-grill-eb02aa23-b51f-42c9-b9c0-8fa331bfb02c.md) establish these decisions:

| Decision | Consequence |
| --- | --- |
| One core executable orchestrates interaction; skills may bring assets/helpers, and fungi may have adjacent optional capabilities | Extensibility must not require rebuilding the core for every new workflow |
| Bundle the simplest, lightest suitable rendering engine | Do not require users to install WebView2 or WebKitGTK; the prior system-webview recommendation is superseded |
| HTML/assets and JSON component descriptions are equally supported in the first release | The JSON authoring path is first-release scope; it should share a renderer and result contract with HTML |
| Automatic delivery is best effort across popular harnesses; open the interface even when it is unavailable | Explicitly report the fallback and let the user say “done”; never substitute agent polling |
| Local agent and UI on the same desktop first | Remote/container-to-desktop transport is deferred |
| Explicit reopen/collect is acceptable after relevant processes exit | Preserve responses; no mandatory login service or automatic reboot recovery in the first release |

The user specifically prioritizes running cleanly out of the box over the smaller download offered by system-webview wrappers. Native operating-system APIs remain necessary; the exact supported desktop baseline still needs verification.

During follow-up research, the user explicitly accepted bundling Chromium/Electron as a worst-case browser-renderer fallback. That initially permitted a larger engine without selecting one. The later Electron + Go decision below settles the default; it does not require maintaining two engines in the first release.

## Second interview and subsequent clarifications

The [second submitted form](journal/2026-10-19-grill-186963ab-2f7b-423c-b40b-97b0afacc6fb.md) selects an open-source renderer for the default distribution. Proprietary renderers may remain optional third-party capabilities, but Sciter and Ultralight no longer qualify for the default.

For packaging, the saved form selects automatic runtime extraction. In chat, while referring specifically to that question, the user instead delegated the decision to the simplest-complexity route. Applying that instruction, the chosen plan is a portable folder/archive containing the fungi CLI and renderer files. This retains one CLI entry point and no installer, while avoiding extraction/cache-management code. Preserve the submitted answer as history; the chat clarification governs the current plan. It relaxes the original literal single-file distribution target.

The user also requires a substantial established following (including GitHub stars or equivalent evidence) and healthy recent renderer development. No numerical star cutoff was specified. Evaluate adoption, actual releases/fixes, contributor breadth, maintenance policy, and any binding used by fungi separately; a popular upstream engine does not establish the health of its wrapper. [Adoption and maintenance research](research/renderer-health.md)

## Accepted Electron + Go direction

On 2026-10-19, after the adoption and maintenance assessment, the user accepted Electron + Go. Go is the core CLI/session engine; a small Electron companion owns the native window and bundled web renderer. This replaces the earlier Go-owned-window idea while preserving Go ownership of orchestration and agent delivery. The choice favors established browser compatibility and straightforward integration over minimum runtime size.

Use the already accepted portable directory/archive distribution. CEF and other engines remain research references, not parallel implementation targets. Selecting the stack does not establish performance, desktop compatibility, or a completed agent return path.

The next milestone is the [first end-to-end proof](proposal.md#first-proof): launch, idle without model polling, submit, save, return to the exact caller, and collect. The third interview below settles its scope and close behavior. The initial API was subsequently implemented as a prototype; see the current contract and verification report.

## Accepted first-milestone scope

The [third submitted form](journal/2026-10-19-grill-689a37f9-b8c4-4f2a-9e54-282234bc2551.md) settles the following:

| Decision | Consequence |
| --- | --- |
| macOS first; Windows/Linux verification follows | Prove the initial loop on this Mac, without claiming other platforms have passed |
| Verify automatic return in Codex, Claude Code, OpenCode, Grok, Kimi and Pi | The answer's notes expand the selected three-harness option to six; Codex alone does not complete this milestone |
| Form plus pasted/selected image attachment | Keep drawing tools and a full screenshot UI outside the initial demo |
| Flat built-in JSON fields first | Defer nested composition and custom JSON components; the user has larger JSON-authoring ideas to discuss later |
| Window close ends the interaction with a distinct `closed` outcome and retains the draft | Never treat close as submission; explicit reopening uses the latest saved draft |
| Include a tiny native helper contract test | Exercise an extension outside Electron without implementing a native widget library |
| Launch input and local UI interaction are enough initially | Defer external live updates into an open session; local sliders/previews still work |

HTML and JSON remain equally supported, sharing submission and result handling. The flat JSON example is a deliberately small first surface, not a commitment to the eventual JSON architecture. Arbitrary local HTML/assets remain the route for custom UI during this milestone.

Read-only inspection confirmed all six commands are installed. Current versions and candidate return mechanisms are recorded in [delivery research](research/delivery.md#first-milestone-harness-inventory). Installation does not establish authentication, callback reliability or live Fungi compatibility. All six are verification targets; any unavailable automatic route remains an explicit gap, with the already accepted manual “done” behavior available rather than model polling.

## Current proposals, not accepted decisions

- Keep session ownership, lifecycle, extension communication, durable results, and delivery in the engine.
- Use a small, separable reference renderer for end-to-end testing.
- Treat native windows, webviews, component libraries, capture interfaces, and terminal interfaces as interchangeable implementations above the engine's protocol.
- Keep high-frequency interactions local. A slider change or drawing stroke should not invoke a model by default.
- Separate storing a valid result from notifying a harness and from confirming that an agent consumed it.
- Use the same rendering bridge for direct HTML and a small JSON component catalog. Catalog composition belongs above the session engine.
- Identify a supported delivery route before promising automatic resumption; otherwise launch with the accepted “done” fallback. Never substitute repeated model polling for an unavailable callback.

## Contract draft on 2026-10-20

The [first implementation contract](contract.md) makes launch, waiting, lifecycle, page/helper APIs, artifact persistence and recovery concrete. It proposes four CLI operations, immutable terminal results, new linked sessions when reopening a draft, and one shared commit path for HTML, flat JSON and native helpers. These are engineering recommendations for review; they do not change the accepted milestone scope or establish implementation evidence.

The user's follow-up asks whether more grilling is needed. Current product choices are sufficient to draft this contract. Another targeted interview is useful only if factual work reveals a consequential tradeoff, such as a required harness setup step. Ordinary command naming, internal transport and test-fixture choices do not need another broad questionnaire.

## Remaining implementation questions

1. Which built-in return mechanism works in each installed harness, with exact originating-session routing and no repeated model checks? The fourth interview rules out harness setup; use the explicit manual fallback where native capabilities are insufficient.
2. Does the draft contract hold up against concrete HTML, flat JSON and native-helper examples, including interrupted submission and exact-session recovery? Refine it from implementation evidence rather than growing a speculative plugin framework.
3. What desktop baselines and packaging changes are needed for later Windows/Linux support? Keep those claims unverified until exercised there.

[Research findings](research/index.md) and the [small authoring proposal](proposal.md) reflect these decisions. Electron + Go and the milestone scope are selected; exact APIs remain proposals. The initial prototype now implements these paths; coverage is recorded in the verification report.

## Implementation checkpoint on 2026-10-20

The user invoked `work-until-grill`, authorizing implementation and verification of the accepted first milestone. Go now owns detached sessions, immutable results, artifacts, drafts, explicit recovery and return delivery. Electron renders local HTML and flat JSON forms through a four-method bridge. A native Go helper exercises the same protocol. The core uses only the Go standard library.

Evidence changed one internal choice: a loopback listener failed under Codex’s sandbox, so process pipes and OS file locks replaced HTTP. The macOS portable folder preserves Electron’s framework symlinks and does not need a separately installed Node runtime. It is a development package, not a notarized public release.

[The feature map](features/index.md) and [first proof report](verification/runs/2026-10-20-first-proof.md) separate observed behavior from remaining work. Five harnesses returned real Fungi helper results; three-minute waits completed in Claude, Grok and Pi. Kimi lacks a configured credential. Longer AFK behavior and all-six native UI journeys remain unverified. The next consequential question concerns opt-in harness integrations where stock tools cannot wait indefinitely; no persistent integration configuration has been changed.

## No-setup delivery decision on 2026-10-20

The [fourth submitted form](journal/2026-10-20-grill-397cc4f4-31e2-4fd5-a08e-30ee197e3297.md) resolves that tradeoff: **require no harness setup; use the manual “done” fallback on unsupported versions.** Do not add a plugin, hook, startup configuration or harness upgrade as a prerequisite for Fungi. Existing built-in completion tools and the already available exact-thread Codex queue remain eligible.

All six installed harnesses remain compatibility targets. For a version or execution mode that cannot sustain a human wait, record its manual fallback explicitly; that is accepted compatibility behavior, not a passing automatic-return test. This clarifies the earlier all-six automatic-return target without claiming that an unsupported version gained a capability.

The user also reports completing Kimi sign-in and requests verification to resume. Retry the actual command; do not inspect credentials or assume authentication succeeded merely from the answer. The existing `work-until-grill` authorization continues through implementation, native UI checks and records.

## Native compatibility checkpoint

The [no-setup follow-up](verification/runs/2026-10-20-no-setup.md) verifies Kimi 2.2.0 after sign-in, complete native form/image returns in Pi, Grok, Kimi and Claude's bounded headless mode, and OpenCode v1's explicitly resumed manual flow. Launch receipts now include `collectEnv`, verified by the core suite and OpenCode's actual later call. No integration setup is required.

The latest portable build is `dist/fungi-darwin-arm64-2026-10-20T08-12-31-161Z`. Its Codex image session `39fdb74018dd13baae7fbc2883699530` completed the final check: the notice reached its originating conversation after the prior turn ended, and one collection using the supplied argv/environment returned the submitted response and matching image. The bounded macOS no-setup proof is complete. Long AFK and sleep/wake behavior, clean-machine distribution, and Windows/Linux runtime verification remain outside the completed proof. The earlier user-edited demo remains open with its draft preserved.

## Skill authoring and longer waits

The user authorized the next work through `work-until-grill`: exercise real skill authoring, longer idle/recovery behavior and portability. They specifically requested the simpler grill-form from yesterday's jstack history while its current version is being reworked. The last committed modification of the relevant assets is `7a7c919` (2026-10-13); those assets are unchanged at yesterday's tip, `32324e47aae7d3ad3f92e6a4aad8966e9451d691`. Use that committed baseline, not the uncommitted editor/review changes.

The Fungi example is now a skill-owned HTML directory and question JSON input. It keeps choices, plain notes and the older rich descriptions, adds one optional image through the existing artifact bridge, and uses Fungi's durable drafts and result delivery. No core or bridge API changes were needed. The native UI has passed close/reopen, abrupt owner-loss recovery, storage-failure retry, submission and cancellation. Linux arm64 core tests also passed in a network-isolated, unprivileged container. [Verification](verification/runs/2026-10-20-skill-proof.md)

The user also suggested Obsidian's Pele. [The assessment](research/pele.md) finds a promising optional diagram asset, with a small dependency-free browser bundle and recent development, but little maintenance history so far. It is not substituted into the committed grill-form baseline.

The two-hour native-helper probe, session `622f59b42e05840dd408622f23c8416a`, completed with one exact-origin Codex notice and one collection. Its [result](evidence/2026-10-20-skill-proof/idle-long-result.json) records 7,200.0148 helper-measured seconds and a UTC interval of 2 h 4 min 17.873 s, spanning the manual sleep/wake check below. The timing artifact matches the returned data and SHA-256; the owner and helper exited. No model polling or additional waiter was used. This verifies the selected helper/owner/Codex path, not arbitrary-duration behavior in every harness.

OrbStack provides Linux core execution but not proof of Linux desktop behavior. A Windows desktop and a clean Mac are not connected to the workspace. The follow-up interview below settles the immediate arrangements; public signing/notarization remains outside the completed prototype.

Questions 18–19 are complete in the Fungi-backed simple form, session `3863d61f14f01136a6d967630a00af91`. Its exact-origin callback was consumed with one collection. The [submitted answers](journal/2026-10-20-grill-desktop-arrangements.md) defer additional desktop checks until last; the user has a Windows machine but wants those checks deferred now. They chose to sleep/wake this Mac manually during the verification session. Do not provision a VM, connect to their Windows machine or trigger sleep as a consequence of this answer.

The sleep/wake fixture, session `1897358a56ceff72064a252004170716`, passed. Its [baseline](evidence/2026-10-20-skill-proof/sleep-wake-baseline.json) records the acknowledged `BEFORE-SLEEP — keep this line` note and reference image. macOS [power events](evidence/2026-10-20-skill-proof/sleep-wake-power.log) confirm clamshell sleep at 18:09:06 +0800 and final lid/HID wake at 18:13:38, with a maintenance dark wake/sleep in between. The exact-origin callback was collected once: the [response](evidence/2026-10-20-skill-proof/sleep-wake-result.json) retains that note, adds the user's `AFTER-WAKE` edit and returns identical image bytes. This verifies one short cycle on the current Mac.

The [skill-authoring and longer-idle proof](verification/runs/2026-10-20-skill-proof.md) is now complete for its selected scenarios. All callbacks in this round have been consumed. Additional desktop checks remain deferred by the user's answer; public distribution and arbitrary-duration reliability are not established. The current portable package is `dist/fungi-darwin-arm64-2026-10-20T09-41-25-905Z`.

## Rich questionnaire feedback

The user tried the updated jstack presentation in a Fungi window. Its [submitted response](evidence/2026-10-20-rich-grill/user-result.json) returned all six items through one callback and one collection; the [verification report](verification/runs/2026-10-20-rich-grill.md) records the tested build and remaining gaps. Mock choices and placeholder comments are demonstration data.

The final note expresses a real concern that UI authoring and handoff may be expensive and hard for models to understand. The user intends eventually to bundle Fungi with jstack and rewrite skills such as grill-form around it once the base is sufficiently stable. This reinforces the existing simplicity requirement. The compatibility adapter used for this experiment is not the intended routine authoring workflow.

The user then explicitly confirmed the authoring requirement: agents should follow robust mechanical instructions while understanding the available grill-form capabilities well enough to author meaningful sessions. Both parts matter: a compact capability reference with working examples and limits, and a fixed procedure for validation, launch, completion and interpreting the result.

The target is question data plus one launch, followed by automatic structured response delivery. Rendering, draft storage, attachment conversion and callback handling should be reusable implementation supplied with the skill/Fungi distribution. An agent should not reconstruct the template adapter or reconcile two collectors for every questionnaire. The skill should describe choices, open feedback, visual comparisons, diagrams, prototype state, region comments and images without requiring agents to read renderer internals. This is a requirement for the eventual integration, not a claim that the current experiment already provides that final authoring experience.

## Foundation work on 2026-10-22

The user renewed `work-until-grill` to pursue a stable, robust, simple and broadly flexible rendering platform. Continue within the accepted Go/Electron, local-first, no-harness-setup scope. Full jstack authoring integration remains deferred; a small reusable authoring guide and independent example are appropriate foundation work.

Concrete verification initially exposed four repairs: asynchronous Blob conversion could reorder mutations; exact-limit JSON could be saved but unreadable because of its newline; malformed/closed helper stdout could leave the owner blocked in process wait; a failed reference-form autosave could prevent a later submit retry without another edit. Repairs retain the existing protocol and bridge.

Renderer support now includes local/Blob workers and WebAssembly, plus streamed local assets with media range requests. The independent before/after comparison example uses ordinary HTML and canvas, exports a PNG and adds no engine API. Its submitted feedback (39% divider, After, “Nice”) is demo data, not a new project requirement. [The verification record](verification/runs/2026-10-22-platform.md) separates runtime observations from broader compatibility claims.

Treat existing protocol-1 meanings as the extension compatibility boundary. Keep the flat JSON schema experimental and grow capabilities through ordinary skill assets/helpers. Cross-platform desktop and clean-machine release verification remain the last outstanding acceptance work; signing/notarization and universal harness resumption are not established by a local macOS pass.

The follow-up Linux run exposed a fifth issue: a helper could emit a valid
result and exit before the launch write completed. The owner now drains its
output after write failure and bounds blocked protocol writes. macOS and Linux
core suites pass after the repair. Portable Windows x64 and Linux x64 bundles
are prepared using pinned Electron downloads. Windows file handling has a
separate directory-sync policy, and lock errors no longer masquerade as live
owners. Actual desktop validation is the next unresolved acceptance boundary;
compiled bundles are not a platform-support claim.

Hosted follow-up testing now exercises the real Electron renderer. macOS and
Windows Server 2025 pass the 21 rendering assertions with an empty child PATH.
Windows required direct access to its inherited input pipe, suppression of
Electron console attachment and tolerance of blank NDJSON lines. The four
page methods and native-helper contract remain unchanged.

All hosted core suites pass. Ubuntu 22.04 completes the renderer workflow and
20 checks, but its virtual display cannot provide WebGL. Ubuntu 24.04 aborts
at Chromium sandbox startup under its default namespace policy. No sandbox
bypass or machine policy change was made. The
[hosted report](verification/runs/2026-10-22-hosted.md) remains **Failed** overall.
This exposed a choice between keeping no-install distribution with
explicit supported Linux environments, allowing an optional Linux installation
path, or deferring Linux UI. The following interview resolves it. Full jstack
authoring integration remains deferred.

## Linux distribution decision

The [submitted answer to question 20](journal/2026-10-22-grill-0c10be35-a07a-4a8d-8d5d-39fda5b2e1bc.md)
selects **keep no-install; support explicitly verified Linux environments**.
The Fungi callback was collected once, and the matching jstack transcript
was validated and finished. There were no additional notes.

Keep the portable folder/archive and bundled Electron runtime. Publish the
tested environments and their capability limits; make unsupported startup
failures diagnosable. This authorizes neither an installer nor sandbox
bypasses, elevated helper setup or global host policy changes. Ubuntu 22.04's
tested workflow is usable without WebGL; Ubuntu 24.04's tested default
sandbox policy remains outside UI support. The earlier broad matrix remains
Failed in its historical report.

The [platform record](platforms.md) distinguishes workflow coverage, optional
GPU/codec availability and physical clean-machine acceptance. Full jstack
authoring integration remains deferred. No new skill/plugin architecture is
needed to apply this decision.

The [baseline follow-up](verification/runs/2026-10-22-baselines.md) passes
those selected scopes on revision `199759a`. Unsupported startup now returns
a saved, distinct failure with the diagnostic log path. The four-method page
API remains unchanged. The foundation is ready for further skill work within
the documented limits. Physical target-machine and clean-machine acceptance
remain the separate distribution gap; no such result is inferred from CI.

The final Windows handoff was opened as Fungi session
`61a2170bba8d62aa0291aec1aad663cc`, using the new Mac portable package and an
armed exact-conversation callback. Question 21 asks for the observed physical
Windows result or an explicit deferral; it does not assume either. The ZIP is
`dist/fungi-win32-x64-2026-10-22T09-50-20-630Z.zip`, with the full portable
folder and `VERIFY.md`. Its SHA-256 is
`413a6526920e3c8d0fbae17d4a79d31fe50a9309472ee74901eb95cf465529c9`.

The callback was collected exactly once and the matching jstack transcript
was validated and finished. The [submitted response](journal/2026-10-22-grill-b8ce5a10-7261-48c1-9d80-ee54e103ec5e.md)
asks for instructions in the wiki: the user will pull on Windows and have the
agent there execute them. It does not report a pass, failure or deferral.
Physical Windows verification therefore remains pending execution there.

The [Windows agent handoff](verification/windows.md) is the entry point. It
specifies `foundation-hardening`, building from tracked source because
portable packages are ignored by Git, the automated and visible checks, and
how to commit evidence back into this wiki. No new questionnaire is pending;
the next action belongs to the Windows agent after the pull. The original
receipt and response remain under
`verification-output/platform-2026-10-22/windows-handoff/` on this Mac.

## Physical Windows handoff on 2026-10-22/23

The user requested execution of the Windows guide and later asked the agent
to continue. The [completed handoff record](verification/runs/windows-2026-10-22-203907.md)
identifies revision `73eb9ad`, the locally built portable package, native
Windows 11 Pro x64 build 26200 and the evidence. The agent operated synthetic
controls through native mouse/keyboard input. Launch, acknowledged draft,
close/reopen with restored values, PNG submission/hash and separate
cancellation all passed with an empty child PATH.

The overall run is **Untested**: 17 core and six Node tests passed, but the
remaining core/asset tests and the full renderer fixture could not create
their prerequisite symlinks under this account. Exact setup failures are
retained; no elevation, security-policy change or weakened fixture was used.
A checkout line-ending mismatch was resolved by applying the branch's
existing LF policy; no product source repair was needed. Node/npm were used
from an isolated portable toolchain for building only.

This completes the visible desktop journey on that machine and records the
remaining automated coverage gap. It does not certify a clean-machine or
downloaded-archive launch, Windows clipboard, WebGL, automatic agent return,
sleep/wake or power-loss behavior. The earlier hosted and Mac reports retain
their original scopes. Full jstack authoring integration remains deferred.
