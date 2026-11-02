# Fungi

<img alt="Fungi icon" src="./docs/assets/icon/fungi-mac.svg" width="34%" align="right">

**Human input for AI agent workflows.**

Let an agent request a review, collect answers or files, and continue when the
response is ready.

Fungi handles sessions, saved responses, files, and completion delivery. You
provide the interface through a **helper**: a program that exchanges JSON with
Fungi over standard input and output.

The engine is one executable with no desktop, webview, or language runtime
required. Helpers supply their own UI and runtime.

[Try it](#try-it) · [Write a helper](docs/manual/helper.md) ·
[Read the manual](docs/manual/start.md)

<br clear="right">

## What you can build

For developers building agent tools and skills, helpers can support workflows
such as:

| Workflow | What the agent gets back |
| --- | --- |
| Review a design, document, or proposed change | A decision, comments, and supporting files |
| Collect requirements or preferences | Structured answers in your chosen JSON format |
| Request screenshots or other local files | Saved files with paths and checksums |

Fungi is headless. Forms, windows, and other interfaces belong to helper
projects.

## Why Fungi

- **Keep responses and files.** Fungi saves them before acknowledging a
  submission. A notification failure does not erase the result.
- **Recover saved progress.** Reopen interrupted or closed sessions from the
  last draft your helper saved.
- **Wait without polling.** Use completion callbacks where supported, a pending
  call where the host allows it, or an explicit manual fallback.
- **Choose your language.** Any program that supports pipes and files can be
  a helper. You define the input and response format.
- **Run on macOS, Windows, and Linux.** Get offline help, machine-readable
  discovery, and managed installation with updates and rollback.

## Try it

**Development preview:** This README describes the unreleased **0.2 engine**.
The [published 0.1 alphas](https://github.com/rengwu/fungi/releases) contain the
older UI and use a different protocol. Build from source below; read the
[migration guide](docs/manual/migration.md) before upgrading an existing integration.

With **Go 1.26** installed, run this example on macOS or Linux:

```sh
git clone https://github.com/rengwu/fungi.git
cd fungi
go build -o ./dist/fungi .
go build -o ./dist/helper ./examples/helper
./dist/fungi open --exec ./dist/helper --delivery wait
```

On Windows, add `.exe` to both output filenames and use those names in the
launch command.

This example completes immediately, with no UI. It saves a draft and a text
file, then returns `outcome: "submitted"`, the response, and the file's saved
path and checksum.

The [Go example](examples/helper/main.go) and [Python example](docs/manual/examples.md)
show the full protocol. Run `./dist/fungi help` for offline docs or
`./dist/fungi doctor` to check session storage.

## Connect your own tool

1. The agent opens a session and passes JSON input to your helper.
2. Your helper collects a response and submits JSON and files.
3. Fungi saves the result and reports completion. The agent reads the response
   and continues.

With your own helper and input file:

```sh
./dist/fungi open --exec ./my-helper --input ./input.json --delivery auto
```

**Read the launch receipt.** Automatic delivery depends on the host. If it
falls back to manual, ask the user to say “done,” then collect once. For a
callback, collect once when notified. Use `wait` only with a supported pending
call. Never poll for a person to finish. See [completion and recovery](docs/manual/delivery.md).

Helpers are trusted local programs with normal OS access. Saved sessions stay
on your machine until explicitly removed.

## Documentation

| Guide | Covers |
| --- | --- |
| [Writing a helper](docs/manual/helper.md) | Messages, drafts, files, and acknowledgments |
| [Complete examples](docs/manual/examples.md) | A runnable Python helper and the Go example |
| [Installation and updates](docs/manual/installation.md) | Packages, managed installs, and rollback |
| [Migration from 0.1](docs/manual/migration.md) | Protocol 2, package format 3, and existing sessions |

## Development

Use Go 1.26 for the engine and Node 22+ for packaging and documentation tools.

```sh
go vet ./...
go test -race ./...
node docs/wiki/_build.mjs --check
```

Omit `-race` for Windows. CI checks macOS arm64, Windows x64, and Ubuntu
22.04/24.04, including installation, updates, and rollback.
[Verification setup](docs/wiki/verification/setup.md) · [Release process](docs/RELEASING.md)

[Share a workflow, report a bug, or propose a helper](https://github.com/rengwu/fungi/issues).
Contributions are welcome.
