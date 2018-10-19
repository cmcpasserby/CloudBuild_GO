package settings

import (
	"github.com/BurntSushi/toml"
	"os"
	"os/user"
	"path"
)

type CliSettings struct {
	ApiKey string `toml:"apikey"`
	OrgId  string `toml:"orgid"`
}

func ParseSettings() (*CliSettings, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}

	settingsFilePath := path.Join(usr.HomeDir, ".cloudbuild")

	f, err := os.Open(settingsFilePath)
	if err != nil {
		return nil, err
	}

	var data CliSettings

	meta, err := toml.DecodeReader(f, &data)
	if err != nil {
		return nil, err
	}
}
