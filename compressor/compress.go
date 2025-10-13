package compressor

import (
	"bytes"
	"compress/zlib"
	"errors"
	"fmt"
	"io"
	"os"
)

// Compresses the file with zlib
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

	// Copy file content to the buffer as well as compressing it
	if _, err = io.Copy(zw, file); err != nil {
		return fmt.Errorf("Can not copy to compression buffer")
	}
	zw.Flush()
	com_file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("Can not create compressed file")
	}
	defer com_file.Close()

	// Copy compressed content from buffer to the file
	if _, err = io.Copy(com_file, &buf); err != nil {
		return fmt.Errorf("Can not copy to compressed file")
	}
	return nil
}

// decompresses the file using a temporary one
func DecompressFile(filePath string) error {
	input, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("cannot open file to decompress: %w", err)
	}
	defer input.Close()

	// Create zlib reader for the file
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

	// Decompress and copy content from zlib reader to temporary file
	if _, err = io.Copy(output, zr); err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
		os.Remove(tmpPath)
		return fmt.Errorf("cannot copy decompressed data: %w", err)
	}

	// Rename temporary file with the actual output file name
	if err = os.Rename(tmpPath, filePath); err != nil {
		return fmt.Errorf("cannot rename decompressed file: %w", err)
	}

	return nil
}
