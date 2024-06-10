package cmd

import (
	"fmt"

	"github.com/joaocgduarte/httpmate/internal/configs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// createCmd represents the createRequest command
var createCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"c"},
	Short:   "Creates a new HTTP request to a collection",
	Long: `Creates a new HTTP request to your collection. You will be prompted
to select all the options.`,
	Run: func(cmd *cobra.Command, args []string) {
		collectionsPath := viper.GetString("collectionDirectory")
		requestConfigs := configs.PromptRequestConfig(collectionsPath)
		requestConfigs.WriteToJSONFile()

		fmt.Println("Request was added to collection")
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createRequestCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createRequestCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
