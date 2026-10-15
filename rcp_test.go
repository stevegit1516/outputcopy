package main

import (
	"os"
	"testing"
)

func TestTransformText(t *testing.T) {
	tests := []struct {
		name   string
		in     string
		trim   bool
		number bool
		want   string
	}{
		{"unchanged", "alpha\nbeta\n", false, false, "alpha\nbeta\n"},
		{"trim", " \nalpha\n ", true, false, "alpha"},
		{"number", "alpha\nbeta\n", false, true, "1 alpha\n2 beta\n"},
		{"trim and number", " alpha \n beta ", true, true, "1 alpha \n2  beta"},
		{"empty", "", false, false, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := transformText(tt.in, tt.trim, tt.number)
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestAddLineNumbers(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"", ""},
		{"one", "1 one"},
		{"one\ntwo", "1 one\n2 two"},
		{"one\ntwo\n", "1 one\n2 two\n"},
	}
	for _, tt := range tests {
		if got := addLineNumbers(tt.in); got != tt.want {
			t.Errorf("addLineNumbers(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestGetEnvInt(t *testing.T) {
	const key = "OUTPUTCOPY_TEST_LIMIT"
	t.Setenv(key, "123")
	if got := getEnvInt(key, 10); got != 123 {
		t.Fatalf("got %d, want 123", got)
	}
	t.Setenv(key, "invalid")
	if got := getEnvInt(key, 10); got != 10 {
		t.Fatalf("invalid value: got %d, want fallback 10", got)
	}
	os.Unsetenv(key)
	if got := getEnvInt(key, 10); got != 10 {
		t.Fatalf("unset value: got %d, want fallback 10", got)
	}
}

func TestAddLineNumbersBlankLines(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"one\n\ntwo\n", "1 one\n2 \n3 two\n"},
		{"\n", "1 \n"},
		{"  spaced  \n", "1   spaced  \n"},
	}
	for _, tt := range tests {
		if got := addLineNumbers(tt.in); got != tt.want {
			t.Errorf("addLineNumbers(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestTransformTrimOnlyOuterWhitespace(t *testing.T) {
	got := transformText(" \n alpha \n beta \n ", true, false)
	want := "alpha \n beta"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
