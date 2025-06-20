# TidyTask 🫧

**TidyTask** is a simple Linux/Mac CLI tool for managing your to-do list, built with Go using Cobra.

## Features ✨
TidyTask is designed for users who want a no-nonsense way to manage their tasks from the terminal 

- **Fast**: Instant access to tasks, no logins or loading screens.
- **Minimal**: Shows only the information you need, no distractions. 
- **Concise**: Short, intuitive commands and powerful flags.
- **Offline**: Fully local, no internet or accounts required.

## Installation 📦
You can install TidyTask by building from source. Homebrew and AUR releases coming soon. 

##### Linux/Mac 🐧🍎
- Ensure the `go` and `git` packages are installed through your system's package manager
- Run the following commands in your terminal:
- `git clone https://github.com/tm-craggs/TidyTask.git`
- `cd tidytask`
- `go build -o tidytask`
- `sudo mv tidytask /usr/local/bin` 

##### Windows 🪟
- Install [Go for Windows](https://go.dev/dl/) and [Git for Windows](https://git-scm.com/downloads/win)
- Run the following commands in Command Prompt or PowerShell:
- `git clone https://github.com/tm-craggs/TidyTask.git`
- `cd tidytask`
- `go build -o tidytask.exe`
- `move .\tidytask.exe C:\Program Files\TidyTask\`
- Ensure the folder is added to PATH

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
Contributions are encouraged! See **Wishlist** for some ideas.

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
I’m a student building free and open source tools for fun. If you’d like to support me, you can do so via Ko-fi. There’s absolutely no pressure or expectation. Your encouragement and feedback mean just as much. <3
[![ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/G2G81GQB6Y)
