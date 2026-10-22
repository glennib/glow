package main

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestGlowFlags(t *testing.T) {
	tt := []struct {
		args  []string
		check func() bool
	}{
		{
			args: []string{"-p"},
			check: func() bool {
				return pager
			},
		},
		{
			args: []string{"-s", "light"},
			check: func() bool {
				return style == "light"
			},
		},
		{
			args: []string{"-w", "40"},
			check: func() bool {
				return width == 40
			},
		},
	}

	for _, v := range tt {
		err := rootCmd.ParseFlags(v.args)
		if err != nil {
			t.Fatal(err)
		}
		if !v.check() {
			t.Errorf("Parsing flag failed: %s", v.args)
		}
	}
}

func TestCLIPreserveNewLines(t *testing.T) {
	oldPager, oldTUI, oldStyle, oldWidth, oldPreserve := pager, tui, style, width, preserveNewLines
	t.Cleanup(func() {
		pager, tui, style, width, preserveNewLines = oldPager, oldTUI, oldStyle, oldWidth, oldPreserve
	})
	pager, tui, style, width = false, false, "notty", 80

	render := func(preserve bool) string {
		preserveNewLines = preserve
		var buf bytes.Buffer
		src := &source{reader: io.NopCloser(strings.NewReader("hello\nworld\n")), URL: "test.md"}
		if err := executeCLI(src, &buf); err != nil {
			t.Fatal(err)
		}
		return buf.String()
	}

	if out := render(false); !strings.Contains(out, "hello world") {
		t.Errorf("expected soft line break to be joined, got %q", out)
	}
	if out := render(true); strings.Contains(out, "hello world") || !strings.Contains(out, "world") {
		t.Errorf("expected soft line break to be preserved, got %q", out)
	}
}
