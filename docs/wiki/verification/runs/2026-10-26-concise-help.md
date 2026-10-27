---
title: Concise CLI documentation
updated: 2026-10-26
result: Passed
---
# Concise CLI documentation

**Passed for the local documentation scope.** The user approved a shorter,
neutral rewrite of all eight help topics while preserving core behavior and
complete examples. Scope: [FEAT-discovery](../../features/discovery.md).
No renderer, session, delivery, installer or formatter implementation changed.

## Build and environment

Tested on macOS arm64 with a standalone core built from `a94e3ffd5dd1625cca724f74f8f40b3db5e894e7`
plus uncommitted manual/test edits. The [build record](../../evidence/2026-10-26-concise-help/build.json)
contains source-file hashes, CLI SHA-256, Go version and hashes for every captured
help output. All topics were invoked outside the repository root with empty
PATH and no adjacent renderer or manuals.

## Results

| Scenario | Expected | Observed | Result |
| --- | --- | --- | --- |
| Offline help and discovery | All topics, aliases and reference routes work; invalid requests return JSON/exit 2 | Existing focused tests passed; [log](../../evidence/2026-10-26-concise-help/tests.log) | Passed |
| Formatting and copyability | Plain terminal text; fenced examples preserved verbatim, including long commands | Formatter tests and all-topic checks passed; all eight captured outputs reviewed | Passed |
| Content preservation | Retain API, outcomes, limits, recovery and installation constraints | Maintainer review against the prior manuals and relevant implementation found no removed core contract; no independent-agent usability test claimed | Passed |
| Installation examples | One physical line per primary shell/PowerShell command | Both blocks checked; commands preserved in [installation output](../../evidence/2026-10-26-concise-help/installation.txt) | Passed |

Captured output: [start](../../evidence/2026-10-26-concise-help/start.txt),
[CLI](../../evidence/2026-10-26-concise-help/cli.txt),
[delivery](../../evidence/2026-10-26-concise-help/delivery.txt),
[HTML](../../evidence/2026-10-26-concise-help/html.txt),
[form](../../evidence/2026-10-26-concise-help/form.txt),
[helper](../../evidence/2026-10-26-concise-help/helper.txt) and
[examples](../../evidence/2026-10-26-concise-help/examples.txt).

## Limits

This verifies documentation output and content, not native UI or cross-platform
runtime behavior. Those implementations are unchanged and were outside the
selected scope. The existing terminal formatter remains the sole formatter;
Markdown is still the shared embedded and packaged source.

No release was published or user installation modified. Temporary build files
were removed after capturing evidence. The next package build will embed the
rewritten manuals.
