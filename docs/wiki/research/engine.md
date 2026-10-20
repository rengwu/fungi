---
title: Engine and extension mechanisms
updated: 2026-10-20
---
# Engine and extension mechanisms

[Research overview](index.md)

## What belongs in the engine

The engine should know who requested an interaction, what was launched, where its result is stored, and how completion reaches that caller. It need not understand a slider, a drawing stroke, a React component, or an interview's decision tree.

The minimal useful responsibilities are:

1. Accept input and identify the exact caller and session.
2. Launch the chosen interface and observe readiness, exit, and failure.
3. Accept explicit submit/cancel events and retain recoverable drafts where supported.
4. Save a result and its artifacts before acknowledging submission.
5. Deliver a completion notice or return a result through a supported host mechanism.
6. Allow collection and recovery by exact session identifier.

Keep this a session runner rather than a general workflow scheduler. A skill can sequence forms, drawing, and review steps. A renderer can implement multiple pages or local branching without asking the model on every interaction.

## Extension mechanisms compared

| Mechanism | Author experience | Useful property | Cost or limitation | Assessment |
| --- | --- | --- | --- | --- |
| HTML file or directory | Write ordinary web content and call submit | Lowest friction for generated or reusable UI | Needs a web renderer; custom code needs a constrained bridge | Preferred common path |
| External executable | Read JSON input; write a structured result | Language independence and separate releases | Extra executable, process lifetime, and trust must be handled | Preferred native escape hatch |
| Compiled-in Go packages | Implement a small internal interface | One selected executable; no runtime plugin loader | Adding native functionality requires a new build | Good for first-party distributions |
| Go `plugin` shared libraries | Match host interfaces and exact builds | Direct calls and shared data | No Windows support; strict toolchain/dependency compatibility | Reject as the general extension contract |
| WebAssembly | Compile a module against host imports | Portable computation with controlled host access | Runtime and ABI to maintain; no native UI appears without host functions | Defer until a concrete need |
| Universal component DSL | Describe every UI using a custom schema | Validatable, discoverable catalog | New authoring language and ongoing adapters for every renderer | Optional renderer, not core requirement |

Go's own documentation lists platform and build-compatibility restrictions and recommends considering IPC. This is stronger evidence than assuming that a Go core should use Go's package named `plugin`. [Go plugin documentation](https://pkg.go.dev/plugin)

HashiCorp's go-plugin demonstrates mature subprocess separation, version handshakes, and reattachment. Its actual startup handshake negotiates a local socket and RPC protocol. That is a useful design reference, but adopting its full RPC interface and generated gRPC machinery is not automatically simpler for a skill author. Process separation also does not sandbox the child from the filesystem or network. [go-plugin](https://github.com/hashicorp/go-plugin), [handshake implementation contract](https://github.com/hashicorp/go-plugin/blob/main/docs/internals.md)

wazero is a pure-Go WebAssembly runtime that avoids a cgo/native-library requirement. Extism demonstrates the host-function approach to granting additional capabilities to modules. Neither gives a UI-free core a browser engine or native window for free. A WASM extension needing windows still needs those capabilities implemented in the host. [wazero](https://wazero.io/), [Extism host functions](https://extism.org/docs/concepts/host-functions/)

## Keep the protocol smaller than the implementation

For ordinary web UI, the author should see input, submit, cancel, and optional draft saving. Session IDs, callback retries, artifact paths, and harness routing stay inside fungi.

For a native process, start with one input envelope and one final result envelope. If draft updates or intermediate events become necessary, add bounded JSON-line events without making a full RPC server mandatory. Standard output carries protocol data; standard error carries bounded diagnostics. A clean process exit without a valid result is not a submission.

Useful versioning is a small protocol version and explicit optional capabilities. Avoid a dependency solver, plugin marketplace, arbitrary callback scripts, and mandatory registries in the first release. An explicit local path should be enough to run authored content.

## Reliability without a large framework

Submission, renderer exit, and agent consumption are different events. A UI may close without submission. A result may be saved after the harness disconnects. A callback may be accepted while its acknowledgement is lost. Represent those differences internally without making authors operate a complicated state machine.

Observable outcomes are submitted, cancelled, closed with draft retained, and failed. A timeout of an agent tool is a transport outcome, not evidence that the human cancelled. The third interview accepted close as a distinct terminal outcome with draft preservation; the [contract draft](../contract.md#result-and-lifecycle) proposes reopening into a new linked session.

Save the response and artifact bytes before returning success to the UI. Associate delivery with a stable result ID. Retry failed notifications locally within explicit bounds; do not repeatedly wake the model. A queued notification means the host accepted a message, not that the agent has read the result. Collection should be repeatable and should not delete data. Cleanup follows an explicit acknowledgement or retention policy.

Exactly-once handling across independent applications cannot be promised just by recording a local flag: a crash can occur after the receiver accepted a message and before fungi recorded success. Use repeatable result retrieval and deduplication by result ID where the host permits it. Where the host lacks an idempotency key, distinguish best-effort notification from confirmed consumption.

## Storage choice

For a small local tool, session directories with JSON records and artifact files are attractive: easy to inspect, recover, and move. But reliable persistence requires more than writing JSON and renaming a file. Go explicitly does not promise atomic rename on non-Unix platforms. SQLite's atomic-commit explanation shows why write ordering and flushes matter under crashes. [Go Rename contract](https://pkg.go.dev/os#Rename), [SQLite atomic commit](https://sqlite.org/atomiccommit.html)

bbolt offers embedded transactions, but its database write lock constrains multiple processes sharing one database. It would be a reasonable option under a single owner, not a reason to introduce a permanent central daemon prematurely. [bbolt limitations](https://github.com/etcd-io/bbolt#caveats--limitations)

Recommendation: keep a narrow persistence interface and choose a proven file-commit implementation or embedded transactional store after testing crash recovery on the supported platforms. Do not build an event-sourcing framework. One authoritative submitted result and modest delivery bookkeeping should be enough initially.

## Proposed verification before expanding UI

| Scenario | Required observation |
| --- | --- |
| Human stays away for hours | No fungi-caused model turns, status tool calls, or repeated notifications |
| User submits twice | One authoritative result; repeated submission returns the existing receipt or an explicit conflict |
| Renderer closes or crashes | Distinct outcome; no fabricated successful response; saved draft survives |
| Engine crashes around submission | UI acknowledges only a committed result; recovery finds a valid result or explicit incomplete state |
| Callback fails or its response is lost | Result remains available; bounded retry; deduplication behavior documented |
| Multiple conversations use fungi | Each receives only its own session result |
| Harness kills a process group | No false success; persistence and supported recovery route remain clear |
| Machine sleeps or reboots | No lost committed data; automatic resumption is claimed only if an actual running/restarted receiver exists |
| Disk full or artifact interrupted | Submission stays unsuccessful; existing data is preserved |
| Invalid or oversized extension output | Bounded memory use; useful failure; no indefinite wait |

These are proposed acceptance cases, not completed tests.
