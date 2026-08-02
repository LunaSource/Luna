package git

import (
	"fmt"
	"math/rand"
	"os/exec"
	"strings"
	"time"
	"unicode"

	"luna/ai"
	"luna/config"
)

// fallbackEmojis is only used when the AI call fails and we can't ask it
// to pick an emoji that matches the commit's intent.
var fallbackEmojis = []string{"✨", "🛠️", "🐛", "🔥", "📝", "🚀", "🔧", "🎨", "🔒", "💄"}

func GenerateCommitMessage(apiKey, diff, filename string, cfg config.Config, includeEmoji bool) string {
	commitMsg := ai.CallOpenRouter(apiKey, cfg.Model, diff, includeEmoji)
	if commitMsg == "" {
		commitMsg = "update " + filename
		if includeEmoji {
			rand.Seed(time.Now().UnixNano())
			commitMsg = fallbackEmojis[rand.Intn(len(fallbackEmojis))] + " " + commitMsg
		}
	}

	checkMsg := stripLeadingEmoji(commitMsg)
	hasPrefix := false
	for _, p := range cfg.CommitPrefixes {
		if strings.HasPrefix(strings.ToLower(checkMsg), strings.ToLower(p)) {
			hasPrefix = true
			break
		}
	}

	if !hasPrefix && len(cfg.CommitPrefixes) > 0 {
		rand.Seed(time.Now().UnixNano())
		prefix := cfg.CommitPrefixes[rand.Intn(len(cfg.CommitPrefixes))]
		commitMsg = insertPrefix(commitMsg, prefix)
	}

	if len(commitMsg) > cfg.MaxCommitLength {
		commitMsg = commitMsg[:cfg.MaxCommitLength-3] + "..."
	}

	return commitMsg
}

// stripLeadingEmoji removes a leading emoji (and any surrounding spaces) so
// the conventional-commit prefix can be detected right after it.
func stripLeadingEmoji(s string) string {
	return strings.TrimLeftFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
}

// insertPrefix adds the missing conventional-commit prefix right after any
// leading emoji, so messages read "<emoji> <type>: <description>".
func insertPrefix(msg, prefix string) string {
	rest := stripLeadingEmoji(msg)
	lead := strings.TrimSuffix(msg, rest)
	return lead + prefix + " " + rest
}

func CommitFile(filename, commitMsg string) (string, error) {
	cmd := exec.Command("git", "commit", "-m", commitMsg, "--", filename)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("error committing %s: %v", filename, err)
	}
	return string(out), nil
}