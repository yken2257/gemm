package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ProtonMail/go-mime"
)

var filename string

// decodeCmd represents the decode command
var decodeCmd = &cobra.Command{
	Use:   "decode",
	Aliases: []string{"dec", "d"},
	Short: "Decode a MIME header",
	Long: `Decode a MIME header. For example:
gemm decode =?ISO-2022-JP?B?GyRCJDMkcyRLJEEkTxsoQg==?=
You can also decode a header from a file:
gemm decode -f test.eml

If the header separator is included, please enclose it in double quotes:
gemm decode "=?ISO-2022-JP?B?GyRCJDMkcyRLJEEkTxsoQg==?= =?ISO-2022-JP?B?GyRCJDMkcyRLJEEkTxsoQg==?="`,
	Version: rootCmd.Version,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.MaximumNArgs(1) (cmd, args); err != nil {
			return err
		}
		if filename == "" && len(args) == 0 {
			return fmt.Errorf("please specify a file or a header to decode")
		}
		if filename != "" && len(args) == 1 {
			return fmt.Errorf("please specify either a file or a header to decode")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// if a file specified
		if filename != "" {
			err := DecodeEmlPrompt(filename)
			if err != nil {
				return err
			}
		}
		// if a argument specified
		if len(args) == 1 {
			str := args[0]
			decoded, err := gomime.DecodeHeader(str)
			if err != nil {
				return err
			}
			fmt.Println(decoded)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(decodeCmd)
	decodeCmd.Flags().StringVarP(&filename, "file", "f", "", "file to decode")
}
