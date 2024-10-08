package cmd

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestEncodeCommand(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		inputStdin     string
		expectOutput   string
		expectError    bool
		expectedErrMsg string
	}{
		{
			name:         "Encode from argument",
			args:         []string{"encode", "こんにちは", "-c", "UTF-8", "-e", "B"},
			expectOutput: "=?UTF-8?b?44GT44KT44Gr44Gh44Gv?=",
			expectError:  false,
		},
		{
			name:         "Encode from stdin",
			args:         []string{"encode", "-c", "UTF-8", "-e", "B"},
			inputStdin:   "こんにちは",
			expectOutput: "=?UTF-8?b?44GT44KT44Gr44Gh44Gv?=",
			expectError:  false,
		},
		{
			name:           "Stdin is piped but charset and encoding are not set",
			args:           []string{"encode"},
			inputStdin:     "こんにちは",
			expectError:    true,
			expectedErrMsg: "charset and encoding must be set when using stdin; e.g. gemm encode -c UTF-8 -e B",
		},
		{
			name:           "Invalid charset",
			args:           []string{"encode", "こんにちは", "-c", "UTF-16", "-e", "B"},
			expectError:    true,
			expectedErrMsg: "charset must be either UTF-8, ISO-2022-JP, or Shift_JIS",
		},
		{
			name:           "Invalid encoding",
			args:           []string{"encode", "こんにちは", "-c", "UTF-8", "-e", "C"},
			expectError:    true,
			expectedErrMsg: "encoding must be either B or Q",
		},
		{
			name:           "Too many arguments",
			args:           []string{"encode", "こんにちは", "こんにちは"},
			expectError:    true,
			expectedErrMsg: "too many arguments; only one arg is allowed",
		},
		{
			name:           "lowercase encoding",
			args:           []string{"encode", "こんにちは", "-c", "UTF-8", "-e", "b"},
			expectOutput:   "=?UTF-8?b?44GT44KT44Gr44Gh44Gv?=",
			expectError:    false,
		},
		{
			name:           "lowercase charset",
			args:           []string{"encode", "こんにちは", "-c", "utf8", "-e", "B"},
			expectOutput:   "=?UTF-8?b?44GT44KT44Gr44Gh44Gv?=",
			expectError:    false,
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
			encodeCmd := EncodeCmd()
			root := &cobra.Command{Use: "gemm"}
			root.AddCommand(encodeCmd)

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
