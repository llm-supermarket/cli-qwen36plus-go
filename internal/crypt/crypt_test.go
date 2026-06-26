package crypt

import (
	"bytes"
	"testing"
)

func TestEncryptDecrypt_NoSalt(t *testing.T) {
	password := "Testpassword1"
	plaintext := []byte("abandon ability able about above absent absorb abstract absurd abuse access accident")

	cipher, err := NewCipher(password, "", EncodingBase32)
	if err != nil {
		t.Fatalf("NewCipher failed: %v", err)
	}

	var encrypted bytes.Buffer
	err = cipher.EncryptFile(bytes.NewReader(plaintext), &encrypted)
	if err != nil {
		t.Fatalf("EncryptFile failed: %v", err)
	}

	var decrypted bytes.Buffer
	err = cipher.DecryptFile(bytes.NewReader(encrypted.Bytes()), &decrypted)
	if err != nil {
		t.Fatalf("DecryptFile failed: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted.Bytes()) {
		t.Errorf("decrypted does not match original\ngot:  %s\nwant: %s", decrypted.Bytes(), plaintext)
	}
}

func TestEncryptDecrypt_WithSalt(t *testing.T) {
	password := "Testpassword1"
	salt := "mycustomsalt"
	plaintext := []byte("hello world with custom salt")

	cipher, err := NewCipher(password, salt, EncodingBase32)
	if err != nil {
		t.Fatalf("NewCipher failed: %v", err)
	}

	var encrypted bytes.Buffer
	err = cipher.EncryptFile(bytes.NewReader(plaintext), &encrypted)
	if err != nil {
		t.Fatalf("EncryptFile failed: %v", err)
	}

	var decrypted bytes.Buffer
	err = cipher.DecryptFile(bytes.NewReader(encrypted.Bytes()), &decrypted)
	if err != nil {
		t.Fatalf("DecryptFile failed: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted.Bytes()) {
		t.Errorf("decrypted does not match original\ngot:  %s\nwant: %s", decrypted.Bytes(), plaintext)
	}
}

func TestEncryptDecrypt_WrongPassword(t *testing.T) {
	password := "Testpassword1"
	wrongPassword := "WrongPassword1"
	plaintext := []byte("secret data")

	cipher, err := NewCipher(password, "", EncodingBase32)
	if err != nil {
		t.Fatalf("NewCipher failed: %v", err)
	}

	var encrypted bytes.Buffer
	err = cipher.EncryptFile(bytes.NewReader(plaintext), &encrypted)
	if err != nil {
		t.Fatalf("EncryptFile failed: %v", err)
	}

	wrongCipher, err := NewCipher(wrongPassword, "", EncodingBase32)
	if err != nil {
		t.Fatalf("NewCipher with wrong password failed: %v", err)
	}

	var decrypted bytes.Buffer
	err = wrongCipher.DecryptFile(bytes.NewReader(encrypted.Bytes()), &decrypted)
	if err == nil {
		t.Error("expected error when decrypting with wrong password, got nil")
	}
}

func TestEncryptDecrypt_DifferentSalts(t *testing.T) {
	password := "Testpassword1"
	salt1 := "salt1"
	salt2 := "salt2"
	plaintext := []byte("test data with different salts")

	cipher1, err := NewCipher(password, salt1, EncodingBase32)
	if err != nil {
		t.Fatalf("NewCipher salt1 failed: %v", err)
	}

	cipher2, err := NewCipher(password, salt2, EncodingBase32)
	if err != nil {
		t.Fatalf("NewCipher salt2 failed: %v", err)
	}

	var encrypted bytes.Buffer
	err = cipher1.EncryptFile(bytes.NewReader(plaintext), &encrypted)
	if err != nil {
		t.Fatalf("EncryptFile failed: %v", err)
	}

	var decrypted bytes.Buffer
	err = cipher2.DecryptFile(bytes.NewReader(encrypted.Bytes()), &decrypted)
	if err == nil {
		t.Error("expected error when decrypting with different salt, got nil")
	}
}

func TestEncryptDecrypt_LargeFile(t *testing.T) {
	password := "Testpassword1"
	plaintext := make([]byte, 200*1024)
	for i := range plaintext {
		plaintext[i] = byte(i % 256)
	}

	cipher, err := NewCipher(password, "", EncodingBase32)
	if err != nil {
		t.Fatalf("NewCipher failed: %v", err)
	}

	var encrypted bytes.Buffer
	err = cipher.EncryptFile(bytes.NewReader(plaintext), &encrypted)
	if err != nil {
		t.Fatalf("EncryptFile failed: %v", err)
	}

	var decrypted bytes.Buffer
	err = cipher.DecryptFile(bytes.NewReader(encrypted.Bytes()), &decrypted)
	if err != nil {
		t.Fatalf("DecryptFile failed: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted.Bytes()) {
		t.Error("decrypted does not match original for large file")
	}
}

func TestFilenameEncryptDecrypt_Base32(t *testing.T) {
	password := "Testpassword1"
	filename := "TEST_FILE.txt"

	cipher, err := NewCipher(password, "", EncodingBase32)
	if err != nil {
		t.Fatalf("NewCipher failed: %v", err)
	}

	encrypted := cipher.EncryptFilename(filename)
	if encrypted == "" {
		t.Fatal("encrypted filename is empty")
	}

	decrypted, err := cipher.DecryptFilename(encrypted)
	if err != nil {
		t.Fatalf("DecryptFilename failed: %v", err)
	}

	if decrypted != filename {
		t.Errorf("decrypted filename mismatch\ngot:  %s\nwant: %s", decrypted, filename)
	}
}

func TestFilenameEncryptDecrypt_Base64(t *testing.T) {
	password := "Testpassword1"
	filename := "TEST_FILE.txt"

	cipher, err := NewCipher(password, "", EncodingBase64)
	if err != nil {
		t.Fatalf("NewCipher failed: %v", err)
	}

	encrypted := cipher.EncryptFilename(filename)
	if encrypted == "" {
		t.Fatal("encrypted filename is empty")
	}

	decrypted, err := cipher.DecryptFilename(encrypted)
	if err != nil {
		t.Fatalf("DecryptFilename failed: %v", err)
	}

	if decrypted != filename {
		t.Errorf("decrypted filename mismatch\ngot:  %s\nwant: %s", decrypted, filename)
	}
}

func TestFilenameEncryptDecrypt_WithPath(t *testing.T) {
	password := "Testpassword1"
	path := "folder/subfolder/file.txt"

	cipher, err := NewCipher(password, "", EncodingBase32)
	if err != nil {
		t.Fatalf("NewCipher failed: %v", err)
	}

	encrypted := cipher.EncryptFilename(path)

	decrypted, err := cipher.DecryptFilename(encrypted)
	if err != nil {
		t.Fatalf("DecryptFilename failed: %v", err)
	}

	if decrypted != path {
		t.Errorf("decrypted path mismatch\ngot:  %s\nwant: %s", decrypted, path)
	}
}

func TestFilenameKnownValues(t *testing.T) {
	password := "Testpassword1"

	cipher, err := NewCipher(password, "", EncodingBase32)
	if err != nil {
		t.Fatalf("NewCipher failed: %v", err)
	}

	encrypted := cipher.EncryptFilename("TEST_FILE.txt")
	expected := "kr9tu4e1da4u3nifdd99g9tf5o"

	if encrypted != expected {
		t.Errorf("encrypted filename mismatch\ngot:  %s\nwant: %s", encrypted, expected)
	}
}

func TestFilenameKnownValues_Base64(t *testing.T) {
	password := "Testpassword1"

	cipher, err := NewCipher(password, "", EncodingBase64)
	if err != nil {
		t.Fatalf("NewCipher failed: %v", err)
	}

	encrypted := cipher.EncryptFilename("TEST_FILE BASE64.txt")
	expected := "Iyxcijgc9bp3o5Y0npW6xqUvwWNcc3MA4SadB0sR6cY"

	if encrypted != expected {
		t.Errorf("encrypted filename mismatch\ngot:  %s\nwant: %s", encrypted, expected)
	}
}

func TestParseEncoding(t *testing.T) {
	tests := []struct {
		input    string
		expected Encoding
		wantErr  bool
	}{
		{"base32", EncodingBase32, false},
		{"", EncodingBase32, false},
		{"base64", EncodingBase64, false},
		{"hex", EncodingHex, false},
		{"invalid", EncodingBase32, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseEncoding(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseEncoding() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("ParseEncoding() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestEncryptedSize(t *testing.T) {
	tests := []struct {
		plainSize int64
		want      int64
	}{
		{0, 32},
		{1, 32 + 16 + 1},
		{65536, 32 + 65552},
		{65537, 32 + 65552 + 16 + 1},
		{131072, 32 + 2*65552},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := EncryptedSize(tt.plainSize)
			if got != tt.want {
				t.Errorf("EncryptedSize(%d) = %d, want %d", tt.plainSize, got, tt.want)
			}
		})
	}
}

func TestIsEncryptedFile(t *testing.T) {
	encrypted := []byte("RCLONE\x00\x00" + "rest of header and data")
	plaintext := []byte("just some plain text")

	if !IsEncryptedFile(encrypted) {
		t.Error("expected IsEncryptedFile to return true for encrypted data")
	}

	if IsEncryptedFile(plaintext) {
		t.Error("expected IsEncryptedFile to return false for plaintext")
	}
}

func TestEncryptDecrypt_EmptyFile(t *testing.T) {
	password := "Testpassword1"
	plaintext := []byte{}

	cipher, err := NewCipher(password, "", EncodingBase32)
	if err != nil {
		t.Fatalf("NewCipher failed: %v", err)
	}

	var encrypted bytes.Buffer
	err = cipher.EncryptFile(bytes.NewReader(plaintext), &encrypted)
	if err != nil {
		t.Fatalf("EncryptFile failed: %v", err)
	}

	var decrypted bytes.Buffer
	err = cipher.DecryptFile(bytes.NewReader(encrypted.Bytes()), &decrypted)
	if err != nil {
		t.Fatalf("DecryptFile failed: %v", err)
	}

	if len(decrypted.Bytes()) != 0 {
		t.Errorf("expected empty decrypted data, got %d bytes", len(decrypted.Bytes()))
	}
}

func TestEncryptDecrypt_UTF8(t *testing.T) {
	password := "Testpassword1"
	plaintext := []byte("Hello 世界 🌍 مرحبا")

	cipher, err := NewCipher(password, "", EncodingBase32)
	if err != nil {
		t.Fatalf("NewCipher failed: %v", err)
	}

	var encrypted bytes.Buffer
	err = cipher.EncryptFile(bytes.NewReader(plaintext), &encrypted)
	if err != nil {
		t.Fatalf("EncryptFile failed: %v", err)
	}

	var decrypted bytes.Buffer
	err = cipher.DecryptFile(bytes.NewReader(encrypted.Bytes()), &decrypted)
	if err != nil {
		t.Fatalf("DecryptFile failed: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted.Bytes()) {
		t.Errorf("decrypted does not match original\ngot:  %s\nwant: %s", decrypted.Bytes(), plaintext)
	}
}
