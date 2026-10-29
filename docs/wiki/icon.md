---
title: Project icon
updated: 2026-10-28
---
# Project icon

## Accepted direction

**Pixel**, chosen in round 11, in the **Larger screen** version chosen in
round 12: a small beige device whose green dot-matrix screen shows a 13×13
pixel spark. The screen fills 77% of the icon's width, the beige bezel is half
its round-11 thickness, and nothing sits on the bezel; the speaker dots and
status light were removed at the user's request.

[icon.museum](https://icon.museum/) is the current visual reference. The icon
does not depict fungi and need not relate to the name. Previews check it as a
GitHub avatar, README header and 16px favicon.

## Files

`docs/assets/icon/` holds the master `fungi.svg`, PNG exports at 1024, 512,
256, 128, 64, 32 and 16px, and a `favicon.ico` with 16, 32 and 48px frames.
Chrome and librsvg render the master identically. Regenerate a PNG with
`rsvg-convert -w 512 -h 512 fungi.svg -o fungi-512.png`. The executable embeds
only `docs/manual/*.md` and packages copy only `docs/manual` and
`docs/PORTABLE.md`, so these files ship in neither.

`fungi-mac.svg` places the same art on the macOS icon grid (824px on a 1024
canvas) with a soft shadow. The README floats it right of the introduction at
34% width, like chartr's README, and clears it before the first code block.

## Open questions

- The docs favicon does not use the icon yet.
- The GitHub organization avatar is uploaded by hand.

## Interview history

All rounds ran on 2026-10-28. Each journal keeps the full question text,
options and the user's notes.

1. [Round 1](journal/2026-10-28-grill-41ad8bda-7e7b-4d13-9cbe-e3d0642e419f.md):
   no mushrooms or literal fungi; the icon is a brand mark first.
2. [Round 2](journal/2026-10-28-grill-940e22d2-0b8e-4aa2-acb7-cbfcffea9d90.md):
   all rejected as too playful; asked for something more serious.
3. [Round 3](journal/2026-10-28-grill-e86152bb-d39e-40f7-834b-5bcc1ec4293a.md):
   asked for a hint of a living organism.
4. [Round 4](journal/2026-10-28-grill-846ecc70-e9cb-49db-8416-65e1f9169b42.md):
   asked to restudy [appicons.store](https://www.appicons.store/).
5. [Round 5](journal/2026-10-28-grill-f296ebf8-c4e9-4c94-8b31-00da532dea03.md):
   nothing memorable; keep it simple, it can be anything.
6. [Round 6](journal/2026-10-28-grill-b6f4accf-9e3d-4a44-9ada-d6fdd581efa1.md):
   liked Pebble and Pearl, asked for different eyes.
7. [Round 7](journal/2026-10-28-grill-d2e5563a-9f64-42f8-8670-2b491a602e66.md):
   restart from scratch with icon.museum as the reference.
8. [Round 8](journal/2026-10-28-grill-5b3e14b5-0e4e-4e6e-8d48-1c6d7168d86e.md):
   Enamel, Frost, Drops and Print rejected.
9. [Round 9](journal/2026-10-28-grill-9d51af1e-8ae9-4ac1-8c79-8aacb8de34f5.md):
   all new designs rejected.
10. [Round 10](journal/2026-10-28-grill-829c423c-d331-4185-8a41-843aa32c577c.md):
    kept designs were steering new ones; restart from scratch.
11. [Round 11](journal/2026-10-28-grill-3db9d420-d7cf-4a8e-aeaf-ee9123f104c1.md):
    Pixel chosen, with a larger screen and nothing on the bezel.
12. [Round 12](journal/2026-10-28-grill-59b957a6-21b0-4adc-88bc-a65a3307a110.md):
    Larger screen chosen, with no further changes.
