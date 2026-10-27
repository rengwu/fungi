---
title: Starting window dimensions
updated: 2026-10-27
result: Passed
---
# Starting window dimensions

**Passed for the selected local Mac checks.** The user requested optional
starting dimensions for skill-authored interfaces. Scope:
[FEAT-interfaces](../../features/interfaces.md) and
[FEAT-sessions](../../features/sessions.md).

## Change and build

`open` and `reopen` accept `--width` and `--height`. Either may be omitted;
defaults remain 820 × 780. Values must be 32-bit integers of at least 420.
Reopen inherits each omitted dimension from the original launch. Native
helpers reject the flags because they own their windows. The existing launch
record and renderer message carry the values; no new page method or dependency
was added. CLI help and HTML/form discovery expose the options.

The [build record](../../evidence/2026-10-27-window-size/BUILD.json) identifies
the macOS arm64 Wails 2.17 package, base commit and binary/helper hashes.
The checkout was dirty, including earlier unrelated work.
[Checks and source hashes](../../evidence/2026-10-27-window-size/checks.json)
identify the code tested. No release or user installation was changed.

## Results

| Scenario | Expected | Observed | Result |
| --- | --- | --- | --- |
| Omitted dimensions | Preserve 820 × 780 defaults | Native page reported 820 × 780 logical pixels at device scale 2 | Passed |
| Both dimensions | Apply 960 × 640 at launch | Native page reported exactly 960 × 640 | Passed |
| One dimension | Preserve the other default | Width-only reported 700 × 780; height-only reported 820 × 600 | Passed |
| Minimum | Accept 420 × 420 | Native page reported exactly 420 × 420 | Passed |
| Reopen | Retain original dimensions and draft; permit overrides | All five cases reopened at identical sizes. A width override changed 960 × 640 to 740 × 640 and retained the original draft | Passed |
| JSON form | Apply size without breaking built-in submission | Form launched with 900 × 600 content; native screenshot includes the title bar. Clicking Send response closed it; one collection returned the exact synthetic note | Passed |
| Invalid flags and compatibility | Reject malformed, missing, undersized and overflowing values, and helper use; accept older records | Go tests returned exit 2 before creating storage, covered HTML/JSON launch parsing and size-less records, and passed the existing session regressions | Passed |

[Native evidence](../../evidence/2026-10-27-window-size/evidence.json) contains
the five cancel/reopen flows and width override, collected through one pending
call per command. The [HTML fixture](../../evidence/2026-10-27-window-size/fixture.html.txt)
reads the real viewport; it does not resize the window or bypass the CLI.
`go vet ./...` and `go test -race ./...` passed (71.450 seconds for the core).

The [JSON fixture](../../evidence/2026-10-27-window-size/form.json),
[receipt](../../evidence/2026-10-27-window-size/form-receipt.json),
[window screenshot](../../evidence/2026-10-27-window-size/form.jpg) and
[submitted result](../../evidence/2026-10-27-window-size/form-result.json)
preserve the native form observation. Its 1800 × 1264 physical-pixel screenshot
corresponds to 900 × 600 content at scale 2 plus a 32-point native title bar.
The UI driver initially timed out while selecting the already-exited test app;
selection succeeded after the disposable form was launched.

## Limits and cleanup

Windows/Linux desktop execution, small-screen OS adjustments and mixed-DPI
displays were not selected for this local run. The maintained smoke driver can
be run on those hosts; no cross-platform qualification is claimed here.
Window decorations vary by platform, and manually resizing does not update
the stored launch dimensions. The source uses the existing Wails size options.

Evidence was copied before removing this run's disposable package and session
directory. Receipt paths identify the original test and are not retained launch
instructions. The user's installed Fungi and other packages were unchanged.
