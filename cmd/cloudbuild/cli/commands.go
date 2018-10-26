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

func PopulateArgs(data interface{}) error {
	v := reflect.Indirect(reflect.ValueOf(data))

	for i := 0; i < v.NumField(); i++ {
		v.Field(i).SetString("Hello")
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

			if err := PopulateArgs(&results); err != nil {
				return err
			}

			fmt.Println(results.ApiKey)
			fmt.Println(results.OrgId)

			qs := make([]*survey.Question, 0, 4)

			// and is valid
			if apiKey, ok := flags["apiKey"]; ok {
				results.ApiKey = apiKey
			} else {
				qs = append(qs, CreateQuestion("apiKey", "API Key"))
			}

			// and is valid
			if orgId, ok := flags["orgId"]; ok {
				results.OrgId = orgId
			} else {
				qs = append(qs, CreateQuestion("orgId", "Organization Id"))
			}

			// and is valid
			if projectId, ok := flags["projectId"]; ok {
				results.ProjectId = projectId
			} else {
				qs = append(qs, CreateQuestion("projectId", "Project Id"))
			}

			// and is valid
			if credId, ok := flags["credId"]; ok {
				results.CredId = credId
			} else {
				qs = append(qs, CreateQuestion("credId", "Credential Id"))
			}

			if err := survey.Ask(qs, &results); err != nil {
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

			qs := make([]*survey.Question, 0, 3)

			// and is valid
			if apiKey, ok := flags["apiKey"]; ok {
				results.ApiKey = apiKey
			} else {
				qs = append(qs, CreateQuestion("apiKey", "API Key"))
			}

			// and is valid
			if orgId, ok := flags["orgId"]; ok {
				results.OrgId = orgId
			} else {
				qs = append(qs, CreateQuestion("orgId", "Organization Id"))
			}

			// and is valid
			if projectId, ok := flags["projectId"]; ok {
				results.ProjectId = projectId
			} else {
				qs = append(qs, CreateQuestion("projectId", "Project Id"))
			}

			if err := survey.Ask(qs, &results); err != nil {
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
