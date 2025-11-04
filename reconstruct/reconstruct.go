package reconstruct

import (
	"bufio"
	bh "github.com/mainak55512/qwe/binaryhandler"
	cp "github.com/mainak55512/qwe/compressor"
	er "github.com/mainak55512/qwe/qwerror"
	utl "github.com/mainak55512/qwe/qweutils"
	tr "github.com/mainak55512/qwe/tracker"
	"io"
	"os"
	"strconv"
	"strings"
)

// Applies previous commits till the commitID supplied on to the base version
func Reconstruct(val tr.Tracker, target string, commitID int) error {
	if strings.HasPrefix(val.Base, "_bin_") {
		if err := bh.RevertBinFile(target, val.Current); err != nil {
			return err
		}
		return nil
	}
	buf := make([]byte, 1024)

	// Decompress the base varient
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

	// Copy the content from base varient to the file
	_, err = io.CopyBuffer(target_content, base_content, buf)
	if err != nil {
		return err
	}
	base_content.Close()

	// Compress the base varient
	if err = cp.CompressFile(".qwe/_object/" + val.Base); err != nil {
		return err
	}
	target_content.Close()

	// if commitID is -2, that means only base varient is needed
	if commitID == -2 {
		return nil
	}

	// Loop through the file versions and apply the changes to the base varient one by one
	for i, elem := range val.Versions {

		// Will stop if the specified commitID is reached; -1 means it will cover all versions
		if commitID != -1 && i > commitID {
			break
		}

		if err = cp.DecompressFile(".qwe/_object/" + elem.UID); err != nil {
			return err
		}

		diff_file, err := os.Open(".qwe/_object/" + elem.UID)
		if err != nil {
			return err
		}
		// defer diff_file.Close()
		diff_scanner := bufio.NewScanner(diff_file)

		base_file, err := os.Open(target)
		if err != nil {
			return err
		}
		// defer base_file.Close()
		base_scanner := bufio.NewScanner(base_file)

		var output string

		// This will retrieve the line number from the commit file,
		// reconstructed file should only have these many lines in it.
		diff_scanner.Scan()

		total_lines, err := strconv.Atoi(diff_scanner.Text())
		if err != nil {
			return err
		}

		// Retrieving a line from commit file
		diff_scanner.Scan()
		idx := 1
		for i := 0; i < total_lines; i++ {

			// Retrieve a line from base file
			base_scanner.Scan()

			// Split the line from the commit file
			// to get the line number where the change occured and the content as string
			comp := strings.Split(diff_scanner.Text(), " @@@ ")
			line_number, _ := strconv.Atoi(comp[0])

			// if current line number of the base file is same as the retrieved line number
			// from the commit file, then decompress the string from the commit file and
			// add it to the output string and scan the next line from the commit file,
			// else add the line from the base file to the output string
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

		// Compress the commit file
		if err = cp.CompressFile(".qwe/_object/" + elem.UID); err != nil {
			return err
		}
		base_file.Close()

		output_content, err := os.Create(target)
		if err != nil {
			return err
		}

		// Write all the changes to the file
		output_writer := bufio.NewWriter(output_content)
		_, err = output_writer.WriteString(output)
		if err != nil {
			return er.BaseWriteErr
		}
		if err = output_writer.Flush(); err != nil {
			return er.OutputWriteErr
		}
		output_content.Close()
	}
	return nil
}
