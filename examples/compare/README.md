# Compare two designs

A complete custom interface using only the four Fungi bridge methods. No
registration, dependencies or build step. Copy this directory, replace its two
images, and launch:

```sh
fungi open ./compare/index.html --delivery auto
```

Optional `--input config.json` accepts `title`, `context`, `before` and `after`.
Image paths are relative to `index.html` and must stay inside its directory.
The sample stretches both images to a 1200×600 comparison; use matching aspect
ratios. Network images are unavailable.

The divider, preference and notes save as drafts. Closing retains the last
acknowledged draft; `fungi reopen SESSION` restores it in a new session.
Submitting returns those values plus a `comparison` PNG artifact exported
from the canvas. Errors keep the form open for retry. The preview is drawn
locally and never invokes an agent while the user adjusts it.

Follow the launch receipt's delivery mode. `auto` can return a manual fallback;
use the user’s “done” message to collect once in that case. See
`fungi describe delivery` for host-managed waiting. Never poll.
