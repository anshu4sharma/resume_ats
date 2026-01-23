package services

import (
	"fmt"
	"io"
	"strings"

	"github.com/anshu4sharma/resume_ats/internal/structs"
	"github.com/anshu4sharma/resume_ats/pkg/utils"
)

type AtsService interface {
	AnalyzeResume(
		file io.Reader,
		fileName string,
	) (*structs.ResumeAnalysisResult, error)
}
type atsService struct {
	logger *utils.Logger
}

func NewAtsService(logger *utils.Logger) AtsService {
	return &atsService{logger: logger}
}

func (s *atsService) AnalyzeResume(
	file io.Reader,
	fileName string,
) (*structs.ResumeAnalysisResult, error) {

	// read directly using go lib
	// content, totalPage, err := utils.ExtractText(filePath)
	// if err != nil || !utils.IsReadableText(content) {
	// 	s.logger.Warnf("fallback to OCR for %s", filename)
	// 	content, totalPage, err = utils.ExtractTextOCR(filePath)
	// 	if err != nil {
	// 		return nil, errors.New("ocr_failed")
	// 	}
	// }

	extraction, err := utils.ExtractTextFromPdfBox(file, fileName)

	if err != nil {
		return nil, fmt.Errorf("extraction failed: %w", err)
	}

	content := extraction.Text

	resumeText := strings.ToLower(content)

	aiRes := structs.ResumeStruct{
		HasProfileSummary:     utils.DetectProfileSummary(resumeText),
		MissingProfileSummary: !utils.DetectProfileSummary(resumeText),

		HasSkills:     utils.DetectSkills(resumeText),
		HasEducation:  utils.DetectEducation(resumeText),
		HasExperience: utils.DetectExperience(resumeText),
		NoExperience:  !utils.DetectExperience(resumeText),

		HasProjectsOrCerts: utils.DetectProjectsOrCerts(resumeText),
		HasAchievements:    utils.DetectAchievements(resumeText),
		HasProofOfWork:     utils.DetectProofOfWork(resumeText),

		HasCodingPlatforms: utils.DetectCodingPlatforms(resumeText),
		HasCodingContests:  utils.DetectCodingContests(resumeText),
		HasLanguages:       utils.DetectLanguages(resumeText),

		MissingKeywords:     utils.DetectMissingKeywords(resumeText),
		ExperienceGap:       utils.DetectExperienceGap(resumeText),
		ComplexSentences:    utils.DetectComplexSentences(resumeText),
		RepeatedActionVerbs: utils.DetectRepeatedActionVerbs(resumeText),
		PassiveLanguage:     utils.DetectPassiveLanguage(resumeText),

		// GoodFormatting:    utils.DetectGoodFormatting(resumeText),
		MultipleFonts:     extraction.MultipleFonts,
		MultiColumnLayout: extraction.MultiColumnLayout,

		PersonalDetailsPresent: utils.DetectPersonalDetails(resumeText),
		// OpenUniversity:         utils.DetectOpenUniversity(resumeText),

		HasEmail:             utils.EmailRegex.MatchString(resumeText),
		HasPhone:             utils.PhoneRegex.MatchString(resumeText),
		HasLinkedIn:          strings.Contains(resumeText, "linkedin"),
		HasGitHubOrPortfolio: utils.DetectGitHubOrPortfolio(resumeText),

		MoreThanTwoPages: extraction.MoreThanTwoPages,
		// LargeFileSize:    utils.IsLargeFile(fileSize),
		IsValidResume: utils.IsValidResume(resumeText),
	}

	score := utils.CalculateResumeScore(aiRes)

	return &structs.ResumeAnalysisResult{
		Score: score,
		Data:  aiRes,
	}, nil
}
