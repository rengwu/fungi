---
title: Agent delivery and waiting
updated: 2026-10-19
---
# Agent delivery and waiting

[Research overview](index.md)

## What fungi can guarantee itself

Fungi can remain idle without invoking an LLM. It can save results and expose event-driven waits. It can avoid repeated progress output, model status checks, and timer-triggered agent turns. It cannot make an arbitrary inactive application receive a message without a supported host mechanism.

Compatibility belongs to the harness, its version, its settings, and the execution environment. A model name such as DeepSeek or Grok does not identify those capabilities. The same model may run inside several different harnesses.

The accepted product policy is: no model polling; best-effort automatic delivery across popular harnesses; durable recovery when delivery is unavailable or interrupted. Open the interface even without an automatic route and clearly ask the user to say “done” after submitting. That continuation collects the exact saved session. [First interview decisions](../planning.md)

## Delivery routes

| Route | Suitable situation | Important limit |
| --- | --- | --- |
| Host-managed background completion | Host owns the command and injects completion into the right conversation | A returned process ID alone is insufficient; exit and cleanup behavior vary |
| Existing-session callback | Harness exposes an authenticated queue or session-message API | Must target the original session; accepted message is not confirmed consumption |
| One pending tool call | Host can keep it pending for the interaction duration | Transport timeout cannot be treated as a human deadline |
| Manual collect after user continuation | No supported wake-up route or an interrupted delivery | Accepted fallback: user says “done”; collect the exact saved session |

Prefer one supported route rather than running multiple delivery attempts that may independently wake the model. Do not use a subagent, repeated sleeps, scheduled model checks, or a stop-hook loop as the waiting mechanism.

## First-milestone harness inventory

The third interview expands the first milestone's automatic-return targets to **Codex, Claude Code, OpenCode, Grok, Kimi and Pi**, all on macOS initially. The user reports all are installed; read-only version/help checks confirmed these commands on 2026-10-19:

| Command | Installed version | First route to investigate |
| --- | --- | --- |
| `codex` | 0.163.0 | Exact-thread queue, already exposed by local CLI help |
| `claude` | 2.1.302 | Host-managed background completion or a silent monitor |
| `opencode` | 1.18.34 | Existing server's exact-session message API; establish instance identity and authentication |
| `grok` | Grok Build 1.0.62, revision 1474f9cbdd91 | Background command completion or one silent monitor |
| `kimi` | Kimi Code 2.0.1 | Background Bash completion with the human-wait timeout handled explicitly |
| `pi` | `@earendil-works/pi-coding-agent` 1.3.0 | One pending built-in shell call, which has no default timeout in the installed source |

The [inventory and source hashes](../attachments/research/harness-inventory-2026-10-19.json) record the evidence boundary. No model prompts were sent, callbacks invoked, accounts checked, or harness settings changed during this inspection. No Fungi handoff has been tested. Installation is not evidence of authentication or usable model access.

Prefer these existing mechanisms before introducing a custom plugin. A single pending tool call can meet the no-polling requirement even though the conversation remains in a tool-wait state; it is different from injecting a message into an already-ended turn. For each route, verify silent waiting, exact caller identity, task lifetime, cancellation and completion behavior. Do not count a manual fallback as a successful automatic-return test, or silently reduce the six-harness target to Codex.

## Current harness evidence

| Harness or interface | Evidence found | What remains unverified |
| --- | --- | --- |
| Codex | Official CLI documents `queue --thread ... --message ...`; local 0.163.0 exposes it. Jstack's existing adapter targets `CODEX_THREAD_ID` and probes support | Current live end-to-end fungi delivery; behavior when caller is unavailable; duplicate acceptance |
| Claude Code | Background Bash exists. Current tool documentation also describes Monitor feeding command output lines or WebSocket events to the agent | Exact completion/wake behavior for each installed mode; a long-lived local UI through exit, memory pressure, and restart |
| OpenCode | Server API provides exact-session asynchronous messages and event streams | Discovering and authenticating to the originating instance without configuration; race behavior with an active turn |
| Gemini CLI | Configuration includes background completion modes; `inject` returns output to the agent, while `silent` is documented as default | Effective settings and task lifetime in an installed version; no automatic settings changes assumed |
| Cursor | Stop/subagent-stop hooks can submit a follow-up when their lifecycle event fires | An externally triggered wake of an already idle local conversation is not established by those hook docs |
| Grok Build | Current project documents background command completion notifications and monitor lines delivered to the conversation; local CLI 1.0.62 is installed | Actual idle-agent resumption, timeouts and caller routing in the installed version |
| Kimi Code CLI | Current tool docs describe automatic background Bash completion notifications and a per-call `disable_timeout` option; local CLI 2.0.1 is installed | Actual idle-agent resumption, task lifetime and cancellation in this version |
| Pi | Installed 1.3.0 shell implementation awaits process completion and has no default timeout; extension APIs also expose turn-triggering messages | Real long human wait and result consumption; extension route only if needed |
| ACP | Defines client-controlled session creation/loading, prompts, updates, and cancellation | Does not itself establish access from an unrelated process to an already open IDE conversation |
| Other/model-named agents | No universal wake API established | Need the actual harness, version, and exposed tools |

Codex's official command description establishes exact-session message targeting; the local CLI check establishes availability in this environment. Jstack separately records earlier successful handoffs. These are different evidence levels. [Codex commands](https://learn.chatgpt.com/docs/developer-commands#codex-queue)

Claude's background tasks have environment-dependent lifetimes and may be stopped when the harness exits. Its Monitor tool suggests another useful event route: fungi could emit one result event after being silent while waiting. Neither should be claimed as verified here. [Claude background Bash](https://code.claude.com/docs/en/interactive-mode#background-bash-commands), [Claude tools](https://code.claude.com/docs/en/tools-reference)

OpenCode exposes `/session/:id/prompt_async` and server event streams. A fungi adapter should reuse an identified instance, rather than launch a separate agent with similar history. [OpenCode server](https://opencode.ai/docs/server/)

Gemini's `tools.shell.backgroundCompletionBehavior` distinguishes `silent`, `inject`, and `notify`; the default is `silent`. This is a concrete example of why finding a background command capability is not sufficient. [Gemini configuration](https://geminicli.com/docs/reference/configuration/)

Cursor's `followup_message` belongs to a hook response triggered by an agent lifecycle event. Repeatedly continuing the model merely to wait would violate fungi's intended behavior. The reviewed docs do not establish a general delayed external callback to local idle chat. [Cursor hooks](https://cursor.com/docs/hooks)

Grok Build documents completion notifications for background terminal commands and a monitor whose output lines become conversation notifications. A Fungi waiter could remain silent until one terminal outcome. This is a candidate route; a visible notification must still be tested for actual model continuation. Recurring prompts and scheduler loops do not fit Fungi's waiting contract. [Grok background tasks](https://github.com/xai-org/grok-build/blob/main/crates/codegen/xai-grok-pager/docs/user-guide/20-background-tasks.md)

Kimi's current Bash documentation describes automatic completion notifications, a ten-minute background timeout by default in interactive mode, and a per-call `disable_timeout` control. The adapter must account for human waits explicitly rather than inherit that default or repeatedly invoke `WaitFor`. Print-mode waiting has separate settings and must not be assumed to match the interactive session. [Kimi tools](https://moonshotai.github.io/kimi-code/en/reference/tools), [background configuration](https://moonshotai.github.io/kimi-code/en/configuration/config-files#background)

Pi's installed `dist/core/tools/bash.js` declares an optional timeout with no default, awaits the child process, and returns the output. This suggests a simple pending-call route without installing an extension; it is source evidence, not a completed human-wait experiment. Pi also supports explicitly loaded local extensions and `sendMessage(..., { triggerTurn: true })`; keep that as an alternative if a concrete workflow needs an ended-turn callback. [Installed source evidence](../attachments/research/harness-inventory-2026-10-19.json), [Pi extension documentation](https://pi.dev/docs/latest/extensions)

Kimi's current repository replaces the archived Python implementation and supports ACP. Old Wire/CLI assumptions should not be copied blindly. ACP is valuable when the client owns the agent connection, but adopting it would not automatically let fungi attach to any existing IDE session. [Kimi Code CLI](https://github.com/MoonshotAI/kimi-code), [legacy project notice](https://github.com/MoonshotAI/kimi-cli), [ACP lifecycle](https://agentclientprotocol.com/protocol/v1/overview)

## MCP does not remove the host problem

MCP is useful for discovery, tool invocation, and structured results. The versioned 2025-11-25 Tasks specification includes asynchronous task results, but status notifications are optional and requestors are advised to poll for updates. That polling can happen inside ordinary host software without involving the model. Exposing `tasks/get` to the model and instructing it to keep checking would recreate the original problem. Do not infer support from the protocol alone; negotiate the actual host's capabilities. [MCP Tasks specification](https://modelcontextprotocol.io/specification/2025-11-25/basic/utilities/tasks)

MCP Apps separately negotiates capabilities such as messages and model-context updates. The host decides what it supports. A UI message, a changed resource, and a resumed model turn are not interchangeable guarantees. [MCP Apps host capabilities](https://apps.extensions.modelcontextprotocol.io/api/interfaces/app.McpUiHostCapabilities.html)

## Lessons from jstack

Local source inspected: `../jstack/skills/grill-form/scripts/grill-form.mjs`, `../jstack/skills/screencap/native/Screencap.m`, `../jstack/scripts/input-completion.mjs`, `../jstack/scripts/completion/codex.mjs`, and their skill/reference documents.

Useful existing properties include exact session paths, validated receipts, preserved drafts, non-destructive collection, and a separate completion adapter. The shared helper probes Codex support, starts one collector, and exits; the remaining collector sends a short notice to the originating thread. `armed` records startup, not guaranteed eventual delivery.

Grill-form still relies on a browser downloading an answer file and a local collector checking for it. Fungi can remove that indirection: the renderer submits directly to the owning session. A local wait loop consumes ordinary CPU/resources, not model tokens, but a direct event bridge also improves latency and avoids download-directory ambiguity.

Jstack's one-attempt/no-auto-restart policy is an intentional safety boundary for those helpers. Fungi should decide its own recovery design rather than inherit that rule accidentally. Retrying delivery can be safe when controlled by the engine and a stable result ID; restarting model turns to check readiness is not.

## Routing and lifecycle rules worth keeping

- Bind the caller identity at launch. Never infer the destination from the newest conversation, newest file, or current foreground window.
- Keep callback configuration outside UI-controlled data. A page should submit an answer, not choose a thread or shell command.
- Invoke adapter commands with an argument array, not an interpolated shell command.
- Store the committed result before notifying. The notification should identify a result to collect rather than insert arbitrary captured content as instructions.
- Distinguish saved, notification accepted, and consumed. Retain artifacts until an explicit acknowledgement or defined retention policy.
- Bound local retries and deduplicate where the harness allows it. Preserve the response when delivery fails.
- Ordinary AFK waiting, sleep, process termination, and reboot are distinct tests. After all relevant processes stop, some later explicit restart or installed service is necessary for delivery to resume.

The first interview accepted explicit reopen/collect after process exit. An always-running service, login installation, or automatic reboot recovery is therefore outside the first-release requirement. Saved responses remain recoverable.

## Prototype findings on 2026-10-20

Live probes now supplement the earlier source inventory. [The first proof report](../verification/runs/2026-10-20-first-proof.md) records actual Fungi results in five harnesses and three-minute waits in Claude, Grok and Pi. Kimi is credential-blocked. These are bounded tests, not unlimited AFK guarantees.

The inspected OpenCode Bash source defines a two-minute default and ten-minute maximum. An ordinary pending Bash call is therefore insufficient for arbitrary AFK waiting. The current development v2 documentation describes a separate background-call capability; it has not been established as an active route in the tested configuration. [Bash source](https://github.com/anomalyco/opencode/blob/dev/packages/core/src/tool/bash.ts), [v2 tool documentation](https://dev.opencode.ai/v2/docs/tools/)

The v2 documentation was rechecked after the native UI tests: its `shell` background mode explicitly has no default timeout and notifies the session on completion. The installed CLI remains 1.18.34 and the observed probe used its `bash` tool. A newer tool contract is a promising future route, not evidence that the installed v1 route already works indefinitely. No harness upgrade was performed.

OpenCode’s documented local plugins and runtime configuration offer an opt-in integration route without editing global settings. This remains an implementation option, not a selected or verified adapter. [Plugins](https://opencode.ai/docs/plugins/), [configuration precedence](https://opencode.ai/docs/config/#precedence-order)

## No-setup follow-up

The [fourth interview](../planning.md#no-setup-delivery-decision-on-2026-10-20) rejects a required integration step. Keep the OpenCode plugin route as research only. For the installed v1, use explicit manual delivery; do not require an upgrade to obtain v2's newer shell behavior.

Claude's current documentation distinguishes local interactive background Bash, which has no time limit, from unattended/headless modes, which do. A text-prompt `-p` run defaults to ten minutes and can request at most two hours without changing configuration. Monitor is also bounded and does not solve arbitrary human waiting. Our three-minute and native-form tests use headless background Bash, so they verify notification/consumption within that mode's limit; they do not establish a long interactive wait. [Claude Bash lifecycle](https://code.claude.com/docs/en/tools-reference#time-limit-for-background-commands), [Monitor limits](https://code.claude.com/docs/en/tools-reference#monitor-tool)

After the user signed in, Kimi reported version 2.2.0. Its live three-minute test used `run_in_background: true` and `disable_timeout: true`, ended with `WAITING`, then automatically continued and read the completed output once. A separate native form-and-image submission followed the same route. The print stream does not expose a separate completion-event payload, so the observed evidence is the unsolicited continuation and actual result consumption; no second prompt was injected. [Kimi Bash controls](https://moonshotai.github.io/kimi-code/en/reference/tools#shell)
