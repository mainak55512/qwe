package cli

import (
	"fmt"
	"os"
	"strconv"

	cm "github.com/mainak55512/qwe/commit"
	"github.com/mainak55512/qwe/diff"
	utl "github.com/mainak55512/qwe/qweutils"
	rv "github.com/mainak55512/qwe/revert"
)

func helpText() {
	fmt.Println("Version: 0.0.2")
	fmt.Println("Available commands:")
	fmt.Println("qwe init")
	fmt.Println("qwe track <file-path>")
	fmt.Println("qwe list <file-path>")
	fmt.Println("qwe commit <file-path> \"<commit message>\"")
	fmt.Println("qwe revert <file-path> <commit-id>")
	fmt.Println("qwe diff <file-path>")
}

func HandleArgs() error {
	command_list := os.Args[1:]

	if len(command_list) == 0 {
		helpText()
	} else {

		switch command_list[0] {
		case "init":
			{
				if len(command_list) != 1 {
					return fmt.Errorf("Init command doesn't take any argument")
				}
				if err := utl.Init(); err != nil {
					return err
				}
			}
		case "track":
			{
				if len(command_list) != 2 {
					return fmt.Errorf("Track command accepts one argument")
				}
				if err := utl.StartTracking(command_list[1]); err != nil {
					return err
				}
			}
		case "commit":
			{
				if len(command_list) != 3 {
					return fmt.Errorf("Commit command accepts two arguments")
				}
				if err := cm.CommitUnit(command_list[1], command_list[2]); err != nil {
					return err
				}
			}
		case "list":
			{
				if len(command_list) != 2 {
					return fmt.Errorf("List command accepts one argument")
				}
				if err := utl.GetCommitList(command_list[1]); err != nil {
					return err
				}
			}
		case "revert":
			{
				if len(command_list) != 3 {
					return fmt.Errorf("Revert command accepts two arguments")
				}
				commitNumber, err := strconv.Atoi(command_list[2])
				if err != nil {
					return fmt.Errorf("Not a valid commit number")
				}
				if err := rv.Revert(commitNumber, command_list[1]); err != nil {
					return err
				}
			}
		case "diff":
			if len(command_list) != 2 {
				return fmt.Errorf("diff command accepts filename as argument")
			}
			if err := diff.Diff(command_list[1]); err != nil {
				return err
			}
		default:
			{
				helpText()
			}
		}
	}
	return nil
}
