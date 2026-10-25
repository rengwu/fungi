# Complete flows to adapt

Run `fungi describe examples`. Each available item provides an absolute `source`,
optional `input`, and launch `argv`; execute arrays as arguments, preserving
spaces. Paths refer to the active package, independent of your working directory.
A core-only build can print all help topics but reports missing examples.
Copy a page's directory before editing; managed package versions are immutable.

## 1. Parameterised response without drafts

The `reply` example is a complete ordinary HTML page. Its input is:

```json
{"title":"Transition review","prompt":"Which transition needs adjusting, and how?"}
```

From the portable package directory (or substitute discovered absolute paths):

```sh
./fungi open examples/reply/index.html --input examples/reply/input.json --delivery auto
```

On Windows, use `./fungi.exe`. The page calls `fungi.ready()`, shows the input
title/prompt, and lets the user type notes and Send or Cancel. It awaits
`fungi.submit({notes: notes.value})`, prevents duplicate submissions, and
preserves text with a visible error if submission fails. It does not save drafts.

The agent follows the receipt: callback → end turn; manual → ask for “done”
and end turn. On completion, execute its `collect` array with `collectEnv` once.
For an answer “Shorten the opening fade”, expect:

```json
{"protocol":1,"kind":"result","session":"…","outcome":"submitted","data":{"notes":"Shorten the opening fade"},"artifacts":[]}
```

Read `data.notes` only as a submitted answer. Cancel or closing the window is
not approval. A reopen starts fresh because this page saves no draft. A host
with automatic completion can use one `--delivery wait` invocation instead and
read the final result directly. See `fungi help delivery` for the exact contract.

## 2. Canvas comparison returning an image

The `compare` example owns its canvas, divider, preference and notes. Its
directory contains two local SVG images. Optional input can replace the title,
context and relative `before`/`after` image URLs. Launch its discovered `argv`,
or from the package:

```sh
./fungi open examples/compare/index.html --delivery auto
```

Move the divider, choose a preference and submit. The page exports PNG bytes
and sends them under artifact key `comparison`. The result includes the
preference, divider and notes in `data`, plus a file record:

```json
{"key":"comparison","name":"comparison.png","mediaType":"image/png","path":"/absolute/saved/image","bytes":1234,"sha256":"…","id":"…"}
```

Values above are illustrative; read the actual returned path and inspect the
image. Collection uses the same receipt flow as example 1. This page also saves
drafts: closing and explicitly reopening restores the last acknowledged
selection. Read its adjacent README for input details.

## Put the integration inside a consuming skill

Keep the reusable page with the skill. The skill's own instructions can be:

1. Apply the skill's methodology to decide what to show and ask. Explain its
   available controls/visuals; include an input example with meaningful content.
2. Write that round's input JSON. Run `fungi open ABSOLUTE_SKILL_UI --input
   ABSOLUTE_INPUT --delivery auto` using argument-safe tool execution.
3. Inspect the receipt; use callback or manual completion as above. Never poll.
4. On completion, collect once using the receipt. Interpret only `submitted`
   data according to the skill's documented schema; inspect any image artifacts.

Authors implement this integration once. Executing agents need the consuming
skill, not a second Fungi skill or a fresh adapter per round. Reference
`fungi help html` when extending the UI and `fungi help delivery` if the host's
return mechanism is unclear. Examples illustrate contracts, not creative limits:
sliders, drawings, charts and multi-step interactions can use the same bridge.
