# Fungi engine package

This package contains one Fungi engine executable, its offline manual, Go
license and build/release metadata. The executable also runs alone. Other
projects provide all UI implementations and executable helpers.

```sh
./fungi help
./fungi doctor
./fungi open --exec /absolute/helper --input /absolute/input.json --delivery auto
```

Use `fungi.exe` on Windows. No webview, graphical desktop, Go, Node or npm is
needed by the engine. A helper brings its own runtime requirements.

Read the receipt's delivery mode. For callback delivery, collect once on
notification; for manual delivery, collect once after the user's “done.”
Use wait mode only with a suitable host pending/background call. Never poll.
See `fungi help helper`, `fungi help delivery` and `fungi help examples`.

This package uses protocol 2 and format 3. Protocol-1 UI sessions need their
original executable and state directory. See `fungi help migration` before
replacing an old installation. No UI companion is maintained in this project.

Portable copies do not self-update. Use `fungi install --archive FILE --sha256
HASH --dir DIRECTORY` to create a managed install, then use its `bin/fungi` for
update/rollback/prune. Keep session state outside the installation. `doctor`
checks storage and file locks, not external helper runtimes.
