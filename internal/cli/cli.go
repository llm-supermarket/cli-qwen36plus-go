package cli

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"

	"github.com/llm-supermarket-org/cli-qwen36plus-go/internal/crypt"
	"golang.org/x/term"
)

const version = "1.0.0"

type Config struct {
	Mode       string
	InputFile  string
	OutputFile string
	Password   string
	Salt       string
	Encoding   string
	Version    bool
}

func Run(cfg Config) error {
	if cfg.Version {
		fmt.Printf("rclone-encrypt %s\n", version)
		return nil
	}

	if cfg.Mode != "encrypt" && cfg.Mode != "decrypt" {
		return fmt.Errorf("mode must be 'encrypt' or 'decrypt', got: %s", cfg.Mode)
	}

	if cfg.InputFile == "" {
		return fmt.Errorf("input file is required (use -i or --input-file)")
	}

	if _, err := os.Stat(cfg.InputFile); os.IsNotExist(err) {
		return fmt.Errorf("input file not found: %s", cfg.InputFile)
	}

	if _, err := crypt.ParseEncoding(cfg.Encoding); err != nil {
		return err
	}

	password := cfg.Password
	salt := cfg.Salt

	if password == "" {
		var err error
		password, err = promptPassword("Enter password: ")
		if err != nil {
			return fmt.Errorf("failed to read password: %w", err)
		}
	} else {
		fmt.Fprintln(os.Stderr, "WARNING: Using --password exposes it in shell history and process listings.")
		fmt.Fprintln(os.Stderr, "Consider using the RCLONE_ENCRYPT_PASSWORD environment variable instead.")
		fmt.Fprintln(os.Stderr, "After use, clear your terminal history (e.g., 'history -c' on bash/zsh).")
	}

	if salt == "" {
		saltVal := os.Getenv("RCLONE_ENCRYPT_SALT")
		if saltVal != "" {
			salt = saltVal
		}
	}

	if salt == "" && term.IsTerminal(int(syscall.Stdin)) {
		s, err := promptPassword("Enter salt (optional, press Enter to skip): ")
		if err != nil {
			return fmt.Errorf("failed to read salt: %w", err)
		}
		salt = strings.TrimSpace(s)
	}

	cipher, err := crypt.NewCipher(password, salt, 0)
	if err != nil {
		return fmt.Errorf("failed to initialize cipher: %w", err)
	}

	inFile, err := os.Open(cfg.InputFile)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inFile.Close()

	var outWriter io.Writer
	if cfg.OutputFile != "" {
		outFile, err := os.Create(cfg.OutputFile)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer outFile.Close()
		outWriter = outFile
	} else {
		outWriter = os.Stdout
	}

	switch cfg.Mode {
	case "encrypt":
		err = cipher.EncryptFile(inFile, outWriter)
	case "decrypt":
		err = cipher.DecryptFile(inFile, outWriter)
	}

	if err != nil {
		return err
	}

	if cfg.OutputFile != "" {
		fmt.Fprintf(os.Stderr, "Successfully %sed: %s -> %s\n", cfg.Mode, cfg.InputFile, cfg.OutputFile)
	}

	return nil
}

func RunFilename(cfg Config) error {
	if cfg.Mode != "encrypt-filename" && cfg.Mode != "decrypt-filename" {
		return fmt.Errorf("mode must be 'encrypt-filename' or 'decrypt-filename', got: %s", cfg.Mode)
	}

	if cfg.InputFile == "" {
		return fmt.Errorf("filename is required (use -i or --input-file)")
	}

	password := cfg.Password
	salt := cfg.Salt

	if password == "" {
		var err error
		password, err = promptPassword("Enter password: ")
		if err != nil {
			return fmt.Errorf("failed to read password: %w", err)
		}
	} else {
		fmt.Fprintln(os.Stderr, "WARNING: Using --password exposes it in shell history and process listings.")
		fmt.Fprintln(os.Stderr, "Consider using the RCLONE_ENCRYPT_PASSWORD environment variable instead.")
		fmt.Fprintln(os.Stderr, "After use, clear your terminal history (e.g., 'history -c' on bash/zsh).")
	}

	if salt == "" {
		saltVal := os.Getenv("RCLONE_ENCRYPT_SALT")
		if saltVal != "" {
			salt = saltVal
		}
	}

	if salt == "" && term.IsTerminal(int(syscall.Stdin)) {
		s, err := promptPassword("Enter salt (optional, press Enter to skip): ")
		if err != nil {
			return fmt.Errorf("failed to read salt: %w", err)
		}
		salt = strings.TrimSpace(s)
	}

	enc, err := crypt.ParseEncoding(cfg.Encoding)
	if err != nil {
		return err
	}

	cipher, err := crypt.NewCipher(password, salt, enc)
	if err != nil {
		return fmt.Errorf("failed to initialize cipher: %w", err)
	}

	switch cfg.Mode {
	case "encrypt-filename":
		result := cipher.EncryptFilename(cfg.InputFile)
		fmt.Println(result)
	case "decrypt-filename":
		result, err := cipher.DecryptFilename(cfg.InputFile)
		if err != nil {
			return fmt.Errorf("failed to decrypt filename: %w", err)
		}
		fmt.Println(result)
	}

	return nil
}

func promptPassword(prompt string) (string, error) {
	if !term.IsTerminal(int(syscall.Stdin)) {
		reader := bufio.NewReader(os.Stdin)
		fmt.Fprint(os.Stderr, prompt)
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(line), nil
	}

	fmt.Fprint(os.Stderr, prompt)
	b, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Fprintln(os.Stderr)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
