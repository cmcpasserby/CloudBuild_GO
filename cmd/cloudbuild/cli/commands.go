package cli

import (
	"flag"
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

var CommandOrder = [...]string{"listCreds"}

var Commands = map[string]Command{
	"listCreds": {
		"listCreds",
		"List all IOS Credentials",
		func() *flag.FlagSet {
			flags := CreateFlagSet("listCreds")
			flags.String("projectId", "", "Project Id")
			return flags
		}(),
		func(flags map[string]string) error {
			return nil
		},
	},
}
