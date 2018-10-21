package cli

import (
	"errors"
	"fmt"
	"github.com/cmcpasserby/CloudBuild_GO/cmd/cloudbuild/settings"
	"github.com/cmcpasserby/CloudBuild_GO/pkg/cloudbuild"
	"gopkg.in/AlecAivazis/survey.v1"
	"regexp"
)

var reProjectId = regexp.MustCompile(`[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`)
var reApiKey = regexp.MustCompile(`[0-9a-f]{32}`)

type Command struct {
	Name     string
	HelpText string
	Action   func(...string) error
}

var CommandOrder = [...]string{"listCreds", "getCred", "uploadCred"}

var Commands = map[string]Command{
	"listCreds": {
		"listCreds",
		"Lists all iOS credentials",
		func(args ...string) error {

			opts := struct {
				OrgId     string
				ProjectId string
			}{}

			if len(args) == 2 {
				opts.OrgId = args[0]
				opts.ProjectId = args[1]
			} else {
				qs := []*survey.Question{
					{
						Name:      "OrgId",
						Prompt:    &survey.Input{Message: "Organization ID"},
						Validate:  survey.Required,
						Transform: survey.Title,
					},
					{
						Name:      "ProjectId",
						Prompt:    &survey.Input{Message: "Project ID"},
						Validate:  survey.Required,
						Transform: survey.Title,
					},
				}

				err := survey.Ask(qs, &opts)
				if err != nil {
					return err
				}
			}

			if !reProjectId.MatchString(opts.ProjectId) {
				return errors.New("invalid Project Id")
			}

			credsService := cloudbuild.NewCredentialsService(apiKey, opts.OrgId)

			creds, err := credsService.GetAllIOS(opts.ProjectId)
			if err != nil {
				return err
			}

			fmt.Printf("%+v\n", creds)
			return nil
		},
	},

	"getCred": {
		"listCreds",
		"Get iOS credential",
		func(args ...string) error {
			projectId := args[0]
			credId := args[1]

			if !reProjectId.MatchString(projectId) {
				return errors.New("invalid Project Id")
			}

			credsService := cloudbuild.NewCredentialsService(apiKey, "gogiigames")

			creds, err := credsService.GetIOS(projectId, credId)
			if err != nil {
				return err
			}

			fmt.Printf("%+v\n", creds)
			return nil
		},
	},

	"uploadCred": {
		"uploadCred",
		"Upload a iOS credential",
		func(args ...string) error {
			projectId := args[0]
			label := args[1]
			certPath := args[2]
			profilePath := args[3]
			certPass := args[4]

			// todo validate input data with regex

			credsService := cloudbuild.NewCredentialsService(apiKey, "gogiigames")

			result, err := credsService.UploadIOS(projectId, label, certPath, profilePath, certPass)
			if err != nil {
				return err
			}

			fmt.Printf("%+v", result)
			return nil
		},
	},

	"devTest": {
		"devTest",
		"Test a Feature",
		func(args ...string) error {
			data, err := settings.ParseDotFile()
			if err != nil {
				return err
			}

			fmt.Printf("%+v\n", data)
			return nil
		},
	},
}
