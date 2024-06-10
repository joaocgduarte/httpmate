package configs

import (
	"os"

	"github.com/joaocgduarte/httpmate/internal/files"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type ApplicationConfigs struct {
	Editor                  string `yaml:"editor"`
	CollectionDirectory     string `yaml:"collectionDirectory"`
	TemporaryFilesDirectory string `yaml:"temporaryFilesDirectory"`
	AlwaysEditBody          bool   `yaml:"alwaysEditBody"`
	AlwaysEditDomain        bool   `yaml:"alwaysEditDomain"`
	AlwaysEditPath          bool   `yaml:"alwaysEditPath"`
	AlwaysEditQueryParams   bool   `yaml:"alwaysEditQueryParams"`
	AlwaysEditHeaders       bool   `yaml:"alwaysEditHeaders"`
	AlwaysEditMethod        bool   `yaml:"alwaysEditMethod"`
	AlwaysEditContentType   bool   `yaml:"alwaysEditContentType"`
	AlwaysEditAll           bool   `yaml:"alwaysEditAll"`
}

func CreateDefaultConfigs(configFilePath, collectionsDirectoryPath, tmpFilesPath string) {
	createConfigFile(configFilePath, collectionsDirectoryPath, tmpFilesPath)
	files.CreateDirectory(collectionsDirectoryPath)
	files.CreateDirectory(tmpFilesPath)
}

func createConfigFile(configFilePath, collectionsDirectoryPath, tmpFilesPath string) {
	if _, err := os.Stat(configFilePath); !os.IsNotExist(err) {
		cobra.CheckErr(err)
		return
	}

	config := ApplicationConfigs{
		Editor:                  "vim",
		CollectionDirectory:     collectionsDirectoryPath,
		TemporaryFilesDirectory: tmpFilesPath,
		AlwaysEditBody:          true,
	}

	yamlData, err := yaml.Marshal(&config)
	cobra.CheckErr(err)

	err = os.WriteFile(configFilePath, yamlData, 0644)
	cobra.CheckErr(err)
}
