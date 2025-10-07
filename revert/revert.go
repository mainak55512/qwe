package revert

import (
	"bufio"
	"encoding/json"
	"fmt"
	cp "github.com/mainak55512/qwe/compressor"
	utl "github.com/mainak55512/qwe/qweutils"
	tr "github.com/mainak55512/qwe/tracker"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func Revert(commitNumber int, filePath string) error {
	tracker, err := tr.GetTracker()
	if err != nil {
		return fmt.Errorf("Can not retrieve Current version of %s", filePath)
	}
	fileId := utl.Hasher(filePath)

	if val, ok := tracker[fileId]; ok {
		target := filePath
		if commitNumber < 0 || commitNumber > len(val.Versions) {
			return fmt.Errorf("Not a valid commit number")
		}

		buf := make([]byte, 1024)
		cp.DecompressFile(".qwe/_object/" + val.Base)
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
			return fmt.Errorf("Copy error!")
		}
		base_content.Close()
		cp.CompressFile(".qwe/_object/" + val.Base)
		target_content.Close()

		for i, elem := range val.Versions {
			if i > commitNumber {
				break
			}
			cp.DecompressFile(".qwe/_object/" + elem.UID)
			diff_file, err := os.Open(".qwe/_object/" + elem.UID)
			if err != nil {
				log.Fatalf("Error opening file: %v", err)
			}
			// defer diff_file.Close()
			diff_scanner := bufio.NewScanner(diff_file)

			base_file, err := os.Open(target)
			if err != nil {
				log.Fatalf("Error opening file: %v", err)
			}
			// defer base_file.Close()
			base_scanner := bufio.NewScanner(base_file)

			var output string
			diff_scanner.Scan()
			total_lines, _ := strconv.Atoi(diff_scanner.Text())
			diff_scanner.Scan()
			for i := 0; i < total_lines; i++ {
				base_scanner.Scan()
				comp := strings.Split(diff_scanner.Text(), " @@@ ")
				line_number, _ := strconv.Atoi(comp[0])
				if i+1 == line_number {
					dec_str, err := utl.ConvStrDec(comp[1])
					if err != nil {
						return err
					}
					output += dec_str + "\n"
					diff_scanner.Scan()
				} else {
					output += base_scanner.Text() + "\n"
				}
			}

			diff_file.Close()
			cp.CompressFile(".qwe/_object/" + elem.UID)
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

		val.Current = val.Versions[commitNumber].UID
		tracker[fileId] = val
		marshalContent, err := json.MarshalIndent(tracker, "", " ")
		if err != nil {
			return fmt.Errorf("Commit unsuccessful!")
		}

		tracker_content, err := os.Create(".qwe/_tracker.qwe")
		if err != nil {
			return err
		}
		defer tracker_content.Close()
		writer := bufio.NewWriter(tracker_content)
		_, err = writer.Write(marshalContent)
		if err != nil {
			log.Fatal("Can not write to base file")
		}
		if err = writer.Flush(); err != nil {
			return fmt.Errorf("Tracker file write error")
		}
	}
	return nil
}
