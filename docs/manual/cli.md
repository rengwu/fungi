# CLI reference

```text
fungi [--help]
fungi help [start|examples|html|form|helper|delivery|cli|installation]
fungi describe [html|form|helper|delivery|installation|examples]
fungi open SOURCE [--input FILE] [--delivery auto|wait|manual]
fungi open --exec HELPER [--input FILE] [--delivery MODE] -- [ARGS...]
fungi collect SESSION
fungi reopen SESSION [--delivery MODE]
fungi --version
fungi doctor
fungi install --archive FILE --sha256 HASH [--dir DIRECTORY]
fungi update [--version vX.Y.Z | --archive FILE --sha256 HASH]
fungi rollback
fungi prune
```

`SOURCE` is a local HTML file or a `.json` built-in form. `--exec` selects a
trusted executable instead; arguments after `--` go directly to it, without
a shell. `--input` reads JSON. Paths are resolved at launch. `SESSION` is the
exact 32-character ID in a receipt. `version` is an alias for `--version`. `help`/`--help` return formatted plain text for terminals; other public commands emit
JSON on stdout. `describe` returns protocol 1 metadata and reference routes;
`describe examples` locates adjacent example files from any working directory.
An unavailable example has `available:false` and no launch `argv`.

Exit codes: **0** success (also `pending`, `cancelled`, `closed`); **1** operational
failure or failed session; **2** invalid input/arguments; **3** specifically
`doctor` with a missing runtime. Always inspect the JSON kind/outcome. Read
`fungi help delivery` before choosing a waiting mechanism.

## Environment and storage

| Variable | Purpose |
| --- | --- |
| `FUNGI_STATE_DIR` | Overrides session storage; use a directory outside a managed installation. Receipts preserve its absolute location in `collectEnv`. |
| `CODEX_THREAD_ID` | Inherited exact conversation ID used by automatic callback detection. Do not fabricate it. |
| `FUNGI_RENDERER_DIR` | Development override for the directory containing the native helper. Normally omit. |

Default sessions live under the OS user config directory: `~/Library/Application
Support/fungi/sessions` on macOS, `%AppData%\fungi\sessions` on Windows, and
`$XDG_CONFIG_HOME/fungi/sessions` or `~/.config/fungi/sessions` on Linux.
Each session retains input, drafts, result, artifacts and a worker diagnostic
log. Collection does not delete them; prune removes eligible installed
versions, not sessions. Keep files needed by active/recoverable sessions.

No background service remains after a completed session. Internal worker
commands and on-disk control files are not authoring APIs.
