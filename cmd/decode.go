package cmd

import (
	"fmt"
	"io"
	"os"

	gomime "github.com/ProtonMail/go-mime"
	"github.com/spf13/cobra"
)

func DecodeCmd() *cobra.Command {
	var filename string

	cmd := &cobra.Command{
		Use:     "decode",
		Aliases: []string{"dec", "d"},
		Short:   "Decode a MIME header",
		Long: `Decode a MIME header. For example:
	gemm decode '=?ISO-2022-JP?B?GyRCJDMkcyRLJEEkTxsoQg==?='
You can also decode a header from a file:
	gemm decode -f test.eml

Please enclose the header with single quotes to prevent unexpected behavior:
	gemm decode '=?ISO-2022-JP?B?GyRCJDMkcyRLJEEkTxsoQg==?= =?ISO-2022-JP?B?GyRCJDMkcyRLJEEkTxsoQg==?='`,
		Version: rootCmd.Version,
		Args: func(cmd *cobra.Command, args []string) error {
			filename, _ := cmd.Flags().GetString("file")

			// 標準入力がパイプされているかどうかをチェック
			stat, err := os.Stdin.Stat()
			if err != nil {
				return fmt.Errorf("failed to stat stdin: %v", err)
			}
			isPiped := (stat.Mode() & os.ModeCharDevice) == 0

			// エラー条件のチェック
			if isPiped {
				if filename != "" || len(args) > 0 {
					return fmt.Errorf("cannot specify a file or arguments when using standard input")
				}
			} else {
				if filename != "" && len(args) > 0 {
					return fmt.Errorf("please specify either a file or a header to decode, not both")
				}

				// if filename == "" && len(args) == 0 {
				// 	return fmt.Errorf("please specify a file, a header to decode, or provide input via stdin")
				// }
			}

			if len(args) > 1 {
				return fmt.Errorf("too many arguments; only one arg is allowed")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var decoded string
			var err error

			if filename != "" {
				err = decodeEmlPrompt(filename)
				if err != nil {
					return fmt.Errorf("failed to decode file '%s': %v", filename, err)
				}
				return nil
			}

			if len(args) == 1 {
				decoded, err = gomime.DecodeHeader(args[0])
				if err != nil {
					return fmt.Errorf("failed to decode header: %v", err)
				}
				fmt.Println(decoded)
				return nil
			}

			// 標準入力からの読み取り
			stats, _ := os.Stdin.Stat()
			if (stats.Mode() & os.ModeCharDevice) != 0 {
				err := selectFuncPrompt(true)
				if err != nil {
					return fmt.Errorf("failed to decode header: %v", err)
				}
				return nil
			}

			data, err := io.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("failed to read from stdin: %v", err)
			}
			decoded, err = gomime.DecodeHeader(string(data))
			if err != nil {
				return fmt.Errorf("failed to decode header from stdin: %v", err)
			}
			fmt.Println(decoded)
			return nil
		},
	}

	cmd.Flags().StringVarP(&filename, "file", "f", "", "file to decode")
	return cmd
}
