package crypt

import (
	"bytes"
	"crypto/aes"
	"crypto/rand"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/rfjakob/eme"
	"golang.org/x/crypto/nacl/secretbox"
	"golang.org/x/crypto/scrypt"
)

const (
	DefaultSaltHex = "a80df43a8fbd0308a7cab83e581f86b1"

	scryptN     = 16384
	scryptR     = 8
	scryptP     = 1
	keySize     = 80
	fileMagic   = "RCLONE\x00\x00"
	fileHeader  = 32
	blockData   = 64 * 1024
	blockHeader = 16
	blockSize   = blockData + blockHeader
)

type Encoding int

const (
	EncodingBase32 Encoding = iota
	EncodingBase64
	EncodingBase32768
	EncodingHex
)

func ParseEncoding(s string) (Encoding, error) {
	switch strings.ToLower(s) {
	case "base32", "":
		return EncodingBase32, nil
	case "base64":
		return EncodingBase64, nil
	case "hex":
		return EncodingHex, nil
	default:
		return EncodingBase32, fmt.Errorf("unknown encoding: %s", s)
	}
}

type Cipher struct {
	dataKey     [32]byte
	nameKey     [32]byte
	nameTweak   [16]byte
	nameCipher  *eme.EMECipher
	fileNameEnc Encoding
}

func NewCipher(password, salt string, enc Encoding) (*Cipher, error) {
	saltBytes, err := hex.DecodeString(DefaultSaltHex)
	if err != nil {
		return nil, fmt.Errorf("invalid default salt: %w", err)
	}
	if salt != "" {
		saltBytes = []byte(salt)
	}

	key, err := scrypt.Key([]byte(password), saltBytes, scryptN, scryptR, scryptP, keySize)
	if err != nil {
		return nil, fmt.Errorf("scrypt failed: %w", err)
	}

	var c Cipher
	copy(c.dataKey[:], key[0:32])
	copy(c.nameKey[:], key[32:64])
	copy(c.nameTweak[:], key[64:80])

	aesBlock, err := aes.NewCipher(c.nameKey[:])
	if err != nil {
		return nil, fmt.Errorf("aes cipher: %w", err)
	}

	c.nameCipher = eme.New(aesBlock)
	c.fileNameEnc = enc
	return &c, nil
}

func (c *Cipher) EncryptFile(in io.Reader, out io.Writer) error {
	var nonce [24]byte
	if _, err := rand.Read(nonce[:]); err != nil {
		return err
	}

	header := make([]byte, fileHeader)
	copy(header[:8], []byte(fileMagic))
	copy(header[8:], nonce[:])
	if _, err := out.Write(header); err != nil {
		return err
	}

	buf := make([]byte, blockData)
	for {
		n, readErr := io.ReadFull(in, buf)
		if n == 0 {
			if readErr == io.EOF || readErr == io.ErrUnexpectedEOF {
				break
			}
			return readErr
		}

		sealed := secretbox.Seal(nil, buf[:n], &nonce, &c.dataKey)
		if _, err := out.Write(sealed); err != nil {
			return err
		}

		for i := range nonce {
			nonce[i]++
			if nonce[i] != 0 {
				break
			}
		}

		if readErr == io.EOF || readErr == io.ErrUnexpectedEOF {
			break
		}
		if readErr != nil {
			return readErr
		}
	}
	return nil
}

func (c *Cipher) DecryptFile(in io.Reader, out io.Writer) error {
	header := make([]byte, fileHeader)
	if _, err := io.ReadFull(in, header); err != nil {
		return fmt.Errorf("read header: %w", err)
	}

	if string(header[:8]) != fileMagic {
		return errors.New("invalid rclone crypt file: bad magic")
	}

	var nonce [24]byte
	copy(nonce[:], header[8:])

	blockBuf := make([]byte, blockSize)
	for {
		n, err := io.ReadFull(in, blockBuf)
		if err == io.EOF {
			break
		}
		if n == 0 {
			break
		}

		opened, ok := secretbox.Open(nil, blockBuf[:n], &nonce, &c.dataKey)
		if !ok {
			return errors.New("failed to decrypt block: wrong password or corrupted data")
		}

		if _, err := out.Write(opened); err != nil {
			return err
		}

		for i := range nonce {
			nonce[i]++
			if nonce[i] != 0 {
				break
			}
		}

		if err == io.EOF {
			break
		}
	}
	return nil
}

func (c *Cipher) EncryptFilename(name string) string {
	parts := strings.Split(name, "/")
	encrypted := make([]string, len(parts))
	for i, part := range parts {
		encrypted[i] = c.encryptSegment(part)
	}
	return strings.Join(encrypted, "/")
}

func (c *Cipher) DecryptFilename(name string) (string, error) {
	parts := strings.Split(name, "/")
	decrypted := make([]string, len(parts))
	for i, part := range parts {
		d, err := c.decryptSegment(part)
		if err != nil {
			return "", err
		}
		decrypted[i] = d
	}
	return strings.Join(decrypted, "/"), nil
}

func (c *Cipher) encryptSegment(name string) string {
	padded := pkcs7Pad([]byte(name), 16)
	encrypted := c.nameCipher.Encrypt(c.nameTweak[:], padded)
	return encodeFilename(encrypted, c.fileNameEnc)
}

func (c *Cipher) decryptSegment(encoded string) (string, error) {
	decoded, err := decodeFilename(encoded, c.fileNameEnc)
	if err != nil {
		return "", err
	}
	decrypted := c.nameCipher.Decrypt(c.nameTweak[:], decoded)
	return pkcs7Unpad(decrypted)
}

func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - (len(data) % blockSize)
	padded := make([]byte, len(data)+padding)
	copy(padded, data)
	for i := len(data); i < len(padded); i++ {
		padded[i] = byte(padding)
	}
	return padded
}

func pkcs7Unpad(data []byte) (string, error) {
	if len(data) == 0 {
		return "", errors.New("empty data")
	}
	padding := int(data[len(data)-1])
	if padding == 0 || padding > len(data) {
		return "", fmt.Errorf("invalid padding: %d", padding)
	}
	for i := len(data) - padding; i < len(data); i++ {
		if data[i] != byte(padding) {
			return "", errors.New("invalid padding")
		}
	}
	return string(data[:len(data)-padding]), nil
}

func encodeFilename(data []byte, enc Encoding) string {
	switch enc {
	case EncodingBase32:
		s := base32.HexEncoding.EncodeToString(data)
		return strings.ToLower(strings.TrimRight(s, "="))
	case EncodingBase64:
		return base64.RawURLEncoding.EncodeToString(data)
	case EncodingHex:
		return hex.EncodeToString(data)
	default:
		s := base32.HexEncoding.EncodeToString(data)
		return strings.ToLower(strings.TrimRight(s, "="))
	}
}

func decodeFilename(s string, enc Encoding) ([]byte, error) {
	switch enc {
	case EncodingBase32:
		upper := strings.ToUpper(s)
		mod := len(upper) % 8
		if mod != 0 {
			upper += strings.Repeat("=", 8-mod)
		}
		return base32.HexEncoding.DecodeString(upper)
	case EncodingBase64:
		return base64.RawURLEncoding.DecodeString(s)
	case EncodingHex:
		return hex.DecodeString(s)
	default:
		upper := strings.ToUpper(s)
		mod := len(upper) % 8
		if mod != 0 {
			upper += strings.Repeat("=", 8-mod)
		}
		return base32.HexEncoding.DecodeString(upper)
	}
}

func EncryptedSize(plainSize int64) int64 {
	blocks, residue := plainSize/blockData, plainSize%blockData
	size := int64(fileHeader) + blocks*blockSize
	if residue != 0 {
		size += blockHeader + residue
	}
	return size
}

func DecryptFilenameOnly(encoded, password, salt, encStr string) (string, error) {
	enc, err := ParseEncoding(encStr)
	if err != nil {
		return "", err
	}
	c, err := NewCipher(password, salt, enc)
	if err != nil {
		return "", err
	}
	return c.DecryptFilename(encoded)
}

func EncryptFilenameOnly(plain, password, salt, encStr string) (string, error) {
	enc, err := ParseEncoding(encStr)
	if err != nil {
		return "", err
	}
	c, err := NewCipher(password, salt, enc)
	if err != nil {
		return "", err
	}
	return c.EncryptFilename(plain), nil
}

func DetectMode(data []byte) string {
	if len(data) >= 8 && string(data[:8]) == fileMagic {
		return "encrypted"
	}
	return "plaintext"
}

func IsEncryptedFile(data []byte) bool {
	return bytes.HasPrefix(data, []byte(fileMagic))
}
