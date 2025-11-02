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

### init
---

**Description**: `init` command initiates a `qwe` repository in the current directory.

**Arguments**: It doesn't take any argument.

**Command**: `qwe init`.

**Example**: `qwe init`.

### track
---

**Description**: `track` command starts tracking the base version of a file.

**Arguments**: It takes `file-path` as the only argument.

**Command**: `qwe track [file-path]`.

**Example**: `qwe track main.go`.

### commit
---

**Description**: `commit` command commits the changes of a file.

**Arguments**: It takes `file-path` and a `commit message` as arguments.

**Command**: `qwe commit [file-path] [commit-message]`.

**Example**: `qwe commit main.go "Example commit"`.

### list
---

**Description**: `list` command lists all the commits of the file.

**Arguments**: It takes `file-path` as the argument.

**Command**: `qwe list [file-path]`.

**Example**: `qwe list main.go`.

### revert
---

**Description**: `revert` command reverts the changes of the file.

**Arguments**: It can take upto `two` arguments: `file-path`, `commit-number`

**Command**: `qwe revert [file-path] [commit-number]`.

**Example**:

- `qwe revert main.go`: this will revert main.go to its latest committed version.

- `qwe revert main.go 2`: this will revert main.go to its `2nd committed version`.

### current
---

**Description**: `current` command shows the details of the current commit of the file.

**Arguments**: It takes `file-path` as the argument.

**Command**: `qwe current [file-path]`.

**Example**: `qwe current main.go`.

### diff
---

**Description**: `diff` command shows the difference between two commits of a file.

**Arguments**: It takes upto `three` arguments. `file-path`, `first-commit-number`, `second-commit-number`

**Command**: `qwe diff [file-path] [first-commit-number] [second-commit-number]`.

**Example**:

- `qwe diff main.go`: this will show the difference between the uncommitted and latest committed version of main.go.

- `qwe diff main.go 1 2`: this will show difference of commit 1 and 2 of main.go

- `qwe diff main.go uncommitted 0`: this will show difference of uncommitted and 0th committed version of main.go.

### rebase
---

**Description**: `rebase` command reverts the file back to its base version (the version from which qwe started tracking).

**Arguments**: It takes `file-path` as the argument.

**Command**: `qwe rebase [file-path]`.

**Example**: `qwe rebase main.go`.

### recover
---

**Description**: `recover` command restores a deleted file if it was earlier tracked by qwe.

**Arguments**: It takes `file-path` as the argument.

**Command**: `qwe recover [file-path]`.

**Example**: `qwe recover main.go`.

### group-init
---

**Description**: `group-init` command initiate a logical group in the repo, creates an initial commit for the group with id 0.

**Arguments**: It takes `group-name` as the argument.

**Command**: `qwe group-init [group-name]`.

**Example**: `qwe group-init new-group`.

### group-track
---

**Description**: `group-track` command starts tracking a file or all files of a directory in a logical group.

**Arguments**: It takes `group-name` and `file-path` or `folder-path` as arguments.

**Command**: `qwe group-track [group-name] [file/folder-path]`.

**Example**:

- `qwe group-track new-group main.go`: this starts tracking `main.go` file in `new-group`.

- `qwe group-track new-group script-folder`: this starts tracking all the files (not from sub-directories) of `script-folder`directory in `new-group`.

### group-commit
---

**Description**: `group-commit` command commits all the changes of all the files tracked by a logical group.

**Arguments**: It takes `group-name` and `commit-message` as arguments.

**Command**: `qwe group-commit [group-name] [commit-message]`.

**Example**: `qwe group-commit new-group "Example commit"`.

### group-list
---

**Description**: `group-list` command lists all the commits of the specified group.

**Arguments**: It takes `group-name` as the argument.

**Command**: `qwe group-list [group-name]`.

**Example**: `qwe group-list new-group`.

### group-revert
---

**Description**: `group-revert` command reverts all the files tracked in a group to a specified commit.

**Arguments**: It takes `group-name` and `commit-number` as arguments.

**Command**: `qwe group-revert [group-name] [commit-number]`.

**Example**: `qwe group-revert new-group 1`.

### groups
---

**Description**: `groups` command lists all the groups present in the qwe repository.

**Arguments**: It doesn't take any argument.

**Command**: `qwe groups`.

**Example**: `qwe groups`.

### group-current
---

**Description**: `group-current` command shows commit details of a specific group including all the tracked files in it.

**Arguments**: It takes upto `two` arguments. `group-name`, `commit-number`

**Command**: `qwe group-current [group-name] [commit-number]`.

**Example**:

- `qwe group-current new-group`: this shows current commit version of `new-group`.

- `qwe group-current new-group 1`: this shows commit details of specified commit number.

