package cmd

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestDecodeCommand(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		inputStdin     string
		expectOutput   string
		expectError    bool
		expectedErrMsg string
	}{
		{
			name:         "Decode from argument",
			args:         []string{"decode", "'=?UTF-8?B?44GT44KT44Gr44Gh44Gv?='"},
			expectOutput: "こんにちは",
			expectError:  false,
		},
		{
			name:         "Decode from stdin",
			args:         []string{"decode"},
			inputStdin:   "'=?UTF-8?B?44GT44KT44Gr44Gh44Gv?='",
			expectOutput: "こんにちは",
			expectError:  false,
		},
		// {
		// 	name:           "No input provided",
		// 	args:           []string{"decode"},
		// 	expectError:    true,
		// 	expectedErrMsg: "please specify a file, a header to decode, or provide input via stdin",
		// },
		{
			name:           "Both file and argument provided",
			args:           []string{"decode", "-f", "test.eml", "=?ISO-2022-JP?B?GyRCJDMkcyRLJEEkTxsoQg==?="},
			expectError:    true,
			expectedErrMsg: "please specify either a file or a header to decode, not both",
		},
		{
			name:           "Too many arguments",
			args:           []string{"decode", "'=?UTF-8?B?44GT44KT44Gr44Gh44Gv?='", "'=?UTF-8?B?44GT44KT44Gr44Gh44Gv?='"},
			expectError:    true,
			expectedErrMsg: "too many arguments; only one arg is allowed",
		},
		{
			name:           "Decode from stdin and file",
			args:           []string{"decode", "-f", "test.eml"},
			inputStdin:     "=?UTF-8?B?44GT44KT44Gr44Gh44Gv?=",
			expectError:    true,
			expectedErrMsg: "cannot specify a file or arguments when using standard input",
		},
		{
			name:           "Decode from stdin and argument",
			args:           []string{"decode", "=?UTF-8?B?44GT44KT44Gr44Gh44Gv?="},
			inputStdin:     "=?UTF-8?B?44GT44KT44Gr44Gh44Gv?=",
			expectError:    true,
			expectedErrMsg: "cannot specify a file or arguments when using standard input",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// キャプチャ用のバッファを作成
			var outputBuf bytes.Buffer
			// エラーメッセージ用のバッファを作成
			var errBuf bytes.Buffer

			// 元のStdoutとStderrを保存
			origStdout := os.Stdout
			origStderr := os.Stderr

			// テスト用にStdoutとStderrをリダイレクト
			rOut, wOut, _ := os.Pipe()
			rErr, wErr, _ := os.Pipe()
			os.Stdout = wOut
			os.Stderr = wErr

			// コマンドのセットアップ
			decodeCmd := DecodeCmd()
			root := &cobra.Command{Use: "gemm"}
			root.AddCommand(decodeCmd)

			// 引数の設定
			root.SetArgs(tt.args[0:])

			var err error

			// 標準入力の設定（必要な場合）
			if tt.inputStdin != "" {
				// 標準入力をモック
				origStdin := os.Stdin
				defer func() { os.Stdin = origStdin }()
				r, w, _ := os.Pipe()
				w.WriteString(tt.inputStdin)
				w.Close()
				os.Stdin = r
			}

			// コマンドの実行
			err = root.Execute()

			// 出力を取得
			wOut.Close()
			wErr.Close()
			outBytes, _ := io.ReadAll(rOut)
			errBytes, _ := io.ReadAll(rErr)
			outputBuf.Write(outBytes)
			errBuf.Write(errBytes)

			// 元のStdoutとStderrに戻す
			os.Stdout = origStdout
			os.Stderr = origStderr

			output := outputBuf.String()
			errorOutput := errBuf.String()

			if tt.expectError {
				assert.Error(t, err)
				if tt.expectedErrMsg != "" {
					assert.Contains(t, errorOutput, tt.expectedErrMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.Contains(t, output, tt.expectOutput)
			}
		})
	}
}
