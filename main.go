package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/llm-supermarket-org/cli-qwen36plus-go/internal/cli"
	"github.com/llm-supermarket-org/cli-qwen36plus-go/internal/crypt"
)

var version = "dev"

func main() {
	cfg, err := parseArgs(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintf(os.Stderr, "Usage: rclone-encrypt [mode] [flags]\n")
		fmt.Fprintf(os.Stderr, "Run 'rclone-encrypt --help' for usage.\n")
		os.Exit(1)
	}

	if cfg.Version {
		fmt.Printf("rclone-encrypt %s\n", version)
		return
	}

	if cfg.Mode == "encrypt-filename" || cfg.Mode == "decrypt-filename" {
		if err := cli.RunFilename(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if err := cli.Run(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func parseArgs(args []string) (cli.Config, error) {
	cfg := cli.Config{
		Mode:     "encrypt",
		Encoding: "base32",
	}

	i := 0
	for i < len(args) {
		arg := args[i]

		switch {
		case arg == "--help", arg == "-h":
			printHelp()
			os.Exit(0)
		case arg == "--version", arg == "-v":
			cfg.Version = true
			return cfg, nil
		case arg == "encrypt", arg == "decrypt":
			cfg.Mode = arg
		case arg == "encrypt-filename", arg == "decrypt-filename":
			cfg.Mode = arg
		case arg == "--password":
			i++
			if i >= len(args) {
				return cfg, fmt.Errorf("--password requires a value")
			}
			cfg.Password = args[i]
		case strings.HasPrefix(arg, "--password="):
			cfg.Password = strings.TrimPrefix(arg, "--password=")
		case arg == "--salt":
			i++
			if i >= len(args) {
				return cfg, fmt.Errorf("--salt requires a value")
			}
			cfg.Salt = args[i]
		case strings.HasPrefix(arg, "--salt="):
			cfg.Salt = strings.TrimPrefix(arg, "--salt=")
		case arg == "--input-file", arg == "-i":
			i++
			if i >= len(args) {
				return cfg, fmt.Errorf("--input-file requires a value")
			}
			cfg.InputFile = args[i]
		case strings.HasPrefix(arg, "--input-file="):
			cfg.InputFile = strings.TrimPrefix(arg, "--input-file=")
		case strings.HasPrefix(arg, "-i="):
			cfg.InputFile = strings.TrimPrefix(arg, "-i=")
		case arg == "--output-file", arg == "-o":
			i++
			if i >= len(args) {
				return cfg, fmt.Errorf("--output-file requires a value")
			}
			cfg.OutputFile = args[i]
		case strings.HasPrefix(arg, "--output-file="):
			cfg.OutputFile = strings.TrimPrefix(arg, "--output-file=")
		case strings.HasPrefix(arg, "-o="):
			cfg.OutputFile = strings.TrimPrefix(arg, "-o=")
		case arg == "--encoding":
			i++
			if i >= len(args) {
				return cfg, fmt.Errorf("--encoding requires a value")
			}
			cfg.Encoding = args[i]
		case strings.HasPrefix(arg, "--encoding="):
			cfg.Encoding = strings.TrimPrefix(arg, "--encoding=")
		default:
			if cfg.InputFile == "" && !strings.HasPrefix(arg, "-") {
				cfg.InputFile = arg
			} else {
				return cfg, fmt.Errorf("unknown flag: %s", arg)
			}
		}
		i++
	}

	if _, err := crypt.ParseEncoding(cfg.Encoding); err != nil {
		return cfg, err
	}

	if envPass := os.Getenv("RCLONE_ENCRYPT_PASSWORD"); envPass != "" && cfg.Password == "" {
		cfg.Password = envPass
	}

	return cfg, nil
}

func printHelp() {
	fmt.Println(`rclone-encrypt - Encrypt and decrypt files using rclone-compatible cryptography

USAGE:
    rclone-encrypt [mode] [flags]

MODES:
    encrypt             Encrypt a file (default)
    decrypt             Decrypt a file
    encrypt-filename    Encrypt a filename
    decrypt-filename    Decrypt a filename

FLAGS:
    -i, --input-file    Input file path (required for encrypt/decrypt)
    -o, --output-file   Output file path (optional, defaults to stdout)
    --password          Password (warns about security risks; prefer env var)
    --salt              Optional salt (default: rclone's built-in salt)
    --encoding          Filename encoding: base32 (default), base64, hex
    -h, --help          Show this help message
    -v, --version       Show version

ENVIRONMENT VARIABLES:
    RCLONE_ENCRYPT_PASSWORD   Password (preferred over --password)
    RCLONE_ENCRYPT_SALT       Salt (optional)

EXAMPLES:
    # Encrypt a file (prompts for password)
    rclone-encrypt encrypt -i secret.txt -o secret.txt.enc

    # Encrypt with password flag
    rclone-encrypt encrypt -i secret.txt -o secret.txt.enc --password mypassword

    # Encrypt with custom salt and base64 filename encoding
    rclone-encrypt encrypt -i secret.txt -o secret.txt.enc --password mypass --salt mysalt --encoding base64

    # Decrypt a file
    rclone-encrypt decrypt -i secret.txt.enc -o secret.txt

    # Use environment variable for password
    export RCLONE_ENCRYPT_PASSWORD=mypassword
    rclone-encrypt encrypt -i secret.txt -o secret.txt.enc

    # Encrypt/decrypt a filename
    rclone-encrypt encrypt-filename -i "TEST_FILE.txt" --password mypass --encoding base32
    rclone-encrypt decrypt-filename -i "kr9tu4e1da4u3nifdd99g9tf5o" --password mypass --encoding base32`)
}
