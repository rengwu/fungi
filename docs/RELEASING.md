# Publishing an engine alpha

An annotated `vX.Y.Z-alpha.N` tag triggers `.github/workflows/release.yml`.
Tag publication requires an explicit maintainer request; local implementation
and verification do not publish a release. The 0.2 engine line uses protocol 2
and package format 3; existing 0.1 alpha downloads remain historical UI builds.

1. Write `docs/releases/VERSION.md` with changes, migration, install commands,
   target platforms and observed limitations. Update current version examples,
   build/check the wiki and commit the reviewed source.
2. Push an annotated tag only when publication is authorized:

   ```sh
   FUNGI_RELEASE_VERSION=vX.Y.Z-alpha.N
   git tag -a "$FUNGI_RELEASE_VERSION" -m "Fungi $FUNGI_RELEASE_VERSION"
   git push origin "$FUNGI_RELEASE_VERSION"
   ```

3. The engine workflow tests Go, packages, offline discovery, helper sessions,
   artifacts, bootstrap, update and rollback on macOS, Windows and Ubuntu.
   It also executes the Ubuntu 22-built archive unchanged on Ubuntu 24.
   No desktop, webview, Xvfb or system package installation is required.
4. Publication selects Mac arm64, Windows amd64 and Ubuntu 22-built Linux amd64
   archives, SHA-256 sidecars and installers. BUILD.json must match the tag
   commit and clean source. The tested bytes are published without rebuilding.
5. Anonymous public-download jobs install the exact published version, check
   doctor, update from GitHub and execute external helper sessions. Preserve
   evidence and record the release result in the wiki.

Alpha versions must be selected explicitly; the default updater/bootstrap
version remains latest stable. Never replace a published version's assets or
move its tag. Unsigned binaries may be subject to host security policy; CI is
not evidence of signing or clean consumer-machine installation.

The scripts are source-controlled maintainer tools requiring Node 22+ and Go
1.26 for builds. End users need neither. Repeat a published engine check with:

```sh
node scripts/public-release-smoke.mjs vX.Y.Z-alpha.N
```

Manual **Publish alpha** runs accept an existing engine version and platform,
skip publication and run only the download check. Use the matching historical
checkout when verifying protocol-1 releases.

Preserve the [protocol 2 contract](manual/helper.md#limits-and-compatibility).
Optional metadata may be added; changing documented names, types or meanings
requires a new protocol. Keep session migration explicit and preserve old data.
