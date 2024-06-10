package cmd

import (
	"os"
	"path/filepath"

	"github.com/joaocgduarte/httpmate/internal/configs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	projectName = "httpmate"
)

// rootCmd represents the base command when called without any subcommands
var (
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "httpmate",
		Short: "Manages your collection of HTTP requests",
		Long: `HTTPMate

httpmate is a command-line tool for managing and executing HTTP requests. It 
supports all standard HTTP methods (GET, POST, PUT, DELETE, etc.), allowing you 
to customize headers, query parameters, and request bodies. `,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		//Run: func(cmd *cobra.Command, args []string) {},
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/httpmate/config.yaml)")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		loadDefaultConfigs()
	}

	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	cobra.CheckErr(err)
}

func loadDefaultConfigs() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	configPath := filepath.Join(home, ".config", projectName)
	err = os.MkdirAll(configPath, 0755)
	cobra.CheckErr(err)

	configFilePath := filepath.Join(configPath, "config.yaml")
	collectionDirectoryPath := filepath.Join(configPath, "collections")
	tmpFilesPath := filepath.Join(configPath, "tmp")
	configs.CreateDefaultConfigs(configFilePath, collectionDirectoryPath, tmpFilesPath)

	viper.AddConfigPath(configPath)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.ReadInConfig()
}
