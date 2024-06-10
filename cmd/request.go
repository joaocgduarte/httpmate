package cmd

import (
	"fmt"
	"net/http"
	"time"

	"github.com/joaocgduarte/httpmate/internal/configs"
	"github.com/joaocgduarte/httpmate/internal/responseprinter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// runCmd represents the request command
var runCmd = &cobra.Command{
	Use:     "run",
	Aliases: []string{"r"},
	Short:   "Performs an HTTP request from your collection",
	Long: `Performs an HTTP request. You will be prompted to choose your request.

You can add the "edit-*" flags to choose whether you want to edit particular
details of the request before you make it. If these flags are true, you will 
be prompted to edit any details you would want. Available "edit" flags are:
    * --edit-body
    * --edit-domain
    * --edit-path
    * --edit-query-params
    * --edit-headers
    * --edit-method
    * --edit-content-type
    * --edit-all

These flags are also available in the default configurations of the application.
If these are set to true, you would always be prompted to edit the parameters
you've chosen.

Example: httpmate r --edit-body`,
	Run: func(cmd *cobra.Command, args []string) {
		reqConfig := configs.PromptNewExistentRequestConfig(
			"What is the request you want to perform?",
			viper.GetString("collectionDirectory"),
		)

		editFlagParser := func(cmdFlag, viperConfig string) bool {
			editFlag, err := cmd.Flags().GetBool(cmdFlag)
			cobra.CheckErr(err)
			return viper.GetBool(viperConfig) || editFlag

		}
		editConfigs := configs.EditRequestFlags{
			EditBody:        editFlagParser("edit-body", "alwaysEditBody"),
			EditDomain:      editFlagParser("edit-domain", "alwaysEditDomain"),
			EditPath:        editFlagParser("edit-path", "alwaysEditPath"),
			EditQueryParams: editFlagParser("edit-query-params", "alwaysEditQueryParams"),
			EditHeaders:     editFlagParser("edit-headers", "alwaysEditHeaders"),
			EditMethod:      editFlagParser("edit-method", "alwaysEditMethod"),
			EditContentType: editFlagParser("edit-content-type", "alwaysEditContentType"),
			EditAll:         editFlagParser("edit-all", "alwaysEditAll"),
		}
		reqConfig.PromptEditConfig(editConfigs)

		printCurl, err := cmd.Flags().GetBool("print-curl")
		if printCurl {
			fmt.Println("cURL equivalent:")
			fmt.Println(reqConfig.ConvertToCurlCommand())
		}

		client := &http.Client{}

		req := reqConfig.BuildHTTPRequest()

		fmt.Println("Request started...")
		startTime := time.Now()
		resp, err := client.Do(req)
		cobra.CheckErr(err)

		responseprinter.PrintHTTPResponse(resp, time.Since(startTime))
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().BoolP("print-curl", "p", false, "Will print a equivalent CURL request")
	runCmd.Flags().BoolP("edit-body", "", false, "If set, you'll be asked to edit the body before making the request")
	runCmd.Flags().BoolP("edit-domain", "", false, "If set, you'll be asked to edit the domain before making the request")
	runCmd.Flags().BoolP("edit-path", "", false, "If set, you'll be asked to edit the path before making the request")
	runCmd.Flags().BoolP("edit-query-params", "", false, "If set, you'll be asked to edit query params before making the request")
	runCmd.Flags().BoolP("edit-headers", "", false, "If set, you'll be asked to edit headers before making the request")
	runCmd.Flags().BoolP("edit-method", "", false, "If set, you'll be asked to edit the method before making the request")
	runCmd.Flags().BoolP("edit-content-type", "", false, "If set, you'll be asked to edit the contentType before making the request")
	runCmd.Flags().BoolP("edit-all", "", false, "If set, you'll be asked to edit all of the configuration before making the request")
}
