package initializer

import (
	"encoding/json"
	"fmt"
	er "github.com/mainak55512/qwe/qwerror"
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
		return er.RepoAlreadyInit
	} else {

		// Create objects directory
		if err := os.MkdirAll(qwePath+"/_object/", os.ModePerm); err != nil {
			return er.RepoInitError
		}
		// Create _tracker.qwe file
		if _, err := os.Create(qwePath + "/_tracker.qwe"); err != nil {
			os.RemoveAll(qwePath)
			return er.RepoInitError
		}
		// Create _group_tracker.qwe file
		if _, err := os.Create(qwePath + "/_group_tracker.qwe"); err != nil {
			os.RemoveAll(qwePath)
			return er.RepoInitError
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

// Initiate a group in a qwe repository
func GroupInit(groupName string) error {

	qwePath := ".qwe"

	if exists := utl.FolderExists(qwePath); !exists {
		return er.RepoNotFound
	}

	// Get group tracker
	_, groupTracker, err := tr.GetTracker(1)
	if err != nil {
		return err
	}

	groupID := utl.Hasher(groupName)
	groupObjectId := "_group_" + utl.Hasher(fmt.Sprintf("%s%d", groupName, time.Now().UnixNano()))

	// Check if group is already tracked
	if _, ok := groupTracker[groupID]; ok {
		return er.GrpAlreadyTracked
	}

	// Instantiate a logical group in the group tracker
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
		return er.CommitUnsuccessful
	}

	// Update the tracker
	if err = tr.SaveTracker(1, marshalContent); err != nil {
		return err
	}
	fmt.Println("Started tracking group ", groupName)
	return nil
}
