---
title: Project icon
updated: 2026-11-19
---
# Project icon

## Accepted direction

**Apparition**, chosen on 2026-11-19 from a reference image the user supplied:
a glowing white-to-lavender ghost with two dark oval eyes, standing in a
pink-to-blue neon portal ring on a near-black squircle lit by a violet halo.
A large four-point sparkle sits upper right, a small one at left and a dot at
lower right. The ghost's scalloped hem fades so the back of the ring shows
through it; the front of the ring passes over the hem. It still reads as a
ghost at 16px.

It replaces **Pixel** (rounds 11 and 12), a beige device whose green
dot-matrix screen showed a pixel spark. The icon does not depict fungi and
need not relate to the name. Previews check it as a GitHub avatar, README
header and 16px favicon.

## Files

`docs/assets/icon/` holds the master `fungi.svg`, PNG exports at 1024, 512,
256, 128, 64, 32 and 16px, and a `favicon.ico` with 16, 32 and 48px PNG
frames. Chromium and WebKit render the master identically. The 1024px PNG is a
headless Chromium screenshot with a transparent background, for example
`headless_shell --headless --default-background-color=00000000
--window-size=1024,1024 --screenshot=fungi-1024.png file://$PWD/fungi.svg`;
the smaller sizes are `sips -z` downscales of it. The headless 0.2 executable
embeds only `docs/manual/*.md`. Its packages copy only `docs/manual` and
`docs/PORTABLE.md` from the documentation tree.

`fungi-mac.svg` places the same art on the macOS icon grid (824px on a 1024
canvas) with a soft shadow. The README floats it right of the introduction at
34% width, like chartr's README, and clears it before the first code block.

## Desktop release

The macOS desktop renderer ships this artwork as `Fungi.icns` starting with
[v0.1.0-alpha.3](https://github.com/rengwu/fungi/releases/tag/v0.1.0-alpha.3).
Updating a temporary alpha.2 installation from GitHub and opening a Fungi
window confirmed that macOS reports the ghost as the running app's icon.
The [release checks](https://github.com/rengwu/fungi/actions/runs/44383046669)
also verify the icon in the installed bundle after installation, update and rollback.

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
