package settings

import (
	"github.com/BurntSushi/toml"
	"os"
	"os/user"
	"path"
)

const dotFileName string = ".cloudbuild"

type CliSettings struct {
	ApiKey string `toml:"apiKey"`
	OrgId  string `toml:"orgId"`
}

func ParseDotFile() (*CliSettings, error) {
	dotPath, err := getFilePath()
	if err != nil {
		return nil, err
	}

	f, err := os.Open(dotPath)
	if os.IsNotExist(err) {
		if err := createDotFile(dotPath); err != nil {
			return nil, err
		}
		return &CliSettings{}, nil
	} else if err != nil {
		return nil, err
	}
	defer f.Close()

	var data CliSettings
	if _, err := toml.DecodeReader(f, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

func createDotFile(dotPath string) error {
	f, err := os.Create(dotPath)
	if err != nil {
		return err
	}
	defer f.Close()

	data := CliSettings{}
	if err := toml.NewEncoder(f).Encode(data); err != nil {
		return err
	}
	return nil
}

func getFilePath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return path.Join(usr.HomeDir, dotFileName), nil
}
