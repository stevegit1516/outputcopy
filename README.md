# rcp

Copy remote text to your local clipboard via OSC52 (works with SSH, tmux, and modern terminals).

rcp is a small CLI tool for copying text from a remote shell session directly into your local clipboard.
It works over SSH and tmux without X forwarding, SCP hacks, or manual selection.

If you spend time on remote servers, this saves friction every day.

---

## What it does

- Copy a file’s contents to your local clipboard
- Copy command output (date, ls, logs, etc.)
- Works over SSH and tmux using OSC52
- Writes only the OSC52 escape sequence to stdout (safe for scripting)

What it does not do:

- Copy huge files (by design)
- Work in terminals without OSC52 support
- Act as a general clipboard manager

---

## Installation

### Prebuilt binaries (recommended)

Download the appropriate binary from GitHub Releases:

```bash
curl -LO [https://github.com/re-verse/rcp/releases/latest/download/rcp-linux-amd64](https://github.com/re-verse/rcp/releases/latest/download/rcp-linux-amd64)
chmod +x rcp-linux-amd64
sudo mv rcp-linux-amd64 /usr/local/bin/rcp
```

Available builds typically include:

- linux-amd64
- linux-arm64
- darwin-amd64
- darwin-arm64

No Go installation is required to run the binary.

### Build from source (Go)

You only need Go to build, not to run.

```bash
# Build with version tag
go build -ldflags "-s -w -X main.Version=1.0.0" -o rcp rcp.go
sudo install -m 0755 rcp /usr/local/bin/rcp
```

Cross-compile from macOS to Linux:

```bash
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X main.Version=1.0.0" -o rcp-linux-amd64 rcp.go
```

---

## Usage

`rcp` reads from stdin by default, or from a file if provided as an argument.

| Flag | Description | Example |
| :--- | :--- | :--- |
| **`-t`** | **Trim whitespace.** Removes leading/trailing whitespace. Essential for secrets. | `cat .env \| rcp -t` |
| **`-n`** | **Add line numbers.** Formats output with line numbers. Great for sharing code. | `rcp -n main.go` |
| **`-p`** | **Primary Selection.** (Linux) Copies to the middle-click buffer. | `date \| rcp -p` |
| **`--limit`** | **Set byte limit.** Override the 100KB safety cap. | `rcp --limit 500000 log.txt` |

### Examples

**Copy a file:**
```bash
rcp file.txt
```

**Copy command output:**
```bash
date | rcp
```

**Copy a clean token (Trim):**
```bash
# Strips the trailing newline so you can paste into password fields
cat token.txt | rcp -t
```

**Share code with context (Line Numbers):**
```bash
rcp -n config.yaml
```

---

## Size limits

By default, rcp refuses to copy more than 100,000 bytes (before base64 encoding).

This is intentional:

- OSC52 payloads are sent through the terminal
- Large payloads are slow and frequently truncated

Override if needed:

```bash
RCOPY_MAX_BYTES=200000 rcp big.txt
```

---

## Terminal support

rcp requires OSC52 clipboard support.

Known to work in:

- iTerm2
- Kitty
- WezTerm
- tmux
- PuTTY (with OSC52 enabled)

### PuTTY users

Enable:

    Window → Selection → Enable OSC 52 clipboard

Then save the session.

---

## Bash version

A Bash implementation is included in:

    scripts/rcp.sh

This version is provided for reference or for environments where installing a compiled binary is not desirable.

The Go binary is the primary, supported implementation.
The Bash version is provided as-is and may lag behind feature-wise.

---

## Exit behavior

- OSC52 escape sequence is written to stdout
- Status and errors are written to stderr

This makes rcp safe to use in pipelines and scripts.

---

## Non-goals

- Copying megabytes of data
- Supporting terminals without OSC52
- Being clever about clipboard history or synchronization

This tool is intentionally small and boring.

---

## License

MIT
