# Experimental flat JSON form

Use this renderer for simple forms, or HTML for custom layouts. The schema
requires a nonempty `title` string and `fields` array; `context` is optional text.

```json
{"title":"Timing feedback","context":"Try a duration, then comment.","fields":[
  {"name":"duration","type":"range","label":"Milliseconds","min":100,"max":1000,"step":50,"value":300},
  {"name":"notes","type":"textarea","label":"Notes","required":true},
  {"name":"reference","type":"image","label":"Optional reference"}
]}
```

```sh
fungi open ./form.json --input ./values.json --delivery auto
```

Optional `--width` and `--height` flags set the initial window dimensions;
see `fungi help cli` for defaults and limits.

## Fields and values

Fields require a unique nonempty `name` and `type`. Optional properties:
`label`, boolean `required`, `value`. Names `__proto__`, `constructor` and
`prototype` are reserved. Unknown schema properties are rejected.

| Type | Values |
| --- | --- |
| `text`, `textarea` | String defaults and input |
| `range` | Numeric `min < max`; positive `step` defaults to 1. Values must be in bounds. |
| `image` | Default absent or null. Images come from selection, paste or restored drafts, not input paths. |

Initial non-image values use draft > input property > field value > empty
(range minimum for ranges). Input must match the field type. This renderer
autosaves; custom HTML manages its own drafts.

Submitted `data` maps names to strings, numbers, or `{"artifact":"FIELD_NAME"}`
image references (null if absent). Image artifacts use the field name as key.
Example: `{"duration":250,"notes":"Faster feels better","reference":null}`.
Completion and errors follow `fungi help delivery`.
