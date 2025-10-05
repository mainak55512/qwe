package revert

import (
	"bufio"
	"encoding/json"
	"fmt"
	utl "github.com/mainak55512/qwe/qweutils"
	tr "github.com/mainak55512/qwe/tracker"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func Revert(commitNumber int, filePath string) error {
	tracker, err := tr.GetTracker()
	if err != nil {
		return fmt.Errorf("Can not retrieve Current version of %s", filePath)
	}
	fileId := utl.Hasher(filePath)
	fileObjectId := utl.Hasher(fmt.Sprintf("%s%d", filePath, time.Now().UnixNano()))

	if val, ok := tracker[fileId]; ok {
		if commitNumber < 0 || commitNumber > len(val.Versions) {
			return fmt.Errorf("Not a valid commit number")
		}
		base_content, _ := os.ReadFile(".qwe/_object/" + val.Base)
		err = os.WriteFile(".qwe/_object/_base_"+fileObjectId, base_content, 0644)
		for i, elem := range val.Versions {
			if i > commitNumber {
				break
			}
			diff_file, err := os.Open(".qwe/_object/" + elem.UID)
			if err != nil {
				log.Fatalf("Error opening file: %v", err)
			}
			defer diff_file.Close()
			diff_scanner := bufio.NewScanner(diff_file)

			base_file, err := os.Open(".qwe/_object/_base_" + fileObjectId)
			if err != nil {
				log.Fatalf("Error opening file: %v", err)
			}
			base_scanner := bufio.NewScanner(base_file)

			var output string
			diff_scanner.Scan()
			total_lines, _ := strconv.Atoi(diff_scanner.Text())
			diff_scanner.Scan()
			// idx := 1
			for i := 0; i < total_lines; i++ {
				base_scanner.Scan()
				comp := strings.Split(diff_scanner.Text(), " @@@ ")
				line_number, _ := strconv.Atoi(comp[0])
				if i+1 == line_number {
					output += utl.ConvStrDec(comp[1]) + "\n"
					diff_scanner.Scan()
				} else {
					output += utl.ConvStrDec(base_scanner.Text()) + "\n"
				}
				// idx++
			}
			// fmt.Println(output)
			err = os.WriteFile(".qwe/_object/_base_"+fileObjectId, []byte(output), 0644)
			if err != nil {
				log.Fatal("Can not write to base file")
			}
			base_file.Close()
		}

		output_content, _ := os.ReadFile(".qwe/_object/_base_" + fileObjectId)
		err = os.WriteFile(filePath, output_content, 0644)
		os.Remove(".qwe/_object/_base_" + fileObjectId)

		val.Current = val.Versions[commitNumber].UID
		tracker[fileId] = val
		marshalContent, err := json.MarshalIndent(tracker, "", " ")
		if err != nil {
			return fmt.Errorf("Commit unsuccessful!")
		}

		if err = os.WriteFile(".qwe/_tracker.qwe", marshalContent, 0644); err != nil {
			return fmt.Errorf("Commit unsuccessful!")
		}
	}
	return nil
}
