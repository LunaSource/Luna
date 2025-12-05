# 🌙 Luna — AI Commit Generator

Luna generates concise commit messages for your using Google Gemini 2.0 Flash API.

## ✨ Features

- **Per-file commits**: one commit for each staged file
- **Gemini 2.0 Flash**: AI-powered summaries from file diffs
- **Conventional prefixes**: adds prefix if missing
- **Length control**: target < 60 chars, configurable max (default 72)
- **Smart filtering**: ignores common binaries and images
- **Optional emojis**: enable with `-e`

## How it works

1. Collects staged files via `git diff --cached --name-only`
2. Sends each diff to Gemini
3. If response lacks known prefix, randomly selects from: `chore:`, `refactor:`, `feat:`, `fix:`, `docs:`, `test:`, `etc..`
4. If `-e` is active, adds random emoji
5. Truncates to `maxCommitLength` and commits with `git commit -m <message> -- <file>`

## Requirements

- Windows
- Git installed and available in PATH
- Google Gemini API key (`https://aistudio.google.com/app/apikey`)

## Installation

### Option A — Use prebuilt binary (`bin/Luna.exe`)

1. Copy `bin/Luna.exe` to a directory, e.g. `C:\Users\user\Luna`
2. Add that folder to system PATH:
   - Press `Win + R`, run `sysdm.cpl`, open "Environment Variables"
   - Edit `Path` variable → "New" → paste folder path
   - Save and reopen terminal

### Option B — Build from source (Go)

```bash
$ go build -o ./bin/Luna.exe main.go
```

Or use the helper script:

```bash
$ ./build.sh
```

## Configuration

Luna reads configuration from project and global files:

- Project: `.lunacfg` (in repository root or nearest parent)
- Global: `.lunarc` (in user home directory)

Priority:

- API key: Global → Project → Default
- Other settings: Project → Default

Default settings (from code):

- `ignoredPatterns`: `*.exe`, `*.dll`, `*.png`, `*.jpg`, `*.jpeg`, `*.gif`, `*.bin`
- `commitPrefixes`: `chore:`, `refactor:`, `feat:`, `fix:`, `docs:`, `test:`
- `maxCommitLength`: `72`
- `defaultEmoji`: `false`

### Set your API key

```bash
$ Luna apikey YOUR_GEMINI_KEY
```

This saves the key in your global `.lunarc`. Reopen terminal after setting.

## Usage

Run Luna inside a Git repository with staged changes.

### Commands and aliases

- `help` | `h`: Show help
- `commit` | `c`: Generate and commit per-file messages
- `apikey <YOUR_KEY>` | `k <YOUR_KEY>`: Set API key
- `config` | `cfg` with subcommands:
  - `init`: Create `.lunacfg` in current directory
  - `show`: Print merged configuration
  - `edit`: Placeholder (not implemented yet)



### Typical flow

```bash
$ git add . && Luna commit # or Luna c
```

### Optional emojis

```bash
$ Luna c -e       # enable emojis in messages
```

## Example output

```
Generating commit for file: src/main.go
Committed src/main.go with message:
🚀 feat: add user authentication system

Generating commit for file: README.md
Committed README.md with message:
📝 docs: update installation instructions
```

## Notes

- Luna skips common binary/image files
- If model returns empty response, fallback is `update <file>`
- Supported prefixes: `feat:`, `fix:`, `docs:`, `refactor:`, `test:`, `chore:`
- `maxCommitLength` is enforced (default 72)

## Troubleshooting

- Error: `Set API key using LunaApikey first`
  - Run `Luna apikey YOUR_KEY` and reopen terminal
- Error running Git commands
  - Ensure you're in a Git repository and Git is installed
- No staged changes
  - Run `git add .` or stage specific files
- API key not working
  - Verify key is valid and has access to Gemini 2.0 Flash

---

Made with ❤️ by hax — version 1.3 (Beta)
