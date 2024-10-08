package cmd

import (
	"fmt"

	"github.com/ProtonMail/go-mime"
  "github.com/manifoldco/promptui"
	"github.com/yken2257/gemm/utils"
)

func selectPromptAction(label string, items []string) (int, string, error) {
	prompt := promptui.Select{
			Label:        label,
			Items:        items,
			CursorPos:    0,
			HideSelected: false,
			HideHelp:     true,
	}

	index, result, err := prompt.Run()
	if err != nil {
			return -1, "", err
	}
	return index, result, nil
}

func selectFuncPrompt(decodeOnly bool) error {
	var funcOptions = []string{"Decode a single header", "Decode an .eml file", "Encode string"}
	if decodeOnly {
		funcOptions = funcOptions[:2]
	}
	_, result, err := selectPromptAction("What do you want to do?", funcOptions)
	if err != nil {
		return err
	}

	switch result {
	case funcOptions[0]:
		return decodeHeaderPrompt()
	case funcOptions[1]:
		result, err := chooseFilePrompt()
		if err != nil {
			return err
		}
		return decodeEmlPrompt(result)
	case funcOptions[2]:
		return encodePrompt("", "", "")
	}
	return nil
}

func decodeHeaderPrompt() error {
	prompt := promptui.Prompt{
		Label:       "Enter the header text",
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
		Label:       "Enter the path to the .eml file",
		HideEntered: false,
	}

	result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return result, nil
}

func decodeEmlPrompt(filename string) error {
	decodedHeaders, err := utils.DecodeHeaders(filename)
	if err != nil {
		return err
	}
	var headerKeys []string
	for key := range decodedHeaders {
		headerKeys = append(headerKeys, key)
	}
	_, choice, err := selectPromptAction("Choose a header to decode", headerKeys)
	if err != nil {
		return err
	}
	decoded := decodedHeaders[choice]

	fmt.Println(decoded)
	return nil
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
		return encodePrompt(textInput, charset, encoding)
	} 
	if charset == "" {
		items := []string{"UTF-8", "ISO-2022-JP", "Shift_JIS"}
		_, charsetInput, err := selectPromptAction("Choose a charset", items)
		if err != nil {
			return err
		}
		return encodePrompt(text, charsetInput, encoding)
	} 
	if encoding == "" {
		items := []string{"B", "Q"}
		_, encodingInput, err := selectPromptAction("Choose an encoding", items)
		if err != nil {
			return err
		}
		return encodePrompt(text, charset, encodingInput)
	}
	
	normalizedCharset := utils.NormalizeCharset(charset)
	encoded, err := utils.EncodeHeader(text, normalizedCharset, encoding)
	if err != nil {
		return err
	}
	fmt.Println(encoded)
	return nil
}