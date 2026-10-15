# outputcopy

Copy the full output text to your local clipboard via OSC52 (works with SSH, tmux, and modern terminals).

`outputcopy` is a small Go-based CLI tool derived from the foundational concepts of `rcp`, designed for copying command output, files, and streamed terminal outputs directly from a remote shell session into your local clipboard.
It works over SSH and tmux without X forwarding, SCP hacks, or manual text selection.

If you spend time on remote servers, this saves friction every day. 

It is ideal for vibecoding or advanced agentic development workflows as you no longer have to onerously highlight batches of text manually. By allowing you to wrap commands and capture full execution blocks programmatically, it enables seamless automated workflows and full-stream batching.

---

## What it does

- Wrap and execute shell commands while capturing both standard output (`stdout`) and error streams (`stderr`) simultaneously.
- Copy complete output streams or file contents to your local clipboard over SSH and tmux using OSC52.
- Write only the OSC52 escape sequence to stdout (keeping pipelines clean and script-safe).

What it does not do:

- Copy huge files (by design)
- Work in terminals without OSC52 support
- Act as a general clipboard manager

---

## How it differs from `rcp` in function

| Feature / Behavior | Original `rcp` | New `outputcopy` |
| :--- | :--- | :--- |
| **Primary Purpose** | Static file copying (treating arguments as file paths to read and copy). | General-purpose command wrapper, stream capturer, and file copier. |
| **Command Execution** | None. Attempting to pass a command resulted in `Error opening file: open [cmd]: no such file or directory`. | Executes any shell command via `exec.Command` and captures its exact execution lifecycle. |
| **Stream Handling** | Limited to single file or basic stdin streams. | Combines and captures **both `stdout` and `stderr`** synchronously into a shared thread-safe buffer without dropping lines or truncating interleaved messages. |
| **Exit Code Propagation** | N/A (file reader only). | Propagates the exact exit code of the wrapped subprocess back to the parent shell. |

---

## Installation

### Prebuilt binaries (recommended)

Download the appropriate binary from GitHub Releases:

```bash
curl -LO [https://github.com/re-verse/outputcopy/releases/latest/download/outputcopy-linux-amd64](https://github.com/re-verse/outputcopy/releases/latest/download/outputcopy-linux-amd64)
chmod +x outputcopy-linux-amd64
sudo mv outputcopy-linux-amd64 /usr/local/bin/outputcopy
