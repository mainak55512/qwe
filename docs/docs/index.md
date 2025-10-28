# Welcome to qwe

`qwe` (pronounced kiwi) makes version control effortless. Track individual files with precision, group them seamlessly, and commit or revert changes individually or together — all in one lightweight, intuitive tool built for speed and simplicity.

## Individual Tracking

`qwe` mainly focuses on file-level version tracking. Users can commit changes of an individual file and revert back to any previous version anytime without affecting other files.

## Grouped Tracking

`qwe` also allows user to commit multiple files as a group while maintaining individual file versions. When you take a snapshot of the group, `qwe` captures the current state of all files within it. This makes rolling back a set of related changes incredibly simple.

## How does it work?

A key design choice in `qwe` is the persistence of file-level tracking, even within a group. This gives you unparalleled flexibility.

Example: Imagine you are tracking files A, B, and C in a group called "Feature-A." You still have the freedom to commit an independent revision for file A alone without affecting the group's snapshot history for B and C.

`qwe` simply groups the current commits of all the files tracked in a group. Hence under-the-hood every version of every single file is preserved. So, when you `group-revert`, `qwe` checks the commits of each file associated with that group commit and reverts the files individually to the associated versions.

This means you can:

- Maintain a clean, unified history for all files in the group (the Group Snapshot).
- Still perform granular, single-file rollbacks or commits outside the group's scope.
This approach ensures that qwe remains the flexible, non-intrusive file revision system that you can rely on.

## Flags

- `qwe` - Shows all available commands
- `init` - Initiates qwe repository
- `group-init` - Initiates a group in a qwe repository
- `groups` - Shows all the groups in the qwe repository
- `track` - Tracks a file
- `group-track` - Tracks a file in a group
- `list` - Lists all the commits of a file
- `group-list` - Lists all commits of a group
- `commit` - Commits a file
- `group-commit` - Commits a group
- `revert` - Reverts a file to a specific version
- `group-revert` - Reverts a group to a specific version
- `current` - Shows details of current commit of a file
- `group-current` - Shows details of current/specific commit of a group
- `rebase` - Reverts a file to its base version
- `recover` - Restores a file if earlier tracked
- `diff` - Shows differences between two commits of a file

## Usage

**Tracking single file**

* `qwe init` - Initiate a qwe repo
* `qwe track notes.txt` - Track 'notes.txt'
* `qwe commit notes.txt "Initial notes added"` - Commit changes of 'note.txt', commitID will be 0
* `qwe list notes.txt` - List all the commits of 'note.txt'
* `qwe revert notes.txt` - Revert 'notes.txt' to the latest committed version
* `qwe revert notes.txt 0` - Revert 'notes.txt' to commitID 0
* `qwe current notes.txt` - Check the current commit on 'notes.txt'
* `qwe diff notes.txt` - Check difference of uncommitted changes with last commit of 'notes.txt'
* `qwe diff notes.txt uncommitted 0` - Check difference of uncommitted changes with the specified commit version of 'notes.txt'
* `qwe diff notes.txt 0 1` - Check difference between two commit versions of 'notes.txt'
* `qwe rebase notes.txt` - Revert 'notes.txt' to the base version (the version from which qwe started tracking)
* `qwe recover some_deleted_file.go` - Restores a deleted file if it was earlier tracked by qwe

**Tracking a group**

* `qwe init` - Initiate a qwe repo
* `qwe group-init new_group` - Initiate a group 'new_group' in the repo, creates commitID 0
* `qwe groups` - List all the groups present in the qwe repository
* `qwe group-track new_group notes.txt` - Add 'notes.txt' to 'new_group' for tracking
* `qwe group-track new_group example.txt` - Add 'example.txt' to 'new_group' for tracking
* `qwe group-commit new_group "Initial commit"` - Commit all the files of 'new_group', commitID will be 1
* `qwe group-commit new_group "Updated commit"` - Commit all the files of 'new_group', commitID will be 2
* `qwe group-list new_group` - Lists all the commits of 'new_group'
* `qwe group-revert new_group 0` - Reverts all files to commitID 0 (the version from which qwe started tracking the group)
* `qwe group-current new_group` - Shows current commit version of 'new_group'
* `qwe group-current new_group 1` - Shows commit details of the specified version of 'new_group'

