---
title: Offline documentation and example discovery
updated: 2026-10-26
---
# Offline documentation and example discovery

**ID: FEAT-discovery**

## Entry points and expected behaviour

`fungi`, `fungi --help` and `fungi help` show a concise capability overview
followed by common starting points. `fungi help TOPIC` prints the complete
Fungi-specific reference as Markdown, embedded in the CLI. It needs no network,
renderer, man viewer or adjacent documentation directory. Invalid topics or
extra arguments produce a JSON error and exit 2.

`fungi describe` retains protocol 1 JSON discovery and exposes documentation
topics. `describe TOPIC` includes a `reference` argument array for its full
manual. `describe examples` reports absolute source/input paths and launch
arrays relative to the executing package, independent of working directory.
Missing example files are marked unavailable and have no runnable `argv`.

Packages include the same `docs/manual/*.md` plus small reply, comparison,
review and flat-form examples. Managed installs resolve the selected version;
updates and rollback must switch discovery accordingly. Copy example directories
before editing. The development wiki, rich skill adapters and third-party
library manuals remain outside the slim package.

The manual distinguishes skill authoring from execution. It offers complete
flows without requiring a separate Fungi skill, plugin registry or framework.
The reply example accepts title/prompt JSON, submits notes, handles cancellation
and errors, and intentionally saves no draft. Draft recovery in other pages is
optional; final result persistence does not depend on it.

## Verification recipes

1. `go test -run TestOfflineHelpAndDiscovery .` builds a standalone core outside
   the checkout and invokes every help topic from another directory with empty
   PATH. Expect exact embedded Markdown, working reference routes, explicit
   missing examples and structured invalid-topic errors.
2. Build two versioned packages and run the existing `scripts/install-smoke.mjs`
   recipe in [setup](../verification/setup.md). It installs under a path with
   spaces, invokes from a different working directory with empty child PATH,
   compares each help topic with packaged Markdown and checks example paths
   before update, after update and after rollback. The ordinary live UI update
   scenario still exercises session preservation independently.
3. Use the selected package's `describe examples` reply invocation with manual
   delivery and isolated session storage. Through the actual native window,
   submit synthetic notes; collect once with its receipt and verify exact text,
   empty artifacts and no draft file. Repeat Cancel and window close. Explicitly
   reopen the closed session and verify empty input text, then cancel it.

An unfamiliar-interface exercise can check whether the documentation transfers
beyond sample copying. Separate an independent fresh-agent result from a
maintainer's own test; do not label one as the other.

## Verification history

- [2026-10-26 offline discovery](../verification/runs/2026-10-26-discovery.md) — local offline/install/native checks; platform CI status and independent-agent limits are explicit.
