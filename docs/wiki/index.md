---
title: Fungi
updated: 2026-10-25
---
# Fungi

Fungi is a prototype portable interaction engine for skill-driven agents. Skills define workflows; fungi runs their interaction sessions and returns human responses to the originating agent.

This wiki records the accepted design, implementation and verification. A Go core and Wails window helper now exist. The latest decision requires no harness setup, with manual fallback on unsupported versions. Native return checks and remaining limits are recorded explicitly.

Wails is now the default development renderer. Native packages use system
webviews and preserve the Go core and skill API. The
[initial assessment](verification/runs/2026-10-25-wails.md) and
[accepted limits](planning.md#accepted-wails-limits-and-next-milestone-on-2026-10-25)
explain the change; distribution integration verification is in progress.
The [first end-to-end proof](proposal.md#first-proof) records the original
Electron/Go scope and automatic-return checks.

- [Requirements and discussion](planning.md)
- [First implementation contract](contract.md)
- [Author a Fungi interface](authoring.md)
- [Verified environments and portable distribution](platforms.md)
- [Command installation and updates](features/installation.md)
- [Wails prototype: rendering fit, size and remaining gaps](verification/runs/2026-10-25-wails.md)
- [Wails follow-up: simpler requirements, Linux media diagnosis and manual recovery](verification/runs/2026-10-25-wails-followup.md)
- [Distribution proof and remaining release checks](verification/runs/2026-10-23-distribution.md)
- [Verified baseline checks and startup diagnostics](verification/runs/2026-10-22-baselines.md)
- [Windows 11 desktop handoff: desktop passed, symlink-dependent checks unexecuted](verification/runs/windows-2026-10-22-203907.md)
- [Feature map](features/index.md)
- [First prototype verification](verification/runs/2026-10-20-first-proof.md)
- [Current no-setup compatibility verification](verification/runs/2026-10-20-no-setup.md)
- [Skill authoring, recovery and longer-idle verification](verification/runs/2026-10-20-skill-proof.md)
- [Renderer flexibility and reliability verification](verification/runs/2026-10-22-platform.md)
- [Hosted platform checks and remaining failures](verification/runs/2026-10-22-hosted.md)
- [Verification setup](verification/setup.md)
- [Windows agent handoff: start here after pulling](verification/windows.md)
- [Research findings and recommendation](research/index.md)
- [Proposed small authoring interface](proposal.md)
- [First interview: submitted answers](journal/2026-10-19-grill-eb02aa23-b51f-42c9-b9c0-8fa331bfb02c.md)
- [Second interview: submitted answers](journal/2026-10-19-grill-186963ab-2f7b-423c-b40b-97b0afacc6fb.md)
- [Third interview: first milestone](journal/2026-10-19-grill-689a37f9-b8c4-4f2a-9e54-282234bc2551.md)
- [Fourth interview: no-setup delivery](journal/2026-10-20-grill-397cc4f4-31e2-4fd5-a08e-30ee197e3297.md)
- [Desktop verification arrangements](journal/2026-10-20-grill-desktop-arrangements.md)
- [Linux distribution decision: submitted answer](journal/2026-10-22-grill-0c10be35-a07a-4a8d-8d5d-39fda5b2e1bc.md)
- [Windows handoff request: submitted answer](journal/2026-10-22-grill-b8ce5a10-7261-48c1-9d80-ee54e103ec5e.md)

The project uses jstack's plan, wiki, grill-form, and rule-lean workflows. The wiki builder is development documentation tooling, not a proposed dependency of the delivered fungi binary.
