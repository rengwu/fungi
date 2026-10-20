---
title: Proposed small authoring interface
updated: 2026-10-20
---
# Proposed small authoring interface

This earlier illustrative design is retained as discussion history. For the implemented commands and JavaScript methods, use the [prototype contract](contract.md) and CLI discovery. The first interview accepted both HTML and JSON authoring in the first release, bundled rendering, adjacent optional capabilities, and manual “done” fallback when automatic delivery cannot be established. The interface details below remain proposals.

The [implementation contract draft](contract.md) now defines the proposed command, result, helper and recovery semantics. This page keeps the short authoring examples and milestone rationale; use the contract for precise behavior.

The user has selected Electron + Go: Go owns the session engine and CLI; a small Electron companion owns the native window and web renderer. Packaging is a portable folder/archive with one fungi CLI entry point and bundled runtime files, following the user's instruction to choose the simplest route. The [third interview](planning.md#accepted-first-milestone-scope) settles the first milestone's scope. Exact APIs and internal implementation sequence remain proposals.

## Author a file, then open it

A basic skill should need only a file like this:

```html
<form id="review">
  <label>Animation duration
    <input name="duration" type="range" min="100" max="1000" step="10">
  </label>
  <label>Notes <textarea name="notes"></textarea></label>
  <button>Submit</button>
</form>
<script type="module">
  const { input } = await fungi.ready();
  const form = document.querySelector('#review');
  form.elements.duration.value = input.duration ?? 300;
  form.onsubmit = async event => {
    event.preventDefault();
    await fungi.submit({
      duration: Number(form.elements.duration.value),
      notes: form.elements.notes.value
    });
  };
</script>
```

The renderer supplies the small bridge. An illustrative invocation is:

```sh
fungi open ./review.html --input ./input.json
```

The agent should get either a supported automatic-delivery acknowledgement with an exact session identifier, a final result through a host-managed pending call, or a notice that the UI is open and the user should say “done” after submitting. It must not get a process ID accompanied by instructions to keep polling. The same exact-session collection path works for automatic notices and manual continuation.

A reference component library can shorten common forms later. The plain HTML path remains valid so novel interfaces never wait for fungi to invent another widget type.

## JSON authoring in the same first release

Support a small declarative form through the same renderer and result bridge. An illustrative file is:

```json
{
  "title": "Animation feedback",
  "fields": [
    { "type": "range", "name": "duration", "label": "Duration", "min": 100, "max": 1000, "value": 300 },
    { "type": "text", "name": "notes", "label": "Notes" }
  ]
}
```

The corresponding invocation could be `fungi open ./review.json`. The renderer turns this into components, collects named values, and submits the same JSON result a hand-authored page would submit. Common fields should not require JavaScript from the skill author. The accepted initial scope is flat built-in fields. Nested layouts and custom JSON component registration are deferred for the user's larger JSON-authoring discussion; do not design that framework into this milestone.

Equal first-release support means both paths are documented, usable, and exercised. It does not require a JSON representation for every possible JavaScript program. Custom components can own their local behavior; the engine should not gain an expression evaluator, workflow DSL, or separate catalog-specific lifecycle.

## A small bridge

| Capability | Meaning |
| --- | --- |
| `ready()` | Read launch input, restored draft if present, and available renderer capabilities |
| `submit(data, attachments?)` | Commit JSON and optional blobs; resolve only after fungi has saved them |
| `cancel()` | Explicitly finish without submitted data |
| `saveDraft(data)` | Optional recovery support for custom interfaces |

Do not expose an arbitrary shell, filesystem, or harness API to page content. A canvas can submit a PNG blob and its editable scene JSON together; fungi supplies artifact references in the returned result. A slider can animate a local preview continuously while the model remains idle.

The native window's close action is handled by the renderer/engine, not by requiring every skill to implement fragile browser-unload behavior. Accepted behavior: report a distinct `closed` outcome, preserve the latest saved draft, and allow explicit reopen. Close must never imply submission. Draft preservation is limited to data already saved; custom pages opt into `saveDraft`.

## Growing an interface

A reusable capability may simply be an `index.html` and local assets. No manifest is needed for the simplest file. Optional metadata can add a name, input schema, output schema, required capabilities, and documentation when discovery makes that useful.

Authors may use plain HTML, Web Components, React, Vue, Svelte, or a canvas library. Prebuilt JavaScript and CSS can ship with a skill or renderer. A library's build-time npm dependency does not require npm on the end user's machine. Ad hoc interfaces should work without a build step.

Actual framework support depends on the chosen bundled engine. Do not advertise a browser library as supported merely because an engine parses HTML and executes JavaScript.

Explicit file paths should always work. Reusable capabilities can live with the skill or in an optional folder near the fungi executable. Named discovery may use a short descriptor when needed; no registry server, installation command, or dependency solver is required. The exact folder name and lookup rules remain implementation choices.

Use composition at the right level: components inside one web renderer can share values and events; a native screenshot panel and a web form compose as workflow steps. Arbitrary third-party native widgets cannot be promised to share one layout merely because they both use fungi.

## Native escape hatch

A native extension should be an ordinary executable: accept input, display its own interface, emit a result. It can be written in any language. The advanced process contract should not be required for someone making a basic HTML form.

Package the Go core and Electron renderer together in a portable directory. A native capability may be a companion executable; basic HTML/JSON authors still need no compiler or runtime installation. Keep renderer launch behind the same small process boundary instead of requiring a Go binding for every engine. No self-extraction mechanism is required for the accepted distribution.

## First proof

**Outcome:** an agent opens a skill-authored interface, the human takes as long as needed, and the exact calling conversation receives a saved response without repeated model tool calls. The accepted platform scope is macOS first. The accepted harness targets are **Codex, Claude Code, OpenCode, Grok, Kimi and Pi**. A successful Codex test alone does not establish completion or universal compatibility.

Inspect and exercise the simplest supported return route in each harness early. A host-managed background completion, exact-session callback, or one pending call can satisfy the waiting requirement when its actual behavior is verified. A process ID with instructions to poll cannot. Keep any route that needs extra setup or cannot be verified visible as a gap; the manual “done” fallback remains usable but is not an automatic-return pass.

1. **Make the small contract concrete.** Use two equivalent examples, plain HTML and a JSON form, to settle launch input, the page bridge, result/artifact shape, and exact-session collection. The result distinguishes submitted, cancelled, closed and failed outcomes. Keep caller identity and delivery configuration outside page-controlled data. An `open` acknowledgement must identify either the supported automatic route or the accepted manual “done” fallback. Help and machine-readable command output must make that distinction discoverable.
2. **Build one complete round trip.** Go starts the session and Electron host; the host renders text, one field, one slider, and submit/cancel, with a pasted/selected image attachment. Both HTML and flat JSON examples use the same submission path. Save the response before acknowledging submission or sending a completion notice. Include a tiny native helper accepting JSON and returning JSON/artifact data through the same engine contract. Implement and verify the return routes for all six named harnesses; use native host capabilities where they suffice. Preserve manual exact-session collection when automatic delivery is unavailable. Do not introduce an adapter registry or broad component system for this milestone.
3. **Verify waiting, failure and packaging.** Record the actual tool/event sequence during a human wait, measure idle resources, then submit and verify correct routing. Exercise two simultaneous sessions, duplicate submit, close/cancel, renderer crash, interrupted delivery and explicit recovery after process exit. Package a portable build and test offline launch with no separately installed Node/npm or browser runtime. Measure archive size and startup. Windows/Linux clean-environment packaging checks follow before claiming cross-platform delivery.

Use a deterministic fake renderer for core failure/routing tests so session correctness can be checked without opening windows. Use the real Electron UI and real harness for the user-visible return path; mocks cannot establish automatic resumption.

The milestone is complete when both authoring paths and the native helper return the expected JSON and artifact references, submitted results survive the exercised failures, outcomes remain distinct, sessions never cross, and each named harness has a verified automatic return path with an observed idle interval producing no recurring model checks. Record the interval, harness version and environment; a short test is not evidence of unlimited-duration reliability. Missing callbacks must leave a collectable result and the documented “done” fallback, while remaining an explicit automatic-delivery gap.

The first useful deliverable is this reliable round trip with tiny HTML, flat JSON and native-helper examples. Custom UI can use ordinary local HTML/assets without changing Go. Drawing tools, native capture UI, rich JSON composition, custom JSON component registration and external live updates are deferred. Broader UI libraries, external catalog protocols, MCP Apps hosting, remote desktops, and WASM can follow when a workflow needs them. The scope is accepted; implementation has not begun. [Research rationale](research/index.md)
