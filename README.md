<div align="center">
  <img src="./assets/qwe.png" alt="qwe Logo">
</div>

# qwe (kiwi) - lightweight, flexible, file-first version/revision control system


**qwe** (pronounced *kiwi*) makes version control effortless.
Track individual files with precision, group them seamlessly, and commit or revert changes individually or together — all in one lightweight, intuitive tool built for speed and simplicity.

![Static Badge](https://img.shields.io/badge/version-control-system?style=for-the-badge&logo=git&logoColor=white&color=blue) ![Static Badge](https://img.shields.io/badge/revision-control-system?style=for-the-badge&logo=git&logoColor=white&color=red) ![Downloads](https://img.shields.io/github/downloads/mainak55512/qwe/total?style=for-the-badge) ![GitHub License](https://img.shields.io/github/license/mainak55512/qwe?style=for-the-badge)

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white) ![GitHub repo size](https://img.shields.io/github/repo-size/mainak55512/qwe?style=for-the-badge)

## Features

- **File-based version control** — track versions of specific files independently.  
- **Easy commits** — record changes with a simple commit message.  
- **Revert anytime** — roll back a single file without affecting others.  
- **Grouped snapshot** — Track multiple files with ease for collective commit and revert.
- **Simple and fast** — minimal setup, no complex repository management.

## How does it work?

`qwe` allows you to track individual files separately or bundle related files into a single, named snapshot for easy tracking and rollback.

**Group Creation:** Create a logical group (e.g., "Project X Assets," "Configuration Files") that contains multiple individual files.

**Unified Tracking:** When you take a snapshot of the group, qwe captures the current state of all files within it. This makes rolling back a set of related changes incredibly simple.

<div align="center">
  <img src="./assets/qwe-diagram.png" alt="qwe diagram Logo">
</div>

## Installation

You can install **qwe** in two ways:

### 1️⃣ Using Go Package Manager
```bash
go install github.com/mainak55512/qwe@latest
```

Make sure your Go environment’s `GOPATH/bin` is added to your system `PATH`.

### 2️⃣ Using Prebuilt Executables
Download the prebuilt binary for your platform from the **[Releases](https://github.com/mainak55512/qwe/releases)** section of this repository and add it to your PATH.

## Example Usage

```bash
qwe init
qwe track notes.txt
qwe commit notes.txt "Initial notes added" // -> commitID 0
qwe commit notes.txt "Updated with new ideas" // -> commitID 1
qwe revert notes.txt 0
qwe group-init new_group
qwe group-track new_group notes.txt example.txt README.md
qwe group-commit new_group "Initial commit" // -> commitID 1
qwe group-commit new_group "Updated commit" // -> commitID 2
qwe group-revert new_group 0 // -> Revert back to base version (to the version from which group tracking started)
```

## Documentation

Full documentation is available at [https://mainak55512.github.io/qwe](https://mainak55512.github.io/qwe/).

## Why qwe?

- Ideal for **independent file tracking** without setting up a full Git repo.  
- Track multiple files while maintaining **individual file versions**.
- Great for **scripts, configs, notes, or documents**.  
- Simple CLI interface — no branching or merging headaches.


## Code of Conduct

Check out code of conduct for `qwe` [here](./CODE_OF_CONDUCT.md).


## Contributions

Contributions are welcome. Please have a look at the [contribution guidelines](./CONTRIBUTING.md) for `qwe` before contributing.


## License

[MIT License © 2025](./LICENSE)
