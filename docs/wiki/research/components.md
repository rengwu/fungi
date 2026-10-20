---
title: Composition and component libraries
updated: 2026-10-19
---
# Composition and component libraries

[Research overview](index.md)

## The useful composability boundary

Let the renderer own layout and interactivity. Let fungi own input, committed output, artifacts, and delivery. This allows a skill to reuse existing UI code instead of rewriting it for a fungi-specific widget tree.

For web interfaces, standard DOM elements and events are the lowest common denominator. Web Components offer reusable controls without requiring all authors to use the same framework. React, Vue, or Svelte can remain choices for larger interfaces. A skill may ship compiled assets without requiring an end user to install Node or a package manager.

No library should be loaded from a CDN by default. Bundle the assets a capability actually uses, including fonts and workers where necessary. Rich component authoring may use a build pipeline, but simple authored HTML should not require one.

The first interview accepts both direct HTML and JSON component descriptions in the first release. Keep one component implementation underneath both paths: a JSON renderer maps descriptions to controls and uses the same submit bridge. It need not make every custom HTML program expressible as JSON. Framework/library compatibility is a test requirement for the selected bundled engine, not something established by this survey.

## Libraries worth reusing

| Need | Candidate | Why it fits | Integration boundary |
| --- | --- | --- | --- |
| Small reusable controls | [Lit](https://lit.dev/docs/) | Builds on Web Components; composable with ordinary HTML | Good optional component implementation, not mandatory author framework |
| Ready-made web controls | [Web Awesome](https://webawesome.com/) | Framework-independent standards-based components | Select core assets; distinguish free components from Pro resources |
| Custom accessible React controls | [React Aria](https://react-aria.adobe.com/) | Keyboard, focus, input behavior, and internationalization | Good for sophisticated React surfaces; accessibility still needs host testing |
| Familiar application components | [shadcn/ui](https://ui.shadcn.com/docs) | Customizable component code and distribution model | Useful in authored bundles; its scaffolding/registry should not become fungi's runtime requirement |
| Schema-generated forms | [JSON Forms](https://jsonforms.io/) | Separates data schema, UI schema, and data; supports multiple frontend frameworks | A form capability above the engine, not the universal workflow format |
| Live parameter tuning | [Tweakpane](https://tweakpane.github.io/docs/) | Compact standalone inputs, monitoring, plugins | Particularly close to the animation-timing use case |
| Object-based drawing | [Fabric.js](https://www.fabricjs.com/) | Editable text, transformations, SVG conversion, serialization | Strong candidate for a modest drawing tool |
| Interactive 2D scenes | [Konva](https://konvajs.org/docs/) | Shape/event model, layers, drag handling, serialization | Useful when interaction logic is custom; restore handlers/assets separately from serialized nodes |
| Freehand stroke primitive | [perfect-freehand](https://github.com/steveruizok/perfect-freehand) | Converts pointer/pressure data into stroke geometry | Small building block; does not supply a complete editor |
| Full sketching interface | [Excalidraw](https://github.com/excalidraw/excalidraw) | Embeddable editor, PNG/SVG and editable-scene output | Larger React-based capability; export both image and scene when useful |
| Extensible canvas SDK | [tldraw](https://tldraw.dev/community/license) | Relevant canvas ecosystem, but production use requires an appropriate license key | Optional choice; not a frictionless default dependency |
| Interactive charts | [Apache ECharts](https://echarts.apache.org/en/feature.html) | Broad chart types, interactions, Canvas/SVG renderers | Return selection/filter state alongside data or an image |
| Declarative analytical charts | [Vega-Lite](https://vega.github.io/vega-lite/) | Concise JSON and composable interactive views | Natural fit for agent-generated chart specifications |
| Editable node diagrams | [React Flow](https://reactflow.dev/) | Custom nodes, edges, pan/zoom, selection | Use when diagrams must be edited, not merely displayed |
| Text-authored diagrams | [Mermaid](https://mermaid.js.org/intro/) | Familiar text syntax and SVG rendering | Preserve source; configure safe rendering; not a general diagram editor |
| Structured rich text | [Tiptap](https://tiptap.dev/docs/editor/getting-started/overview) | Modular ProseMirror-based editor and custom nodes | Useful for annotated feedback; distinguish open-source and paid extensions |
| Terminal interaction | [Huh](https://github.com/charmbracelet/huh) | Forms and accessible standard-prompt mode | Optional visible-terminal renderer |

The recommendation is to make these usable, not ship all of them. The first reference UI needs ordinary form controls. Add a library when it saves substantial implementation work for an actual capability.

## Agent-oriented UI formats

### A2UI

A2UI separates declarative component descriptions from client-owned renderers and catalogs. It supports multiple rendering technologies and incremental JSON. This is useful when a skill wants a constrained catalog with predictable output. It is not a substitute for a local process owner, durable result storage, or an idle-conversation callback. Keep it as a possible renderer input format. [A2UI](https://a2ui.org/)

### json-render

json-render lets an application define typed components and actions, then map generated JSON to real components. It provides one approach for fungi's accepted JSON authoring path. The cost is maintaining that catalog, action semantics, bindings, and renderer integration. Compare it with a smaller form-specific adapter before adopting the full framework. Direct HTML remains equally supported. [json-render documentation](https://json-render.dev/docs)

### MCP Apps

MCP Apps packages an HTML interface as an MCP resource and defines communication with a supporting host, commonly through an isolated iframe. This is a valuable optional delivery surface when the user wants the interface inside chat. A native fungi window is a different host; compatibility requires implementing the relevant bridge and capabilities. Merely exposing an MCP tool does not make every harness an MCP Apps host. [MCP Apps overview](https://modelcontextprotocol.io/extensions/apps/overview)

### AG-UI and Remote DOM

AG-UI standardizes agent/frontend event streams, shared state, and human interruptions. It is broader than fungi's initial local interaction round trip. Use an adapter if a real integration needs it, rather than building an agent runtime to satisfy the protocol. [AG-UI](https://docs.ag-ui.com/introduction)

Shopify's Remote DOM mirrors a restricted component tree from a separate JavaScript environment to a host. It is relevant if fungi eventually runs mutually untrusted component authors in one surface. It also requires host/remote plumbing and component restrictions; that complexity is not justified for a single trusted skill-authored page. [Remote DOM](https://github.com/Shopify/remote-dom)

## What should travel back to the agent

A useful result preserves meaning, not only appearance. A form returns typed values; a chart returns selected data and filters; a drawing can return a PNG plus editable scene data; annotated screenshots return each image paired with its note. Fungi should store these as JSON plus artifact references, with a compact summary for the agent.

An image path is not equivalent to the model seeing the image. CLI clients need an image-reading tool or a harness attachment mechanism. MCP can express image/resource content, but the host still controls supported modalities. This belongs in capability discovery and skill instructions.

Keep local changes local: dragging a node, moving a slider, typing, and drawing should not create model turns. Submission, explicit requests for assistance, or workflow-defined milestones may do so. This is the distinction that makes rich interaction compatible with the AFK requirement.
