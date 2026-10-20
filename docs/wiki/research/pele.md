---
title: Pele as an optional diagram asset
updated: 2026-10-20
---
# Pele as an optional diagram asset

The user suggested [Obsidian's Pele](https://github.com/obsidianmd/pele) during the simpler grill-form port. It renders Mermaid text to SVG; it is not a notes editor or native window engine. Its MIT-licensed package declares no runtime dependencies. It targets Mermaid 12.1.0 syntax with its own layout, so visual equivalence should not be assumed. [Package metadata](https://registry.npmjs.org/pele/1.7.0)

## Fit and maintenance

The GitHub API snapshot on 2026-10-20 recorded 2,553 stars, 77 forks, creation on 2026-10-05 and a last push on 2026-10-19. The contributor listing attributed 195 commits to `kepano` and one to `uzuraDev`. Version 1.7.0 was released on 2026-10-19. These are strong early interest and activity signals, with a short history and concentrated maintenance. They do not establish mature long-term support. [Repository metadata](https://api.github.com/repos/obsidianmd/pele), [contributors](https://api.github.com/repos/obsidianmd/pele/contributors), [release commit](https://github.com/obsidianmd/pele/commit/8603bf3313add9143cd417015d4a6bd6c522dada)

Pele fits the skill-assets layer: vendor a pinned browser file and call its API from ordinary page JavaScript. Neither the Go runner nor Electron's bridge needs to know about diagrams. `render()` returns SVG synchronously; `mount()` adds container resizing, source updates and zoom. Individual diagram imports are available if bundle size matters. [API](https://pele.run/api)

For untrusted source, retain application-controlled options, disable links and images, bound source/output/edges, and display errors as text. Supplied icon markup and internal link handling remain the application's responsibility. Synchronous rendering alone does not impose a time limit; upstream recommends a terminable worker for that. [Security contract](https://pele.run/security)

## Local probe

Downloaded the pinned npm 1.7.0 tarball and verified its published SHA-512 integrity before reading the browser/module distributions. The browser file is 683,259 bytes, or 254,846 with local gzip. The older grill-form's Mermaid Tiny file is 2,811,019 bytes, or 752,544 gzip. These are measured files, not a rendering-speed benchmark or a statement of equivalent coverage.

Three simple flowchart, sequence and state fixtures produced nonempty SVG with positive dimensions under Node 22.19.0, without a DOM or additional packages. Invalid syntax and a deliberately tiny source limit were rejected. This does not establish browser interaction, complex-diagram compatibility or visual quality. [Probe evidence](../evidence/2026-10-20-skill-proof/pele.json)

Recommendation: keep Pele as a strong candidate for the next diagram-focused skill experiment. Preserve Mermaid Tiny in this requested historical baseline so the migration tests isolate the Fungi lifecycle changes. No new renderer abstraction or core dependency is justified.
