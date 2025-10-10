package diff

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	tw "text/tabwriter"
	"time"

	cp "github.com/mainak55512/qwe/compressor"
	utl "github.com/mainak55512/qwe/qweutils"
	tr "github.com/mainak55512/qwe/tracker"
)

type Changes struct {
	Prev string
	Curr string
}

// Diff shows changes between the last committed and current version of a tracked file.
func Diff(filePath string) error {
	tracker, err := tr.GetTracker()
	if err != nil {
		return fmt.Errorf("Can not retrieve Tracker, err: %s", err)
	}
	fileId := utl.Hasher(filePath)
	fileObjectId := utl.Hasher(fmt.Sprintf("%s%d", filePath, time.Now().UnixNano()))

	if val, ok := tracker[fileId]; ok {
		target := ".qwe/_object/_diff_" + fileObjectId

		new_file, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("Error opening file: %v", err)
		}
		defer new_file.Close()

		buf := make([]byte, 1024)
		if err = cp.DecompressFile(".qwe/_object/" + val.Base); err != nil {
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

		for _, elem := range val.Versions {
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

		current_file, err := os.Open(target)
		if err != nil {
			return fmt.Errorf("Error opening file: %v", err)
		}
		// defer current_file.Close()
		current_scanner := bufio.NewScanner(current_file)

		new_scanner := bufio.NewScanner(new_file)

		var diff_content []Changes

		line := 0
		for new_scanner.Scan() {
			line++
			current_scanner.Scan()
			if !bytes.Equal(current_scanner.Bytes(), new_scanner.Bytes()) {
				diff_content = append(diff_content, Changes{
					Prev: fmt.Sprintf("- %d %s", line, current_scanner.Text()),
					Curr: fmt.Sprintf("+ %d %s", line, new_scanner.Text()),
				})
			}
		}
		// diff_content = fmt.Sprintf("%d\n%s", line, diff_content)
		current_file.Close()
		os.Remove(target)
		if len(diff_content) == 0 {
			fmt.Println("No Change!")
		} else {
			w := new(tw.Writer)
			w.Init(os.Stdout, 0, 0, 0, ' ', tw.TabIndent)
			fmt.Printf("===Start Diff view===\n\n")
			fmt.Fprintln(w, "Previous\t\tCurrent")
			fmt.Fprintln(w, "--------\t\t-------")
			for _, elem := range diff_content {
				fmt.Fprintf(w, "%s  \t\t%s\n", elem.Prev, elem.Curr)
			}
			w.Flush()
			fmt.Printf("\n===End of Diff===")
		}
	}
	return nil
}
