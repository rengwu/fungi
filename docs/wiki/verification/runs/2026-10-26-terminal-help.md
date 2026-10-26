---
title: Terminal help and one-line installation
updated: 2026-10-26
result: Passed
---
# Terminal help and one-line installation

**Passed for the terminal-help and installation-documentation scope.**
The alpha is published and all 11 release-workflow jobs passed.
Scope: [FEAT-discovery](../../features/discovery.md) and the documented
[installation entry points](../../features/installation.md). The user supplied
a screenshot of raw Markdown in the installed CLI and asked for terminal
formatting, while requiring future installation instructions to be one line.

## Build and changes

Release tag `v0.1.0-alpha.2` points to source
`68e239cdf6f4972d7001f95f810c373e222fe83a`. The
[release workflow](https://github.com/rengwu/fungi/actions/runs/40289692224)
builds and verifies the actual archives before publication. The
[local build record](../../evidence/2026-10-26-terminal-help/local-build.json)
identifies the standalone CLI and Go version.

Help formats the existing embedded manual as plain terminal text. Headings and
inline formatting markers are removed, two-column tables align, prose wraps at
80 columns, and fenced examples keep their content with a two-space indent.
Long commands remain one physical line. It needs no new dependency, ANSI escape
sequences, pager or adjacent manual. `describe` retains JSON output. Source
Markdown still ships for offline reading.

README, current manual and release notes provide one-line shell and PowerShell
installation examples. Shell examples include current-terminal PATH setup;
persistent shell configuration remains explicit. The maintainer release recipe
records the one-line requirement for future releases.

## Results

| Scenario | Result | Observation |
| --- | --- | --- |
| Local formatting and offline discovery | Passed | Focused tests cover readable columns/prose, literal code content, long commands, every help topic without PATH tools/docs/renderer, and invalid-topic JSON/exit 2; [test output](../../evidence/2026-10-26-terminal-help/tests.log) |
| Actual terminal and redirected output | Passed | Help was inspected through a PTY and saved as plain text; [overview](../../evidence/2026-10-26-terminal-help/start.txt), [HTML contracts](../../evidence/2026-10-26-terminal-help/html.txt) and [installation](../../evidence/2026-10-26-terminal-help/installation.txt) have no raw Markdown headings/tables/fences |
| Copyable install instructions | Passed | README, manual and both alpha release-note pages were checked: each primary shell/PowerShell install block contains one physical line |
| Platform package/install/update checks | Passed | Mac, Windows and both Ubuntu foundation jobs passed, including all eight formatted help topics and preserved code examples before update, after update and after rollback |
| Shared Linux archive | Passed | The Ubuntu 22-built archive passed unchanged on Ubuntu 24, including formatted help during install/update/rollback |
| Published alpha and public-download checks | Passed | Three archives and both installer scripts are public with SHA-256 sidecars; anonymous installation, update, example discovery and the existing renderer fixture passed on all four hosted baselines |

The [completed CI record](../../evidence/2026-10-26-terminal-help/ci.json)
contains every job and step conclusion. The
[release inventory](../../evidence/2026-10-26-terminal-help/release.json) records
published assets and checksums, while [public build identities](../../evidence/2026-10-26-terminal-help/public-builds.json)
confirm the exact tagged source was installed. Detailed public checks are kept
for [Mac](../../evidence/2026-10-26-terminal-help/macos-latest-public.json),
[Windows](../../evidence/2026-10-26-terminal-help/windows-latest-public.json),
[Ubuntu 22](../../evidence/2026-10-26-terminal-help/ubuntu-22.04-public.json) and
[Ubuntu 24](../../evidence/2026-10-26-terminal-help/ubuntu-24.04-public.json).

The user's installed alpha.1 was not modified. The supported update command is
`fungi update --version v0.1.0-alpha.2`. Temporary standalone-test files and
downloaded CI scratch artifacts were removed after preserving evidence. The
source manual and published release notes use one-line installation examples;
the earlier alpha's immutable package is not rewritten.

## Limits

This is a formatter for the manual's controlled syntax, not a general Markdown
library. Code lines are deliberately not wrapped. Renderer implementation,
callbacks, runtime installation and session lifecycle are unchanged. Existing
unsigned-alpha, macOS 26+ and media qualification limits still apply; this change
does not establish clean consumer-machine or long sleep/wake acceptance.
