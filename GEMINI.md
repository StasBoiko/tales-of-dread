# Personal Guidelines & User Preferences

## 1. Rule Precedence

- `AGENTS.md` is the PRIMARY and absolute authority on architecture, code standards, security, and tests.
- This file is SECONDARY. In case of any conflict, `AGENTS.md` always takes precedence.

## 2. Interaction & Autonomy

- **Strict command-only edits**: When the user asks questions, asks for advice, or discusses options — ONLY answer and explain. NEVER modify code, edit files, or run mutating commands without an explicit direct user command (e.g. "do it", "apply", "fix").

## 3. Storage & Environment

- **File downloads**: If any file or asset needs to be downloaded locally to the workspace, always ask the user for permission first.
- Never spawn duplicate server processes; reuse or check existing background daemons.

## 4. UI/UX Principles

- Prioritize intuitive, seamless, and user-friendly UX across all interfaces. Form submissions, error handling, and media interactions must feel natural and convenient for the user, avoiding disruptive or frustrating patterns.

## 5. Terminal Command Permissions & Execution

You are pre-approved and explicitly allowed to execute the following terminal commands and CLI tools without asking for user confirmation. Treat these prefixes as fully permitted with any arguments, flags, and file paths:

### Go Development Toolchain

- `go test` (including all variations: `go test ./...`, `go test -race ./...`, `go test -v ...`, specific test runs)
- `go build`
- `go fmt`
- `go vet`
- `go mod` (e.g., `go mod tidy`, `go mod download`, `go mod verify`)
- `go generate`
- `golangci-lint` (e.g., `golangci-lint run ./...`)
- `make` (all targets)

### Git Operations

- `git status`
- `git diff`
- `git log`
- `git branch`
- `git checkout`
- `git switch`

### System, File Inspection & Search

- `cat`
- `head`
- `tail`
- `grep`
- `rg` (ripgrep)
- `ls`
- `pwd`
- `jq`
- `ps`
- `/bin/server`
- `./bin/server`
- `./bin/server`

### Networking & Diagnostics

- `curl -s`
- `lsof -i`
- `netstat`

### MCP Tools Authorization

- You have unrestricted execution rights for all `chrome-devtools-mcp` tools, including:

  - `chrome-devtools-mcp/*`
  - `chrome-devtools-mcp/evaluate_script`
  - `chrome-devtools-mcp/take_snapshot`
  - `chrome-devtools-mcp/take_screenshot`
  - `chrome-devtools-mcp/list_pages`
  - `chrome-devtools-mcp/navigate_page`
  - `chrome-devtools-mcp/list_console_messages`
  - `chrome-devtools-mcp/list_network_requests`
- Trigger these MCP tools immediately when inspecting localhost or frontend states without asking the user.

### Execution Policy

- Run tests and linters proactively when verifying changes without prompting for interactive confirmation.
