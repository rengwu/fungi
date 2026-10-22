# Fungi portable package

Fungi runs local interfaces authored by agent skills. This package includes
the Go CLI and Electron runtime for one OS/CPU. Keep the directory together.
No Go, Node, npm or separate browser is needed to run it on a supported host.

```sh
./fungi --version
./fungi describe html
./fungi open examples/review.html
```

On Windows use `fungi.exe`. For custom workflows, supply your own HTML/assets
or a trusted native helper. `fungi describe` lists the CLI's capabilities.

To create a managed installation from a downloaded release archive:

```sh
./fungi install --archive PATH_TO_ARCHIVE --sha256 EXPECTED_SHA256
```

The default stable command is `~/.fungi/bin/fungi` (`fungi.exe` on Windows).
Add its directory to PATH if needed. Managed installations support
`fungi update`, `fungi rollback` and `fungi prune`. Read
`fungi describe installation` for details. Keep session data outside the
installation; `FUNGI_STATE_DIR` can select a separate location.

This remains a development preview. Public release publication, signing,
notarization and clean consumer-machine acceptance are not established by
creating an archive. An updater needs a published release or an explicit
local archive/checksum. Linux support depends on the documented host baseline;
Fungi does not install OS libraries or change sandbox policy.

Documentation, skills, source and verification records:
[rengwu/fungi](https://github.com/rengwu/fungi).
