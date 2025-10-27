package recover

import (
	"fmt"
	er "github.com/mainak55512/qwe/qwerror"
	utl "github.com/mainak55512/qwe/qweutils"
	res "github.com/mainak55512/qwe/reconstruct"
	tr "github.com/mainak55512/qwe/tracker"
)

// Restores a deleted file if it was earlier tracked by qwe
func Recover(filePath string) error {

	// Check if the file is present before recovering
	if exists := utl.FileExists(filePath); exists {
		return er.FileExists
	}

	// Get tracker details
	tracker, _, err := tr.GetTracker(0)
	if err != nil {
		return err
	}
	fileId := utl.Hasher(filePath)
	val, ok := tracker[fileId]
	if !ok {
		return er.FileNotTracked
	}

	target := filePath

	// Reconstruct the file all the way to the latest version
	if err = res.Reconstruct(val, target, -1); err != nil {
		return err
	}
	fmt.Println("Successfully recovered", filePath)
	return nil
}
