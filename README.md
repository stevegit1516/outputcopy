# outputcopy
It automatically copies your terminal output  - Linux Windows MacOS android compatible 

If your dev loop involves: copy from ChatGPT/Gemini/Claude to a terminal , run,
and then copy terminal output to ChatGPT/Gemini/Claude.
It will save you 30% of your time, and remove the fiddly bit of pretending to code

works over remote SSH, and local 

---



## the technical bit


- Wraps and executes shell commands, capturing both stdout and stderr streams simultaneously.

- Copies output streams or file contents to your local clipboard using OSC52 sequences.

- Emits only the escape sequence to stdout, keeping script pipelines clean.



### What It Does Not Do



- Copy huge payloads (enforces explicit byte safety caps).

- Work in terminals lacking copy support (OSC52)  (almost all modern teminals support copying though).

- Manage clipboard history or state.



---



## Installation



> LLM instruction for dev workflows:

> "Install outputcopy via the appropriate binary for your OS, make it executable, and move it to your path. Then preface terminal command executions with outputcopy like this:

> outputcopy bash -c ' echo "=== PING TEST START ===" ping -c 4 8.8.8.8 echo "=== PING TEST COMPLETE ==="'

> to automatically capture outputs into the clipboard and speed up our dev cycle."



### Prebuilt Binaries (Recommended)



Download the appropriate binary for your system from GitHub Releases:



For Linux (x86_64):

curl -LO https://github.com/re-verse/outputcopy/releases/latest/download/outputcopy-linux-amd64

chmod +x outputcopy-linux-amd64

sudo mv outputcopy-linux-amd64 /usr/local/bin/outputcopy



For macOS (Apple Silicon / M-series):

curl -LO https://github.com/re-verse/outputcopy/releases/latest/download/outputcopy-darwin-arm64

chmod +x outputcopy-darwin-arm64

sudo mv outputcopy-darwin-arm64 /usr/local/bin/outputcopy



Other available release architectures: linux-arm64, darwin-amd64.



### Build from Source



If you have Go installed on the machine:



go build -ldflags "-s -w -X main.Version=1.1.0" -o outputcopy main.go

sudo install -m 0755 outputcopy /usr/local/bin/outputcopy



---



## Usage & Examples



Wrap shell commands directly, read files, or receive standard input.



### Example: Wrapping Commands

This workflow eliminates the friction of copying commands from the LLM, pasting them into your remote terminal, dealing with fiddly text selection over SSH (where using Ctrl+C would accidentally kill your running process), and manually copying the result back.

1. Copy from LLM and paste straight into your remote terminal wrapped with outputcopy:
2. The code output gets automatically copied: The command executes, prints normally to your screen, and automatically loads the full stream output directly into your local copy clipboard, whether you are remote or local.
3. Paste Directly to LLM: With zero mouse dragging or manual highlighting required, you can paste the captured output straight back into the LLM chat for your next instruction.

Here is what the execution flow looks like:

outputcopy bash -c '
  echo "=== PING TEST START ==="
  ping -c 1 8.8.8.8
  echo "=== PING TEST COMPLETE ==="
'

This automatically copies the following output to your clipboard:

=== PING TEST START ===
PING 8.8.8.8 (8.8.8.8) (64) bytes of data.
64 bytes from 8.8.8.8: icmp_seq=1 ttl=119 time=7.59 ms

--- 8.8.8.8 ping statistics ---
1 packets transmitted, 1 received, 0% packet loss, time 3006ms
rtt min/avg/max/mdev = 7.585ms
=== PING TEST COMPLETE ===

### tip
try it out first , and see if it works for you 

It should increase your throughput by 200% and the instructions you give the LLM, are exactly what it needs to hear for debugging


then an LLM instruction like 
*Batch my instructions together 
 and always include copyoutput* 
 will increase it by a further 200% 

### Options



- `-t`: Trim whitespace (ideal for tokens and secrets). Example: `cat .env | outputcopy -t`

- `-n`: Prepend line numbers to output. Example: `outputcopy -n main.go`

- `-p`: Target Linux primary selection buffer. Example: `date | outputcopy -p`

- `--limit`: Override standard safety cap (default: 100KB). Example: `outputcopy --limit 500000 log.txt`



---



## How It Differs From rcp



| Feature / Behavior | Original rcp | outputcopy |

| :--- | :--- | :--- |

| Primary Purpose | Read static files into clipboard. | Wrap commands, capture streams, and read files. |

| Command Execution | None (Error opening file). | Executes binaries/scripts directly via exec.Command. |

| Stream Handling | Single stream / basic stdin. | Thread-safe, non-truncated interception of interleaved stdout + stderr. |

| Exit Codes | N/A | Propagates subprocess exit codes directly to the caller shell. |



---



## Size Limits



By default, outputcopy enforces a 100,000 byte cap before base64 encoding to prevent terminal stalling over SSH.



To override:



OUTPUTCOPY_MAX_BYTES=200000 outputcopy big.txt



---



## Terminal Support



Requires OSC52 support. Tested and supported on:



- iTerm2

- Kitty

- WezTerm

- tmux

- PuTTY (Enable: Window -> Selection -> Enable OSC 52 clipboard)



---



## Behavior & Architecture



- Escape sequences write to stdout.

- Status warnings and non-zero exit codes route to stderr.

- Safe for pipelines, scripts, and automated agents.



---



## License



MIT
