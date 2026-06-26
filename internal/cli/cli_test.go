package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRun_EncryptDecrypt_NoSalt(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "input.txt")
	outputFile := filepath.Join(tmpDir, "encrypted.bin")
	decryptFile := filepath.Join(tmpDir, "decrypted.txt")

	original := []byte("abandon ability able about above absent")
	if err := os.WriteFile(inputFile, original, 0644); err != nil {
		t.Fatalf("failed to write input file: %v", err)
	}

	cfg := Config{
		Mode:       "encrypt",
		InputFile:  inputFile,
		OutputFile: outputFile,
		Password:   "Testpassword1",
		Encoding:   "base32",
	}

	if err := Run(cfg); err != nil {
		t.Fatalf("Run encrypt failed: %v", err)
	}

	cfg = Config{
		Mode:       "decrypt",
		InputFile:  outputFile,
		OutputFile: decryptFile,
		Password:   "Testpassword1",
		Encoding:   "base32",
	}

	if err := Run(cfg); err != nil {
		t.Fatalf("Run decrypt failed: %v", err)
	}

	decrypted, err := os.ReadFile(decryptFile)
	if err != nil {
		t.Fatalf("failed to read decrypted file: %v", err)
	}

	if !bytes.Equal(original, decrypted) {
		t.Errorf("decrypted does not match original\ngot:  %s\nwant: %s", decrypted, original)
	}
}

func TestRun_EncryptDecrypt_WithSalt(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "input.txt")
	outputFile := filepath.Join(tmpDir, "encrypted.bin")
	decryptFile := filepath.Join(tmpDir, "decrypted.txt")

	original := []byte("test with custom salt")
	if err := os.WriteFile(inputFile, original, 0644); err != nil {
		t.Fatalf("failed to write input file: %v", err)
	}

	cfg := Config{
		Mode:       "encrypt",
		InputFile:  inputFile,
		OutputFile: outputFile,
		Password:   "Testpassword1",
		Salt:       "mysalt123",
		Encoding:   "base32",
	}

	if err := Run(cfg); err != nil {
		t.Fatalf("Run encrypt failed: %v", err)
	}

	cfg = Config{
		Mode:       "decrypt",
		InputFile:  outputFile,
		OutputFile: decryptFile,
		Password:   "Testpassword1",
		Salt:       "mysalt123",
		Encoding:   "base32",
	}

	if err := Run(cfg); err != nil {
		t.Fatalf("Run decrypt failed: %v", err)
	}

	decrypted, err := os.ReadFile(decryptFile)
	if err != nil {
		t.Fatalf("failed to read decrypted file: %v", err)
	}

	if !bytes.Equal(original, decrypted) {
		t.Errorf("decrypted does not match original\ngot:  %s\nwant: %s", decrypted, original)
	}
}

func TestRun_EncryptDecrypt_Base64(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "input.txt")
	outputFile := filepath.Join(tmpDir, "encrypted.bin")
	decryptFile := filepath.Join(tmpDir, "decrypted.txt")

	original := []byte("test with base64 encoding")
	if err := os.WriteFile(inputFile, original, 0644); err != nil {
		t.Fatalf("failed to write input file: %v", err)
	}

	cfg := Config{
		Mode:       "encrypt",
		InputFile:  inputFile,
		OutputFile: outputFile,
		Password:   "Testpassword1",
		Encoding:   "base64",
	}

	if err := Run(cfg); err != nil {
		t.Fatalf("Run encrypt failed: %v", err)
	}

	cfg = Config{
		Mode:       "decrypt",
		InputFile:  outputFile,
		OutputFile: decryptFile,
		Password:   "Testpassword1",
		Encoding:   "base64",
	}

	if err := Run(cfg); err != nil {
		t.Fatalf("Run decrypt failed: %v", err)
	}

	decrypted, err := os.ReadFile(decryptFile)
	if err != nil {
		t.Fatalf("failed to read decrypted file: %v", err)
	}

	if !bytes.Equal(original, decrypted) {
		t.Errorf("decrypted does not match original\ngot:  %s\nwant: %s", decrypted, original)
	}
}

func TestRunFilename_EncryptDecrypt_Base32(t *testing.T) {
	cfg := Config{
		Mode:     "encrypt-filename",
		InputFile: "TEST_FILE.txt",
		Password: "Testpassword1",
		Encoding: "base32",
	}

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := RunFilename(cfg)
	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("RunFilename encrypt failed: %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)
	encrypted := strings.TrimSpace(buf.String())

	if encrypted != "kr9tu4e1da4u3nifdd99g9tf5o" {
		t.Errorf("encrypted filename mismatch\ngot:  %s\nwant: kr9tu4e1da4u3nifdd99g9tf5o", encrypted)
	}

	cfg = Config{
		Mode:     "decrypt-filename",
		InputFile: encrypted,
		Password: "Testpassword1",
		Encoding: "base32",
	}

	r, w, _ = os.Pipe()
	os.Stdout = w

	err = RunFilename(cfg)
	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("RunFilename decrypt failed: %v", err)
	}

	buf.Reset()
	buf.ReadFrom(r)
	decrypted := strings.TrimSpace(buf.String())

	if decrypted != "TEST_FILE.txt" {
		t.Errorf("decrypted filename mismatch\ngot:  %s\nwant: TEST_FILE.txt", decrypted)
	}
}

func TestRunFilename_EncryptDecrypt_Base64(t *testing.T) {
	cfg := Config{
		Mode:     "encrypt-filename",
		InputFile: "TEST_FILE BASE64.txt",
		Password: "Testpassword1",
		Encoding: "base64",
	}

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := RunFilename(cfg)
	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("RunFilename encrypt failed: %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)
	encrypted := strings.TrimSpace(buf.String())

	if encrypted != "Iyxcijgc9bp3o5Y0npW6xqUvwWNcc3MA4SadB0sR6cY" {
		t.Errorf("encrypted filename mismatch\ngot:  %s\nwant: Iyxcijgc9bp3o5Y0npW6xqUvwWNcc3MA4SadB0sR6cY", encrypted)
	}

	cfg = Config{
		Mode:     "decrypt-filename",
		InputFile: encrypted,
		Password: "Testpassword1",
		Encoding: "base64",
	}

	r, w, _ = os.Pipe()
	os.Stdout = w

	err = RunFilename(cfg)
	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("RunFilename decrypt failed: %v", err)
	}

	buf.Reset()
	buf.ReadFrom(r)
	decrypted := strings.TrimSpace(buf.String())

	if decrypted != "TEST_FILE BASE64.txt" {
		t.Errorf("decrypted filename mismatch\ngot:  %s\nwant: TEST_FILE BASE64.txt", decrypted)
	}
}

func TestRun_MissingInputFile(t *testing.T) {
	cfg := Config{
		Mode:     "encrypt",
		Password: "Testpassword1",
	}

	err := Run(cfg)
	if err == nil {
		t.Error("expected error for missing input file, got nil")
	}
}

func TestRun_InvalidMode(t *testing.T) {
	cfg := Config{
		Mode:      "invalid",
		InputFile: "test.txt",
		Password:  "Testpassword1",
	}

	err := Run(cfg)
	if err == nil {
		t.Error("expected error for invalid mode, got nil")
	}
}

func TestRun_InvalidEncoding(t *testing.T) {
	cfg := Config{
		Mode:      "encrypt",
		InputFile: "test.txt",
		Password:  "Testpassword1",
		Encoding:  "invalid",
	}

	err := Run(cfg)
	if err == nil {
		t.Error("expected error for invalid encoding, got nil")
	}
}
