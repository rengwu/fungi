---
title: Fungi
updated: 2026-10-22
---
# Fungi

Fungi is a prototype portable interaction engine for skill-driven agents. Skills define workflows; fungi runs their interaction sessions and returns human responses to the originating agent.

This wiki records the accepted design, implementation and verification. A Go core and Electron reference renderer now exist. The latest decision requires no harness setup, with manual fallback on unsupported versions. Native return checks and remaining limits are recorded explicitly.

Electron + Go is the accepted direction. The [first end-to-end proof](proposal.md#first-proof) targets macOS and automatic return in six harnesses, with small HTML/flat-JSON forms, image attachments and a native helper contract test.

- [Requirements and discussion](planning.md)
- [First implementation contract](contract.md)
- [Author a Fungi interface](authoring.md)
- [Feature map](features/index.md)
- [First prototype verification](verification/runs/2026-10-20-first-proof.md)
- [Current no-setup compatibility verification](verification/runs/2026-10-20-no-setup.md)
- [Skill authoring, recovery and longer-idle verification](verification/runs/2026-10-20-skill-proof.md)
- [Renderer flexibility and reliability verification](verification/runs/2026-10-22-platform.md)
- [Hosted platform checks and remaining failures](verification/runs/2026-10-22-hosted.md)
- [Verification setup](verification/setup.md)
- [Research findings and recommendation](research/index.md)
- [Proposed small authoring interface](proposal.md)
- [First interview: submitted answers](journal/2026-10-19-grill-eb02aa23-b51f-42c9-b9c0-8fa331bfb02c.md)
- [Second interview: submitted answers](journal/2026-10-19-grill-186963ab-2f7b-423c-b40b-97b0afacc6fb.md)
- [Third interview: first milestone](journal/2026-10-19-grill-689a37f9-b8c4-4f2a-9e54-282234bc2551.md)
- [Fourth interview: no-setup delivery](journal/2026-10-20-grill-397cc4f4-31e2-4fd5-a08e-30ee197e3297.md)
- [Desktop verification arrangements](journal/2026-10-20-grill-desktop-arrangements.md)

The project uses jstack's plan, wiki, grill-form, and rule-lean workflows. The wiki builder is development documentation tooling, not a proposed dependency of the delivered fungi binary.
