# henge (変化)

A universal CLI data transformation tool, inspired by [CyberChef](https://gchq.github.io/CyberChef/).

Encode, decode, hash, format, and convert data without leaving the terminal. No more googling `base64` flags or `shasum` syntax.

```bash
echo "aGVsbG8=" | henge                 # Auto-detect: "looks like Base64" → hello
echo "hello" | henge encode base64      # → aGVsbG8=
echo '{"z":3,"a":1}' | henge format json # → pretty-printed, sorted JSON
cat config.yaml | henge convert json    # → YAML to JSON
echo "hello" | henge hash sha256        # → 2cf24dba5fb0a30e...
```

## Why henge?

| Problem | Today | With henge |
|---------|-------|------------|
| Decode Base64 | `base64 -d` or `base64 -D`? Depends on OS | `henge decode base64` |
| SHA256 hash | `shasum -a 256` or `sha256sum`? | `henge hash sha256` |
| JSON to YAML | Install `yq`, learn its syntax | `henge convert yaml` |
| "What is this data?" | Try multiple tools manually | `echo "data" \| henge` |
| Image to Base64 | `base64 < image.png` + manual Data URI | `henge encode image logo.png --data-uri` |

One tool. One syntax. Works with pipes.

## Install

### From source (requires Go 1.22+)

```bash
go install github.com/HituziANDO/henge@latest
```

### Build from repository

```bash
git clone https://github.com/HituziANDO/henge.git
cd henge
go build -o henge .
```

## Usage

### Auto-detect (default)

Just pipe data into `henge` with no arguments. It detects the format and applies the most useful transformation:

```bash
echo "aGVsbG8=" | henge              # Base64 → hello
echo '{"a":1}' | henge              # JSON → pretty-printed
echo "hello%20world" | henge        # URL-encoded → hello world
echo "68656c6c6f" | henge           # Hex → hello
echo "name: henge" | henge          # YAML → JSON
```

Detection priority: JSON → Base64 → YAML → URL encoding → Hex

### Encode / Decode

```bash
# Base64
echo "hello" | henge encode base64        # → aGVsbG8=
echo "aGVsbG8=" | henge decode base64     # → hello

# URL encoding
echo "hello world" | henge encode url     # → hello+world
echo "hello%20world" | henge decode url   # → hello world

# Hex
echo "hello" | henge encode hex           # → 68656c6c6f
echo "68656c6c6f" | henge decode hex      # → hello
```

### Image Base64

Encode image files to Base64 strings and decode them back. Useful for embedding images in HTML/CSS or API payloads.

```bash
# Image → Base64
henge encode image logo.png               # → iVBORw0KGgo...

# Image → Data URI (for HTML/CSS embedding)
henge encode image logo.png --data-uri    # → data:image/png;base64,iVBORw0KGgo...

# Wrap output at 76 characters (for email)
henge encode image photo.jpg --wrap 76

# Save Base64 to file
henge encode image logo.png -o encoded.txt

# Base64 → Image
henge decode image encoded.txt -o restored.png

# Data URI → Image (auto-detected)
echo "data:image/png;base64,iVBORw..." | henge decode image -o output.png

# Round-trip: encode then decode
henge encode image logo.png | henge decode image -o copy.png
```

Supported formats: PNG, JPEG, GIF, WebP, BMP, SVG, ICO

### Hash

```bash
echo -n "hello" | henge hash md5       # → 5d41402abc4b2a76b9719d911017c592
echo -n "hello" | henge hash sha1      # → aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d
echo -n "hello" | henge hash sha256    # → 2cf24dba5fb0a30e26e83b2ac5b9e29e...
echo -n "hello" | henge hash sha512    # → 9b71d224bd62f3785d96d46ad3ea3d73...
```

### Format (pretty-print)

```bash
# JSON
echo '{"z":3,"a":1,"b":2}' | henge format json
# {
#   "a": 1,
#   "b": 2,
#   "z": 3
# }

# Custom indent
echo '{"a":1}' | henge format json --indent 4

# Compact (minify)
echo '{ "a" : 1 }' | henge format json -c    # → {"a":1}

# YAML
echo "name: henge" | henge format yaml

# XML
echo '<root><a>b</a></root>' | henge format xml
```

### Convert (format transformation)

```bash
# JSON → YAML
echo '{"name":"henge","version":"1.0"}' | henge convert yaml
# name: henge
# version: "1.0"

# YAML → JSON
echo "name: henge" | henge convert json
# {
#   "name": "henge"
# }

# JSON → TOML
echo '{"name":"henge","version":"1.0"}' | henge convert toml

# CSV → JSON
echo -e "name,age\nAlice,30\nBob,25" | henge convert json
# [
#   {"age":"30","name":"Alice"},
#   {"age":"25","name":"Bob"}
# ]

# Explicit input format (override auto-detection)
cat data.txt | henge convert json --from yaml
```

### File input

```bash
# As argument
henge format json data.json
henge hash sha256 file.txt

# Via redirect
henge convert yaml < config.json

# Via pipe
cat data.json | henge format json
```

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--output <file>` | `-o` | Write output to file instead of stdout |
| `--from <format>` | `-f` | Specify input format (override auto-detection) |
| `--compact` | `-c` | Compact output (no indentation) |
| `--no-newline` | `-n` | Do not append trailing newline |
| `--help` | `-h` | Show help |
| `--version` | `-v` | Show version |

## Shell Completion

```bash
# Bash
henge completion bash > /etc/bash_completion.d/henge

# Zsh
henge completion zsh > "${fpath[1]}/_henge"

# Fish
henge completion fish > ~/.config/fish/completions/henge.fish
```

## Pipe-friendly

henge follows Unix philosophy. Compose it with other tools:

```bash
# Extract a field with jq, then convert to YAML
curl -s https://api.example.com/data | jq '.config' | henge convert yaml

# Hash a downloaded file
curl -sL https://example.com/release.tar.gz | henge hash sha256

# Decode a JWT payload (Base64 URL-encoded)
echo "$JWT" | cut -d. -f2 | henge decode base64

# Chain encode/decode
echo "hello" | henge encode base64 | henge decode base64
```

## Command Reference

```
henge [file]                          Auto-detect and transform
henge auto [input]                    Auto-detect (explicit alias)

henge encode base64 [input]           Base64 encode
henge encode url [input]              URL percent-encode
henge encode hex [input]              Hex encode
henge encode image <file>             Image to Base64 (--data-uri, --wrap)

henge decode base64 [input]           Base64 decode
henge decode url [input]              URL percent-decode
henge decode hex [input]              Hex decode
henge decode image [input] -o <file>  Base64 to image file

henge hash md5 [input]                MD5 hash
henge hash sha1 [input]               SHA-1 hash
henge hash sha256 [input]             SHA-256 hash
henge hash sha512 [input]             SHA-512 hash

henge format json [input]             Pretty-print JSON
henge format yaml [input]             Pretty-print YAML
henge format xml [input]              Pretty-print XML

henge convert json [input]            Convert to JSON
henge convert yaml [input]            Convert to YAML
henge convert toml [input]            Convert to TOML
```

## Development

```bash
# Build
go build -o henge .

# Run tests
go test ./...

# Run with verbose output
go test ./... -v
```

### Cross-platform build with GoReleaser

[GoReleaser](https://goreleaser.com/) is used to build binaries for multiple platforms.

**Supported platforms:**

| OS | Arch | Target |
|----|------|--------|
| macOS | x86_64 | Intel Mac |
| macOS | arm64 | Apple Silicon |
| Linux | x86_64 | Intel / AMD |
| Linux | arm64 | ARM64 |
| Windows | x86_64 | Intel / AMD |

**Snapshot build (local, without tag):**

```bash
goreleaser build --snapshot --clean
```

Binaries are output to `dist/`.

**Release build (requires a Git tag and `GITHUB_TOKEN`):**

```bash
git tag v{X.Y.Z}
git push origin v{X.Y.Z}
```

## Tech Stack

| Component | Choice | Reason |
|-----------|--------|--------|
| Language | Go | Single binary, cross-platform, fast startup |
| CLI Framework | [Cobra](https://github.com/spf13/cobra) | Used by kubectl, docker, gh |
| YAML | [gopkg.in/yaml.v3](https://pkg.go.dev/gopkg.in/yaml.v3) | Standard Go YAML library |
| TOML | [BurntSushi/toml](https://github.com/BurntSushi/toml) | De facto Go TOML library |
| CSV, JSON, Hash | Go standard library | No external dependencies needed |

## License

MIT
