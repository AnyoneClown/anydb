# AnyDB

`anydb` is a CLI tool for managing your databases.

It allows you to `view` table contents, `configure` database connections, and perform `backups`.
## Features

- Display tables and their contents
- Configure database connections
- Backup your database(Currently in progress...)
## Installation

### Using Golang

You can install `anydb` directly using Go:

```sh
go install github.com/AnyoneClown/anydb@latest
```

### Using Linux

1. Download the latest release from the [GitHub Releases](https://github.com/AnyoneClown/anydb/releases) page.

2. Extract the downloaded archive:

    ```sh
    tar -xzf anydb-linux-amd64.tar.gz
    ```

3. Move the binary to `/usr/local/bin`:

    ```sh
    sudo mv anydb /usr/local/bin/
    ```

4. Ensure `/usr/local/bin` is in your PATH. Add the following line to your `~/.bashrc` or `~/.zshrc`:

    ```sh
    export PATH=$PATH:/usr/local/bin
    ```

5. Reload your shell configuration:

    ```sh
    source ~/.bashrc  # or source ~/.zshrc
    ```

### Windows installation

1. Download the latest release from the [GitHub Releases](https://github.com/AnyoneClown/anydb/releases) page.

2. Extract the downloaded archive.

3. Move the `anydb.exe` file to a directory of your choice, for example, `C:\Program Files\anydb`.

4. Add the directory to your PATH:
    - Open the Start Search, type in "env", and select "Edit the system environment variables".
    - In the System Properties window, click on the "Environment Variables" button.
    - In the Environment Variables window, find the `Path` variable in the "System variables" section, and click "Edit".
    - Click "New" and add the path to the directory where you placed `anydb.exe` (e.g., `C:\Program Files\anydb`).
    - Click "OK" to close all the windows.

5. Open a new command prompt and type `anydb` to verify the installation.
## Roadmap

### Version 1.0

- [x] Basic CLI structure
- [x] View table contents
- [x] Configure database connections
- [ ] Backup functionality (In progress)

### Future Plans

- [ ] Add support for more database drivers
- [ ] Implement restore functionality
- [ ] Improve error handling and logging
- [ ] Add more configuration options

## Acknowledgments

- [Cobra](https://github.com/spf13/cobra) for the CLI framework
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) for the TUI framework
- [sqlx](https://github.com/jmoiron/sqlx) for SQL extensions