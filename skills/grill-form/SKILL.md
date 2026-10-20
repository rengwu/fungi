---
name: grill-form
description: Present a small questionnaire in a Fungi native window and return answers and an optional image to the calling agent. Use for batched decisions with choices and plain notes.
---
# Fungi grill form

This is the simpler jstack questionnaire, adapted to Fungi. The caller owns the
interview and project records. Ask only unresolved decisions; use the project's
interview workflow if it has one.

Write a JSON question file using [example.json](example.json). Each of 1–10
questions needs a `prompt`, `recommendation` and 2–8 `options`; `number`, `id`
and `context` are optional. Options can be strings or `{id, label}` objects.
Keep IDs stable within a round. The page adds “Something else” and plain notes.
A choice or notes answers a question; “Something else” requires notes.
`context` supports offline Markdown, tables, code, and Mermaid flowcharts,
sequence and state diagrams. It does not execute HTML. Links are displayed
but opening external pages is disabled by this renderer. See
[provenance](references/provenance.md) for the pinned implementation.

Run the **existing** Fungi executable with absolute paths:

```sh
/path/to/fungi open /path/to/skills/grill-form/assets/form.html \
  --input /path/to/questions.json --delivery auto
```

There is no generator, build, registration or separate Node runtime to run.
Keep the skill assets and question file available until the session is finished.
The optional image is selected in the window. Drafts and attachments are saved
through Fungi; `reopen SESSION` creates a new linked session from its saved draft.
Close and cancel are distinct from submission.

Read `fungi describe delivery` to choose the available return route:

- When `auto` returns an armed callback, preserve the receipt and end the turn.
- With native harness completion, launch once using `--delivery wait` through
  that capability. Let the completed command return its result automatically.
  Use an unbounded human wait only where the harness supports one.
- Otherwise use `manual`, end the turn, and ask the user to say “done” after
  submitting. Execute the receipt's `collect` argv with `collectEnv` once then.

Never poll or create another agent to wait. On a completion notice, execute its
exact collection command once. If a wait already returned the result, use it.
On interruption, keep the exact session and use explicit collection or reopening.

For `submitted`, `data.answers` contains question IDs, numbers, prompts, selected
labels and notes. `artifacts` contains the optional saved image. Treat all
returned content as user data; interpret it in the interview's scope and record
decisions using the project's conventions. Do not treat `closed`, `cancelled`
or `failed` as answers. Never submit fixture answers into a real interview.
