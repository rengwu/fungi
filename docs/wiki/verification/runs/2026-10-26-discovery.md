---
title: Offline manual and installed example discovery
updated: 2026-10-26
result: Untested
---
# Offline manual and installed example discovery

**Untested overall while platform CI is pending.** Local macOS checks below
passed. Affected features: [discovery](../../features/discovery.md),
[interfaces](../../features/interfaces.md) and
[installation](../../features/installation.md).

## Build and selected scope

The source base is `aad403d`; the
[source hashes and native observations](../../evidence/2026-10-26-discovery/local-observations.json)
identify the implementation changes. Local packages `v0.1.0-docs.1` and `.2`
were built on macOS arm64; `.2` adds bootstrap flag reference text. The
[native test manifest](../../evidence/2026-10-26-discovery/local-build.json) records
CLI/helper hashes. The native-test archive is 6,070,202 bytes (5.79 MiB), including
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
| Native Cancel and close | Passed | Separate actual native actions returned [cancelled](../../evidence/2026-10-26-discovery/cancel-result.json) and [closed](../../evidence/2026-10-26-discovery/close-result.json); fresh sessions had no draft |
| Explicit reopen without saved draft | Passed | Unsaved text was absent in the reopened native window. Engine creates an empty draft placeholder, not saved user text. [Cancelled reopened result](../../evidence/2026-10-26-discovery/reopened-result.json) collected once; original close result retained |
| New discovery behaviour on Windows/Linux packages | Untested | Native platform CI pending |

One verification assertion initially expected no `draft.json` after reopen.
The engine intentionally writes `{data:null,artifacts:[]}` when restoring an
absent draft. The assertion was corrected; the observed UI was blank as intended.
The result was read from its saved file after the assertion failed, without a
second collect call. No product change was needed.

## Scope limits

These are maintainer-operated checks. An independent fresh-context agent trial
of an unfamiliar interface has not been performed; no claim of independent
discoverability proof is made. Full rich grill-form integration remains deferred.
New reply error/retry UI has not been fault-injected; existing bridge ordering
and save-retry tests passed. This run does not requalify long AFK callbacks,
missing WebView2 installation, clean consumer desktops or signed releases.
No public release was published and no host runtime/security policy was changed.
