# TidyTask 🫧

**TidyTask** is a simple CLI tool for managing your to-do list, built with Go using Cobra.

<br>

## Features ✨
TidyTask is designed for users who want a no-nonsense way to manage their tasks from the terminal 

- **Fast**: Instant access to tasks, no logins or loading screens.
- **Minimal**: Shows only the information you need, no distractions. 
- **Concise**: Short, intuitive commands and powerful flags.
- **Offline**: Fully local, no internet or accounts required.

<br>

## Installation 📦


#### Linux 🐧
- Download archive

  `wget https://github.com/tm-craggs/tidytask/archive/refs/tags/v1.0.0.tar.gz`
 
- Extract and install
  
  `tar -xzf v1.0.0.tar.gz`
  
  `chmod +x tidytask`
  
  `sudo mv tidytask /usr/local/bin`
  
- Optional: Remove source archive
  
  `rm -rf v1.0.0.tar.gz`

<br>

#### Mac 🍎

The easiest way to install TidyTask on Mac is to use **Homebrew**

- Install [Homebrew](https://brew.sh/)
- Add the tap:
  
  `brew tap tm-craggs/tidytask https://github.com/tm-craggs/homebrew-tidytask.git`
  
- Install:

- `brew install tm-craggs/tidytask/tidytask`

<br>

## Usage 🚀
Once installed, you can start using TidyTask directly from your terminal.

<br>

#### Add

To add a task run:
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

<br>

#### List

To view your to-do list use:
```
tidytask list
```

You can use flags to just view certain types of task:
```
tidytask list --priority
```

<br>

#### Complete/Remove/Reopen

These commands are formatted the same way the same way.
- `complete` → Marks a task as complete
- `remove` → Removes a task from to-do list
- `reopen` → Marks completed tasks as incomplete

<br>

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

<br>

### Search
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

<br>

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

<br>

## Contributions 🤝
Contributions are encouraged! See **Wishlist** for some ideas. I will respond to pull requests and issues as soon as possible.

<br>

## Wishlist 💭
- Add support for a configuration file so that users can control the colour scheme and column layout.
- Add a basic Terminal User Interface, designed to enhance navigation without interfering with users who 
prefer the core command line.
- Implement an action log table and multi-level undo

<br>

## Libraries Used 📚
- [Cobra](https://github.com/spf13/cobra.git)
- [go-sqlite3](https://github.com/mattn/go-sqlite3.git)
- [TableWriter for Go](https://github.com/olekukonko/tablewriter.git)
- [termenv](https://github.com/muesli/termenv)

<br>

## Support ❤️
TidyTask may be a simple tool but I've put a lot of time and love into building it, if you find it helpful I'm so glad!

I’m a student building free and open source tools for fun. If you would like to support my work, you can do so on Ko-Fi. 

There is absolutely no pressure or obligation, your kind words mean just as much to me.  <3

[![ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/G2G81GQB6Y)
