package cli

import (
	"encoding/json"
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

type Command struct {
	Name     string
	HelpText string
	Flags    *flag.FlagSet
	Action   func(flags map[string]string) error
}

var certIdRe = regexp.MustCompile(`[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`)

func populateArgs(flags map[string]string, data interface{}, credsService *cloudbuild.CredentialsService) error {
	v := reflect.Indirect(reflect.ValueOf(data))
	tt := v.Type()
	fCount := v.NumField()

	qs := make([]*survey.Question, 0, fCount)

	hasInteractiveCert := false

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
			} else if fType == "filePath" {
				promptType = &survey.Input{Message: fmt.Sprintf("%s (absoulte path, can drag and drop)", fName)}
			} else if fType == "certId" {
				hasInteractiveCert = true

				creds, err := credsService.GetAllIOS()
				if err != nil {
					return err // maybe fallback on manual text input instead of error
				}

				options := make([]string, 0, len(creds))

				for _, cred := range creds {
					options = append(options, fmt.Sprintf("%s {%s}", cred.Label, cred.Id))
				}

				promptType = &survey.Select{
					Message:  fName,
					Options:  options,
					PageSize: 10,
				}
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

	if hasInteractiveCert {
		for i := 0; i < fCount; i++ {
			fType := tt.Field(i).Tag.Get("type")
			if fType != "certId" {
				continue
			}

			oldValue := v.Field(i).String()
			v.Field(i).SetString(certIdRe.FindString(oldValue))
		}
	}

	return nil
}

var CommandOrder = [...]string{"getCred", "listCreds", "updateCred", "uploadCred", "deleteCred", "listProjects", "config"}

func prettyPrint(data interface{}) {
	if s, err := json.MarshalIndent(data, "", "    "); err == nil {
		fmt.Println(string(s))
		return
	}

	fmt.Printf("%+v\n", data)
}

var Commands = map[string]Command{

	"getCred": {
		"getCred",
		"Get IOS Credential Detials",
		func() *flag.FlagSet {
			flags := CreateFlagSet("getCred")
			flags.String("credId", "", "Credential Id")
			return flags
		}(),
		func(flags map[string]string) error {
			results := struct {
				ApiKey string `survey:"apiKey"`
				OrgId  string `survey:"orgId"`
				CredId string `survey:"credId" type:"credId"`
			}{}

			credsService := cloudbuild.NewCredentialsService(flags["apiKey"], flags["orgId"])

			if err := populateArgs(flags, &results, credsService); err != nil {
				return err
			}

			cred, err := credsService.GetIOS(results.CredId)
			if err != nil {
				return err
			}

			prettyPrint(cred)

			return nil
		},
	},

	"listCreds": {
		"listCreds",
		"List all IOS Credentials",
		func() *flag.FlagSet {
			return CreateFlagSet("listCreds")
		}(),
		func(flags map[string]string) error {
			// parse args and settings, and question if needed
			results := struct {
				ApiKey string `survey:"apiKey"`
				OrgId  string `survey:"orgId"`
			}{}

			if err := populateArgs(flags, &results, nil); err != nil {
				return err
			}

			credsService := cloudbuild.NewCredentialsService(results.ApiKey, results.OrgId)
			creds, err := credsService.GetAllIOS()
			if err != nil {
				return err
			}

			prettyPrint(creds)

			return nil
		},
	},

	"updateCred": {
		"updateCred",
		"Update a IOS Credential",
		func() *flag.FlagSet {
			flags := CreateFlagSet("updateCred")
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
				CertId      string `survey:"certId" type:"certId"`
				Label       string `survey:"label"`
				CertPath    string `survey:"certPath" type:"filePath"`
				ProfilePath string `survey:"profilePath" type:"filePath"`
				CertPass    string `survey:"certPass" type:"password"`
			}{}

			credsService := cloudbuild.NewCredentialsService(flags["apiKey"], flags["orgId"])

			if err := populateArgs(flags, &results, credsService); err != nil {
				return err
			}

			cred, err := credsService.UpdateIOS(results.CertId, results.Label, results.CertPath, results.ProfilePath, results.CertPass)
			if err != nil {
				return err
			}

			prettyPrint(cred)

			return nil
		},
	},

	"uploadCred": {
		"uploadCred",
		"Upload a IOS Credential",
		func() *flag.FlagSet {
			flags := CreateFlagSet("uploadCred")
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
				Label       string `survey:"label"`
				CertPath    string `survey:"certPath" type:"filePath"`
				ProfilePath string `survey:"profilePath" type:"filePath"`
				CertPass    string `survey:"certPass" type:"password"`
			}{}

			credsService := cloudbuild.NewCredentialsService(flags["apiKey"], flags["orgId"])

			if err := populateArgs(flags, &results, credsService); err != nil {
				return err
			}

			cred, err := credsService.UploadIOS(results.Label, results.CertPath, results.ProfilePath, results.CertPass)
			if err != nil {
				return err
			}

			prettyPrint(cred)

			return nil
		},
	},

	"deleteCred": {
		"deleteCred",
		"Delete a IOS Credential",
		func() *flag.FlagSet {
			flags := CreateFlagSet("deleteCred")
			flags.String("credId", "", "Credential Id")
			return flags
		}(),
		func(flags map[string]string) error {
			results := struct {
				ApiKey string `survey:"apiKey"`
				OrgId  string `survey:"orgId"`
				CertId string `survey:"certId" type:"certId"`
			}{}

			credsService := cloudbuild.NewCredentialsService(flags["apiKey"], flags["orgId"])

			if err := populateArgs(flags, &results, credsService); err != nil {
				return err
			}

			resp, err := credsService.DeleteIOS(results.CertId)
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

			if err := populateArgs(flags, &results, nil); err != nil {
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
