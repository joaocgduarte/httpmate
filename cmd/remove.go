package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joaocgduarte/httpmate/internal/configs"
	"github.com/joaocgduarte/httpmate/internal/prompts"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Removes a specific request",
	Long: `Removes the configurations of a specific request with this command
command. 

If you provide a collection, using the flag --collection, you would be able
to remove only the request only inside a particullar collection

Example: httpmate remove --collection "collection name"

You can also specify the request which you want to remove directly, using the
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

		confirm := prompts.ConfirmPrompt("Are you sure?")
		if !confirm {
			fmt.Println("Request was not deleted")
			return
		}

		err = os.Remove(filepath.Join(reqConfig.Collection, fmt.Sprintf("%s.json", reqConfig.RequestName)))
		cobra.CheckErr(err)
		fmt.Println("Request was removed successfully")
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)

	removeCmd.Flags().StringP("collection", "c", "", "Collection to which you will be prompted to remove a request")
	removeCmd.Flags().StringP("request", "r", "", "Specify the request which you want to remove, without being prompted")
}
