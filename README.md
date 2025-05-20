# TidyTask

**TidyTask** is a simple Linux/Mac CLI tool for managing your to-do list, built with Go using Cobra.

> This app is currently in development and is not yet ready for use.
I'm working on it in my free time, with a first release targeted for July 2025.

### Design Philosophy
TidyTask is designed for users who want a no-nonsense way to manage their tasks from the terminal 

- **Fast**: No login or loading times. Tasks are instantly accessible through a small set of simple commands.
- **Minimal**: Shows only the information you need, nothing more. 
- **Concise**: Commands are simple and quick to type. Flags allow for fine-grained control, but are kept to a minimum.
- **Offline**: Everything stays on your machine and works without internet. No account needed.

### Contributions
Contributions are encouraged and will be welcome after first release. See **Wishlist** for future ideas.

### Pre-Release TODO
- Add `search` command
- Add confirmation prompts for destructive tasks
- Store empty dates as `NULL` , instead of empty strings
- Add `completion_date` to tasks table
- Use completion date to evaluate if due dates have been met 
- Improve documentation and CLI `--help` output
- General code and documentation cleanup
- Package for release and ensure README has installation instructions for supported platforms

### Wishlist (Post-Release Ideas)
- Add support for a configuration file, so that users can control colour scheme and column layout.
- Add a basic Terminal User Interface, designed to enhance navigation without interfering with users who 
prefer the core command line. 
- Flatpak release