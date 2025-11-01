package qweutils

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	er "github.com/mainak55512/qwe/qwerror"
	"io"
	"io/fs"
	"os"
	"unicode"
)

// Encodes strings to base64
func ConvStrEnc(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

// Retrieves strings from base64
func ConvStrDec(str string) (string, error) {
	dec_str, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", fmt.Errorf("Failed to decode!")
	}
	return string(dec_str), nil
}

// creates Hash of a given string and returns first 32 characters as string
func Hasher(str string) string {
	hasher := sha256.New()
	hasher.Write([]byte(str))
	hashByte := hasher.Sum(nil)
	return hex.EncodeToString(hashByte)[:32]
}

// Checks if a folder exists
func FolderExists(path string) bool {
	info, err := os.Stat(path)
	if err == nil {
		if info.IsDir() {
			return true
		}
		return false
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false
	}
	return false
}

// Check if a file exists
func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if err == nil {
		return true
	}
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return false
}

// Checks if the file is a text type or binary type, returns true if binary type
func CheckBinFile(filePath string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, er.InvalidFile
	}
	defer file.Close()

	buffer := make([]byte, 1024)

	size, err := io.ReadFull(file, buffer)
	if err != nil && !(errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF)) {
		return false, err
	}
	for i := 0; i < size; i++ {
		runeValue := rune(buffer[i])
		if buffer[i] == 0 && !unicode.IsSpace(runeValue) && !unicode.IsPrint(runeValue) {
			return true, nil
		}
	}
	return false, nil
}
