package utils

import (
	"fmt"
	"net/mail"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/anshu4sharma/resume_ats/internal/structs"
)

const (
	BaselineScore = 40
	MaxScore      = 100
	MinScore      = 1
	MaxPenalty    = 40
)

func CalculateResumeScore(r structs.ResumeStruct) int {
	score := BaselineScore
	penalty := 0

	// --------------------
	// Positive Scoring
	// --------------------
	if r.HasProfileSummary {
		score += 8
	}
	if r.HasExperience {
		score += 10
	}
	if r.HasSkills {
		score += 8
	}
	if r.HasEducation {
		score += 6
	}
	if r.HasEmail && r.HasPhone {
		score += 4
	}
	if r.HasLinkedIn {
		score += 4
	}
	if r.HasGitHubOrPortfolio {
		score += 8
	}
	if r.HasProjectsOrCerts {
		score += 6
	}
	if r.HasAchievements {
		score += 4
	}
	if r.HasCodingPlatforms {
		score += 3
	}
	if r.HasCodingContests {
		score += 2
	}
	if r.HasLanguages {
		score += 2
	}
	if r.GoodFormatting {
		score += 3
	}
	if r.TopInstitute {
		score += 2
	}
	if r.HasProofOfWork {
		score += 8
	}

	// --------------------
	// Negative Penalties
	// --------------------
	if r.MissingProfileSummary {
		penalty += 8
	}
	if r.NoExperience {
		penalty += 10
	}
	if r.MissingProofOfWork {
		penalty += 8
	}
	if r.MissingKeywords {
		penalty += 8
	}
	if r.MissingEducation {
		penalty += 6
	}
	if r.PassiveLanguage {
		penalty += 4
	}
	if r.ExperienceGap {
		penalty += 5
	}
	if r.MissingSocialPresence {
		penalty += 5
	}
	if r.GrammarIssues {
		penalty += 4
	}
	if r.AILanguageDetected {
		penalty += 4
	}
	if r.ComplexSentences {
		penalty += 3
	}
	if r.RepeatedActionVerbs {
		penalty += 2
	}
	if r.MultipleFonts {
		penalty += 2
	}
	if r.MultiColumnLayout {
		penalty += 2
	}
	if r.MoreThanTwoPages {
		penalty += 3
	}
	if r.LargeFileSize {
		penalty += 2
	}
	if r.PersonalDetailsPresent {
		penalty += 3
	}
	// if r.OpenUniversity {
	// 	penalty += 2
	// }

	// --------------------
	// Penalty Normalization
	// --------------------
	if penalty > MaxPenalty {
		penalty = MaxPenalty
	}

	score -= penalty

	// --------------------
	// Clamp Final Score
	// --------------------
	if score > MaxScore {
		return MaxScore
	}
	if score < MinScore {
		return MinScore
	}

	return score
}

var (
	EmailRegex = regexp.MustCompile(`(?i)\b[a-z0-9._%+\-]+@[a-z0-9\-]+(\.[a-z0-9\-]+)*\.[a-z]{2,}\b`)
	PhoneRegex = regexp.MustCompile(`(\+?\d{1,3})?[\s\-]?\(?\d{2,4}\)?[\s\-]?\d{3,4}[\s\-]?\d{4}`)
)

func DetectProfileSummary(text string) bool {
	headers := []string{
		"summary",
		"profile",
		"professional summary",
		"about me",
	}

	for _, h := range headers {
		if strings.Contains(text, h) {
			return true
		}
	}

	// fallback: first paragraph heuristic
	paras := strings.Split(text, "\n\n")
	if len(paras) > 0 && len(strings.Fields(paras[0])) >= 30 {
		return true
	}

	return false
}
func DetectSkills(text string) bool {
	keywords := []string{
		"skills",
		"technical skills",
		"technologies",
		"tools",
		"frameworks",
	}

	for _, k := range keywords {
		if strings.Contains(text, k) {
			return true
		}
	}
	return false
}
func DetectEducation(text string) bool {
	keywords := []string{
		"education",
		"bachelor",
		"master",
		"b.tech",
		"m.tech",
		"degree",
		"university",
		"college",
	}

	for _, k := range keywords {
		if strings.Contains(text, k) {
			return true
		}
	}
	return false
}

func DetectMissingProfileSummary(text string) bool {
	return !DetectProfileSummary(text)
}

func DetectProofOfWork(text string) bool {
	signals := []string{
		"github.com",
		"gitlab.com",
		"bitbucket.org",
		"live demo",
		"deployed at",
		"vercel.app",
		"netlify.app",
	}
	text = NormalizeText(text)
	for _, s := range signals {
		if strings.Contains(text, s) {
			return true
		}
	}
	return false
}
func DetectMissingKeywords(text string) bool {
	coreKeywords := []string{
		"api",
		"backend",
		"frontend",
		"database",
		"microservices",
		"cloud",
		"testing",
	}

	count := 0
	for _, k := range coreKeywords {
		if strings.Contains(text, k) {
			count++
		}
	}

	return count < len(coreKeywords)/2
}

func DetectAchievements(text string) bool {
	keywords := []string{
		"achievement",
		"award",
		"recognition",
		"winner",
		"rank",
		"honor",
	}

	for _, k := range keywords {
		if strings.Contains(text, k) {
			return true
		}
	}
	return false
}

func DetectCodingPlatforms(text string) bool {
	platforms := []string{
		"leetcode",
		"codeforces",
		"codechef",
		"hackerrank",
		"atcoder",
	}

	for _, p := range platforms {
		if strings.Contains(text, p) {
			return true
		}
	}
	return false
}
func DetectCodingContests(text string) bool {
	keywords := []string{
		"contest",
		"competition",
		"hackathon",
		"challenge",
	}

	for _, k := range keywords {
		if strings.Contains(text, k) {
			return true
		}
	}
	return false
}

func DetectLanguages(text string) bool {
	keywords := []string{
		"languages",
		"english",
		"hindi",
		"spanish",
		"french",
	}

	for _, k := range keywords {
		if strings.Contains(text, k) {
			return true
		}
	}
	return false
}

func DetectMultipleFonts(text string) bool {
	return strings.Count(text, "font") > 2
}

func DetectMultiColumnLayout(text string) bool {
	lines := strings.Split(text, "\n")
	tabHeavy := 0

	for _, l := range lines {
		if strings.Count(l, "  ") >= 5 {
			tabHeavy++
		}
	}

	return tabHeavy > len(lines)/5
}

func DetectPersonalDetails(text string) bool {
	patterns := []string{
		"date of birth",
		"dob",
		"age",
		"gender",
		"marital status",
		"father name",
		"mother name",
		"pet name",
		"cast",
		"religion",
	}

	for _, p := range patterns {
		if strings.Contains(text, p) {
			return true
		}
	}
	return false
}

func DetectNoExperience(text string) bool {
	return !DetectExperience(text)
}

func DetectExperience(text string) bool {
	keywords := []string{
		"experience",
		"work experience",
		// "employment",
		// "intern",
		// "engineer",
		// "developer",
		// "manager",
		// "qa",
	}

	for _, k := range keywords {
		if strings.Contains(text, k) {
			return true
		}
	}
	return false
}

func DetectExperienceGap(text string) bool {
	yearRegex := regexp.MustCompile(`(19|20)\d{2}`)
	yearsStr := yearRegex.FindAllString(text, -1)

	if len(yearsStr) < 2 {
		return false
	}

	years := make([]int, 0, len(yearsStr))
	for _, y := range yearsStr {
		if v, err := strconv.Atoi(y); err == nil {
			years = append(years, v)
		}
	}

	sort.Ints(years)

	for i := 1; i < len(years); i++ {
		if years[i]-years[i-1] >= 3 {
			return true
		}
	}
	return false
}

func DetectComplexSentences(text string) bool {
	sentences := regexp.MustCompile(`[.!?]`).Split(text, -1)

	var totalWords int
	var complexCount int

	for _, s := range sentences {
		words := strings.Fields(s)
		if len(words) == 0 {
			continue
		}
		totalWords += len(words)

		if len(words) > 25 || strings.Count(s, ",") >= 2 {
			complexCount++
		}
	}

	if len(sentences) == 0 {
		return false
	}

	return complexCount > len(sentences)/3
}

func DetectRepeatedActionVerbs(text string) bool {
	lines := strings.Split(text, "\n")
	verbCount := map[string]int{}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		words := strings.Fields(line)
		if len(words) == 0 {
			continue
		}

		verb := (words[0])
		verbCount[verb]++
	}

	for _, c := range verbCount {
		if c >= 4 {
			return true
		}
	}
	return false
}

func DetectGoodFormatting(text string) bool {
	lines := strings.Split(text, "\n")
	var bullets int
	var longLines int

	for _, l := range lines {
		l = strings.TrimSpace(l)

		if strings.HasPrefix(l, "-") || strings.HasPrefix(l, "•") {
			bullets++
		}
		if len(l) > 120 {
			longLines++
		}
	}

	return bullets >= 5 && longLines == 0
}

func DetectPassiveLanguage(text string) bool {
	passivePhrases := []string{
		"was responsible for",
		"was involved in",
		"was tasked with",
		"worked on",
	}

	for _, p := range passivePhrases {
		if strings.Contains(text, p) {
			return true
		}
	}
	return false
}

func DetectProjectsOrCerts(text string) bool {
	keywords := []string{
		"project",
		"certification",
		"certificate",
		"codedamn",
		"coursera",
		"udemy",
		"aws",
	}

	for _, k := range keywords {
		if strings.Contains(text, k) {
			return true
		}
	}
	return false
}

func DetectOpenUniversity(text string) bool {
	return strings.Contains(text, "open university") ||
		strings.Contains(text, "distance learning") ||
		strings.Contains(text, "online university") ||
		strings.Contains(text, "ignou")
}

func DetectGitHubOrPortfolio(text string) bool {
	return strings.Contains(text, "github.com") ||
		strings.Contains(text, "portfolio")
}

func DetectEmailPresent(text string) bool {
	normalized := NormalizeForEmail(text)

	emails := EmailRegex.FindAllString(normalized, -1)
	if len(emails) == 0 {
		return false
	}

	for _, e := range emails {
		if _, err := mail.ParseAddress(e); err == nil {
			return true
		}
	}
	return false
}

func NormalizeForEmail(text string) string {
	s := strings.ToLower(text)
	s = regexp.MustCompile(`\s+`).ReplaceAllString(s, " ")
	s = regexp.MustCompile(`\s*@\s*`).ReplaceAllString(s, "@")
	s = regexp.MustCompile(`\s*\.\s*`).ReplaceAllString(s, ".")
	return s
}

func IsValidResume(text string) bool {
	fmt.Println(!DetectEmailPresent(text), !DetectSectionHeaders(text), !DetectMinimumWordCount(text), "IsValidResume")
	if !DetectEmailPresent(text) {
		return false
	}
	if !DetectMinimumWordCount(text) {
		return false
	}
	if !DetectSectionHeaders(text) {
		return false
	}

	return true
}

func NormalizeText(text string) string {
	replacements := []string{
		"•", " ",
		",", " ",
		".", " ",
		"—", " ",
		"-", " ",
		"+", " ",
		"/", " ",
		"(", " ",
		")", " ",
		"\n", " ",
		"\t", " ",
	}

	r := strings.NewReplacer(replacements...)
	normalized := r.Replace(text)

	normalized = strings.Join(strings.Fields(normalized), " ")

	return normalized
}

func DetectMinimumWordCount(text string) bool {
	normalized := NormalizeText(text)
	wc := len(strings.Fields(normalized))

	if len(text) > 500 && wc < 30 {
		return false
	}

	return wc >= 80
}

func DetectSectionHeaders(text string) bool {
	headers := []string{
		"experience", "education", "skills", "projects",
		"certifications", "summary", "profile", "objective",
		"competencies", "achievements",
	}

	count := 0
	for _, h := range headers {
		if strings.Contains(text, h) {
			count++
		}
	}
	return count >= 2
}

func DetectDatesPresent(text string) bool {
	yearRegex := regexp.MustCompile(`(19|20)\d{2}`)
	return len(yearRegex.FindAllString(text, -1)) >= 2
}

func DetectEssayStyle(text string) bool {
	lines := strings.Split(text, "\n")
	longParas := 0

	for _, l := range lines {
		if len(strings.Fields(l)) > 40 {
			longParas++
		}
	}
	return longParas > len(lines)/3
}
