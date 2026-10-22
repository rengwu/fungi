---
title: Command installation and update proof
updated: 2026-10-23
result: Untested
---
# Command installation and update proof

**Untested overall while hosted execution is pending.** The selected local
Mac scenarios passed. This run covers [installation](../../features/installation.md)
and preservation of [sessions](../../features/sessions.md) through updates.
It does not establish a public release or clean consumer-machine support.

## Build and environment

Local base `c768218`, with implementation edits; [source hashes](../../evidence/2026-10-23-distribution/source-hashes.json)
record the dirty source checkpoint. macOS 27.0 arm64, Go 1.26.4, Node 22.19.0
and Electron 44.9.0. [Package A](../../evidence/2026-10-23-distribution/package-a.json)
and [package B](../../evidence/2026-10-23-distribution/package-b.json) identify
the actual executable hashes and sizes. Their versions are `v0.1.0-install.1`
and `v0.1.0-install.2`. Later packaging changes increase archive compression
and substitute a compact portable README; those two archives predate these
packaging-only changes. Hosted checks will identify the committed source.

All local interactions use synthetic fixtures and isolated session/install
directories. No real questionnaire is submitted, no agent callback is queued,
no release is published and no shell profile/PATH is changed. Installed CLI
children use an empty PATH; the shell bootstrap uses standard OS tools.

## Scenarios

| Scenario | Result | Observation |
| --- | --- | --- |
| Core regression and archive tests | Passed | [Race-enabled log](../../evidence/2026-10-23-distribution/go-tests.log): 20 root tests and two archive test groups; Go vet and seven Node tests also passed |
| Managed A → B → C, recovery and cleanup | Passed | Actual CLI/helper processes; stable receipt path, failed checksum preserves A, concurrent mutation rejected, running/recoverable A retained, draft restored through C, old result immutable, artifact retained, rollback and later pruning succeed |
| Invalid downloads/archive paths | Passed | Oversize, partial and missing HTTP responses rejected; ZIP traversal/absolute/drive paths, duplicate files, escaping links and TAR hardlinks rejected; normal symlinks/modes survive Unix round trips |
| Shell bootstrap on development Mac | Passed | [Download evidence](../../evidence/2026-10-23-distribution/mac-install.json): exactly one matching archive plus checksum fetched from a local fixture server, checksum checked, space-containing destination installed |
| Open Electron session across update | Passed | A acknowledges draft and PNG, B activates, A submits, original stable receipt collects the matching saved answer; [PNG](../../evidence/2026-10-23-distribution/submitted.png) hash verified, rollback selects A |
| Installed renderer regression | Passed | [All 21 checks](../../evidence/2026-10-23-distribution/mac-renderer.json), including WebGL, cancel/reopen, restored Blob, submission and hash through the stable installed command |
| Hosted macOS, Windows and Linux install/update paths | Untested | CI execution pending; Ubuntu 24.04's installation check will intentionally omit UI |

The first unsafe-archive test attempt reused a ZIP header object and failed in
the test writer. A fresh duplicate header corrected the fixture. The first
managed-path assertion also exposed macOS `/var` versus `/private/var` aliases;
package-reference matching now resolves existing symlinks before retaining
assets. The rerun passes. No failed product assertion was simply removed.

## Limits and follow-up

The local HTTP fixture proves the bootstrap download path; updates use an
explicit local archive, not a live GitHub release. HTTP failures are tested
separately. A saved receipt was collected through the stable command; actual
agent queue delivery during an update is not claimed here. UI actions in this
run are automated page fixtures, not native human keyboard/mouse input.

Hosted/physical Windows and Linux execution, fresh-shell PATH integration,
downloaded/quarantined consumer packages, signing/notarization, release
publication, interruption during selector replacement and power-loss recovery
remain unverified. Existing platform and harness limits remain in force.
Verification was performed by the implementing agent; no independent review
is claimed.
