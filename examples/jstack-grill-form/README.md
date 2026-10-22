# Rich jstack questionnaire demo

This experiment uses the current sibling jstack skill's prepared HTML and
generator. It keeps its presentation and replaces browser downloads with
Fungi draft/submission calls. It does not modify jstack or the simpler bundled
Fungi skill.

Generate a temporary form with the updated jstack:

```sh
node ../jstack/skills/grill-form/scripts/grill-form.mjs create \
  --project "$PWD" --temporary < examples/jstack-grill-form/questions.json
```

Use the returned HTML path as `GENERATED_FORM`. Choose a new output path outside
the jstack session directory, so its collector still owns only its original
files. The adapter checks the expected template boundaries and stops if they
have changed.

```sh
node examples/jstack-grill-form/adapt.mjs GENERATED_FORM OUTPUT_HTML
/path/to/current/fungi open OUTPUT_HTML --delivery auto
```

Use a build containing inline-frame support. Only Fungi owns completion here:
do not start jstack's wait/completion commands or a second waiter. On the exact
Fungi notice, collect once with the returned argv and environment.

The six mock questions cover plain choices, visual cards, an animation slider
and easing selector, an open layout critique with region comments, Markdown
with a Pele diagram, and free text with optional pasted images. Partial
feedback is allowed. Mock responses are not project decisions.

Result data preserves jstack's session ID, specification digest, answer IDs,
prototype states and region comments. Inline images become Fungi artifacts:
each image keeps its numeric ID and note offset, and gains `artifact` and
`mimeType` fields instead of embedding base64 in result JSON. To validate with
jstack's model/collector, restore each image's `dataUrl` from its matching
artifact bytes before passing the submission to the original session.

The same mapping saves/restores drafts. Fungi's JSON and artifact limits still
apply; a jstack form that reaches its larger theoretical JSON limit may exceed
Fungi's limit. The UI reports persistence failure and retains the editable form.
This adapter is an example, not a new plugin framework or a complete parity
claim for every jstack feature.
