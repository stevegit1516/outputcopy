package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

var Version = "1.0.0"

const DefaultLimit = 100000

func main() {
	trim := flag.Bool("t", false, "Trim whitespace before copying")
	lineNums := flag.Bool("n", false, "Add line numbers")
	primary := flag.Bool("p", false, "Use Linux Primary selection")
	ver := flag.Bool("version", false, "Print version")
	fileMode := flag.Bool("file", false, "Copy a file instead of running a command")
	maxBytes := flag.Int("limit", getEnvInt("RCOPY_MAX_BYTES", DefaultLimit), "Maximum clipboard bytes")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: outputcopy [flags] COMMAND [ARGS...]")
		fmt.Fprintln(os.Stderr, "       outputcopy --file FILE")
		fmt.Fprintln(os.Stderr, "       COMMAND | outputcopy")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *ver {
		fmt.Printf("outputcopy version %s\n", Version)
		return
	}
	if *maxBytes < 1 {
		fmt.Fprintln(os.Stderr, "[outputcopy] limit must be positive")
		os.Exit(2)
	}

	args := flag.Args()
	if *fileMode {
		if len(args) != 1 {
			flag.Usage()
			os.Exit(2)
		}
		f, err := os.Open(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, "[outputcopy]", err)
			os.Exit(1)
		}
		defer f.Close()
		copyInput(f, *maxBytes, *trim, *lineNums, *primary)
		return
	}
	if len(args) > 0 {
		os.Exit(runCommand(args, *maxBytes, *trim, *lineNums, *primary))
	}

	stat, err := os.Stdin.Stat()
	if err == nil && stat.Mode()&os.ModeCharDevice != 0 {
		flag.Usage()
		os.Exit(2)
	}
	copyInput(os.Stdin, *maxBytes, *trim, *lineNums, *primary)
}

type captureWriter struct {
	mu        sync.Mutex
	buf       bytes.Buffer
	limit     int
	truncated bool
}

func (w *captureWriter) Write(p []byte) (int, error) {
	n, err := os.Stdout.Write(p)
	w.mu.Lock()
	defer w.mu.Unlock()

	remaining := w.limit - w.buf.Len()
	take := n
	if take > remaining {
		take = remaining
	}
	if take > 0 {
		_, _ = w.buf.Write(p[:take])
	}
	if n > take {
		w.truncated = true
	}
	if err == nil && n != len(p) {
		err = io.ErrShortWrite
	}
	return n, err
}

func runCommand(args []string, limit int, trim, lineNums, primary bool) int {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	captured := &captureWriter{limit: limit}
	cmd.Stdout, cmd.Stderr = captured, captured
	err := cmd.Run()

	code := 0
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			code = exitErr.ExitCode()
			if code < 0 {
				code = 1
			}
		} else {
			fmt.Fprintln(os.Stderr, "[outputcopy] cannot start command:", err)
			return 127
		}
	}

	captured.mu.Lock()
	text := transformText(captured.buf.String(), trim, lineNums)
	truncated := captured.truncated
	captured.mu.Unlock()

	if text != "" {
		emitClipboard(text, primary)
	}
	if truncated {
		fmt.Fprintf(os.Stderr, "\n[outputcopy] Clipboard limited to %d bytes; terminal output is complete.\n", limit)
	}
	return code
}

func copyInput(input io.Reader, limit int, trim, lineNums, primary bool) {
	content, err := io.ReadAll(io.LimitReader(input, int64(limit)+1))
	if err != nil {
		fmt.Fprintln(os.Stderr, "[outputcopy] read error:", err)
		os.Exit(1)
	}
	truncated := len(content) > limit
	if truncated {
		content = content[:limit]
	}
	text := transformText(string(content), trim, lineNums)
	if text != "" {
		emitClipboard(text, primary)
	}
	if truncated {
		fmt.Fprintf(os.Stderr, "[outputcopy] Clipboard limited to %d bytes.\n", limit)
	}
}

func transformText(text string, trim, lineNums bool) string {
	if trim {
		text = strings.TrimSpace(text)
	}
	if lineNums {
		text = addLineNumbers(text)
	}
	return text
}

func emitClipboard(text string, primary bool) {
	target := "c"
	if primary {
		target = "p"
	}
	encoded := base64.StdEncoding.EncodeToString([]byte(text))
	fmt.Printf("\033]52;%s;%s\a", target, encoded)
}

func addLineNumbers(text string) string {
	lines := strings.Split(text, "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	digits := len(strconv.Itoa(len(lines)))
	format := fmt.Sprintf("%%%dd %%s", digits)
	var out strings.Builder
	out.Grow(len(text) + len(lines)*(digits+2))
	for i, line := range lines {
		if i > 0 {
			out.WriteByte('\n')
		}
		fmt.Fprintf(&out, format, i+1, line)
	}
	if strings.HasSuffix(text, "\n") {
		out.WriteByte('\n')
	}
	return out.String()
}

func getEnvInt(key string, fallback int) int {
	if val, ok := os.LookupEnv(key); ok {
		if n, err := strconv.Atoi(val); err == nil {
			return n
		}
	}
	return fallback
}
