# Committed baseline

The UI layout, question normalization and rich description renderer come from
jstack commit `7a7c919` (2026-10-13), the last committed change to these files
available before 2026-10-20. The repository's tip as of the end of 2026-10-19 was
`32324e47aae7d3ad3f92e6a4aad8966e9451d691`; its grill-form files match this baseline.
The ongoing uncommitted review/editor rework was not used.

Fungi replaces browser downloads, sessionStorage and a separate collector with
`ready`, `saveDraft`, `submit` and `cancel`. Questions arrive through `--input`;
answers and an optional image return through the normal result envelope.
The engine and bridge need no grill-form-specific API.

This preserves ordinary choices, notes and offline rich descriptions. It does
not include the newer visual choices, prototypes, threaded comments or rich
notes editor. One optional shared image exercises Fungi's existing attachment
contract without adding a review model.

The [bundled renderer record](../assets/vendor/README.md) pins versions and
hashes. Their licenses remain beside the JavaScript. Mermaid Tiny is retained
for this baseline; Pele is assessed separately, not silently substituted.
