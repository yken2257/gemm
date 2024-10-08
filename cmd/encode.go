package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yken2257/gemm/utils"
)

func EncodeCmd() *cobra.Command {
	var charset string
	var encoding string

	cmd := &cobra.Command{
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
			charset, _ := cmd.Flags().GetString("char")
			encoding, _ := cmd.Flags().GetString("enc")

			// if stdin is piped, check if both charset and encoding are set
			stat, err := os.Stdin.Stat()
			if err != nil {
				return fmt.Errorf("failed to stat stdin: %v", err)
			}
			isPiped := (stat.Mode() & os.ModeCharDevice) == 0
			if isPiped && (charset == "" || encoding == "") {
				return fmt.Errorf("charset and encoding must be set when using stdin; e.g. gemm encode -c UTF-8 -e B")
			}

			// charset must be either UTF-8, ISO-2022-JP, or Shift_JIS
			if charset != "" {
				normalizedCharset := utils.NormalizeCharset(charset)
				if _, valid := utils.ValidCharsets[normalizedCharset]; !valid {
					return fmt.Errorf("charset must be either UTF-8, ISO-2022-JP, or Shift_JIS")
				}
			}
			// encoding must be either B or Q (case-insensitive)
			if encoding != "" {
				encoding = strings.ToUpper(encoding)
				if encoding != "B" && encoding != "Q" {
					return fmt.Errorf("encoding must be either B or Q")
				}
			}
			if len(args) > 1 {
				return fmt.Errorf("too many arguments; only one arg is allowed")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var text string
			var err error

			// stdin is piped
			stat, err := os.Stdin.Stat()
			if err != nil {
				return fmt.Errorf("failed to stat stdin: %v", err)
			}
			isPiped := (stat.Mode() & os.ModeCharDevice) == 0

			if isPiped {
				data, err := io.ReadAll(os.Stdin)
				if err != nil {
					return fmt.Errorf("failed to read from stdin: %v", err)
				}
				text = string(data)
			} else if len(args) == 1 {
				text = args[0]
			} else if len(args) == 0 {
				text = ""
			}
			return encodePrompt(text, charset, encoding)
		},
	}
	cmd.Flags().StringVarP(&charset, "char", "c", "", "charset; UTF-8, ISO-2022-JP, Shift_JIS")
	cmd.Flags().StringVarP(&encoding, "enc", "e", "", "encoding; B, Q")
	
	return cmd
}
