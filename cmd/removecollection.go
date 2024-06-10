package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joaocgduarte/httpmate/internal/files"
	"github.com/joaocgduarte/httpmate/internal/prompts"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// removeCollectionCmd represents the removeCollection command
var removeCollectionCmd = &cobra.Command{
	Use: "remove-collection",
	Aliases: []string{
		"rc",
	},
	Short: "Removes a collection and all of it's requests",
	Run: func(cmd *cobra.Command, args []string) {
		collectionDir := viper.GetString("collectionDirectory")

		listItems := files.GetSubDirectories(collectionDir)
		toRemove := prompts.Select("Which collection do you want to remove?", listItems)
		confirm := prompts.ConfirmPrompt("Are you sure?")

		if !confirm {
			fmt.Println("Exiting...")
			return
		}

		err := os.RemoveAll(filepath.Join(collectionDir, toRemove))
		cobra.CheckErr(err)
		fmt.Println("Collection removed successfully")
	},
}

func init() {
	rootCmd.AddCommand(removeCollectionCmd)
}
