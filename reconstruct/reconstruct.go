package reconstruct

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	// bh "github.com/mainak55512/qwe/binaryhandler"
	cp "github.com/mainak55512/qwe/compressor"
	er "github.com/mainak55512/qwe/qwerror"
	utl "github.com/mainak55512/qwe/qweutils"
	tr "github.com/mainak55512/qwe/tracker"
)

const (
	LastVersion = -1 // Last version - all commits
	BaseVersion = -2 // Only use base version, no commits
)

// Applies previous commits till the commitID supplied on to the base version
func Reconstruct(val tr.Tracker, target string, commitID int) error {
	// if strings.HasPrefix(val.Base, "_bin_") {
	// 	if err := bh.RevertBinFile(target, val.Current); err != nil {
	// 		return err
	// 	}
	// 	return nil
	// }
	buf := make([]byte, 1024)
	targetInfo, err := os.Lstat(target)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if err == nil && targetInfo.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("refusing to reconstruct symlink %s", target)
	}

	// Build the complete result beside the target and replace it only after
	// every object has been read successfully.
	targetDir := filepath.Dir(target)
	tmp, err := os.CreateTemp(targetDir, ".qwe-reconstruct-*")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)

	fileMode := os.FileMode(0600)
	if targetInfo != nil {
		fileMode = targetInfo.Mode().Perm()
	}
	if err := tmp.Chmod(fileMode); err != nil {
		tmp.Close()
		return err
	}

	// Decompress the base varient
	if err := cp.DecompressFile(".qwe/_object/" + val.Base); err != nil {
		tmp.Close()
		return err
	}

	base_content, err := os.Open(".qwe/_object/" + val.Base)
	if err != nil {
		tmp.Close()
		return err
	}

	// Copy the content from base varient to the file
	_, err = io.CopyBuffer(tmp, base_content, buf)
	baseCloseErr := base_content.Close()
	if err != nil {
		tmp.Close()
		return err
	}
	if baseCloseErr != nil {
		tmp.Close()
		return baseCloseErr
	}

	// Compress the base varient
	if err = cp.CompressFile(".qwe/_object/" + val.Base); err != nil {
		tmp.Close()
		return err
	}

	// if commitID is -2, that means only base varient is needed
	if commitID == BaseVersion {
		if err := tmp.Sync(); err != nil {
			tmp.Close()
			return err
		}
		if err := tmp.Close(); err != nil {
			return err
		}
		return os.Rename(tmpPath, target)
	}

	// Loop through the file versions and apply the changes to the base varient one by one
	for i, elem := range val.Versions {

		// Will stop if the specified commitID is reached; -1 means it will cover all versions
		if commitID != LastVersion && i > commitID {
			break
		}

		if err = cp.DecompressFile(".qwe/_object/" + elem.UID); err != nil {
			tmp.Close()
			return err
		}

		diff_file, err := os.Open(".qwe/_object/" + elem.UID)
		if err != nil {
			tmp.Close()
			return err
		}
		// defer diff_file.Close()
		diff_scanner := bufio.NewScanner(diff_file)

		if err := tmp.Close(); err != nil {
			diff_file.Close()
			return err
		}
		base_file, err := os.Open(tmpPath)
		if err != nil {
			diff_file.Close()
			return err
		}
		// defer base_file.Close()
		base_scanner := bufio.NewScanner(base_file)

		var output strings.Builder

		// This will retrieve the line number from the commit file,
		// reconstructed file should only have these many lines in it.
		if !diff_scanner.Scan() {
			diff_file.Close()
			base_file.Close()
			return fmt.Errorf("invalid revision %s: missing line count", elem.UID)
		}
		total_lines, err := strconv.Atoi(diff_scanner.Text())
		if err != nil {
			diff_file.Close()
			base_file.Close()
			return err
		}
		if total_lines < 0 {
			diff_file.Close()
			base_file.Close()
			return fmt.Errorf("invalid revision %s: negative line count", elem.UID)
		}

		// Retrieving a line from commit file
		hasDiff := diff_scanner.Scan()
		for idx := 1; idx <= total_lines; idx++ {
			baseLine := ""
			if base_scanner.Scan() {
				baseLine = base_scanner.Text()
			}

			if hasDiff {
				comp := strings.SplitN(diff_scanner.Text(), " @@@ ", 2)
				if len(comp) != 2 {
					diff_file.Close()
					base_file.Close()
					return fmt.Errorf("invalid revision %s: malformed line", elem.UID)
				}
				lineNumber, err := strconv.Atoi(comp[0])
				if err != nil || lineNumber < idx {
					diff_file.Close()
					base_file.Close()
					return fmt.Errorf("invalid revision %s: invalid line number", elem.UID)
				}
				if lineNumber == idx {
					decoded, err := utl.ConvStrDec(comp[1])
					if err != nil {
						diff_file.Close()
						base_file.Close()
						return err
					}
					output.WriteString(decoded)
					hasDiff = diff_scanner.Scan()
				} else {
					output.WriteString(baseLine)
				}
			} else {
				output.WriteString(baseLine)
			}
			output.WriteByte('\n')
		}
		if hasDiff || diff_scanner.Err() != nil || base_scanner.Err() != nil {
			diff_file.Close()
			base_file.Close()
			if diff_scanner.Err() != nil {
				return diff_scanner.Err()
			}
			if base_scanner.Err() != nil {
				return base_scanner.Err()
			}
			return fmt.Errorf("invalid revision %s: changes exceed line count", elem.UID)
		}

		diff_file.Close()

		// Compress the commit file
		if err = cp.CompressFile(".qwe/_object/" + elem.UID); err != nil {
			return err
		}
		base_file.Close()

		next, err := os.CreateTemp(targetDir, ".qwe-reconstruct-*")
		if err != nil {
			return err
		}
		nextPath := next.Name()

		// Write all the changes to the file
		output_writer := bufio.NewWriter(next)
		_, err = output_writer.WriteString(output.String())
		if err != nil {
			next.Close()
			os.Remove(nextPath)
			return er.BaseWriteErr
		}
		if err = output_writer.Flush(); err != nil {
			next.Close()
			os.Remove(nextPath)
			return er.OutputWriteErr
		}
		if err = next.Chmod(fileMode); err != nil {
			next.Close()
			os.Remove(nextPath)
			return err
		}
		if err = next.Close(); err != nil {
			os.Remove(nextPath)
			return err
		}
		if err = os.Rename(nextPath, tmpPath); err != nil {
			os.Remove(nextPath)
			return err
		}
		tmp, err = os.OpenFile(tmpPath, os.O_RDWR, fileMode)
		if err != nil {
			return err
		}
	}

	if err := tmp.Sync(); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmpPath, target)
}
