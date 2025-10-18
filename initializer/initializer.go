package initializer

import (
	"encoding/json"
	"fmt"
	utl "github.com/mainak55512/qwe/qweutils"
	tr "github.com/mainak55512/qwe/tracker"
	"os"
	"time"
)

// Initiates qwe repository
func Init() error {
	qwePath := ".qwe"

	// Check if qwe is already initialized
	if exists := utl.FolderExists(qwePath); exists {
		return fmt.Errorf("Repository already initiated!")
	} else {

		// Create objects directory
		if err := os.MkdirAll(qwePath+"/_object/", os.ModePerm); err != nil {
			return fmt.Errorf("Can not initiate repository!")
		}
		// Create _tracker.qwe file
		if _, err := os.Create(qwePath + "/_tracker.qwe"); err != nil {
			os.RemoveAll(qwePath)
			return fmt.Errorf("Can not initiate repository!")
		}
		// Create _group_tracker.qwe file
		if _, err := os.Create(qwePath + "/_group_tracker.qwe"); err != nil {
			os.RemoveAll(qwePath)
			return fmt.Errorf("Can not initiate repository!")
		}
		// Initialize the tracker with '{}'
		if err := tr.SaveTracker(0, []byte("{}")); err != nil {
			return err
		}
		// Initialize the group tracker with '{}'
		if err := tr.SaveTracker(1, []byte("{}")); err != nil {
			return err
		}
	}
	fmt.Println("QWE initiated")
	return nil
}

func GroupInit(groupName string) error {
	qwePath := ".qwe"
	if exists := utl.FolderExists(qwePath); !exists {
		return fmt.Errorf("No qwe repository found!")
	}
	_, groupTracker, err := tr.GetTracker(1)
	if err != nil {
		return err
	}
	groupID := utl.Hasher(groupName)
	groupObjectId := "_group_" + utl.Hasher(fmt.Sprintf("%s%d", groupName, time.Now().UnixNano()))

	if _, ok := groupTracker[groupID]; ok {
		return fmt.Errorf("Group is already being tracked!")
	}
	groupTracker[groupID] = tr.GroupTracker{
		GroupName:    groupName,
		Current:      groupObjectId,
		VersionOrder: []string{groupObjectId},
		Versions: map[string]tr.GroupVersionDetails{
			groupObjectId: {
				CommitMessage: "Initial Tracking",
				Files:         map[string]tr.FileDetails{},
			},
		},
	}

	marshalContent, err := json.MarshalIndent(groupTracker, "", " ")
	if err != nil {
		return fmt.Errorf("Commit unsuccessful!")
	}

	// Update the tracker
	if err = tr.SaveTracker(1, marshalContent); err != nil {
		return err
	}
	fmt.Println("Started tracking group ", groupName)
	return nil
}
