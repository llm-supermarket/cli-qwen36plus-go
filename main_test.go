package main

import (
	"os"
	"testing"
)

func TestParseArgs_DefaultMode(t *testing.T) {
	args := []string{"-i", "test.txt"}
	cfg, err := parseArgs(args)
	if err != nil {
		t.Fatalf("parseArgs failed: %v", err)
	}
	if cfg.Mode != "encrypt" {
		t.Errorf("expected mode 'encrypt', got '%s'", cfg.Mode)
	}
	if cfg.InputFile != "test.txt" {
		t.Errorf("expected input file 'test.txt', got '%s'", cfg.InputFile)
	}
}

func TestParseArgs_DecryptMode(t *testing.T) {
	args := []string{"decrypt", "-i", "test.enc"}
	cfg, err := parseArgs(args)
	if err != nil {
		t.Fatalf("parseArgs failed: %v", err)
	}
	if cfg.Mode != "decrypt" {
		t.Errorf("expected mode 'decrypt', got '%s'", cfg.Mode)
	}
}

func TestParseArgs_Password(t *testing.T) {
	args := []string{"-i", "test.txt", "--password", "mypass"}
	cfg, err := parseArgs(args)
	if err != nil {
		t.Fatalf("parseArgs failed: %v", err)
	}
	if cfg.Password != "mypass" {
		t.Errorf("expected password 'mypass', got '%s'", cfg.Password)
	}
}

func TestParseArgs_PasswordEquals(t *testing.T) {
	args := []string{"-i", "test.txt", "--password=mypass"}
	cfg, err := parseArgs(args)
	if err != nil {
		t.Fatalf("parseArgs failed: %v", err)
	}
	if cfg.Password != "mypass" {
		t.Errorf("expected password 'mypass', got '%s'", cfg.Password)
	}
}

func TestParseArgs_Salt(t *testing.T) {
	args := []string{"-i", "test.txt", "--salt", "mysalt"}
	cfg, err := parseArgs(args)
	if err != nil {
		t.Fatalf("parseArgs failed: %v", err)
	}
	if cfg.Salt != "mysalt" {
		t.Errorf("expected salt 'mysalt', got '%s'", cfg.Salt)
	}
}

func TestParseArgs_OutputFile(t *testing.T) {
	args := []string{"-i", "test.txt", "-o", "output.bin"}
	cfg, err := parseArgs(args)
	if err != nil {
		t.Fatalf("parseArgs failed: %v", err)
	}
	if cfg.OutputFile != "output.bin" {
		t.Errorf("expected output file 'output.bin', got '%s'", cfg.OutputFile)
	}
}

func TestParseArgs_Encoding(t *testing.T) {
	args := []string{"-i", "test.txt", "--encoding", "base64"}
	cfg, err := parseArgs(args)
	if err != nil {
		t.Fatalf("parseArgs failed: %v", err)
	}
	if cfg.Encoding != "base64" {
		t.Errorf("expected encoding 'base64', got '%s'", cfg.Encoding)
	}
}

func TestParseArgs_EncryptFilename(t *testing.T) {
	args := []string{"encrypt-filename", "-i", "TEST_FILE.txt"}
	cfg, err := parseArgs(args)
	if err != nil {
		t.Fatalf("parseArgs failed: %v", err)
	}
	if cfg.Mode != "encrypt-filename" {
		t.Errorf("expected mode 'encrypt-filename', got '%s'", cfg.Mode)
	}
}

func TestParseArgs_DecryptFilename(t *testing.T) {
	args := []string{"decrypt-filename", "-i", "kr9tu4e1da4u3nifdd99g9tf5o"}
	cfg, err := parseArgs(args)
	if err != nil {
		t.Fatalf("parseArgs failed: %v", err)
	}
	if cfg.Mode != "decrypt-filename" {
		t.Errorf("expected mode 'decrypt-filename', got '%s'", cfg.Mode)
	}
}

func TestParseArgs_InvalidEncoding(t *testing.T) {
	args := []string{"-i", "test.txt", "--encoding", "invalid"}
	_, err := parseArgs(args)
	if err == nil {
		t.Error("expected error for invalid encoding, got nil")
	}
}

func TestParseArgs_EnvPassword(t *testing.T) {
	os.Setenv("RCLONE_ENCRYPT_PASSWORD", "envpass")
	defer os.Unsetenv("RCLONE_ENCRYPT_PASSWORD")

	args := []string{"-i", "test.txt"}
	cfg, err := parseArgs(args)
	if err != nil {
		t.Fatalf("parseArgs failed: %v", err)
	}
	if cfg.Password != "envpass" {
		t.Errorf("expected password from env 'envpass', got '%s'", cfg.Password)
	}
}

func TestParseArgs_FlagOverridesEnv(t *testing.T) {
	os.Setenv("RCLONE_ENCRYPT_PASSWORD", "envpass")
	defer os.Unsetenv("RCLONE_ENCRYPT_PASSWORD")

	args := []string{"-i", "test.txt", "--password", "flagpass"}
	cfg, err := parseArgs(args)
	if err != nil {
		t.Fatalf("parseArgs failed: %v", err)
	}
	if cfg.Password != "flagpass" {
		t.Errorf("expected flag password to override env, got '%s'", cfg.Password)
	}
}
