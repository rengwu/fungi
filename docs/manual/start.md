# Fungi — headless sessions for agent workflows

Fungi runs external helpers, saves JSON responses and file artifacts, and
returns completion to the calling agent. Other projects provide all interfaces
and helper runtimes. The engine needs no webview or graphical desktop.

## Start

```sh
fungi open --exec ./helper --input ./input.json --delivery auto
fungi describe
fungi help helper
fungi doctor
```

Pass helper arguments after `--`. Input and result schemas belong to the helper.
The standalone executable includes this manual; no adjacent files are required.
For a complete runnable protocol example, use `fungi help examples`.

## Completion

Inspect the receipt: `auto` may fall back to `manual`.

| Delivery | Action |
| --- | --- |
| `callback` | End the turn; collect once when notified. |
| `manual` | Ask the user to say “done”; end the turn, then collect once. |
| `wait` | Use only if the host keeps one call pending or automatically delivers completion. Read the returned result directly. |

Use the receipt's collection command and environment. Never poll `collect`.
Check `outcome`: `submitted`, `cancelled`, `closed` or `failed`. Draft saving
is optional; submitted results are saved independently.

## Reference

| Command | Reference |
| --- | --- |
| `fungi help helper` | Protocol 2 messages, acknowledgments and limits |
| `fungi help delivery` | Completion, results and explicit recovery |
| `fungi help cli` | Commands, storage and exit codes |
| `fungi help examples` | Complete headless helper example |
| `fungi help installation` | Engine packages, updates and rollback |
| `fungi help migration` | Moving from protocol 1 and bundled UI packages |
