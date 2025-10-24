package compressor

import (
	"bytes"
	"compress/zlib"
	"errors"
	er "github.com/mainak55512/qwe/qwerror"
	"io"
	"os"
)

// Compresses the file with zlib
func CompressFile(filePath string) error {
	var buf bytes.Buffer
	file, err := os.Open(filePath)
	if err != nil {
		return er.CompOpenErr
	}
	defer file.Close()
	zw, err := zlib.NewWriterLevel(&buf, zlib.BestCompression)
	if err != nil {
		return er.CompBufInitErr
	}
	defer zw.Close()

	// Copy file content to the buffer as well as compressing it
	if _, err = io.Copy(zw, file); err != nil {
		return er.BufCopyErr
	}
	zw.Flush()
	com_file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer com_file.Close()

	// Copy compressed content from buffer to the file
	if _, err = io.Copy(com_file, &buf); err != nil {
		return er.BufCopyErr
	}
	return nil
}

// decompresses the file using a temporary one
func DecompressFile(filePath string) error {
	input, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer input.Close()

	// Create zlib reader for the file
	zr, err := zlib.NewReader(input)
	if err != nil {
		return er.DecompBufInitErr
	}
	defer zr.Close()

	tmpPath := filePath + ".tmp"
	output, err := os.Create(tmpPath)
	if err != nil {
		return err
	}

	defer output.Close()

	// Decompress and copy content from zlib reader to temporary file
	if _, err = io.Copy(output, zr); err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
		os.Remove(tmpPath)
		return er.BufCopyErr
	}

	// Rename temporary file with the actual output file name
	if err = os.Rename(tmpPath, filePath); err != nil {
		return err
	}

	return nil
}
