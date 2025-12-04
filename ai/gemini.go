package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const API_URL = "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent"

type Part struct {
	Text string `json:"text"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type RequestBody struct {
	Contents []Content `json:"contents"`
}

type Candidate struct {
	Content Content `json:"content"`
}

type Response struct {
	Candidates []Candidate `json:"candidates"`
}

func CallGemini(apiKey, diff string) string {
	body := RequestBody{
		Contents: []Content{
			{
				Parts: []Part{
					{Text: fmt.Sprintf(
						"Generate a highly polished, professional, one-line commit message following GitHub's Conventional Commit specification. "+
							"Use ONLY these allowed types: feat, fix, chore, refactor, docs, test, perf, ci, build, style.\n"+
							"Format the commit exactly as: prefix(scope): message, where scope is optional.\n"+
							"The message must:\n"+
							"- Be under 60 characters\n"+
							"- Be crystal clear and specific\n"+
							"- Use active voice and descriptive wording\n"+
							"- Never include trailing punctuation\n"+
							"- Never include quotes or extra commentary\n"+
							"- Precisely summarize the intent of the diff, not its details\n"+
							"Here is the diff:\n%s",
						diff,
					)},
				},
			},
		},
	}

	jsonData, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", API_URL+"?key="+apiKey, bytes.NewReader(jsonData))
	if err != nil {
		return ""
	}
	req.Header.Set("Content-Type", "application/json")

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

	if len(response.Candidates) > 0 && len(response.Candidates[0].Content.Parts) > 0 {
		return strings.TrimSpace(response.Candidates[0].Content.Parts[0].Text)
	}
	return ""
}
