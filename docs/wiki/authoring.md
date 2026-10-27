---
title: Author a Fungi interface
updated: 2026-10-27
---
# Author a Fungi interface

The maintained user-facing contract now ships with Fungi. Run `fungi help`
for capabilities and common starting points, `fungi help html` for the page API,
and `fungi help delivery` for launch, completion, results and explicit recovery.
These commands work offline without adjacent docs or a renderer; the same
Markdown files live in the package's `docs/manual/` directory.

`fungi describe` stays machine-readable. Each topic has a `reference` argument
array; `fungi describe examples` locates the active package's examples from any
working directory. Copy an example directory before adapting it. The small
reply flow has no draft; comparison demonstrates PNG return and optional
recovery. Their full input → launch → response flows are in `fungi help examples`.

A consuming skill owns methodology, UI choices, input schema and result
interpretation. It should include its prepared UI and repeatable launch/collect
instructions. Fungi needs no registration, plugin manifest or second procedural
skill. Libraries remain local to the skill; examples do not limit composition.

The Compatibility section of `fungi help html` defines the protocol 1 promise:
documented page methods, session commands and message names/types/meanings are
preserved; additive metadata is allowed. Read advertised limits, ignore unknown
metadata, and check protocol/kind/outcome before interpreting a response. The
promise excludes Wails internals, on-disk layout, help prose and the experimental
flat form schema. It does not require a plugin SDK or a new authoring layer.

The [interface feature](features/interfaces.md) records expected rendering
behaviour; [discovery](features/discovery.md) records the offline entry points.
The [platform record](platforms.md) retains build-specific verification and gaps.
Workflow skills and jstack integration are developed in the separate
`jstack-fungi-test` project. This repository maintains the engine, public
contract and generic examples.
