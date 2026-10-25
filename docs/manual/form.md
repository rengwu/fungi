# Experimental flat JSON form

For custom composition use HTML. This small built-in renderer accepts a
`title` string, optional `context` string and `fields` array:

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

Each field requires unique nonempty `name` and `type` (`text`, `textarea`,
`range`, `image`). Optional common fields: `label`, boolean `required`, `value`.
Names `__proto__`, `constructor` and `prototype` are reserved. Unknown schema
properties are rejected. Text defaults are strings. Ranges require numeric
`min < max`; optional positive `step` defaults to 1 and `value` must be within
bounds. Image defaults must be absent or null.

Initial non-image values use saved draft > input property by name > field
value > empty (range minimum for ranges). Provide input values of the same
types. Images come from file selection/paste or restored draft bytes, not a
path in input. This renderer autosaves; custom HTML controls its own drafts.

Submitted `data` maps names to text strings, range numbers and image
`{"artifact":"FIELD_NAME"}` references (or null). Image artifacts use the
field name as key. With the example, `data` might be
`{"duration":250,"notes":"Faster feels better","reference":null}`.
Launch receipts, completion and error handling are the same as HTML;
see `fungi help delivery`.
