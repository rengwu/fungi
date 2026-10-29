# External helper protocol

```sh
fungi open --exec ./my-helper --input ./input.json --delivery auto -- custom-arg
```

Helpers are trusted programs with ordinary OS access. Use any language that
supports pipes and files; supply its runtime. All UI implementations live in other projects. Fungi does not install helper
runtimes or ship a compiler.

## Launch and readiness

Use UTF-8 newline-delimited JSON over stdin/stdout. Reserve stdout for protocol
messages and stderr for diagnostics; flush each message. Fungi sends this launch
object (abbreviated values):

```json
{"protocol":2,"session":"…","input":null,"draft":null,"draftArtifacts":[],"outputDir":"/absolute/output","limits":{"jsonBytes":1048576,"artifactBytes":20971520,"totalArtifactBytes":104857600},"delivery":{"mode":"manual"}}
```

Reject unsupported protocols. Read launch input, initialize the helper and emit
`{"type":"ready"}` within 25 seconds. No human timeout follows readiness.
Keep reading acknowledgments while the UI is active: Fungi can stop a helper
that does not read a protocol message within five seconds.
Ignore unknown metadata in supported messages; protocol 2 permits additive fields.

## Drafts and results

Optional draft:

```json
{"id":"save-1","type":"draft","data":{"notes":"In progress"},"artifacts":[]}
```

Submission after writing `note.txt` under `outputDir`:

```json
{"id":"finish-1","type":"result","outcome":"submitted","data":{"notes":"Ready"},"artifacts":[{"key":"note","name":"note.txt","path":"note.txt","mediaType":"text/plain"}]}
```

Other outcomes: `cancelled`, `closed`, `failed`. For failure, include
`error:{"code":"your_code","message":"Useful explanation"}`. Each draft/result
receives an ack with matching `id`:

```json
{"id":"save-1","type":"ack","ok":true}
{"id":"finish-1","type":"ack","ok":true,"result":{"protocol":2,"kind":"result","session":"…","outcome":"submitted","data":{"notes":"Ready"},"artifacts":[]}}
```

Successful result acks contain the complete saved result; artifact details
are omitted above. Failure ack:
`{"id":"save-1","type":"ack","ok":false,"error":{"code":"…","message":"…"}}`.
Wait for acknowledgment before reporting success or exiting. Preserve input
for retry on failure. Exiting without a terminal result produces `failed`.

Each draft/result replaces all data and attachments. Artifact paths must be
relative to and confined within `outputDir`; keys must be unique and nonempty,
and names nonempty. Fungi copies bytes into durable storage before acknowledgment.
To retain a restored `draftArtifacts` record, copy its absolute saved path's
contents into the new `outputDir` and include the file in the next message.

## Limits and compatibility

Read the advertised `limits`: currently 1 MiB per JSON envelope, 20 MiB per
artifact and 100 MiB of retained draft/result artifacts per session. Keys are
nonempty, unique and at most 200 UTF-8 bytes; names are nonempty and at most
255 bytes. Binary files use artifact paths rather than JSON-encoded bytes.

Protocol 2 preserves documented message names, types and meanings, session
commands, receipts and results. Optional fields may be added; ignore unknown
metadata. A breaking change requires a new protocol. Helper input/result schemas
belong to their providers. Internal commands and storage layout are not APIs.

Fungi launches helpers without a shell. They are trusted executables with normal
OS access; this is not an application sandbox. The agent callback target is
captured outside helper-controlled data, and CODEX_THREAD_ID is removed from
the helper environment. Remote UI transport is the helper project's concern.

Result interpretation: `fungi help delivery`.
The source repository includes a protocol fixture at `examples/helper/main.go`.
