# CLI reference

```text
fungi [--help]
fungi help [start|helper|delivery|cli|examples|installation|migration]
fungi describe [helper|delivery|installation|examples|migration]
fungi open --exec HELPER [--input FILE] [--delivery auto|wait|manual] -- [ARGS...]
fungi collect SESSION
fungi reopen SESSION [--delivery auto|wait|manual]
fungi --version
fungi doctor
fungi install --archive FILE --sha256 HASH [--dir DIRECTORY]
fungi update [--version vX.Y.Z | --archive FILE --sha256 HASH]
fungi rollback
fungi prune
```

`--exec` launches a trusted executable without a shell. Arguments after `--`
are passed literally. Helper names resolve through PATH; resolved executable
paths and input JSON are captured at launch. Use absolute helper/asset paths
when working directories may change. `SESSION` is the receipt's exact
32-character ID. `version` aliases `--version`.

All interfaces, assets, layout and window sizing belong to the external helper.
Direct HTML/JSON-form sources and engine window flags are unsupported.

Help prints terminal text; other public commands emit protocol 2 JSON.
`describe TOPIC` includes its reference command. `describe examples` points to
the embedded runnable example; no UI or example executables are installed.

`doctor` checks atomic reads/writes and exclusive locking in the session
storage directory. It removes its temporary probe files, launches no helper,
and does not validate the runtime of an external helper.

## Exit codes

| Code | Meaning |
| --- | --- |
| 0 | Success, including `pending`, `cancelled` and `closed` |
| 1 | Operational failure, incompatible session or failed helper session |
| 2 | Invalid input or arguments |

Inspect JSON `kind` and `outcome`; exit zero alone does not mean submission.

## Environment and storage

| Variable | Purpose |
| --- | --- |
| `FUNGI_STATE_DIR` | Session storage override, outside managed installs; its absolute path is returned in `collectEnv`. |
| `CODEX_THREAD_ID` | Inherited exact conversation ID for optional callback detection. Do not fabricate. |

Default session directories:

| Platform | Path |
| --- | --- |
| macOS | `~/Library/Application Support/fungi/sessions` |
| Windows | `%AppData%\fungi\sessions` |
| Linux | `$XDG_CONFIG_HOME/fungi/sessions`, or `~/.config/fungi/sessions` |

Sessions retain input, drafts, results, artifacts and a worker log. `collect`
preserves them. `prune` removes eligible installed versions, not session data.
Keep external helpers and assets needed by active or recoverable sessions.
No background service remains after completion. Internal worker commands and
storage layout are implementation details.
