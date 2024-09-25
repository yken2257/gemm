package cmd

import (
	"os"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/manifoldco/promptui"
	"github.com/ProtonMail/go-mime"
	"github.com/yken2257/gemm/utils"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gemm",
	Short: "Encode and decode MIME headers",
	Long: `Encode and decode MIME headers. Simply run the command and follow the prompts.
You can also run the command with arguments. For example:
gemm decode "=?ISO-2022-JP?B?GyRCJDMkcyRLJEEkTxsoQg==?="
gemm decode -f test.eml
gemm encode "こんにちは"`,
	Version: "0.1.0",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := selectFuncPrompt()
		if err != nil {
			return err
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.go-hello.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}

func selectFuncPrompt() error {
	prompt := promptui.Select{
		Label:     "What do you want to do?",
		Items:     []string{"Decode a single header", "Decode an .eml file", "Encode string"},
		CursorPos: 0,
		HideSelected: false,
		HideHelp: true,
	}

	_, result, err := prompt.Run() 

	if err != nil {
		return err
	}

	switch result {
	case "Decode a single header":
		err := decodeHeaderPrompt()
		if err != nil {
			return err
		}
	case "Decode an .eml file":
		result, err := chooseFilePrompt()
		if err != nil {
			return err
		}
		err = DecodeEmlPrompt(result)
		if err != nil {
			return err
		}
	case "Encode string":
		fmt.Printf("Encode string\n")
	}
	return nil
}

func decodeHeaderPrompt() error {
	prompt := promptui.Prompt{
		Label: "Enter the header text",
		HideEntered: false,
	}

	result, err := prompt.Run()

	if err != nil {
		return err
	}

	decoded, err := gomime.DecodeHeader(result)

	if err != nil {
		return err
	}

	fmt.Println(decoded)
	return nil
}

func chooseFilePrompt() (string, error) {
	prompt := promptui.Prompt{
		Label: "Enter the path to the .eml file",
		HideEntered: false,
	}

	result, err := prompt.Run()

	if err != nil {
		return "", err
	}

	return result, nil
}

func DecodeEmlPrompt(filename string) error {
	decodedHeaders, err := utils.DecodeHeaders(filename)
	if err != nil {
		return err
	}
	var headerKeys []string
	for key := range decodedHeaders {
		headerKeys = append(headerKeys, key)
	}
	prompt := promptui.Select{
		Label:     "Select the header to decode",
		Items:     headerKeys,
		CursorPos: 0,
		HideSelected: false,
		HideHelp: true,
	}
	_, choice, err := prompt.Run()
	if err != nil {
		return err
	}
	decoded := decodedHeaders[choice]
	
	fmt.Println(decoded)
	return nil
}
