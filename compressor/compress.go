package compressor

import (
	"bytes"
	"compress/zlib"
	"errors"
	"fmt"
	"io"
	"os"
)

func CompressFile(filePath string) error {
	var buf bytes.Buffer
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("Can not open file to compress")
	}
	defer file.Close()
	zw, err := zlib.NewWriterLevel(&buf, zlib.BestCompression)
	if err != nil {
		return fmt.Errorf("Can not initialize compression buffer")
	}
	defer zw.Close()
	if _, err = io.Copy(zw, file); err != nil {
		return fmt.Errorf("Can not copy to compression buffer")
	}
	zw.Flush()
	com_file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("Can not create compressed file")
	}
	defer com_file.Close()
	if _, err = io.Copy(com_file, &buf); err != nil {
		return fmt.Errorf("Can not copy to compressed file")
	}
	return nil
}

func DecompressFile(filePath string) error {
	input, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("cannot open file to decompress: %w", err)
	}
	defer input.Close()

	zr, err := zlib.NewReader(input)
	if err != nil {
		return fmt.Errorf("cannot initialize decompressor: %w", err)
	}
	defer zr.Close()

	tmpPath := filePath + ".tmp"
	output, err := os.Create(tmpPath)
	if err != nil {
		return fmt.Errorf("cannot create temporary decompressed file: %w", err)
	}

	defer output.Close()

	if _, err = io.Copy(output, zr); err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
		os.Remove(tmpPath)
		return fmt.Errorf("cannot copy decompressed data: %w", err)
	}

	if err = os.Rename(tmpPath, filePath); err != nil {
		return fmt.Errorf("cannot rename decompressed file: %w", err)
	}

	return nil
}
