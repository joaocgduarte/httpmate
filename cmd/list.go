package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joaocgduarte/httpmate/internal/files"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "Lists all of the requests available in the system",
	Long: `You can list all of the requests available in the system with this 
command. 

If you provide a collection, using the flag --collection, you would be able
to list only the requests inside a particular collection

Example: httpmate l --collection "collection name"

You can also provide --collection-only flag, which will list you only the 
collections but not any request`,
	Run: func(cmd *cobra.Command, args []string) {
		collectionDir := viper.GetString("collectionDirectory")

		chosenCollection, err := cmd.Flags().GetString("collection")
		cobra.CheckErr(err)

		if chosenCollection != "" {
			collectionDir = filepath.Join(collectionDir, chosenCollection)
		}

		collectionsOnly, err := cmd.Flags().GetBool("collection-only")
		cobra.CheckErr(err)

		var listItems []string
		if !collectionsOnly {
			listItems = files.GetFilesFromDirectoryWithoutExtension(collectionDir)
		} else {
			listItems = files.GetSubDirectories(collectionDir)
		}
		cobra.CheckErr(err)

		if len(listItems) == 0 {
			cobra.CompError("There are no available requests in collection")
			os.Exit(-1)
		}

		for _, item := range listItems {
			fmt.Println(item)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringP("collection", "c", "", "Collection to which you will be prompted to inspect a request")
	listCmd.Flags().BoolP("collection-only", "", false, "Collection to which you will be prompted to inspect a request")
}
