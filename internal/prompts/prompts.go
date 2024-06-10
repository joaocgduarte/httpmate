package prompts

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func SelectWithAdd(label, addLabel string, options []string) string {
	index := -1
	var result string
	var err error

	for index < 0 {
		prompt := promptui.SelectWithAdd{
			Label:    label,
			Items:    options,
			AddLabel: addLabel,
		}

		index, result, err = prompt.Run()
		cobra.CheckErr(err)

		if index == -1 {
			options = append(options, result)
			break
		}
	}

	return result
}

func Prompt(label string) string {
	prompt := promptui.Prompt{
		Label: label,
	}

	result, err := prompt.Run()
	cobra.CheckErr(err)

	return result
}

func PromptWithDefault(label, defaultValue string) string {
	prompt := promptui.Prompt{
		Label:   label,
		Default: defaultValue,
	}

	result, err := prompt.Run()
	cobra.CheckErr(err)

	return result
}

func Select(label string, options []string) string {
	prompt := promptui.Select{
		Label: label,
		Items: options,
		Searcher: func(input string, index int) bool {
			return fuzzy.MatchFold(input, options[index])
		},
	}

	_, result, err := prompt.Run()
	cobra.CheckErr(err)
	return result
}

func PromptWhileConfirm(confirmLabel, promptKeyLabel, promptValueLabel string) map[string]string {
	result := map[string]string{}

	for true {
		confirmRes := ConfirmPrompt(confirmLabel)
		if !confirmRes {
			break
		}

		result[Prompt(promptKeyLabel)] = Prompt(promptValueLabel)
	}

	return result
}

func ConfirmPrompt(confirmLabel string) bool {
	prompt := promptui.Prompt{
		Label: fmt.Sprintf("%s %s", confirmLabel, "(Y/n)"),
	}

	confirmRes, err := prompt.Run()
	cobra.CheckErr(err)
	parsedConfirm := strings.Trim(strings.ToLower(confirmRes), " ")

	return parsedConfirm != "n"
}

func TextEditorPrompt(
	editor,
	filename,
	temporaryFilesDir,
	startingContent string,
) string {
	// Create the file
	filepath := filepath.Join(temporaryFilesDir, filename)
	file, err := os.Create(filepath)
	cobra.CheckErr(err)

	if startingContent != "" {
		_, err = file.WriteString(startingContent)
		cobra.CheckErr(err)
	}
	file.Close()

	// Run editor on the file
	cmd := exec.Command(editor, filepath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	cobra.CheckErr(err)

	content, err := os.ReadFile(filepath)
	cobra.CheckErr(err)
	err = os.Remove(filepath)
	cobra.CheckErr(err)

	return string(content)
}
