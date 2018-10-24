package cli

import (
	"gopkg.in/AlecAivazis/survey.v1"
)

func CreateQuestion(name, message string) *survey.Question {
	return &survey.Question{
		Name:      name,
		Prompt:    &survey.Input{Message: message},
		Validate:  survey.Required,
		Transform: survey.ToLower,
	}
}
