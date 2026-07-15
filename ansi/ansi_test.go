package ansi_test

import (
	"testing"

	"github.com/romanSPB15/tui-compose/v3/ansi"
)

func TestStrip(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no ANSI",
			input:    "Hello, World!",
			expected: "Hello, World!",
		},
		{
			name:     "simple color",
			input:    "\033[31mHello\033[0m",
			expected: "Hello",
		},
		{
			name:     "RGB foreground",
			input:    "\033[38;2;255;0;0mRed\033[0m",
			expected: "Red",
		},
		{
			name:     "background color",
			input:    "\033[48;2;0;0;255mBlue bg\033[0m",
			expected: "Blue bg",
		},
		{
			name:     "multiple sequences",
			input:    "\033[1;32mBold green\033[0m and \033[44mblue bg\033[0m",
			expected: "Bold green and blue bg",
		},
		{
			name:     "escape sequences with parameters",
			input:    "\033[?1006h", // SGR mouse
			expected: "",
		},
		{
			name:     "CSI sequences",
			input:    "\033[K", // clear line
			expected: "",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only escape",
			input:    "\033[0m",
			expected: "",
		},
		{
			name:     "Russian text",
			input:    "\033[31mПривет\033[0m",
			expected: "Привет",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ansi.Strip(tt.input)
			if got != tt.expected {
				t.Errorf("Strip(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestFind(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedLen   int
		expectedClean string
		checkSequence func(t *testing.T, matches []ansi.AnsiMatch)
	}{
		{
			name:          "no ANSI",
			input:         "Hello, World!",
			expectedLen:   0,
			expectedClean: "Hello, World!",
		},
		{
			name:          "single sequence",
			input:         "\033[31mHello\033[0m",
			expectedLen:   2,
			expectedClean: "Hello",
			checkSequence: func(t *testing.T, matches []ansi.AnsiMatch) {
				if matches[0].Seq != "\033[31m" {
					t.Errorf("expected first sequence %q, got %q", "\033[31m", matches[0].Seq)
				}
				if matches[0].Index != 0 {
					t.Errorf("expected first index 0, got %d", matches[0].Index)
				}
				if matches[1].Seq != "\033[0m" {
					t.Errorf("expected second sequence %q, got %q", "\033[0m", matches[1].Seq)
				}
				if matches[1].Index != 10 {
					t.Errorf("expected second index 10, got %d", matches[1].Index)
				}
			},
		},
		{
			name:          "RGB with Unicode",
			input:         "\033[38;2;255;0;0mКрасный\033[0m",
			expectedLen:   2,
			expectedClean: "Красный",
			checkSequence: func(t *testing.T, matches []ansi.AnsiMatch) {
				if matches[0].Index != 0 {
					t.Errorf("expected first index 0, got %d", matches[0].Index)
				}
				if matches[1].Index != 22 {
					t.Errorf("expected second index 22, got %d", matches[1].Index)
				}
			},
		},
		{
			name:          "only sequences",
			input:         "\033[31m\033[0m",
			expectedLen:   2,
			expectedClean: "",
		},
		{
			name:          "empty string",
			input:         "",
			expectedLen:   0,
			expectedClean: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches, clean := ansi.Find(tt.input)
			if len(matches) != tt.expectedLen {
				t.Errorf("expected %d matches, got %d", tt.expectedLen, len(matches))
			}
			if clean != tt.expectedClean {
				t.Errorf("clean string = %q, want %q", clean, tt.expectedClean)
			}
			if tt.checkSequence != nil && len(matches) > 0 {
				tt.checkSequence(t, matches)
			}
		})
	}
}

func TestFindEdgeCases(t *testing.T) {
	input := "a\033[31mb\033[0mc"
	matches, clean := ansi.Find(input)
	if len(matches) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(matches))
	}
	if clean != "abc" {
		t.Errorf("clean = %q, want %q", clean, "abc")
	}

	if matches[0].Index != 1 {
		t.Errorf("first match index = %d, want 1", matches[0].Index)
	}
	if matches[1].Index != 7 {
		t.Errorf("second match index = %d, want 7", matches[1].Index)
	}
}
