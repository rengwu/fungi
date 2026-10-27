---
title: Repository cleanup and verification instructions
updated: 2026-10-27
result: Passed
---
# Repository cleanup and verification instructions

**Passed for the selected cleanup, argument-validation and local native checks.**
Scope: [discovery](../../features/discovery.md) and
[interfaces](../../features/interfaces.md). The user moved workflow-skill
experimentation to `jstack-fungi-test`, requested removal of `skills/`, and
approved audit items 1–8. Optional research/proposal archiving and removal of
the existing `dist/release-v0.1.0-clean.2/` were not selected.

## Changes

Removed the orphaned questionnaire model test and the four-file jstack HTML
adapter example. Generic examples, the built-in JSON form, core runtime and
verification drivers remain. Updated skill references, current Wails recovery
guidance, Windows checkout instructions and release examples/defaults. The
renderer smoke driver now requires one explicit executable path; it no longer
defaults to the removed root executable.

The engine contract now links to the embedded manual for API details and keeps
the engine invariants. Current decisions and references appear before historical
chronology. Original decisions, journals and evidence are preserved; original
planning headings retain their anchors. No release, commit or push was performed.

## Build and environment

The [build record](../../evidence/2026-10-27-repo-cleanup/build.json) identifies
base `a94e3ffd5dd1625cca724f74f8f40b3db5e894e7`, dirty source fingerprints,
macOS arm64, Go 1.26.4, Node 22.19.0 and Wails 2.17.0. A fresh development
package was built from the working tree; its CLI/helper hashes identify the
tested binaries. Existing edits and concurrent window-sizing work were
preserved. Those runtime changes are outside this cleanup's verification scope.

The renderer driver used its own disposable state/assets and cleared child
PATH and caller routing. The native idle helper used a separate state directory
and no callback. These fixtures submit synthetic responses without a human wait.

## Results

| Scenario | Expected | Observed | Result |
| --- | --- | --- | --- |
| Remaining Node suite | Removing the skill leaves no broken test import | All four bridge tests passed; [output](../../evidence/2026-10-27-repo-cleanup/node-tests.log) | Passed |
| Smoke-driver argument handling | Exactly one binary is required before fixture creation | Missing and extra positional arguments exited 1 with the usage message; [checks](../../evidence/2026-10-27-repo-cleanup/checks.json) | Passed |
| Fresh native package | Complete CLI/helper reports an available runtime | Build succeeded and [doctor](../../evidence/2026-10-27-repo-cleanup/doctor.json) reported WKWebView available | Passed |
| Explicit-path renderer recipe | Cancel, reopen and submit through the real bridge; preserve the PNG | All 21 checks passed with no exclusions; [renderer evidence](../../evidence/2026-10-27-repo-cleanup/renderer.json) and [returned PNG](../../evidence/2026-10-27-repo-cleanup/renderer-canvas.png) | Passed |
| Standalone idle-helper recipe | Five-second native wait returns the marker and timing artifact | Submitted after 5.011 seconds; JSON and artifact SHA-256 matched; [helper evidence](../../evidence/2026-10-27-repo-cleanup/idle-helper.json) | Passed |
| Documentation integrity | Current recipes avoid removed files; history and links survive | Removed-path scan, planning-history comparison, report metadata check, wiki build/check and diff whitespace check passed | Passed |
| Offline reader | Rebuilt local HTML exposes current guidance and the audit report | Opened in Helium; current-reference navigation and cleanup-report link rendered | Passed |

The initial full Node suite failure after deleting `skills/` was the expected
orphan import in `scripts/grill-form.test.mjs`. Removing that test repairs the
harness; no product behavior was changed to make it pass.

## Limits and cleanup

This run does not qualify Windows/Linux desktop input, published downloads,
window-sizing behavior, clipboard input, a new crash-recovery scenario or long
AFK intervals. Recovery guidance cites the already recorded Wails full-helper
exit proof; earlier Electron results retain their original scope. Release and
Windows instructions were edited, not executed on external hosts.

Only this run's newly built package and disposable renderer/helper state were
removed after preserving evidence; the [cleanup record](../../evidence/2026-10-27-repo-cleanup/cleanup.json)
confirms no matching processes remained and the historical record files were
unchanged. The pre-existing distribution, historical
research/proposal pages, user decisions and earlier verification records remain.
