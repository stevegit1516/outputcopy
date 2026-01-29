package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Version can be injected at build time: go build -ldflags "-X main.Version=1.0.0"
var Version = "1.0.0"

// Default limit prevents blowing up the terminal (100KB)
// Override with env var RCOPY_MAX_BYTES
const DefaultLimit = 100000

func main() {
	// Flags
	var (
		trim      = flag.Bool("t", false, "Trim leading/trailing whitespace")
		lineNums  = flag.Bool("n", false, "Add line numbers to the output")
		primary   = flag.Bool("p", false, "Use Linux Primary selection (middle-click) instead of Clipboard")
		ver       = flag.Bool("version", false, "Print version and exit")
		maxBytes  int
	)

	// Custom usage message for 1.0 polish
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: rcp [flags] [file]\n\n")
		fmt.Fprintf(os.Stderr, "Copy input to the remote clipboard via OSC52.\n")
		fmt.Fprintf(os.Stderr, "Reads from stdin if no file is specified.\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}

	flag.IntVar(&maxBytes, "limit", getEnvInt("RCOPY_MAX_BYTES", DefaultLimit), "Max bytes to copy (env RCOPY_MAX_BYTES)")
	flag.Parse()

	if *ver {
		fmt.Printf("rcp version %s\n", Version)
		os.Exit(0)
	}

	// 1. Determine Input Source
	var input io.Reader
	var err error

	if flag.NArg() > 0 {
		// Read from file argument
		f, err := os.Open(flag.Arg(0))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		input = f
	} else {
		// Read from Stdin
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) != 0 {
			// User ran 'rcp' with no args and no pipe. Show help instead of hanging.
			flag.Usage()
			os.Exit(1)
		}
		input = os.Stdin
	}

	// 2. Read Content (up to limit)
	// We read everything first to apply transforms (Trim, LineNums)
	// For massive files, the limit kicks in during read to save memory.
	reader := io.LimitReader(input, int64(maxBytes))
	content, err := io.ReadAll(reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	if len(content) == 0 {
		return // Nothing to copy
	}

	// 3. Process Content
	text := string(content)

	if *trim {
		text = strings.TrimSpace(text)
	}

	if *lineNums {
		text = addLineNumbers(text)
	}

	// 4. Encode & Output OSC52
	target := "c" // Default: Clipboard
	if *primary {
		target = "p" // Linux Primary (Middle-click)
	}

	encoded := base64.StdEncoding.EncodeToString([]byte(text))
	
	// \033]52;target;content\a
	fmt.Printf("\033]52;%s;%s\a", target, encoded)

	// Optional: Print a friendly stderr message if we hit the limit
	if len(content) == maxBytes {
		fmt.Fprintf(os.Stderr, "[rcp] Warning: Input truncated at %d bytes.\n", maxBytes)
	}
}

// Helper to number lines efficiently
func addLineNumbers(text string) string {
	lines := strings.Split(text, "\n")
	// Handle trailing newline edge case (split creates empty string at end)
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	// Calculate padding width based on line count
	digitCount := len(strconv.Itoa(len(lines)))
	format := fmt.Sprintf("%%%dd %%s", digitCount) // e.g. "%3d %s"

	var sb strings.Builder
	// Pre-allocate approximate size to avoid re-allocations
	sb.Grow(len(text) + (len(lines) * (digitCount + 2)))

	for i, line := range lines {
		if i > 0 {
			sb.WriteByte('\n')
		}
		fmt.Fprintf(&sb, format, i+1, line)
	}
	
	// Preserve original trailing newline if it existed
	if strings.HasSuffix(text, "\n") {
		sb.WriteByte('\n')
	}

	return sb.String()
}

// Helper to get env var as int
func getEnvInt(key string, defaultVal int) int {
	if val, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}
