# rclone-encrypt-qwen36plus
A small CLI tool that encrypts and decrypts using the rclone encryption defaults. 

Rclone uses a custom salt if no salt is provided, which this tool will use by default. A few similar tools:

- https://github.com/rclone/rclone
- https://github.com/mcolatosti/rclonedecrypt
- https://github.com/br0kenpixel/rclone-rcc
- @fyears/rclone-crypt

Rclone encryption uses: 
- NaCl SecretBox (XSalsa20 + Poly1305) for the file contents.
- AES256 for the filenames.
- scrypt for keymaterial.

## Installation

**Homebrew (macOS/Linux)**

```bash
brew tap llm-supermarket/cli-qwen36plus-go https://github.com/llm-supermarket/cli-qwen36plus-go
brew install rclone-encrypt
```

**Scoop (Windows)**

```bash
scoop bucket add rclone-encrypt https://github.com/llm-supermarket/cli-qwen36plus-go
scoop install rclone-encrypt
```

## Usage

### Basic example

```bash
# Encrypt a file (prompts for password)
rclone-encrypt encrypt -i secret.txt -o secret.txt.enc

# Decrypt a file
rclone-encrypt decrypt -i secret.txt.enc -o secret.txt
```

### With password flag

```bash
# Encrypt with --password (warns about security risks)
rclone-encrypt encrypt -i secret.txt -o secret.txt.enc --password mypassword

# Decrypt with --password
rclone-encrypt decrypt -i secret.txt.enc -o secret.txt --password mypassword
```

> **Security warning**: Using `--password` exposes it in shell history and process listings. 
> Consider using the `RCLONE_ENCRYPT_PASSWORD` environment variable instead, and clear your 
> terminal history after use (e.g., `history -c` on bash/zsh).

### With environment variables

```bash
# Set password via environment variable (preferred)
export RCLONE_ENCRYPT_PASSWORD=mypassword
rclone-encrypt encrypt -i secret.txt -o secret.txt.enc

# Optional: set salt via environment variable
export RCLONE_ENCRYPT_SALT=mycustomsalt
rclone-encrypt encrypt -i secret.txt -o secret.txt.enc
```

### With custom salt

```bash
# Encrypt with a custom salt
rclone-encrypt encrypt -i secret.txt -o secret.txt.enc --password mypass --salt mysalt

# Decrypt with the same salt
rclone-encrypt decrypt -i secret.txt.enc -o secret.txt --password mypass --salt mysalt
```

### Filename encryption

```bash
# Encrypt a filename (default base32 encoding)
rclone-encrypt encrypt-filename -i "TEST_FILE.txt" --password mypass
# Output: kr9tu4e1da4u3nifdd99g9tf5o

# Decrypt a filename
rclone-encrypt decrypt-filename -i "kr9tu4e1da4u3nifdd99g9tf5o" --password mypass
# Output: TEST_FILE.txt
```

### Custom filename encoding

```bash
# Use base64 encoding for filenames
rclone-encrypt encrypt-filename -i "TEST_FILE.txt" --password mypass --encoding base64

# Decrypt with base64 encoding
rclone-encrypt decrypt-filename -i "encrypted-base64-filename" --password mypass --encoding base64
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-i`, `--input-file` | *(required)* | Input file path (or filename for filename modes) |
| `-o`, `--output-file` | stdout | Output file path (optional for encrypt/decrypt) |
| `--password` | *(prompt)* | Password (warns about security risks) |
| `--salt` | rclone default | Optional salt for key derivation |
| `--encoding` | `base32` | Filename encoding: `base32`, `base64`, `hex` |
| `-h`, `--help` | | Show help message |
| `-v`, `--version` | | Show version |

## Modes

| Mode | Description |
|------|-------------|
| `encrypt` | Encrypt a file (default) |
| `decrypt` | Decrypt a file |
| `encrypt-filename` | Encrypt a filename only |
| `decrypt-filename` | Decrypt a filename only |

## Environment Variables

| Variable | Description |
|----------|-------------|
| `RCLONE_ENCRYPT_PASSWORD` | Password (preferred over `--password`) |
| `RCLONE_ENCRYPT_SALT` | Optional salt |

## Details

### Encryption

rclone-encrypt uses the same encryption scheme as rclone's crypt remote:

- **Key derivation**: scrypt with N=16384, r=8, p=1, producing 80 bytes of key material
- **File contents**: NaCl SecretBox (XSalsa20 + Poly1305) with 64KB chunks
- **Filenames**: EME (ECB-Mix-ECB) wide-block mode with AES-256
- **Salt**: Uses rclone's default salt (`a80df43a8fbd0308a7cab83e581f86b1`) if not provided

### File format

Encrypted files start with a 32-byte header:
- 8 bytes: Magic string `RCLONE\x00\x00`
- 24 bytes: Random nonce for the first block

Each subsequent block is encrypted with SecretBox, with the nonce incremented for each 64KB chunk.

### Filename encoding

Filenames can be encoded using different schemes:

- **base32** (default): Modified hex base32, lowercase, no padding. Compatible with all remotes.
- **base64**: URL-safe base64 without padding. Shorter but case-sensitive.
- **hex**: Hexadecimal encoding. Longest but universally compatible.

## Building from Source

Requires Go 1.25+.

```bash
git clone https://github.com/llm-supermarket/cli-qwen36plus-go
cd cli-qwen36plus-go
go build -o rclone-encrypt .
```

### Running tests

```bash
go test ./...
```

## Releases

Pushing a `vX.Y.Z` tag triggers the [Build and Release workflow](.github/workflows/build-release.yml), which cross-compiles binaries for Linux and macOS (amd64/arm64) and Windows (amd64), publishes a GitHub Release, and updates the Scoop manifest (`rclone-encrypt.json`) and Homebrew formula (`Formula/rclone-encrypt.rb`) in this repo.
