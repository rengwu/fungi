# Fungi — local interfaces for agent workflows

Skills define workflows and interfaces. Fungi opens them, manages sessions,
and saves JSON responses and file artifacts for the calling agent.

Use local HTML/CSS/JavaScript and libraries, or a native executable helper.
Web interfaces support forms, live controls, custom UI, diagrams, Canvas/SVG,
WebGL, workers and WebAssembly. Graphics and media support depend on the system
webview.

## Start

```sh
fungi open ./ui/index.html --input ./input.json --delivery auto
fungi describe           # JSON capabilities and reference commands
fungi describe examples  # Installed example paths and launch arguments
```

## Completion

Check the launch receipt: `auto` may fall back to `manual`.

| Delivery | Action |
| --- | --- |
| `callback` | End the turn; collect once when notified. |
| `manual` | Ask the user to say “done”; end the turn, then collect once. |
| `wait` | Use only if the host keeps one call pending or automatically delivers its completion. Read the returned result directly. |

Use the receipt's collection command and environment. Never poll `collect`.
Check the result's outcome: `submitted`, `cancelled`, `closed` or `failed`.
Draft saving is optional; submitted results are saved independently.

## Reference

| Command | Reference |
| --- | --- |
| `fungi help html` | HTML page API and rendering limits |
| `fungi help delivery` | Completion, results, attachments and recovery |
| `fungi help cli` | Commands, arguments, storage and exit codes |
| `fungi help examples` | Complete workflows to adapt |
| `fungi help form` | Experimental JSON forms |
| `fungi help helper` | Native helper protocol |
| `fungi help installation` | Installation, updates and runtime requirements |

Help works offline without a renderer. The same Markdown ships in `docs/manual/`.

Read the relevant references when authoring or changing a skill integration.
For an existing integration, follow the skill's invocation instructions.
