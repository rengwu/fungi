# Complete flows to adapt

`fungi describe examples` returns absolute `source`, optional `input`, and
launch `argv` for each available example. Execute arrays as arguments to preserve
spaces. Paths resolve from the active package, regardless of working directory.
Core-only builds provide help but report missing examples. Copy an example's
directory before editing; managed versions are immutable.

## 1. Response without drafts

The `reply` HTML page accepts:

```json
{"title":"Transition review","prompt":"Which transition needs adjusting, and how?"}
```

From the package directory (or use discovered absolute paths):

```sh
./fungi open examples/reply/index.html --input examples/reply/input.json --delivery auto
```

Use `./fungi.exe` on Windows. The page reads title/prompt from `fungi.ready()`
and offers notes, Send and Cancel. It awaits `fungi.submit({notes: notes.value})`,
prevents duplicate submissions and preserves text with a visible error on
failure. It saves no drafts.

For a callback receipt, end the turn and await notification. For manual delivery,
ask the user to say “done” and end the turn. On either cue, execute the receipt's
`collect` array once with `collectEnv`. Never poll. A submitted answer produces:

```json
{"protocol":1,"kind":"result","session":"…","outcome":"submitted","data":{"notes":"Shorten the opening fade"},"artifacts":[]}
```

Read `data.notes` only for `submitted`. Cancellation or window close is not
approval; reopening starts fresh. If the host supports a pending call or
automatic background completion, `--delivery wait` returns the result directly.
See `fungi help delivery` for outcomes, artifacts and recovery.

## 2. Canvas comparison returning an image

The `compare` page includes a canvas, divider, preference, notes and two local
SVG images. Optional input overrides title, context and relative `before`/`after`
image URLs. Run its discovered `argv`, or from the package:

```sh
./fungi open examples/compare/index.html --delivery auto
```

Move the divider, choose a preference and submit. The result contains preference,
divider and notes in `data`, plus a PNG artifact with key `comparison`:

```json
{"key":"comparison","name":"comparison.png","mediaType":"image/png","path":"/absolute/saved/image","bytes":1234,"sha256":"…","id":"…"}
```

Inspect the image at the returned path; values above are illustrative.
Collect as in example 1. This page saves drafts; explicit reopening restores the
last acknowledged selection. Its adjacent README documents input fields.

## Integrate with a skill

Keep the reusable page with the skill. Its instructions should cover:

1. Methodology: what to show and ask, available controls/visuals, and a meaningful
   input example.

2. Invocation: write the round's input JSON and execute the following with
   argument-safe tool calls:

```sh
fungi open ABSOLUTE_SKILL_UI --input ABSOLUTE_INPUT --delivery auto
```

3. Completion: follow the receipt's callback or manual route; collect once on
   notification or “done.” Never poll.

4. Interpretation: check the outcome, interpret `submitted` data using the
   skill's schema, and inspect image artifacts.

Implement the integration once; agents follow the skill without another Fungi
skill or per-round adapter. Use `fungi help html` when extending the UI.
Sliders, drawings, charts and multi-step interactions use the same bridge;
these examples do not limit possible workflows. Keep frequent interactions
inside the page; send the agent the finished response.
