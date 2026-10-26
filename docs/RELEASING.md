# Publishing an alpha

An annotated `vX.Y.Z-alpha.N` tag triggers `.github/workflows/release.yml`.
Only maintainers should push release tags. The workflow publishes automatically
after its checks pass; stable releases are intentionally not automated yet.

1. Write `docs/releases/VERSION.md`, including install commands, changes,
   initial targets and material limitations. Update version-pinned README/manual
   examples and rebuild the wiki. Commit all release source and notes.
2. Push an annotated alpha tag pointing to that exact commit:

   ```sh
   git tag -a v0.1.0-alpha.1 -m 'Fungi v0.1.0-alpha.1'
   git push origin v0.1.0-alpha.1
   ```

3. Check the **Publish alpha** run. It reuses the foundation workflow to test
   the core, bridge, native packages, installer, update/rollback and Linux
   runtime setup. The exact Ubuntu 22 package must also pass on Ubuntu 24.
4. Publication selects only Mac arm64, Windows amd64 and the Ubuntu 22-built
   Linux amd64 archive, with SHA-256 sidecars and both installer scripts.
   It checks each BUILD.json against the tag's commit and requires clean source.
   The GitHub CLI uploads all assets before publishing the prerelease. There
   is no rebuilding between verification and publication.
5. The remaining jobs anonymously download the published installer, install to
   a temporary path with spaces, run doctor, fetch the same published version
   through `fungi update`, discover examples and exercise the real renderer.
   They do not modify persistent PATH or install runtimes on a user's machine.
   Preserve their evidence and record the release run in the wiki.

The release stays a prerelease and is not marked Latest. Bootstrap and updater
commands must select an alpha explicitly. Their default continues to mean
latest stable. This avoids introducing an update-channel service or registry.

Do not overwrite a published version's assets or move its tag: installed
versions are immutable. A correction needs a new alpha number. If upload fails
and leaves a draft, inspect it before removing that unpublished draft and
retrying. If post-publication checks fail, record the failure on the release
and fix it in the next alpha; do not silently replace the download.

These releases are unsigned and unnotarized. Clean consumer-machine setup and
physical desktop qualification are separate from CI. No signing keys, paid
services, OS policy changes or automatic runtime installs are part of this flow.

To repeat the public-download check on a machine with a ready desktop/runtime:

```sh
node scripts/public-release-smoke.mjs v0.1.0-alpha.1
# On Linux CI or a virtual display:
xvfb-run -a node scripts/public-release-smoke.mjs v0.1.0-alpha.1
```

The driver requires Node only for verification; users need neither Node nor Go.
It saves evidence under `verification-output/public-release-*` and deletes its
temporary installation. It opens synthetic test windows, never user interviews.

To check an existing download on a hosted platform with the current verification
driver, manually run **Publish alpha**, selecting its version and platform.
Manual runs skip validation/build/publication and perform only the public-download
check. For example:

```sh
gh workflow run release.yml --ref main -f version=v0.1.0-alpha.1 -f platform=macos-latest
```
