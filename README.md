# TidyTask 🫧

**TidyTask** is a simple CLI tool for managing your to-do list, built with Go using Cobra.

## Features ✨
TidyTask is designed for users who want a no-nonsense way to manage their tasks from the terminal 

- **Fast**: Instant access to tasks, no logins or loading screens.
- **Minimal**: Shows only the information you need, no distractions. 
- **Concise**: Short, intuitive commands and powerful flags.
- **Offline**: Fully local, no internet or accounts required.

## Installation 📦
You can install TidyTask using precompiled binaries or building from source.


#### Linux 🐧
- Download archive
  - Download from [releases](https://github.com/tm-craggs/tidytask/releases)
  - Download using wget: `wget https://github.com/tm-craggs/tidytask/releases/download/untagged-abc31a2cef31ac408b1e/tidytask-linux-amd64.tar.gz`
 
- Extract and install
  - Extract binary: `tar -xzf tidytask-linux-amd64.tar.gz`
  - Grant execution permissions: `chmod +x tidytask`
  - Move to binaries folder: `sudo mv tidytask /usr/local/bin`
 
- Run
  - You can now run `tidytask` from anywhere 


#### Mac 🍎
- Download the archive from [releases](https://github.com/tm-craggs/tidytask/releases)
  - For Intel: `tidytask-mac-amd64.tar.gz`
  - For Apple Silicon (M1/M2): `tidytask-mac-arm64.tar.gz`
 
- Extract and install
  - Extract binary: `tar -xzf tidytask-mac-<arch>.tar.gz`
  - Grant execution permissions: `chmod +x tidytask`
  - Move to binaries folder: `sudo mv tidytask /usr/local/bin`

- Run
  - You can now run `tidytask` from anywhere 


#### Windows 🪟
- Download `tidytask-windows-amd64.zip` from [releases](https://github.com/tm-craggs/tidytask/releases)
- Extract tidytask.exe
- Double-click to run, or use from the terminal. `.\tidytask.exe`
- Optional: Add the containing folder to PATH to run `tidytask` from anywhere

#### Build from Source 🔧
If you prefer, or the precompiled binaries don't work, you can install by compiling from source.

- Ensure `go` is installed and usable from the command line

- Download Source Code:
  - From [releases](https://github.com/tm-craggs/tidytask/releases)
  - Using git: `git clone https://github.com/tm-craggs/tidytask`

- Compile binary
  - Enter the repo: `cd tidytask`
  - Compile binary: `go build -o tidytask`

- Install
  - Move to binaries folder: `mv tidytask /usr/local/bin`
  - Optional, delete source code: `cd ..` `rm -rf tidytask`

## Usage 🚀
Once installed, you can start using TidyTask directly from your terminal.

#### Add

To add a task, run:
```
tidytask add "First Task"
```

Tasks can be given a due date using --due:
```
tidytask add "Finish Homework" --due 2025-06-01
```

Tasks can also be marked as high priority using --priority:
```
tidytask add "Submit Essay" --due 2025-06-25 --priority
```
#### List

To view your to-do list use:
```
tidytask list
```

You can use flags to just view certain types of task:
```
tidytask list --priority
```

#### Complete/Remove/Reopen
These commands are formatted the same way the same way.
- `complete` → Marks a task as complete
- `remove` → Removes a task from to-do list
- `reopen` → Marks completed tasks as incomplete

To complete task numbers 1, 2 and 3:
```
tidytask complete 1 2 3
```

These commands support batch operations through the --all flag:
```
tidytask remove --all
```

The --all flag can be used with constrictions to target specific types of task
```
tidytask reopen --all --priority
```

#### Search
The search command displays all tasks that match a certain keyword. By default, it searches all fields.
```
tidytask search "homework"
```

You can specify with flags which fields to search. For example, this is a search for tasks with the number 2025 in their due date:
```
tidytask search 2025 --due
```

Like with list, you can use flags just view certain types of task. For example, this is a search for all priority tasks with essay in the title:
```
tidytask search essay --title --complete
```

#### Undo

To reverse the previous action, run:
```
tidytask undo
```

#### Help
You can view all commands and general information, use:
```
tidytask --help
```
To view the help page for a specific command use:
```
tidytask [command] --help
```

## Contributions 🤝
Contributions are encouraged! See **Wishlist** for some ideas. I will respond to pull requests and issues as soon as possible.

## Wishlist 💭
- Add support for a configuration file so that users can control the colour scheme and column layout.
- Add a basic Terminal User Interface, designed to enhance navigation without interfering with users who 
prefer the core command line.
- Implement an action log table and multi-level undo

## Libraries Used 📚
- [Cobra](https://github.com/spf13/cobra.git)
- [go-sqlite3](https://github.com/mattn/go-sqlite3.git)
- [TableWriter for Go](https://github.com/olekukonko/tablewriter.git)
- [termenv](https://github.com/muesli/termenv)

## Support ❤️
If TidyTask has been helpful to you, I’m really glad to hear it!
I’m a student building free and open source tools for fun, if you’d like to support me you can do so via Ko-fi. There’s absolutely no pressure or expectation. Your encouragement and feedback mean just as much. <3

[![ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/G2G81GQB6Y)
