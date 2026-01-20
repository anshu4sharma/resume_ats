package structs

type ResumeStruct struct {
	// Positive signals
	HasProfileSummary    bool `json:"has_profile_summary"`
	HasExperience        bool `json:"has_experience"`
	HasSkills            bool `json:"has_skills"`
	HasEducation         bool `json:"has_education"`
	HasEmail             bool `json:"has_email"`
	HasPhone             bool `json:"has_phone"`
	HasLinkedIn          bool `json:"has_linkedin"`
	HasGitHubOrPortfolio bool `json:"has_github_or_portfolio"`
	HasProjectsOrCerts   bool `json:"has_projects_or_certs"`
	HasAchievements      bool `json:"has_achievements"`
	HasCodingPlatforms   bool `json:"has_coding_platforms"`
	HasCodingContests    bool `json:"has_coding_contests"`
	HasLanguages         bool `json:"has_languages"`
	GoodFormatting       bool `json:"good_formatting"`
	TopInstitute         bool `json:"top_institute"`

	// Negative signals
	MissingProfileSummary  bool `json:"missing_profile_summary"`
	NoExperience           bool `json:"no_experience"`
	MissingProofOfWork     bool `json:"missing_proof_of_work"`
	MissingKeywords        bool `json:"missing_keywords"`
	MissingEducation       bool `json:"missing_education"`
	PassiveLanguage        bool `json:"passive_language"`
	ExperienceGap          bool `json:"experience_gap"`
	MissingSocialPresence  bool `json:"missing_social_presence"`
	GrammarIssues          bool `json:"grammar_issues"`
	AILanguageDetected     bool `json:"ai_language_detected"`
	ComplexSentences       bool `json:"complex_sentences"`
	RepeatedActionVerbs    bool `json:"repeated_action_verbs"`
	MultipleFonts          bool `json:"multiple_fonts"`
	MultiColumnLayout      bool `json:"multi_column_layout"`
	MoreThanTwoPages       bool `json:"more_than_two_pages"`
	LargeFileSize          bool `json:"large_file_size"`
	PersonalDetailsPresent bool `json:"personal_details_present"`
	OpenUniversity         bool `json:"open_university"`
	HasProofOfWork         bool `json:"has_proof_of_work"`
}



type ResumeAnalysisResult struct {
	Score int          `json:"score"`
	Data  ResumeStruct `json:"data"`
}
