# 🌙 Luna — AI Commit Generator

Luna generates concise commit messages using the OpenRouter API, letting you bring your own key and pick any supported model.

## ✨ Features

- **Interactive file picker**: choose which changed files to stage, right from the TUI
- **Per-file commits**: one commit for each selected file
- **OpenRouter-powered**: AI summaries from file diffs, using the model you choose
- **Conventional prefixes**: adds prefix if missing
- **AI-picked emojis**: with `-e`, the model chooses a gitmoji-style emoji that matches the commit's intent
- **Length control**: target < 60 chars, configurable max (default 72)
- **Smart filtering**: ignores common binaries and images

## How it works

1. Lists every changed file via `git status --porcelain` (modified, staged, untracked), minus ignored files
2. You pick which files to include: `↑`/`↓` to move, `space` to select, `a` to select all, `enter` to confirm
3. Stages the selected files (`git add`), then processes them one by one
4. Sends each file's diff to the configured OpenRouter model
5. If `-e` is active, the model prefixes the message with a matching emoji; if the response lacks a known prefix, one is randomly selected from: `chore:`, `refactor:`, `feat:`, `fix:`, `docs:`, `test:`, `etc..`
6. Truncates to `maxCommitLength` and commits with `git commit -m <message> -- <file>`

## Requirements

- Windows
- Git installed and available in PATH
- OpenRouter API key (`https://openrouter.ai/keys`)
- A [Nerd Font](https://www.nerdfonts.com/) set as your terminal font, so the interface icons render correctly

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

- API key and model: Global → Project → Default
- Other settings: Project → Default

Default settings (from code):

- `ignoredPatterns`: `*.exe`, `*.dll`, `*.png`, `*.jpg`, `*.jpeg`, `*.gif`, `*.bin`
- `commitPrefixes`: `chore:`, `refactor:`, `feat:`, `fix:`, `docs:`, `test:`
- `maxCommitLength`: `72`
- `defaultEmoji`: `false`

### Set your API key and model

```bash
$ Luna apikey YOUR_OPENROUTER_KEY
$ Luna model openai/gpt-4o-mini
```

This saves the key and model in your global `.lunarc`. Reopen terminal after setting. Browse available models at [openrouter.ai/models](https://openrouter.ai/models).

## Usage

Run Luna inside a Git repository with any changes (staged, unstaged or untracked) — you'll pick which files to stage from the TUI.

### Commands and aliases

- `help` | `h`: Show help
- `commit` | `c`: Generate and commit per-file messages
- `apikey <YOUR_KEY>` | `k <YOUR_KEY>`: Set API key
- `model <MODEL_ID>` | `m <MODEL_ID>`: Set OpenRouter model
- `config` | `cfg` with subcommands:
  - `init`: Create `.lunacfg` in current directory
  - `show`: Print merged configuration
  - `edit`: Placeholder (not implemented yet)



### Typical flow

```bash
$ Luna commit # or Luna c
```

Use `↑`/`↓` to move, `space` to select files, `a` to select all, `enter` to stage and commit them.

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

- Error: `Set API key using 'luna apikey' first`
  - Run `Luna apikey YOUR_KEY` and reopen terminal
- Error: `Set model using 'luna model' first`
  - Run `Luna model MODEL_ID` (see [openrouter.ai/models](https://openrouter.ai/models))
- Error running Git commands
  - Ensure you're in a Git repository and Git is installed
- No changes found
  - Nothing shows up in `git status`; modify or create a file first
- API key not working
  - Verify key is valid and the selected model is accessible on your OpenRouter account

---

Made with ❤️ by hax — version 1.3 (Beta)
