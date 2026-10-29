---
title: Author a Fungi helper
updated: 2026-10-28
---
# Author a Fungi helper

Other projects provide all UI implementations and runtimes. This repository
owns the headless session engine, protocol examples and engine verification.
Workflow-skill integration lives in `fungi-skills`.

Use `fungi help helper` for the authoritative protocol 2 reference and
`fungi help examples` for a complete Python standard-library helper. The source
also includes `examples/helper/main.go`; it exercises the protocol without UI.

A helper reads one launch object from stdin, initializes itself, and writes
`{"type":"ready"}` to stdout. It can then send optional draft snapshots and one
terminal result, waiting for each matching acknowledgment. Use stderr for logs.
Input and result schemas belong to the helper. Persist binary data by writing
files under the supplied output directory and sending relative artifact paths.

```sh
fungi open --exec /absolute/helper --input /absolute/input.json --delivery auto -- helper-argument
```

Use absolute paths for helper assets that must survive working-directory
changes. Keep the executable and assets available for active sessions and
explicit reopen. Fungi restores acknowledged drafts and durable attachments;
helpers implement their own UI state restoration and retries.

Follow the receipt: callback/manual completion permits one collection, while a
supported host wait returns the result directly. Never poll a human. Distinguish
submitted, cancelled, closed and failed outcomes before interpreting data.

Direct HTML/JSON launches, the Fungi page bridge and engine window flags were
removed in protocol 2. See `fungi help migration`. UI projects may provide
whatever page API and controls they need above the external helper contract.

[Session contract](contract.md) · [External interfaces](features/interfaces.md) ·
[Verification setup](verification/setup.md)
