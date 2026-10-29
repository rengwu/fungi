# Migration to the headless engine

Fungi 0.2 development uses protocol 2 and engine package format 3. This is an
explicit breaking change from the protocol-1 bundled UI releases. No 0.2 public
release is implied by a development build.

## Interface ownership

Other projects own all UI implementations, rendering, assets and runtimes.
Fungi no longer supplies Wails, webviews, HTML hosting, a JavaScript page bridge,
JSON forms or UI examples. There is no optional renderer maintained here.

Replace direct page/form launches with an external executable helper:

```sh
fungi open --exec /absolute/helper --input /absolute/input.json --delivery auto -- /absolute/page.html
```

The final page argument is illustrative: its meaning belongs to that helper.
The engine neither opens nor validates it. `--width`, `--height` and
`FUNGI_RENDERER_DIR` have been removed. UI projects supply their own page APIs
and window controls.

## Helper protocol

The launch object now advertises `protocol:2`. Its unused `source` field is
removed; pass helper-specific paths in arguments or input. Readiness,
draft/result messages, acknowledgments, delivery, limits and artifact records
retain their existing shapes. Helpers must explicitly support protocol 2;
ignore unknown optional metadata in supported messages. See `fungi help helper`.

`doctor` now checks engine storage and file locking. Installer runtime-setup
flags and the renderer fields in package manifests have been removed.
`describe examples` points to the embedded helper example, with no installed
UI paths or launch arrays.

## Existing sessions and installations

Use the original protocol-1 executable and original state directory to collect
or reopen protocol-1 sessions. The new engine rejects them without rewriting
their files. Preserve their renderer, helpers and assets while still needed.
The existing on-disk records are not automatically converted.

Test engine packages in a separate installation directory first. The old updater
cannot install format 3; run the new extracted engine's `install` command.
Engine update/rollback accepts format 3 only. Keep the previous UI installation
separate if you need to return to it. No release assets are silently replaced.
