package reconstruct

import (
	"bufio"
	"fmt"
	cp "github.com/mainak55512/qwe/compressor"
	utl "github.com/mainak55512/qwe/qweutils"
	tr "github.com/mainak55512/qwe/tracker"
	"io"
	"os"
	"strconv"
	"strings"
)

func Reconstruct(val tr.Tracker, target string, commitID int) error {
	buf := make([]byte, 1024)
	if err := cp.DecompressFile(".qwe/_object/" + val.Base); err != nil {
		return err
	}
	base_content, err := os.Open(".qwe/_object/" + val.Base)
	if err != nil {
		return err
	}
	target_content, err := os.Create(target)
	if err != nil {
		return err
	}

	_, err = io.CopyBuffer(target_content, base_content, buf)
	if err != nil {
		return err
	}
	base_content.Close()
	if err = cp.CompressFile(".qwe/_object/" + val.Base); err != nil {
		return fmt.Errorf("Compression error")
	}
	target_content.Close()

	for i, elem := range val.Versions {
		if commitID != -1 && i > commitID {
			break
		}
		if err = cp.DecompressFile(".qwe/_object/" + elem.UID); err != nil {
			return err
		}
		diff_file, err := os.Open(".qwe/_object/" + elem.UID)
		if err != nil {
			return fmt.Errorf("Error opening file: %v", err)
		}
		// defer diff_file.Close()
		diff_scanner := bufio.NewScanner(diff_file)

		base_file, err := os.Open(target)
		if err != nil {
			return fmt.Errorf("Error opening file: %v", err)
		}
		// defer base_file.Close()
		base_scanner := bufio.NewScanner(base_file)

		var output string
		diff_scanner.Scan()
		total_lines, _ := strconv.Atoi(diff_scanner.Text())
		diff_scanner.Scan()
		idx := 1
		for i := 0; i < total_lines; i++ {
			base_scanner.Scan()
			comp := strings.Split(diff_scanner.Text(), " @@@ ")
			line_number, _ := strconv.Atoi(comp[0])
			if idx == line_number {
				dec_str, err := utl.ConvStrDec(comp[1])
				if err != nil {
					return err
				}
				output += dec_str + "\n"
				diff_scanner.Scan()
			} else {
				output += base_scanner.Text() + "\n"
			}
			idx++
		}

		diff_file.Close()
		if err = cp.CompressFile(".qwe/_object/" + elem.UID); err != nil {
			return err
		}
		base_file.Close()

		output_content, err := os.Create(target)
		if err != nil {
			return err
		}

		output_writer := bufio.NewWriter(output_content)
		_, err = output_writer.WriteString(output)
		if err != nil {
			return fmt.Errorf("Can not write to base file")
		}
		if err = output_writer.Flush(); err != nil {
			return fmt.Errorf("Output file write error")
		}
		output_content.Close()
	}
	return nil
}
