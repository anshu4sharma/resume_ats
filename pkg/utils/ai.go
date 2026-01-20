package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/anshu4sharma/resume_ats/internal/structs"
	"github.com/openai/openai-go/v3"
)

func AnalyzeResumeWithAi(resumeText string) (*structs.ResumeStruct, error) {
	fmt.Println("before trimming resume text", len(resumeText))
	resumeText = trimForLLM(resumeText, 8000)
	fmt.Println("after trimming resume text", len(resumeText))

	prompt := buildAiPrompt(resumeText)

	client := openai.NewClient()
	chatCompletion, err := client.Chat.Completions.New(
		context.TODO(),
		openai.ChatCompletionNewParams{
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage(
					"You are a deterministic ATS semantic analysis engine. " +
						"You must output valid JSON only and follow instructions exactly.",
				),
				openai.UserMessage(prompt),
			},
			Model:       openai.ChatModelGPT4o,
			Temperature: openai.Float(0),
		},
	)

	if err != nil {
		return &structs.ResumeStruct{}, err
	}

	content := chatCompletion.Choices[0].Message.Content
	fmt.Println(content, "AI Result")

	result, err := parseAiRestoJSON(content)
	return &result, err
}

func parseAiRestoJSON(raw string) (structs.ResumeStruct, error) {
	clean := sanitizeJSON(raw)

	var result structs.ResumeStruct
	err := json.Unmarshal([]byte(clean), &result)

	if err != nil {
		return structs.ResumeStruct{}, err
	}

	return result, nil
}

func buildAiPrompt(resume string) string {
	return fmt.Sprintf(`
TASK:
Analyze the resume text and return a structured ATS signal report.

RULES:
- Output ONE valid JSON object only
- No explanations, no markdown, no extra text
- Do not omit any fields
- If a signal cannot be confidently inferred, return false

OUTPUT SCHEMA (must match exactly):

{
  "has_profile_summary": false,
  "has_experience": false,
  "has_skills": false,
  "has_education": false,
  "has_projects_or_certs": false,
  "has_achievements": false,
  "has_coding_platforms": false,
  "has_coding_contests": false,
  "has_languages": false,
  "good_formatting": false,
  "top_institute": false,
  "missing_profile_summary": false,
  "no_experience": false,
  "missing_proof_of_work": false,
  "missing_keywords": false,
  "missing_education": false,
  "passive_language": false,
  "experience_gap": false,
  "missing_social_presence": false,
  "grammar_issues": false,
  "ai_language_detected": false,
  "complex_sentences": false,
  "repeated_action_verbs": false,
  "personal_details_present": false,
  "open_university": false
}

RESUME TEXT:
<<<
%s
>>>
`, resume)
}

func sanitizeJSON(raw string) string {
	clean := strings.TrimSpace(raw)
	clean = strings.TrimPrefix(clean, "```json")
	clean = strings.TrimPrefix(clean, "```")
	clean = strings.TrimSuffix(clean, "```")
	clean = strings.TrimSpace(clean)

	if idx := strings.LastIndex(clean, "}"); idx != -1 {
		clean = clean[:idx+1]
	}
	return clean
}

func trimForLLM(text string, maxChars int) string {
	if len(text) <= maxChars {
		return text
	}
	return text[:maxChars]
}
