# Complete helper example

This headless Python example echoes input, saves a draft and returns a text
artifact. It uses only Python's standard library; Python is a requirement of
this example, not of Fungi. All human-facing UI belongs to another project.

Save as `helper.py`:

```python
import json
import pathlib
import sys

launch = json.loads(sys.stdin.readline())
if launch["protocol"] != 2:
    raise RuntimeError("Unsupported Fungi protocol")

def send(message):
    print(json.dumps(message), flush=True)

def persist(message):
    send(message)
    ack = json.loads(sys.stdin.readline())
    if ack.get("id") != message["id"] or not ack.get("ok"):
        raise RuntimeError(str(ack))
    return ack

send({"type": "ready"})
state = launch["draft"] if launch["draft"] is not None else launch["input"]
persist({"id": "save", "type": "draft", "data": state})
pathlib.Path(launch["outputDir"], "note.txt").write_text("Example artifact", encoding="utf-8")
persist({"id": "finish", "type": "result", "outcome": "submitted", "data": state,
         "artifacts": [{"key": "note", "name": "note.txt", "path": "note.txt",
                        "mediaType": "text/plain"}]})
```

Launch with an input file containing any JSON:

```sh
fungi open --exec python3 --input ./input.json --delivery wait -- /absolute/helper.py
```

Use your Python executable's actual name or absolute path (`python.exe` on many
Windows installations). This example completes immediately; an interactive
helper should use `auto` unless the host supports a pending human wait.

Expect `protocol:2`, `kind:result`, `outcome:submitted`, the original JSON in
`data`, and a `note` artifact with an absolute saved path and SHA-256. No model
callback is sent by this `wait` example. For manual/callback sessions, follow
the receipt and collect once on completion; see `fungi help delivery`.

The source checkout also includes `examples/helper/main.go`:

```sh
go build -o /tmp/fungi-helper ./examples/helper
fungi open --exec /tmp/fungi-helper --delivery wait
```

These examples exercise the protocol without a UI. UI projects implement the
same pipe contract and define their own input, response and asset schemas.
