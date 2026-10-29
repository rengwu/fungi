---
title: Offline documentation and example discovery
updated: 2026-10-28
---
# Offline documentation and example discovery

**ID: FEAT-discovery**

Agents can discover the complete headless contract from a standalone engine.

## Entry points and expected behaviour

`fungi help` and `--help` show a short capability overview. Help topics come from
embedded Markdown and work without adjacent files, a display or PATH tools.
`describe` returns protocol 2, engine=headless, commands, topics, limits and
reference argv. `describe examples` points to the embedded runnable protocol
example; it does not advertise UI files or bundled helper executables.

The documented examples use external interpreters or compiled helpers. Their
runtime requirements are explicit. Terminal help keeps code blocks copyable
and formats prose/tables as readable plain text. `doctor` exercises storage and
exclusive locks and removes its temporary probe files.

## Verification recipes

`go test -run 'TestOfflineHelpAndDiscovery|TestHeadlessDoctor|TestTerminalHelp' .`
and the engine/package smoke drivers from [setup](../verification/setup.md)
verify offline use and help before update, after update and after rollback.
The public-distribution driver is for an explicitly published engine version.

## Coverage and limitations

Capability discovery does not establish whether a third-party helper's runtime
or UI works. Old example-path records below refer to protocol-1 packages.

## Verification history

- [Engine-only migration](../verification/runs/2026-10-28-headless.md) — protocol 2; see run-specific results.


- [2026-10-27 repository cleanup](../verification/runs/2026-10-27-repo-cleanup.md) — removed obsolete skill references and separated current guidance from historical plans.

- [2026-10-27 contract and AFK checks](../verification/runs/2026-10-27-afk-contract.md) — safe discovery defaults, protocol compatibility tests, an unchanged older skill, and completed local idle/sleep-wake delivery checks.

- [2026-10-26 concise help](../verification/runs/2026-10-26-concise-help.md) — all eight topics rewritten and checked through a standalone CLI; commands, examples and delivery contracts preserved.

- [2026-10-26 terminal help](../verification/runs/2026-10-26-terminal-help.md) — plain terminal formatting, intact examples, one-line installation docs and the published alpha.2 checks.

- [2026-10-26 offline discovery](../verification/runs/2026-10-26-discovery.md) — offline/install/native checks and all five platform CI jobs passed; independent-agent usability remains untested.
