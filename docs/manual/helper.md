# Native executable protocol

```sh
fungi open --exec ./my-helper --input ./input.json --delivery auto -- custom-arg
```

Implement in any language that reads/writes pipes and files. Ship whatever
runtime your helper needs. Fungi does not install it. Helpers are trusted
programs with ordinary OS access.

Transport: UTF-8 newline-delimited JSON, one message per line. Read stdin;
write only protocol messages to stdout, diagnostics to stderr. Flush after
each message. Fungi sends this launch object (values abbreviated):

```json
{"protocol":1,"session":"…","input":null,"draft":null,"draftArtifacts":[],"source":"","outputDir":"/absolute/output","limits":{"jsonBytes":1048576,"artifactBytes":20971520,"totalArtifactBytes":104857600},"delivery":{"mode":"manual"}}
```

Reject unsupported protocols. Read launch input promptly; show the interface
and emit `{"type":"ready"}` within 25 seconds. There is no human timeout after
readiness. A helper that fails to read a protocol message within five seconds
can be stopped; keep reading acknowledgements while the UI is active.

Optional draft:

```json
{"id":"save-1","type":"draft","data":{"notes":"In progress"},"artifacts":[]}
```

Submission after writing `note.txt` under `outputDir`:

```json
{"id":"finish-1","type":"result","outcome":"submitted","data":{"notes":"Ready"},"artifacts":[{"key":"note","name":"note.txt","path":"note.txt","mediaType":"text/plain"}]}
```

Other terminal outcomes: `cancelled`, `closed`, `failed`. For failed, include
`error:{"code":"your_code","message":"Useful explanation"}`. Draft/result
messages receive an ack with matching `id`:

```json
{"id":"save-1","type":"ack","ok":true}
{"id":"finish-1","type":"ack","ok":true,"result":{"protocol":1,"kind":"result","session":"…","outcome":"submitted","data":{"notes":"Ready"},"artifacts":[]}}
```

The second example omits artifact details for brevity; actual successful
result acks contain the complete saved result. Failure ack:
`{"id":"save-1","type":"ack","ok":false,"error":{"code":"…","message":"…"}}`.
Wait for ack before reporting success or exiting; preserve user input for
retry on failure. Exiting without a terminal result becomes `failed`.

Each draft/result replaces data and the attachment list. Attachment paths are
relative to `outputDir`, confined to it; use unique nonempty keys and names.
Fungi copies bytes into durable artifacts before acknowledgement. Restored
`draftArtifacts` are saved artifact records with absolute paths; to retain one,
copy its bytes into the new `outputDir` and include it in the next message.
The same limits and result interpretation apply as in `fungi help html` and
`fungi help delivery`. The source repository's `examples/helper/main.go` is a
small executable protocol fixture; no helper compiler is shipped to end users.
