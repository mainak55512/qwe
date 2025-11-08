package binaryhandler

import (
	"bytes"
	"errors"
	"fmt"
	cp "github.com/mainak55512/qwe/compressor"
	er "github.com/mainak55512/qwe/qwerror"
	utl "github.com/mainak55512/qwe/qweutils"
	"io"
	"os"
	"time"
	"unicode"
)

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

func CheckBinDiff(file_one, file_two string) (bool, error) {
	file_1, err := os.Open(file_one)
	if err != nil {
		return false, err
	}
	defer file_1.Close()

	file_2, err := os.Open(file_two)
	if err != nil {
		return false, err
	}
	defer file_2.Close()

	for {
		buff1 := make([]byte, 8192)
		buff2 := make([]byte, 8192)

		_, err1 := file_1.Read(buff1)
		_, err2 := file_2.Read(buff2)

		if err1 != nil || err2 != nil {
			if err1 == io.EOF && err2 == io.EOF {
				return true, nil
			} else if err1 == io.EOF || err2 == io.EOF {
				return false, nil
			} else {
				return false, fmt.Errorf("Error while comparing!")
			}
		}

		if !bytes.Equal(buff1, buff2) {
			return false, nil
		}
	}
}

func RevertBinFile(filePath, fileObjID string) error {
	commitFile := ".qwe/_object/" + fileObjID
	src, err := os.Open(commitFile)
	if err != nil {
		return err
	}
	dest, err := os.Create(filePath)
	if err != nil {
		return err
	}
	if _, err = io.Copy(dest, src); err != nil {
		return err
	}
	if err = cp.DecompressFile(filePath); err != nil {
		return err
	}
	return nil
}

func CommitBinFile(filePath, lastCommit string) (string, error) {
	src, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	fileObjID := "_bin_" + utl.Hasher(fmt.Sprintf("%s%d", filePath, time.Now().UnixNano()))
	target := ".qwe/_object/" + fileObjID
	dest, err := os.Create(target)
	if err != nil {
		return "", err
	}

	lastCommittedFile, err := os.Open(".qwe/_object/" + lastCommit)
	if err != nil {
		return "", err
	}
	defer lastCommittedFile.Close()

	if _, err := io.Copy(dest, lastCommittedFile); err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
		return "", err
	}

	if err = cp.DecompressFile(target); err != nil {
		return "", err
	}

	isEq, err := CheckBinDiff(target, filePath)
	if err != nil {
		return "", err
	}

	if !isEq {
		dest, err = os.Create(target)
		if err != nil {
			return "", err
		}
		if _, err := io.Copy(dest, src); err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
			return "", err
		}
		if err = cp.CompressFile(target); err != nil {
			return "", err
		}
	} else {
		os.Remove(target)
		return "", er.NoFileOrDiff
	}

	return fileObjID, nil
}
