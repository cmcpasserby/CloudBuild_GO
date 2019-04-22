package cli

import (
	"flag"
	"fmt"
	"github.com/cmcpasserby/CloudBuild_GO/cmd/cloudbuild/settings"
	"github.com/cmcpasserby/CloudBuild_GO/pkg/cloudbuild"
	"gopkg.in/AlecAivazis/survey.v1"
	"os"
	"os/exec"
	"reflect"
	"regexp"
)

var (
	reProjectId = regexp.MustCompile(`[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`)
	reApiKey    = regexp.MustCompile(`[0-9a-f]{32}`)
)

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
		fType := tt.Field(i).Tag.Get("type")
		if fName == "" {
			fName = tt.Field(i).Name
		}

		if val, ok := flags[fName]; ok {
			v.Field(i).SetString(val)
		} else {
			var promptType survey.Prompt

			if fType == "password" {
				promptType = &survey.Password{Message: fName}
			} else {
				promptType = &survey.Input{Message: fName}
			}

			qs = append(qs, &survey.Question{
				Name:     fName,
				Prompt:   promptType,
				Validate: survey.Required,
			})
		}
	}

	if err := survey.Ask(qs, data); err != nil {
		return err
	}

	return nil
}

var CommandOrder = [...]string{"getCred", "listCreds", "updateCred", "uploadCred", "deleteCred", "listProjects", "config"}

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

	"updateCred": {
		"updateCred",
		"Update a IOS Credential",
		func() *flag.FlagSet {
			flags := CreateFlagSet("updateCred")
			flags.String("projectId", "", "Project Id")
			flags.String("certId", "", "Certificate Id")
			flags.String("label", "", "Label")
			flags.String("certPath", "", "Certificate Path")
			flags.String("profilePath", "", "Provisioning Profile Path")
			flags.String("certPass", "", "Certificate password")
			return flags
		}(),
		func(flags map[string]string) error {
			results := struct {
				ApiKey      string `survey:"apiKey"`
				OrgId       string `survey:"orgId"`
				ProjectId   string `survey:"projectId"`
				CertId      string `survey:"certId"`
				Label       string `survey:"label"`
				CertPath    string `survey:"certPath"`
				ProfilePath string `survey:"profilePath"`
				CertPass    string `survey:"certPass" type:"password"`
			}{}

			if err := PopulateArgs(flags, &results); err != nil {
				return err
			}

			credsService := cloudbuild.NewCredentialsService(results.ApiKey, results.OrgId)
			cred, err := credsService.UpdateIOS(results.ProjectId, results.CertId, results.Label, results.CertId, results.ProfilePath, results.CertPass)
			if err != nil {
				return err
			}

			fmt.Printf("%+v\n", cred)

			return nil
		},
	},

	"uploadCred": {
		"uploadCred",
		"Upload a IOS Credential",
		func() *flag.FlagSet {
			flags := CreateFlagSet("uploadCred")
			flags.String("projectId", "", "Project Id")
			flags.String("label", "", "Label")
			flags.String("certPath", "", "Certificate Path")
			flags.String("profilePath", "", "Provisioning Profile Path")
			flags.String("certPass", "", "Certificate password")
			return flags
		}(),
		func(flags map[string]string) error {
			results := struct {
				ApiKey      string `survey:"apiKey"`
				OrgId       string `survey:"orgId"`
				ProjectId   string `survey:"projectId"`
				Label       string `survey:"label"`
				CertPath    string `survey:"certPath"`
				ProfilePath string `survey:"profilePath"`
				CertPass    string `survey:"certPass" type:"password"`
			}{}

			if err := PopulateArgs(flags, &results); err != nil {
				return err
			}

			credsService := cloudbuild.NewCredentialsService(results.ApiKey, results.OrgId)
			cred, err := credsService.UploadIOS(results.ProjectId, results.Label, results.CertPath, results.ProfilePath, results.CertPass)
			if err != nil {
				return err
			}

			fmt.Printf("%+v\n", cred)

			return nil
		},
	},

	"deleteCred": {
		"deleteCred",
		"Delete a IOS Credential",
		func() *flag.FlagSet {
			flags := CreateFlagSet("deleteCred")
			flags.String("projectId", "", "Project Id")
			flags.String("credId", "", "Credential Id")
			return flags
		}(),
		func(flags map[string]string) error {
			results := struct {
				ApiKey    string `survey:"apiKey"`
				OrgId     string `survey:"orgId"`
				ProjectId string `survey:"projectId"`
				CertId    string `survey:"certId"`
			}{}

			if err := PopulateArgs(flags, &results); err != nil {
				return err
			}

			credsService := cloudbuild.NewCredentialsService(results.ApiKey, results.OrgId)
			resp, err := credsService.DeleteIOS(results.ProjectId, results.CertId)
			if err != nil {
				return err
			}

			fmt.Println(resp.Status)

			return nil
		},
	},

	"listProjects": {
		"listProjects",
		"List Projects On CloudBuild",
		func() *flag.FlagSet {
			flags := CreateFlagSet("listProjects")
			return flags
		}(),
		func(flags map[string]string) error {
			results := struct {
				ApiKey string `survey:"apiKey"`
				OrgId  string `survey:"orgId"`
			}{}

			if err := PopulateArgs(flags, &results); err != nil {
				return err
			}

			projectService := cloudbuild.NewProjectsService(results.ApiKey, results.OrgId)
			projects, err := projectService.ListAll()
			if err != nil {
				return err
			}

			for _, proj := range projects {
				fmt.Printf("Name: %s || Id: %s\n", proj.Name, proj.Guid)
			}

			return nil
		},
	},

	"config": {
		"config",
		"Edit config file",
		func() *flag.FlagSet {
			return flag.NewFlagSet("config", flag.ExitOnError)
		}(),
		func(flags map[string]string) error {
			dotFilePath, err := settings.GetFilePath()
			if err != nil {
				return err
			}

			if _, err := os.Stat(dotFilePath); os.IsNotExist(err) {
				if err := settings.CreateDotFile(dotFilePath); err != nil {
					return err
				}
			}

			cmd := exec.Command("vim", dotFilePath)

			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
				return err
			}

			return nil
		},
	},
}
