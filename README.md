# TidyTask

**TidyTask** is a simple Linux/Mac CLI tool for managing your to-do list, built with Go using Cobra.

> This app is currently in development. Core functionality is in place, however needs polish before release

### Design Philosophy
TidyTask is designed for users who want a no-nonsense way to manage their tasks from the terminal 

- **Fast**: No login or loading times. Tasks are instantly accessible through a small set of simple commands.
- **Minimal**: Shows only the information you need, nothing more. 
- **Concise**: Commands are simple and quick to type. Flags allow for fine-grained control, but are kept to a minimum.
- **Offline**: Everything stays on your machine and works without an internet connection. No account needed.

### Contributions
Contributions are encouraged and will be welcome after the first release. See **Wishlist** for future ideas.

### Pre-Release TODO
- Rework reopen and complete to function the same as remove
- Ensure all necessary commands only take in one input, and throw error otherwise
- Finish documentation
- Package for release

### Wishlist (Post-Release Ideas)
- Add support for a configuration file so that users can control the colour scheme and column layout.
- Add a basic Terminal User Interface, designed to enhance navigation without interfering with users who 
prefer the core command line.
- Implement an action log table and multi-level undo