package cmd

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/joaocgduarte/httpmate/internal/configs"
	"github.com/joaocgduarte/httpmate/internal/responseprinter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// inspectCmd represents the inspect command
var inspectCmd = &cobra.Command{
	Use:     "inspect",
	Aliases: []string{"i"},
	Short:   "Inpects the details of a specific request",
	Long: `Inspects the configurations of a specific request with this command
command. 

If you provide a collection, using the flag --collection, you would be able
to inspect only the requests inside a particullar collection

Example: httpmate i --collection "collection name"

You can also specify the request which you want to inspect directly, using the
--request flag. With this, you won't be prompted to select a request.`,
	Run: func(cmd *cobra.Command, args []string) {
		collectionDir := viper.GetString("collectionDirectory")

		chosenCollection, err := cmd.Flags().GetString("collection")
		cobra.CheckErr(err)

		if chosenCollection != "" {
			collectionDir = filepath.Join(collectionDir, chosenCollection)
		}

		specifiedRequest, err := cmd.Flags().GetString("request")
		cobra.CheckErr(err)

		var reqConfig *configs.RequestConfig
		if specifiedRequest == "" {
			reqConfig = configs.PromptNewExistentRequestConfig(
				"What is the request you want to perform?",
				collectionDir,
			)
		} else {
			reqConfig = configs.NewRequestConfigFromFilePath(filepath.Join(collectionDir, fmt.Sprintf("%s.json", specifiedRequest)))
		}

		requestDetails, err := json.MarshalIndent(reqConfig, "", "    ")
		cobra.CheckErr(err)

		fmt.Println("Details:")
		responseprinter.PrintJSON(requestDetails)
	},
}

func init() {
	rootCmd.AddCommand(inspectCmd)

	inspectCmd.Flags().StringP("collection", "c", "", "Collection to which you will be prompted to inspect a request")
	inspectCmd.Flags().StringP("request", "r", "", "Specify the request which you want to inspect, without being prompted")
}
