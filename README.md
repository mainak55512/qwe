# 🥝 qwe — A Lightweight File-Level Version Control System

**qwe** (pronounced *kiwi*) is a simple yet powerful version control system that tracks **individual files**, not entire projects.  
Unlike Git, which manages repositories as a whole, qwe provides a more granular approach — perfect for quick file-level tracking, experimentation, or standalone scripts.

---

## 🚀 Features

- 📁 **File-based version control** — track versions of specific files independently.  
- 💾 **Easy commits** — record changes with a simple commit message.  
- 🔙 **Revert anytime** — roll back a single file without affecting others.  
- 🧾 **Simple and fast** — minimal setup, no complex repository management.

---

## ⚙️ Installation

You can install **qwe** in two ways:

### 1️⃣ Using Go Package Manager
```bash
go install github.com/mainak55512/qwe@latest
```

Make sure your Go environment’s `GOPATH/bin` is added to your system `PATH`.

### 2️⃣ Using Prebuilt Executables
Download the prebuilt binary for your platform from the **[Releases](https://github.com/mainak55512/qwe/releases)** section of this repository and add it to your PATH.

---

## 🧰 Commands

| Command | Description |
|----------|-------------|
| `qwe init` | Initialize qwe in the current directory |
| `qwe track <file-path>` | Start tracking a file |
| `qwe list <file-path>` | List all commits for the specified file |
| `qwe commit <file-path> "<commit message>"` | Commit changes to the file with a message |
| `qwe revert <file-path> <commit-id>` | Revert the file to a previous version |

---

## 💡 Example Usage

```bash
qwe init
qwe track notes.txt
qwe commit notes.txt "Initial notes added"
qwe commit notes.txt "Updated with new ideas"
qwe list notes.txt
qwe revert notes.txt 2
```

---

## 🧩 Why qwe?

- Ideal for **independent file tracking** without setting up a full Git repo.  
- Great for **scripts, configs, notes, or documents**.  
- Simple CLI interface — no branching or merging headaches.

---

## 🏗️ Future Plans

- Add diff support to compare file versions  
- Enable remote file sync  
- Provide optional GUI interface  

---

## 📜 License

MIT License © 2025
