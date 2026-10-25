---
title: Offline manual and installed example discovery
updated: 2026-10-26
result: Passed
---
# Offline manual and installed example discovery

**Passed for the implemented offline discovery and native example scope.**
Local native checks and all five CI jobs passed. Independent fresh-agent
usability remains outside this proof. Affected features: [discovery](../../features/discovery.md),
[interfaces](../../features/interfaces.md) and
[installation](../../features/installation.md).

## Build and selected scope

The source base is `aad403d`; the
[source hashes and native observations](../../evidence/2026-10-26-discovery/local-observations.json)
identify the implementation changes. Local packages `v0.1.0-docs.1` and `.2`
were built on macOS arm64; `.2` adds bootstrap flag reference text. The
[native test manifest](../../evidence/2026-10-26-discovery/local-build.json) records
CLI/helper hashes. [Platform CI](https://github.com/rengwu/fungi/actions/runs/40219931322)
tests implementation `9572d80ae41a41a7916f74f4d80b1f4f581b2593`, including a final
saveDraft return-description correction after the local packages. The native
example and renderer code are unchanged between those builds. The native-test archive is 6,070,202 bytes (5.79 MiB), including
the offline manual and four small examples. The renderer implementation is
unchanged by this documentation/discovery work.

## Observed results

| Scenario | Result | Observation |
| --- | --- | --- |
| Core vet/race tests and existing Node suite | Passed | `go vet ./...`, `go test -race ./...` and all six `node --test` cases passed locally |
| Standalone offline help and reference routing | Passed | New executable integration test invokes every topic from an unrelated directory with empty PATH and no adjacent docs/renderer; checks exact embedded text, entry points, JSON reference routes, missing examples and invalid-topic exit 2 |
| Installed docs before/after update and rollback | Passed | [Install evidence](../../evidence/2026-10-26-discovery/local-install.json): real bootstrap, install path with spaces, unrelated cwd and empty child PATH. All eight help topics match packaged Markdown; all four example paths select the active version at each stage |
| Existing live session through update | Passed | Same install driver retains draft, updates while UI is active, submits PNG and verifies bytes/hash, then rolls back |
| Native reply submit | Passed | Actual Wails window displayed input title/prompt. Typed synthetic notes, clicked Send, collected once. [Result](../../evidence/2026-10-26-discovery/reply-result.json) matches exactly, with no artifacts or draft file |
| Native failed write and retry | Passed | Made only the test session unwritable. Send showed permission denied and retained the typed text, with no result committed. Restored permissions and retried; [result](../../evidence/2026-10-26-discovery/retry-result.json) preserves exact text. [Fault observations](../../evidence/2026-10-26-discovery/retry-observation.json) |
| Native Cancel and close | Passed | Separate actual native actions returned [cancelled](../../evidence/2026-10-26-discovery/cancel-result.json) and [closed](../../evidence/2026-10-26-discovery/close-result.json); fresh sessions had no draft |
| Explicit reopen without saved draft | Passed | Unsaved text was absent in the reopened native window. Engine creates an empty draft placeholder, not saved user text. [Cancelled reopened result](../../evidence/2026-10-26-discovery/reopened-result.json) collected once; original close result retained |
| Workflow variation beyond shipped samples | Passed | Maintainer authored a [priority-order page source](../../evidence/2026-10-26-discovery/variation.html.txt) with [input](../../evidence/2026-10-26-discovery/variation-input.json). Native move-up changed the visible order; Send returned [ordered data](../../evidence/2026-10-26-discovery/variation-result.json) and [SVG source](../../evidence/2026-10-26-discovery/priority.svg.txt). Artifact hash and XML labels match the submitted order; no engine change or draft |
| Discovery on all four native CI platforms | Passed | [macOS](../../evidence/2026-10-26-discovery/macos-ci-install.json), [Windows](../../evidence/2026-10-26-discovery/windows-ci-install.json), [Ubuntu 22](../../evidence/2026-10-26-discovery/ubuntu22-ci-install.json), [Ubuntu 24](../../evidence/2026-10-26-discovery/ubuntu24-ci-install.json): embedded and packaged docs match; all example argv resolve inside the selected version before/after update and rollback, with an unrelated cwd, spaces and empty child PATH. Ordinary rendering and existing runtime checks also passed within accepted exclusions |
| Same Ubuntu 22 package on Ubuntu 24 | Passed | [Compatibility evidence](../../evidence/2026-10-26-discovery/compatibility-ci-install.json): the unchanged older-baseline archive passes the same help/example/update/rollback assertions. Existing renderer/runtime checks also passed within accepted media exclusions |

One verification assertion initially expected no `draft.json` after reopen.
The engine intentionally writes `{data:null,artifacts:[]}` when restoring an
absent draft. The assertion was corrected; the observed UI was blank as intended.
The result was read from its saved file after the assertion failed, without a
second collect call. No product change was needed.

[CI completion record](../../evidence/2026-10-26-discovery/ci-run.json) identifies
all five successful jobs. The final compatibility runner spent several minutes
installing its system runtime before successfully running the checks; no Fungi
fix or rerun was required. The wiki reader and report metadata checks passed.
All six local fixture sessions reached terminal results, and the temporary
write-failure permissions were restored.

## Scope limits

These are maintainer-operated checks. An independent fresh-context agent trial
of an unfamiliar interface has not been performed; no claim of independent
discoverability proof is made. Full rich grill-form integration remains deferred.
Existing bridge ordering and save-retry tests also passed. This run does not requalify long AFK callbacks,
missing WebView2 installation, clean consumer desktops or signed releases.
No public release was published and no host runtime/security policy was changed.
