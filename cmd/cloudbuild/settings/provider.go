package settings

import "os/user"

type CliSettings struct {
	ApiKey string `toml:"apikey"`
	OrgId  string `toml:"orgid"`
}

func ParseSettings() (*CliSettings, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}

	usd := usr.HomeDir
}
