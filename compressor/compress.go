package compressor

import (
	"bytes"
	"compress/zlib"
	"io"
	"os"
)

func CompressFile(filePath string) error {
	var buf bytes.Buffer
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	zw, err := zlib.NewWriterLevel(&buf, zlib.BestCompression)
	if err != nil {
		return err
	}
	defer zw.Close()
	if _, err = io.Copy(zw, file); err != nil {
		return err
	}
	zw.Flush()
	file.Close()
	com_file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	if _, err = io.Copy(com_file, &buf); err != nil {
		return err
	}
	com_file.Close()
	return nil
}

func DecompressFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	zr, err := zlib.NewReader(file)
	if err != nil {
		return err
	}
	file.Close()
	op_decompress, err := os.Create(filePath)
	if err != nil {
		return err
	}
	if _, err = io.Copy(op_decompress, zr); err != nil {
		return err
	}
	zr.Close()
	op_decompress.Close()
	return nil
}
