---
title: External helper interfaces
updated: 2026-10-28
---
# External helper interfaces

**ID: FEAT-interfaces**

The accepted engine-only scope delegates all UI implementation to other
projects. Fungi executes helpers and manages the session protocol.

## Entry points and expected behavior

`fungi open --exec HELPER --input FILE --delivery MODE -- ARGS...` launches a
trusted executable without a shell. Arguments are literal; input is arbitrary
JSON. Protocol 2 launch messages contain input, draft, draftArtifacts, outputDir,
limits and public delivery metadata. Helpers send ready, draft and result
messages and consume matching acknowledgments. No UI executable or page bridge
is bundled or maintained here.

Direct HTML/JSON-form launch and engine window flags fail with invalid_arguments.
Old help/describe html/form topics are withdrawn. `fungi help migration` explains
the break. Other projects can implement web/native/other presentation above
this process contract. They supply runtimes, assets and input handling.

## Verification recipe

Follow [setup](../verification/setup.md), then run the packaged engine smoke
and `go test -run 'TestHelper|TestEngine|TestUnlaunchable' .`. The smoke launches
an external Node fixture with an empty child PATH and no display environment;
Node is a verifier/helper dependency, not an engine requirement. Verify literal
JSON, isolated sessions, saved files, close/reopen, and rejected UI entry points.

## Coverage and limitations

No visual behavior is qualified by this engine feature. External UI projects
must verify their own controls and rendering. Remote transports are outside
this contract. Historical HTML/JSON evidence below applies to protocol 1.

## Verification history

- [Engine-only migration](../verification/runs/2026-10-28-headless.md) — protocol 2; see run-specific results.


- [2026-10-27 repository cleanup](../verification/runs/2026-10-27-repo-cleanup.md) — removed workflow-specific examples/tests, corrected historical crash guidance, and verified the explicit-path smoke driver against a fresh Wails package.

- [2026-10-27 starting window size](../verification/runs/2026-10-27-window-size.md) — native Mac defaults, custom and minimum sizes, single-axis flags, reopen inheritance/override, and JSON form submission.

- [2026-10-27 contract and AFK checks](../verification/runs/2026-10-27-afk-contract.md) — current Wails bridge checks, older skill compatibility and completed two-hour idle/sleep-wake tests.

- [2026-10-26 offline discovery](../verification/runs/2026-10-26-discovery.md) — installed help, examples, update/rollback and native throwaway reply.

- [2026-10-25 Wails integration](../verification/runs/2026-10-25-wails-integration.md) — normal packaged launches, native input and recovery, artifact return, and updates with an active UI.

- [2026-10-25 Wails follow-up](../verification/runs/2026-10-25-wails-followup.md) — accepted simpler requirements, confirmed Linux custom-scheme media loader failure, and local full-host-exit draft recovery.

- [2026-10-25 Wails assessment](../verification/runs/2026-10-25-wails.md) — historical replacement-equivalence proof. Later follow-up and integration runs resolved its distribution work and recorded the accepted native-access, storage and Linux-media limits.

- [2026-10-22/23 Windows desktop handoff](../verification/runs/windows-2026-10-22-203907.md) — agent-operated native input, visible draft restoration and comparison PNG/hash passed with an empty child PATH; the full renderer fixture remained unexecuted because symlink creation lacked privilege. Aggregate Untested.

- [2026-10-22 verified baselines](../verification/runs/2026-10-22-baselines.md) — unchanged authoring API, explicit WebGL availability and diagnosable unsupported startup.

- [2026-10-22 hosted checks](../verification/runs/2026-10-22-hosted.md) — real Electron execution, Windows pipe repair and Linux sandbox/graphics limits.

- [2026-10-22 platform verification](../verification/runs/2026-10-22-platform.md) — call ordering, rendering primitives, saved/exported images, rich-form retry and a separate comparison example.

- [2026-10-20 rich jstack demo](../verification/runs/2026-10-20-rich-grill.md) — inline previews, native draft checks and a completed user questionnaire with region feedback and exact-origin return.

- [2026-10-20 skill proof](../verification/runs/2026-10-20-skill-proof.md) — historical grill-form port, native draft/image recovery, owner loss, retry and cancellation.
- [2026-10-20 first proof](../verification/runs/2026-10-20-first-proof.md) — includes asset-loader and clipboard-fixture corrections, native artifact submission, and crash recovery.
- [2026-10-20 no-setup compatibility](../verification/runs/2026-10-20-no-setup.md) — native form/image submission and originating-harness consumption.
