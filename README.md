# Glow

Render markdown on the CLI, with _pizzazz_!

## About this fork

This is a fork of [charmbracelet/glow](https://github.com/charmbracelet/glow).
It differs from upstream in these ways:

- **CLI fixes**
  - `--config <file>` is read. Upstream declares the flag but only loads the
    config from the default locations.
  - A display mode flag (`-p` or `-t`) overrides the mode set in the config or
    environment, instead of failing with `cannot use both pager and tui`
    ([charmbracelet/glow#1068](https://github.com/charmbracelet/glow/issues/1068)).
  - `--pager=false` and `--tui=false` are respected.
  - `preserveNewLines` (`-n`) applies to CLI output as well as the TUI
    ([charmbracelet/glow#829](https://github.com/charmbracelet/glow/issues/829)).
- **Builds and releases**: CI runs on GitHub-hosted runners without paid
  tooling. Releases are published on [this fork's releases page][releases] for
  Linux, macOS and Windows on x86_64 and arm64, with checksums and build
  provenance attestations. There are no package-manager, deb/rpm/apk or
  container image builds.
- **Go toolchain**: tracks the latest Go release.

<p align="center">
    <img src="https://stuff.charm.sh/glow/glow-banner-github.gif" alt="Glow Logo">
    <a href="https://github.com/glennib/glow/releases"><img src="https://img.shields.io/github/release/glennib/glow.svg" alt="Latest Release"></a>
    <a href="https://pkg.go.dev/github.com/charmbracelet/glow?tab=doc"><img src="https://godoc.org/github.com/golang/gddo?status.svg" alt="GoDoc"></a>
    <a href="https://github.com/glennib/glow/actions/workflows/ci.yml"><img src="https://github.com/glennib/glow/actions/workflows/ci.yml/badge.svg" alt="Build Status"></a>
    <a href="https://goreportcard.com/report/github.com/charmbracelet/glow"><img src="https://goreportcard.com/badge/charmbracelet/glow" alt="Go ReportCard"></a>
</p>

<p align="center">
    <img src="https://github.com/user-attachments/assets/c2246366-f84b-4847-b431-32a61ca07b74" width="800" alt="Glow UI Demo">
</p>

## What is it?

Glow is a terminal based markdown reader designed from the ground up to bring
out the beauty—and power—of the CLI.

Use it to discover markdown files, read documentation directly on the command
line. Glow will find local markdown files in subdirectories or a local
Git repository.

## Installation

Builds of this fork are published on [its releases page][releases]. The
package managers and `go install` below install upstream Glow.

### Download a release

Binaries for Linux, macOS and Windows on x86_64 and arm64 are available on the
[releases][releases] page, packaged with shell completions and a manpage. Each
release includes a `checksums.txt` and a build provenance attestation, which
you can verify with the [GitHub CLI](https://cli.github.com/):

```bash
gh attestation verify glow_<version>_Linux_x86_64.tar.gz --repo glennib/glow
```

### Package Manager (upstream)

```bash
# macOS or Linux
brew install glow
```

```bash
# macOS (with MacPorts)
sudo port install glow
```

```bash
# Arch Linux (btw)
pacman -S glow
```

```bash
# Void Linux
xbps-install -S glow
```

```bash
# Nix shell
nix-shell -p glow --command glow
```

```bash
# FreeBSD
pkg install glow
```

```bash
# Solus
eopkg install glow
```

```bash
# Windows (with Chocolatey, Scoop, or Winget)
choco install glow
scoop install glow
winget install charmbracelet.glow
```

```bash
# Android (with termux)
pkg install glow
```

```bash
# Ubuntu (Snapcraft)
sudo snap install glow
```

```bash
# Debian/Ubuntu
sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://repo.charm.sh/apt/gpg.key | sudo gpg --dearmor -o /etc/apt/keyrings/charm.gpg
echo "deb [signed-by=/etc/apt/keyrings/charm.gpg] https://repo.charm.sh/apt/ * *" | sudo tee /etc/apt/sources.list.d/charm.list
sudo apt update && sudo apt install glow
```

```bash
# Fedora/RHEL
echo '[charm]
name=Charm
baseurl=https://repo.charm.sh/yum/
enabled=1
gpgcheck=1
gpgkey=https://repo.charm.sh/yum/gpg.key' | sudo tee /etc/yum.repos.d/charm.repo
sudo yum install glow
```

### Go (upstream)

Or just install upstream Glow with `go`:

```bash
go install charm.land/glow/v3@latest
```

This fork keeps upstream's module path, so `go install` can't fetch it. To use
this fork from source, build it as below.

### Build (requires Go 1.21+)

```bash
git clone https://github.com/glennib/glow.git
cd glow
go build
```

[releases]: https://github.com/glennib/glow/releases

## The TUI

Simply run `glow` without arguments to start the textual user interface and
browse local. Glow will find local markdown files in the
current directory and below or, if you’re in a Git repository, Glow will search
the repo.

Markdown files can be read with Glow's high-performance pager. Most of the
keystrokes you know from `less` are the same, but you can press `?` to list
the hotkeys.

## The CLI

In addition to a TUI, Glow has a CLI for working with Markdown. To format a
document use a markdown source as the primary argument:

```bash
# Read from file
glow README.md

# Read from stdin
echo "[Glow](https://github.com/charmbracelet/glow)" | glow -

# Fetch README from GitHub / GitLab
glow github.com/charmbracelet/glow

# Fetch markdown from HTTP
glow https://host.tld/file.md
```

### Word Wrapping

The `-w` flag lets you set a maximum width at which the output will be wrapped:

```bash
glow -w 60
```

### Paging

CLI output can be displayed in your preferred pager with the `-p` flag. This defaults
to the ANSI-aware `less -r` if `$PAGER` is not explicitly set.

### Styles

You can choose a style with the `-s` flag. When no flag is provided `glow` tries
to detect your terminal's current background color and automatically picks
either the `dark` or the `light` style for you.

```bash
glow -s [dark|light]
```

Alternatively you can also supply a custom JSON stylesheet:

```bash
glow -s mystyle.json
```

For additional usage details see:

```bash
glow --help
```

Check out the [Glamour Style Section](https://github.com/charmbracelet/glamour/blob/master/styles/gallery/README.md)
to find more styles. Or [make your own](https://github.com/charmbracelet/glamour/tree/master/styles)!

## The Config File

If you find yourself supplying the same flags to `glow` all the time, it's
probably a good idea to create a config file. Run `glow config`, which will open
it in your favorite $EDITOR. Alternatively you can manually put a file named
`glow.yml` in the default config path of you platform. If you're not sure where
that is, please refer to `glow --help`.

Here's an example config:

```yaml
# style name or JSON path (default "auto")
style: "light"
# mouse wheel support (TUI-mode only)
mouse: true
# use pager to display markdown
pager: true
# at which column should we word wrap?
width: 80
# show all files, including hidden and ignored.
all: false
# show line numbers (TUI-mode only)
showLineNumbers: false
# preserve newlines in the output
preserveNewLines: false
```

## Contributing

See [contributing][contribute].

[contribute]: https://github.com/charmbracelet/glow/contribute

## Feedback

We’d love to hear your thoughts on this project. Feel free to drop us a note!

- [Twitter](https://twitter.com/charmcli)
- [The Fediverse](https://mastodon.social/@charmcli)
- [Discord](https://charm.sh/chat)

## License

[MIT](https://github.com/charmbracelet/glow/raw/master/LICENSE)

---

Part of [Charm](https://charm.sh).

<a href="https://charm.sh/"><img alt="The Charm logo" src="https://stuff.charm.sh/charm-badge.jpg" width="400"></a>

Charm热爱开源 • Charm loves open source
