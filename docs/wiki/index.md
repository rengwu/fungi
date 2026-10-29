---
title: Fungi
updated: 2026-10-28
---
# Fungi

Fungi is a headless session engine for skill-driven agents. Go owns sessions,
durable drafts/results/artifacts, recovery and completion delivery. Other
projects supply all interfaces, executable helpers and their runtimes.

The 0.2 development line uses protocol 2 and engine package format 3. Previously
published 0.1 alpha releases contain the old UI. No new release is implied by
local source changes. Workflow-skill work lives in `fungi-skills`.

## Current reference

- [Current decisions and remaining work](planning.md#current-decisions)
- [Author a Fungi helper](authoring.md)
- [Session engine contract](contract.md)
- [Feature map](features/index.md)
- [Verified environments and portable distribution](platforms.md)
- [Command installation and updates](features/installation.md)
- [Verification setup](verification/setup.md)
- [Portable desktop verification](verification/desktop.md)
- [Windows verification handoff](verification/windows.md)
- [Project icon](icon.md)

## Recent verification

- [Engine-only distribution and protocol migration](verification/runs/2026-10-28-headless.md)

- [Callback across an intervening chat message](verification/runs/2026-10-27-interleaved-message.md)

- [Repository cleanup and verification instructions](verification/runs/2026-10-27-repo-cleanup.md)
- [Starting window dimensions](verification/runs/2026-10-27-window-size.md)
- [Authoring contract and AFK verification](verification/runs/2026-10-27-afk-contract.md)
- [Concise CLI documentation](verification/runs/2026-10-26-concise-help.md)
- [Terminal help and alpha.2 installation](verification/runs/2026-10-26-terminal-help.md)
- [First public alpha distribution](verification/runs/2026-10-26-release.md)
- [Stale-artifact and Electron cleanup](verification/runs/2026-10-26-cleanup.md)
- [Wails integration: default packages, runtime setup and updates](verification/runs/2026-10-25-wails-integration.md)

Each report identifies its tested build, scope and remaining gaps. Hosted checks
and local passes do not qualify clean consumer machines or every harness mode.

## Historical verification

These reports retain their original renderer and environment scope.

- [Wails prototype: rendering fit, size and initial gaps](verification/runs/2026-10-25-wails.md)
- [Wails follow-up: Linux media and manual recovery](verification/runs/2026-10-25-wails-followup.md)
- [Electron distribution proof](verification/runs/2026-10-23-distribution.md)
- [Electron baseline checks and startup diagnostics](verification/runs/2026-10-22-baselines.md)
- [Windows 11 Electron desktop handoff](verification/runs/windows-2026-10-22-203907.md)
- [First prototype verification](verification/runs/2026-10-20-first-proof.md)
- [October 20 no-setup compatibility](verification/runs/2026-10-20-no-setup.md)
- [Historical skill authoring, recovery and longer-idle proof](verification/runs/2026-10-20-skill-proof.md)
- [Electron renderer flexibility and reliability](verification/runs/2026-10-22-platform.md)
- [Electron hosted platform checks](verification/runs/2026-10-22-hosted.md)

## Decision and research history

- [Decision chronology](planning.md#decision-history)
- [Research findings and dated recommendations](research/index.md)
- [Original authoring proposal and first-proof plan](proposal.md)
- [Documentation interview](journal/2026-10-25-grill-e3b1de69-4014-457d-b32b-d23570cce284.md)
- [First interview](journal/2026-10-19-grill-eb02aa23-b51f-42c9-b9c0-8fa331bfb02c.md)
- [Second interview](journal/2026-10-19-grill-186963ab-2f7b-423c-b40b-97b0afacc6fb.md)
- [Third interview: first milestone](journal/2026-10-19-grill-689a37f9-b8c4-4f2a-9e54-282234bc2551.md)
- [Fourth interview: no-setup delivery](journal/2026-10-20-grill-397cc4f4-31e2-4fd5-a08e-30ee197e3297.md)
- [Desktop verification arrangements](journal/2026-10-20-grill-desktop-arrangements.md)
- [Linux distribution decision](journal/2026-10-22-grill-0c10be35-a07a-4a8d-8d5d-39fda5b2e1bc.md)
- [Windows handoff request](journal/2026-10-22-grill-b8ce5a10-7261-48c1-9d80-ee54e103ec5e.md)

The project uses jstack's planning and wiki workflows. The wiki builder is
development documentation tooling and is not shipped in the Fungi package.
