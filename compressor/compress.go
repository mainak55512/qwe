package compressor

import (
	"bytes"
	"compress/zlib"
	"io"
	"os"
)

func CompressFile(filePath string) {
	var buf bytes.Buffer
	file, _ := os.Open(filePath)
	zw, _ := zlib.NewWriterLevel(&buf, zlib.BestCompression)
	defer zw.Close()
	io.Copy(zw, file)
	zw.Flush()
	file.Close()
	com_file, _ := os.Create(filePath)
	io.Copy(com_file, &buf)
	com_file.Close()
}

func DecompressFile(filePath string) {
	file, _ := os.Open(filePath)
	zr, _ := zlib.NewReader(file)
	file.Close()
	op_decompress, _ := os.Create(filePath)
	io.Copy(op_decompress, zr)
	zr.Close()
	op_decompress.Close()
}
