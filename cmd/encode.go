package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/manifoldco/promptui"
	"github.com/yken2257/gemm/utils"
)

var charset string
var encoding string

// encodeCmd represents the encode command
var encodeCmd = &cobra.Command{
	Use:   "encode",
	Aliases: []string{"enc", "e"},
	Short: "Encode a string",
	Long: `Encode a string. Simply run the command and follow the prompts.
You can also run the command with arguments. For example:
gemm encode "こんにちは"
With flags:
gemm encode "こんにちは" -c ISO-2022-JP -e Q`,
	Version: rootCmd.Version,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.MaximumNArgs(1) (cmd, args); err != nil {
			return err
		}
		// charset must be either UTF-8, ISO-2022-JP, or Shift_JIS
		if charset != "" && charset != "UTF-8" && charset != "ISO-2022-JP" && charset != "Shift_JIS" {
			return fmt.Errorf("charset must be either UTF-8, ISO-2022-JP, or Shift_JIS")
		}
		// encoding must be either B or Q
		if encoding != "" && encoding != "B" && encoding != "Q" {
			return fmt.Errorf("encoding must be either B or Q")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var text string
		if len(args) == 0 {
			text = ""
		} else {
			text = args[0]
		}
		err := encodePrompt(text, charset, encoding)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(encodeCmd)
	encodeCmd.Flags().StringVarP(&charset, "char", "c", "", "charset; UTF-8, ISO-2022-JP, Shift_JIS")
	encodeCmd.Flags().StringVarP(&encoding, "enc", "e", "", "encoding; B, Q")
}

func encodePrompt(text, charset, encoding string) error {
	if text == "" {
		prompt := promptui.Prompt{
			Label: "Enter text to encode",
			HideEntered: false,
		}
		textInput, err := prompt.Run()
		if err != nil {
			return err
		}
		err = encodePrompt(textInput, charset, encoding)
		if err != nil {
			return err
		}
	} else if charset == "" {
		prompt := promptui.Select{
			Label:     "Choose a charset",
			Items:     []string{"UTF-8", "ISO-2022-JP", "Shift_JIS"},
			CursorPos: 0,
			HideSelected: false,
			HideHelp: true,
		}
		_, charsetInput, err := prompt.Run()
		if err != nil {
			return err
		}
		err = encodePrompt(text, charsetInput, encoding)
		if err != nil {
			return err
		}
	} else if encoding == "" {
		prompt := promptui.Select{
			Label:     "Choose an encoding",
			Items:     []string{"B", "Q"},
			CursorPos: 0,
			HideSelected: false,
			HideHelp: true,
		}
		_, encodingInput, err := prompt.Run()
		if err != nil {
			return err
		}
		err = encodePrompt(text, charset, encodingInput)
		if err != nil {
			return err
		}
	} else {
		encoded, err := utils.EncodeHeader(text, charset, encoding)
		if err != nil {
			return err
		}
		fmt.Println(encoded)
	}
	return nil
}
