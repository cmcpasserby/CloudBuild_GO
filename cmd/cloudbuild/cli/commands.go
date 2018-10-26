package cli

import (
	"flag"
	"fmt"
	"github.com/cmcpasserby/CloudBuild_GO/pkg/cloudbuild"
	"gopkg.in/AlecAivazis/survey.v1"
	"reflect"
	"regexp"
)

var reProjectId = regexp.MustCompile(`[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`)
var reApiKey = regexp.MustCompile(`[0-9a-f]{32}`)

type Command struct {
	Name     string
	HelpText string
	Flags    *flag.FlagSet
	Action   func(flags map[string]string) error
}

func PopulateArgs(flags map[string]string, data interface{}) error {
	v := reflect.Indirect(reflect.ValueOf(data))
	tt := v.Type()
	fCount := v.NumField()

	qs := make([]*survey.Question, 0, fCount)

	for i := 0; i < fCount; i++ {
		fName := tt.Field(i).Tag.Get("survey")
		if fName == "" {
			fName = tt.Field(i).Name
		}

		if val, ok := flags[fName]; ok {
			v.Field(i).SetString(val)
		} else {
			qs = append(qs, CreateQuestion(fName, fName))
		}
	}

	if err := survey.Ask(qs, data); err != nil {
		return err
	}

	return nil
}

var CommandOrder = [...]string{"listCreds", "getCred"}

var Commands = map[string]Command{

	"getCred": {
		"getCred",
		"Get IOS Credential Detials",
		func() *flag.FlagSet {
			flags := CreateFlagSet("getCred")
			flags.String("projectId", "", "Project Id")
			flags.String("credId", "", "Credential Id")
			return flags
		}(),
		func(flags map[string]string) error {
			results := struct {
				ApiKey    string `survey:"apiKey"`
				OrgId     string `survey:"orgId"`
				ProjectId string `survey:"projectId"`
				CredId    string `survey:"credId"`
			}{}

			if err := PopulateArgs(flags, &results); err != nil {
				return err
			}

			credsService := cloudbuild.NewCredentialsService(results.ApiKey, results.OrgId)
			cred, err := credsService.GetIOS(results.ProjectId, results.CredId)
			if err != nil {
				return err
			}

			fmt.Printf("%+v\n", cred)

			return nil
		},
	},

	"listCreds": {
		"listCreds",
		"List all IOS Credentials",
		func() *flag.FlagSet {
			flags := CreateFlagSet("listCreds")
			flags.String("projectId", "", "Project Id")
			return flags
		}(),
		func(flags map[string]string) error {
			// parse args and settings, and question if needed
			results := struct {
				ApiKey    string `survey:"apiKey"`
				OrgId     string `survey:"orgId"`
				ProjectId string `survey:"projectId"`
			}{}

			if err := PopulateArgs(flags, &results); err != nil {
				return err
			}

			credsService := cloudbuild.NewCredentialsService(results.ApiKey, results.OrgId)
			creds, err := credsService.GetAllIOS(results.ProjectId)
			if err != nil {
				return err
			}

			fmt.Printf("%+v\n", creds)

			return nil
		},
	},
}
