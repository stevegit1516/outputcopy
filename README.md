# outputcopy



Copy full command output text directly to your local clipboard via OSC52 (works over SSH, tmux, and modern terminals).



`outputcopy` is a lightweight Go CLI tool evolved from `rcp`. It wraps commands, processes files, and captures live terminal streams from remote shell sessions without requiring X forwarding, SCP, or manual text selection.



It eliminates friction during remote server management, vibecoding, and agentic LLM workflows by allowing complete execution blocks to be programmatically captured and pasted into models or documents.



## What It Does



- Wraps and executes shell commands, capturing both `stdout` and `stderr` streams simultaneously.

- Copies output streams or file contents to your local clipboard using OSC52 sequences.

- Emits only the escape sequence to `stdout`, keeping script pipelines clean.



### What It Does Not Do



- Copy huge payloads (enforces explicit byte safety caps).

- Work in terminals lacking OSC52 support.

- Manage clipboard history or state.



## How It Differs From `rcp`



| Feature / Behavior | Original `rcp` | `outputcopy` |

| --- | --- | --- |

| **Primary Purpose** | Read static files into clipboard. | Wrap commands, capture streams, and read files. |

| **Command Execution** | None (`Error opening file`). | Executes binaries/scripts directly via `exec.Command`. |

| **Stream Handling** | Single stream / basic stdin. | Thread-safe, non-truncated interception of interleaved `stdout` + `stderr`. |

| **Exit Codes** | N/A | Propagates subprocess exit codes directly to the caller shell. |



## Installation



### Prebuilt Binary



curl -LO https://github.com/re-verse/outputcopy/releases/latest/download/outputcopy-linux-amd64

chmod +x outputcopy-linux-amd64

sudo mv outputcopy-linux-amd64 /usr/local/bin/outputcopy



Available architectures: linux-amd64, linux-arm64, darwin-amd64, darwin-arm64.



### Build from Source



go build -ldflags "-s -w -X main.Version=1.1.0" -o outputcopy main.go

sudo install -m 0755 outputcopy /usr/local/bin/outputcopy



## Usage & Examples



Wrap shell commands directly, read files, or receive standard input.



### Example: Wrapping Execution Streams



Capture complex multi-line output directly to clipboard without touch interaction:



outputcopy bash -c '

  echo "=== PING TEST START ==="

  ping -c 4 8.8.8.8

  echo "=== PING TEST COMPLETE ==="

'



### Options



| Flag | Description | Example |

| --- | --- | --- |

| **`-t`** | Trim whitespace (ideal for tokens and secrets). | `cat .env \| outputcopy -t` |

| **`-n`** | Prepend line numbers to output. | `outputcopy -n main.go` |

| **`-p`** | Target Linux primary selection buffer. | `date \| outputcopy -p` |

| **`--limit`** | Override standard safety cap (default: 100KB). | `outputcopy --limit 500000 log.txt` |



## Size Limits



By default, `outputcopy` enforces a 100,000 byte cap before base64 encoding to prevent terminal stalling over SSH.



To override:



RCOPY_MAX_BYTES=200000 outputcopy big.txt



## Terminal Support



Requires OSC52 support. Tested and supported on:



- iTerm2

- Kitty

- WezTerm

- tmux

- PuTTY (Enable: Window -> Selection -> Enable OSC 52 clipboard)



## Behavior & Architecture



- Escape sequences write to `stdout`.

- Status warnings and non-zero exit codes route to `stderr`.

- Safe for pipelines, scripts, and automated agents.



## License



MIT
