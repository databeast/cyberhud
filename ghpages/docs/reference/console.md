# Console Protocol

The CyberHUD daemon (`cyberhudd`) exposes a control interface via a Unix domain socket. The `cyberhudctl` CLI tool and any custom tooling communicate with the daemon through this socket using a line-oriented text protocol.

## Connection

### Socket Path

The default socket path is:

```
/run/cyberhudd/console.sock
```

The path is configurable in two ways:

- **Daemon side**: Set the `"socket"` field in the JSON configuration file passed to `cyberhudd`.
- **Client side**: Use the `-socket` flag when invoking `cyberhudctl`:

```bash
cyberhudctl -socket /path/to/custom.sock status
```

### Socket Permissions

The daemon creates the socket with mode `0600`. If a `cyberhud` system group exists, the socket is widened to `0660` and ownership is set to that group, allowing any member of the group to issue commands.

### Connecting

Clients connect using a standard Unix domain socket (`AF_UNIX`, `SOCK_STREAM`). In Go:

```go
conn, err := net.Dial("unix", "/run/cyberhudd/console.sock")
```

Upon successful connection, the daemon immediately sends a greeting line:

```
OK cyberhud daemon ready
```

The client should read and discard (or verify) this greeting before sending commands.

## Protocol Format

The protocol is **line-oriented UTF-8 text**. Each message (request or response) is terminated by a newline character (`\n`).

### Request Format

```
<verb> [<sub-command>] [<args>]\n
```

- Commands are case-insensitive.
- Arguments are space-separated.
- Empty lines are silently ignored.

### Response Format

Every response begins with a status prefix on the first line:

| Prefix | Meaning |
|--------|---------|
| `OK`   | Command succeeded. Payload follows the prefix on the same or subsequent lines. |
| `ERR`  | Command failed. An error message follows the prefix. |

Single-line response examples:

```
OK stemma_devices=2 gpio_pins=4 display_regions=3
ERR unknown verb
```

Multi-line responses place `OK` on the first line, followed by indented data lines:

```
OK
  region=main.0 name="main" controller=st7789 mode=clock modes=clock,system,ticker
  region=left-aux.0 name="left-aux" controller=ssd1680 mode=thermal modes=thermal,gpio
```

### Session Lifecycle

1. Client connects via Unix socket.
2. Daemon sends greeting: `OK cyberhud daemon ready\n`
3. Client sends one or more command lines.
4. Daemon replies with a response for each command.
5. Client sends `quit` or `exit` to close the session gracefully.
6. Daemon responds `OK bye\n` and closes the connection.

Alternatively, the client may simply close the connection at any time.

## Command Overview

The following verbs are accepted by the protocol:

| Verb | Description |
|------|-------------|
| `status` | One-line daemon summary (device counts, region count) |
| `gpio` | GPIO pin control sub-commands |
| `stemma` | STEMMA QT / QWIIC device queries |
| `display` | Display region and mode control |
| `policy` | Policy store queries |
| `freeze` | Persist configuration to disk |
| `config` | Runtime configuration dump |
| `help` | Command metadata queries |
| `quit` / `exit` | Close the connection |

For full command syntax and usage examples, see the [CLI Reference](cli.md).

## Error Responses

Error responses always begin with `ERR` followed by a human-readable message:

```
ERR unknown verb
ERR usage: display set <region> <mode> [key=value ...]
ERR unknown region "foo.0"; available: main.0, left-aux.0, right-aux.0
ERR unsupported mode "badmode" for region main.0
ERR speed: must be in [0.1, 10.0]
```

Policy validation errors include the field name and the rejection reason, making it straightforward to identify which parameter was invalid.

## Timeouts

The `cyberhudctl` client uses a configurable timeout (default: 2 seconds) for socket read/write operations via the `-timeout` flag:

```bash
cyberhudctl -timeout 5s display regions
```

The daemon itself does not enforce idle timeouts on connections — clients may hold a connection open indefinitely.

## Multi-Command Sessions

A single socket connection can be reused for multiple sequential commands. The client sends each command line, reads its response, and then sends the next. The `cyberhudctl` tool supports this via semicolon-delimited multi-command syntax:

```bash
cyberhudctl region main.0 ';' mode clock ';' config format=24h
```

Each semicolon-separated group is resolved and sent as a separate protocol command over the same connection. Execution stops on the first error response.

## Implementation Notes

- Each client connection is served in its own goroutine.
- The server removes any stale socket file on startup before binding.
- The `Dial` and `SendCommand` helpers in the `runtime/console` package provide convenient programmatic access for tests and tooling.
