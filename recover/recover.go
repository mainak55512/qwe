package recover

import (
	"fmt"
	utl "github.com/mainak55512/qwe/qweutils"
	res "github.com/mainak55512/qwe/reconstruct"
	tr "github.com/mainak55512/qwe/tracker"
)

func Recover(filePath string) error {

	// Check if the file is present before recovering
	if exists := utl.FileExists(filePath); exists {
		return fmt.Errorf("File already exists")
	}

	// Get tracker details
	tracker, _, err := tr.GetTracker(0)
	if err != nil {
		return fmt.Errorf("Can not retrieve Current version of %s", filePath)
	}
	fileId := utl.Hasher(filePath)
	val, ok := tracker[fileId]
	if !ok {
		return fmt.Errorf("Error parsing tracker!")
	}

	target := filePath

	// Reconstruct the file all the way to the latest version
	if err = res.Reconstruct(val, target, -1); err != nil {
		return err
	}
	fmt.Println("Successfully recovered", filePath)
	return nil
}
