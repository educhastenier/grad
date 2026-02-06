# grad

A CLI tool that generates and executes Gradle commands from file or folder paths. Designed for multi-module Gradle projects, it saves you from manually constructing `./gradlew :module:submodule:task --tests "..."` commands.

## Features

- **Smart task detection** - Automatically selects the right Gradle task based on file type:
  - `*IT.java` files run `integrationTest`
  - `*Test.java` files run `test`
  - Directories run `build`
- **Flexible input** - Pass a path as an argument or let it read from your clipboard
- **File search** - Pass just a class name and it finds the file in the current directory tree
- **Auto-execution** - Generates and runs the command in one step (can be disabled)
- **Clipboard integration** - Optionally copies the generated command to clipboard

## Installation

Requires Go 1.24+.

```bash
go install grad@latest
```

Or build from source:

```bash
git clone https://github.com/educhastenier/grad.git
cd grad
go build
```

## Usage

```bash
# Pass a file path directly
grad subscription/project-runtime/src/test/java/com/company/domain/MyServiceTest.java
# => ./gradlew -PcreateTestReports :subscription:project-runtime:test --tests "com.company.domain.MyServiceTest"

# Integration test (auto-detected from *IT.java suffix)
grad src/test/java/com/example/service/MyServiceIT.java
# => ./gradlew -PcreateTestReports :integrationTest --tests "com.example.service.MyServiceIT"

# Build a module
grad community/some-module/sub-project/
# => ./gradlew -PcreateTestReports :some-module:sub-project:build

# Read path from clipboard (no argument)
grad

# Just print the command, don't execute
grad -n some/path/MyTest.java

# Override the task
grad -t customTask some/path/MyTest.java

# Copy generated command to clipboard
grad -c some/path/MyTest.java
```

## Flags

| Flag                  | Short | Description                                 |
|-----------------------|-------|---------------------------------------------|
| `--verbose`           | `-v`  | Enable verbose output with debug info       |
| `--copy-to-clipboard` | `-c`  | Copy the generated command to the clipboard |
| `--no-execute`        | `-n`  | Print the command without executing it      |
| `--task <task>`       | `-t`  | Override the auto-detected Gradle task      |
| `--help`              | `-h`  | Show help                                   |

## Configuration

Create a `config.yaml` in the current directory or in `$HOME/.grad/`:

```yaml
verbose: false
copy-to-clipboard: false
no-execute: false
task: integrationTest
# Shell for command execution (auto-detected from $SHELL if not set)
# shell: /bin/zsh
```

See [config.example.yaml](config.example.yaml) for the full template.

Command-line flags take precedence over configuration file values.

## License

[Apache License 2.0](LICENSE)
