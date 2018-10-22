package cli

import (
	"flag"
	"github.com/cmcpasserby/CloudBuild_GO/cmd/cloudbuild/settings"
)

func CreateFlagSet(name string) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ExitOnError)
	fs.String("apiKey", "", "Api Key")
	fs.String("orgId", "", "Organization Id")
	return fs
}

func ParseFlags(set *flag.FlagSet, args []string) (map[string]string, error) {
	data, err := settings.ParseDotFile()
	if err != nil {
		return nil, err
	}

	if err := set.Parse(args); err != nil {
		return nil, err
	}

	flagMap := make(map[string]string)

	set.Visit(func(flag *flag.Flag) {
		flagMap[flag.Name] = flag.Value.String()
	})

	// apply from dot settings if not defined as flags
	if _, ok := flagMap["apiKey"]; !ok {
		flagMap["apiKey"] = data.ApiKey
	}

	if _, ok := flagMap["orgId"]; !ok {
		flagMap["orgId"] = data.OrgId
	}

	return flagMap, nil
}
