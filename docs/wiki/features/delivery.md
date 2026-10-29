---
title: Returning to an agent
updated: 2026-10-28
---
# Returning to an agent

**ID: FEAT-delivery**

The [accepted requirement](../planning.md) is automatic response delivery without repeated model checks or AFK token expenditure. Fungi preserves results independently of the delivery route.

## Entry points and expected outcomes

- `--delivery auto`: bind an inherited exact `CODEX_THREAD_ID` when `codex queue` is supported; otherwise explicitly report manual delivery.
- `--delivery wait`: keep one command pending through a process pipe, with no Fungi human deadline. The harness must retain that call or supply background completion.
- `--delivery manual`: return an opened receipt; after the user says “done,” collect the exact session once.

The callback contains the executable, exact session and state-directory environment, not user response content. Delivery bookkeeping distinguishes `armed`, `attempting`, `accepted` and `unknown`. A notification failure never erases a committed result. An uncertain attempt is not automatically retried. Acceptance does not prove the originating agent consumed the notice.

## Verification recipe

Run `go test -run 'TestCallback|TestManual|TestWait|TestDisconnected' .` and the
packaged engine smoke from [setup](../verification/setup.md). These check result
preservation after an uncertain notification, manual receipts, pending calls
and disconnected launchers. The automated fixtures have finite bounds and do
not send real model callbacks.

A real callback check must launch the protocol-2 idle helper from the exact
conversation, end the turn, receive the notification and collect once. Host
consumption is distinct from an accepted queue command. Do not poll.

## Coverage and limitations

The headless run records new engine checks. Real end-of-turn callbacks,
all-harness resumption and sleep/wake need their own protocol-2 evidence;
historical UI observations below are not new passes. No host configuration is
required; unsupported automatic routes remain explicit manual fallbacks.

## Verification history

- [Engine-only migration](../verification/runs/2026-10-28-headless.md) — protocol 2; previous reports retain their original build scope.

- [2026-10-27 intervening chat message](../verification/runs/2026-10-27-interleaved-message.md) — human-message survival and a separate mid-tool submission; both queued callbacks resumed this conversation after turn end without “done.”

- [2026-10-27 contract and AFK checks](../verification/runs/2026-10-27-afk-contract.md) — current Wails package and Codex 0.164.0: short callback, actual sleep/wake and two-hour idle all consumed once without polling.

- [2026-10-20 skill proof](../verification/runs/2026-10-20-skill-proof.md) — native skill return after one macOS sleep/wake cycle and automatic return from a two-hour helper both passed.
- [2026-10-20 first proof](../verification/runs/2026-10-20-first-proof.md) — versioned probes, failures, corrections and explicit limits.
- [2026-10-20 no-setup compatibility](../verification/runs/2026-10-20-no-setup.md) — native form/image journeys, Kimi after sign-in, and OpenCode's manual fallback.

## Historical protocol-1 observations

The [October 27 local callback test](../verification/runs/2026-10-27-interleaved-message.md)
returned a human response after an intervening chat message and a separate
timed response submitted during a running tool call. In both cases the queued
notice resumed the same conversation after its active turn ended, without
“done,” polling or restarting a waiter. This observation applies to the tested
Codex callback route, not to a host-managed `wait` call that the host cancels.

| Tested harness | Invocation for a human interaction | Evidence and limits |
| --- | --- | --- |
| Codex 0.164.0 | `auto` with the current Wails package | Ten-second submission, user-controlled Mac sleep/wake and a two-hour page timer all resumed the exact conversation after turn end. One collection per notice verified JSON and PNG. The long timer spanned two sleep/wake cycles and 2 h 7 min 47 s of wall time |
| Codex 0.163.0 | `auto` with inherited exact thread and available queue command | Helper, native-close and packaged form/image callbacks consumed; later native sleep/wake and two-hour helper callbacks also consumed once |
| Claude Code 2.1.302 | `wait` through native background Bash in a suitable local interactive session | Headless form/image and three-minute completion passed; headless/cloud execution is bounded, so use manual for an unrestricted human wait there |
| OpenCode 1.18.34 | `manual`; end turn; collect once after “done” | Native image flow passed; ordinary v1 Bash is bounded |
| Grok 1.0.62 | `wait` through `run_terminal_command` with `block_until_ms: 0`; end turn | Native form/image and three-minute completion passed; retrieve completed output once after notification |
| Kimi 2.2.0 | `wait` through background Bash with `disable_timeout: true`; end turn | Native form/image and three-minute completion passed after user sign-in |
| Pi 1.3.0 | `wait` through one pending Bash call, timeout omitted | Native form/image and three-minute completion passed |

This is a versioned compatibility record, not a universal promise about another mode or future harness version. `auto` does not set another harness's tool arguments; skills choose the native route or the explicit fallback.

## Historical verification recipe

Use a disposable conversation and the native helper, with an isolated state directory. `helper wait 180` stays silent for three minutes before returning a real Fungi artifact/result. Allow only the fixture command; do not change persistent harness settings.

For Claude, use native background Bash, end the turn, and observe the completion event followed by result consumption. For Grok, use `block_until_ms: 0`, end the turn, and retrieve completed output once after notification. For Pi, omit the Bash timeout and retain one call. Record actual arguments: in the completed Claude test the model supplied a 600,000 ms timeout, so that test does not establish arbitrary AFK duration.

For Codex and OpenCode, ten-second pending-call probes returned real Fungi results. Codex requires permission to write the session directory. Canonicalize macOS `/var` paths before using them as an OpenCode working directory. These short tests do not prove an ended-turn callback.

For the Codex callback, launch from the intended live conversation, end the turn, submit the UI, and verify that the exact conversation consumes its notice. Check `delivery.json` separately. After the implementation turn ended, the queued notice for native-helper session `a4e20e6a9ce120f2ccf7461a5ae22b0a` reached the originating conversation. One collection returned its submitted data and artifact. This verifies callback consumption for that fixture, including a notice queued while the original turn was still active. Its older notice omitted the custom state directory, which was recovered from the known launch context.

The corrected notice for native-form session `f2b35d71fb650809ac354bdd593df8fd` subsequently reached the same originating conversation with its full executable, argument array and state-directory environment. One collection using those supplied values returned `closed` without answer data or artifacts. Both callback formats have now been consumed; the second verifies the corrected self-contained collection instructions.

The initial Kimi probe stopped before tool execution because its selected provider had no credential. After the user completed sign-in, Kimi 2.2.0 passed a three-minute background wait with `disable_timeout: true`, then a separate native form/image submission. Both continued after WAITING and read completed output once.

## Historical coverage and limitations

All six harnesses have returned actual Fungi results. Native form/image automatic return passed in Codex, Pi, Grok, Kimi and Claude's bounded headless mode. OpenCode v1 passed the explicitly resumed manual flow. Codex consumed the final portable image callback after the prior turn ended; one collection using the supplied argv/environment returned the submitted response and matching image. [Current report](../verification/runs/2026-10-20-no-setup.md)

No OpenCode callback/plugin adapter is required or implemented. A bounded shell route is not presented as an indefinite one. One user-controlled macOS sleep/wake cycle subsequently returned its native skill response through the exact-origin Codex callback, collected once. A separate two-hour native-helper run also completed through the same callback route without polling; its timing artifact matched the saved result. Arbitrary-duration and broader sleep/wake reliability remain unverified.

The [fourth interview](../planning.md#no-setup-delivery-decision-on-2026-10-20) settled the setup boundary: use built-in capabilities only, with explicit manual fallback on unsupported versions. Do not require an integration or harness upgrade. A short bounded command still does not qualify as an arbitrary human wait.
