package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

// saveOptions restores the globals that validateOptions sets when the test ends.
func saveOptions(t *testing.T) {
	t.Helper()
	oldWidth, oldMouse, oldPager, oldTUI := width, mouse, pager, tui
	oldAll, oldPreserve, oldLineNumbers, oldStyle := showAllFiles, preserveNewLines, showLineNumbers, style
	t.Cleanup(func() {
		width, mouse, pager, tui = oldWidth, oldMouse, oldPager, oldTUI
		showAllFiles, preserveNewLines, showLineNumbers, style = oldAll, oldPreserve, oldLineNumbers, oldStyle
	})
}

func TestConfigFlag(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "custom.yml")
	if err := os.WriteFile(path, []byte("width: 42\npreserveNewLines: true\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name      string
		args      []string
		wantWidth uint
		wantError bool
	}{
		{name: "default config", wantWidth: 80},
		{name: "config flag", args: []string{"--config", path}, wantWidth: 42},
		{name: "missing config file", args: []string{"--config", filepath.Join(dir, "missing.yml")}, wantError: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			saveOptions(t)
			oldConfigFile := configFile
			t.Cleanup(func() { configFile = oldConfigFile })
			viper.Reset()
			t.Cleanup(viper.Reset)
			viper.SetDefault("style", "auto")
			viper.SetDefault("width", 80)

			cmd := &cobra.Command{Use: "glow"}
			cmd.Flags().StringVar(&configFile, "config", "", "")
			if err := cmd.ParseFlags(tt.args); err != nil {
				t.Fatal(err)
			}

			err := validateOptions(cmd)
			if tt.wantError {
				if err == nil {
					t.Fatal("expected an error for a missing config file")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if width != tt.wantWidth {
				t.Errorf("width=%d; want %d", width, tt.wantWidth)
			}
			if wantPreserve := tt.wantWidth == 42; preserveNewLines != wantPreserve {
				t.Errorf("preserveNewLines=%v; want %v", preserveNewLines, wantPreserve)
			}
		})
	}
}

func TestDisplayModePrecedence(t *testing.T) {
	tests := []struct {
		name      string
		config    string
		env       map[string]string
		args      []string
		wantPager bool
		wantTUI   bool
		wantError bool
	}{
		{name: "default modes"},
		{name: "configured pager", config: "pager: true", wantPager: true},
		{name: "configured tui", config: "tui: true", wantTUI: true},
		{name: "tui flag overrides configured pager", config: "pager: true", args: []string{"-t"}, wantTUI: true},
		{name: "long tui flag overrides configured pager", config: "pager: true", args: []string{"--tui=true"}, wantTUI: true},
		{name: "tui flag overrides env pager", env: map[string]string{"GLOW_PAGER": "true"}, args: []string{"-t"}, wantTUI: true},
		{name: "pager flag overrides configured tui", config: "tui: true", args: []string{"-p"}, wantPager: true},
		{name: "disabled tui keeps configured pager", config: "pager: true", args: []string{"--tui=false"}, wantPager: true},
		{name: "disabled pager", config: "pager: true", args: []string{"--pager=false"}},
		{name: "disabled pager with tui", config: "pager: true", args: []string{"--pager=false", "-t"}, wantTUI: true},
		{name: "explicit modes conflict", args: []string{"-p", "-t"}, wantPager: true, wantTUI: true, wantError: true},
		{name: "configured modes conflict", config: "pager: true\ntui: true", wantPager: true, wantTUI: true, wantError: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			saveOptions(t)
			t.Setenv("GLOW_PAGER", "")
			t.Setenv("GLOW_TUI", "")
			for k, v := range tt.env {
				t.Setenv(k, v)
			}
			viper.Reset()
			t.Cleanup(viper.Reset)
			viper.SetDefault("style", "auto")
			viper.SetConfigType("yaml")
			if err := viper.ReadConfig(strings.NewReader(tt.config)); err != nil {
				t.Fatal(err)
			}
			viper.SetEnvPrefix("glow")
			viper.AutomaticEnv()

			cmd := &cobra.Command{Use: "glow"}
			cmd.Flags().BoolP("pager", "p", false, "")
			cmd.Flags().BoolP("tui", "t", false, "")
			_ = viper.BindPFlag("pager", cmd.Flags().Lookup("pager"))
			_ = viper.BindPFlag("tui", cmd.Flags().Lookup("tui"))
			if err := cmd.ParseFlags(tt.args); err != nil {
				t.Fatal(err)
			}

			err := validateOptions(cmd)
			if tt.wantError {
				if err == nil || err.Error() != "cannot use both pager and tui" {
					t.Fatalf("expected display mode conflict, got %v", err)
				}
			} else if err != nil {
				t.Fatal(err)
			}
			if pager != tt.wantPager || tui != tt.wantTUI {
				t.Errorf("pager=%v, tui=%v; want pager=%v, tui=%v", pager, tui, tt.wantPager, tt.wantTUI)
			}
		})
	}
}
