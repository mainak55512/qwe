package binaryhandler

import (
	"errors"
	"fmt"
	cp "github.com/mainak55512/qwe/compressor"
	utl "github.com/mainak55512/qwe/qweutils"
	"io"
	"os"
	"time"
)

func CommitBinFile(filePath string) (string, error) {
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
	if _, err := io.Copy(dest, src); err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
		return "", err
	}
	if err = cp.CompressFile(target); err != nil {
		return "", err
	}
	return fileObjID, nil
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
