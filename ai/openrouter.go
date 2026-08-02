package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const API_URL = "https://openrouter.ai/api/v1/chat/completions"

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type RequestBody struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Choice struct {
	Message Message `json:"message"`
}

type Response struct {
	Choices []Choice `json:"choices"`
}

func CallOpenRouter(apiKey, model, diff string, includeEmoji bool) string {
	instructions := "Generate a highly polished, professional, one-line commit message following GitHub's Conventional Commit specification. " +
		"Use ONLY these allowed types: feat:, fix:, chore:, refactor:, docs:, test:, perf:, ci:, build:, style:. " +
		"The message must:\n" +
		"- Be under 60 characters\n" +
		"- Be crystal clear and specific\n" +
		"- Use active voice and descriptive wording\n" +
		"- Never include trailing punctuation\n" +
		"- Never include quotes or extra commentary\n" +
		"- Precisely summarize the intent of the diff, not its details\n"

	if includeEmoji {
		instructions += "- Prefix the entire message with exactly one emoji that best matches the commit's intent, " +
			"following the gitmoji convention (examples: ✨ feat, 🐛 fix, 🔥 remove code, 📝 docs, ♻️ refactor, " +
			"✅ test, 🎨 style, ⚡️ perf, 👷 ci, 📦 build, 🔧 chore), followed by a space, then the conventional type\n" +
			"- Output format: \"<emoji> <type>: <description>\"\n"
	} else {
		instructions += "- Do not include any emoji\n"
	}

	prompt := fmt.Sprintf("%sHere is the diff:\n%s", instructions, diff)

	body := RequestBody{
		Model: model,
		Messages: []Message{
			{Role: "user", Content: prompt},
		},
	}

	jsonData, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", API_URL, bytes.NewReader(jsonData))
	if err != nil {
		return ""
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)

	var response Response
	if err := json.Unmarshal(data, &response); err != nil {
		return ""
	}

	if len(response.Choices) > 0 {
		return strings.TrimSpace(response.Choices[0].Message.Content)
	}
	return ""
}
