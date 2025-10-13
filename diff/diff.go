package diff

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
	"time"

	utl "github.com/mainak55512/qwe/qweutils"
	res "github.com/mainak55512/qwe/reconstruct"
	tr "github.com/mainak55512/qwe/tracker"
)

type Changes struct {
	Prev string
	Curr string
}

func Diff(filePath, commitID1Str, commitID2Str string) error {

	// Only allow if both are either empty or non-empty
	if !((commitID1Str == "") == (commitID2Str == "")) {
		return fmt.Errorf("Argument number missmatch")
	}

	// Get details from _tracker.qwe
	tracker, err := tr.GetTracker()
	if err != nil {
		return fmt.Errorf("Can not retrieve Tracker, err: %s", err)
	}

	fileId := utl.Hasher(filePath)
	fileObjectId := utl.Hasher(fmt.Sprintf("%s%d", filePath, time.Now().UnixNano()))

	// Check if file is being tracked
	val, ok := tracker[fileId]
	if !ok {
		return fmt.Errorf("File is not tracked!")
	}

	// Will run if no commit id is passed or both commit id is passed and first one is 'uncommitted'
	if (commitID1Str == "" && commitID2Str == "") || commitID1Str == "uncommitted" {

		target := ".qwe/_object/_diff_" + fileObjectId

		if commitID2Str != "" {
			commitID, err := strconv.Atoi(commitID2Str)
			if err != nil {
				return err
			}

			// Reconstruct till the specified commit id
			if err = res.Reconstruct(val, target, commitID); err != nil {
				return err
			}
		} else {
			// Reconstruct till the current version
			if err = res.Reconstruct(val, target, -1); err != nil {
				return err
			}
		}

		// As commitID1Str is either empty or 'uncommitted', need to compare uncommited changes of the file
		// with the latest committed version or with the version specified by the commitID2Str
		new_file, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("Error opening file: %v", err)
		}
		defer new_file.Close()

		current_file, err := os.Open(target)
		if err != nil {
			return fmt.Errorf("Error opening file: %v", err)
		}

		current_scanner := bufio.NewScanner(current_file)

		new_scanner := bufio.NewScanner(new_file)

		var diff_content []Changes

		// Check the differences
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

		current_file.Close()
		os.Remove(target)
		if len(diff_content) == 0 {
			fmt.Println("No Change!")
		} else {
			fmt.Printf("===Start Diff view===\n\n")
			for _, elem := range diff_content {
				fmt.Println(elem.Prev + "\n" + elem.Curr)
				fmt.Println()
			}
			fmt.Printf("\n===End of Diff===")
		}
	} else {

		// This part will execute if both commitIDs are supplied through the command line
		var commit1, commit2 int
		if commitID1Str != "" {
			commit1, err = strconv.Atoi(commitID1Str)
			if err != nil {
				return err
			}
		}
		if commitID2Str != "" {
			commit2, err = strconv.Atoi(commitID2Str)
			if err != nil {
				return err
			}
		}
		src := ".qwe/_object/_diff_src_" + fileObjectId
		dest := ".qwe/_object/_diff_dest_" + fileObjectId

		// reconstruct till first commitID
		if err = res.Reconstruct(val, src, commit1); err != nil {
			return err
		}

		// Reconstruct till second commitID
		if err = res.Reconstruct(val, dest, commit2); err != nil {
			return err
		}
		new_file, err := os.Open(dest)
		if err != nil {
			return fmt.Errorf("Error opening file: %v", err)
		}
		defer new_file.Close()

		current_file, err := os.Open(src)
		if err != nil {
			return fmt.Errorf("Error opening file: %v", err)
		}
		// defer current_file.Close()
		current_scanner := bufio.NewScanner(current_file)

		new_scanner := bufio.NewScanner(new_file)

		var diff_content []Changes

		// Compare for changes
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

		// Remove temporary files
		os.Remove(src)
		os.Remove(dest)

		if len(diff_content) == 0 {
			fmt.Println("No Change!")
		} else {
			fmt.Printf("===Start Diff view===\n\n")
			for _, elem := range diff_content {
				fmt.Println(elem.Prev + "\n" + elem.Curr)
				fmt.Println()
			}
			fmt.Printf("\n===End of Diff===")
		}
	}
	return nil
}
