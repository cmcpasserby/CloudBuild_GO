package cli

import (
	"github.com/cmcpasserby/CloudBuild_GO/cmd/cloudbuild/settings"
	"gopkg.in/AlecAivazis/survey.v1"
)

func CreateQuestions(settings *settings.CliSettings, questions []*survey.Question) []*survey.Question {
	qs := make([]*survey.Question, 0, len(questions)+2)

	if settings.ApiKey == "" {
		qs = append(qs, CreateQuestion("ApiKey", "API Key"))
	}
	if settings.OrgId == "" {
		qs = append(qs, CreateQuestion("OrgId", "Organizaion ID"))
	}

	return append(qs, questions...)
}

func CreateQuestion(name, message string) *survey.Question {
	return &survey.Question{
		Name:      name,
		Prompt:    &survey.Input{Message: message},
		Validate:  survey.Required,
		Transform: survey.Title,
	}
}
