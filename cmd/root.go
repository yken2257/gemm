package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gemm",
	Short: "Encode and decode MIME headers",
	Long: `Encode and decode MIME headers. Simply run the command and follow the prompts.
You can also run the command with arguments. For example:
	gemm decode '=?ISO-2022-JP?B?GyRCJDMkcyRLJEEkTxsoQg==?='
	gemm decode -f test.eml
	gemm encode こんにちは`,
	Version: "0.2.0",
	RunE: func(cmd *cobra.Command, args []string) error {
		return selectFuncPrompt(false)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(EncodeCmd())
	rootCmd.AddCommand(DecodeCmd())
}
