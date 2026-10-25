# Fungi portable package

Fungi runs local interfaces authored by agent skills. This package includes
the Go CLI and Wails window helper for one OS/CPU. Keep the directory together.
No Go, Node or npm is needed. The OS supplies WKWebView (macOS), WebView2
(Windows) or WebKitGTK 4.1 (Linux). Run `fungi doctor` to check availability.

```sh
./fungi --version
./fungi help
./fungi help examples
./fungi describe examples
```

On Windows use `fungi.exe`. For custom workflows, supply your own HTML/assets
or a trusted native helper. [Start here](docs/manual/start.md), or read any
reference offline with `fungi help TOPIC`. `fungi describe` lists capabilities;
`fungi describe examples` returns absolute example paths and launch arguments
from any working directory. Draft saving is optional for custom HTML.

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
normal launches never install OS libraries or change sandbox policy. The
bootstrap offers explicit runtime setup with `--setup-runtime` on supported
Ubuntu 22.04/24.04, or `-SetupRuntime` on Windows. Other distributions need
their own WebKitGTK 4.1 setup instructions. Linux direct media URLs are a known
limitation; browser storage isolation is best effort and pages have native access.

Documentation, skills, source and verification records:
[rengwu/fungi](https://github.com/rengwu/fungi).
