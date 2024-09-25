package utils

import (
	"os"
	"fmt"
	"net/mail"
	"strings"
	"mime"
	"golang.org/x/text/encoding/japanese"

	"github.com/ProtonMail/go-mime"
)

func DecodeHeaders(filename string) (map[string]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	byteData := make([]byte, 1000000)
	count, err := file.Read(byteData)
	if err != nil {
		return nil, err
	}
	message := string(byteData[:count])
	mm, err := mail.ReadMessage(strings.NewReader(message))
	if err != nil {
		return nil, err
	}
	header := mm.Header
	decodedHeaders := make(map[string]string)

	for key, value := range header {
		if containsEncodedWord(value[0]) {
			decoded, err := gomime.DecodeHeader(value[0])
			if err != nil {
				return nil, err
			}
			decodedHeaders[key] = decoded
		}
	}
	// if no encoded word found, raise an error
	if len(decodedHeaders) == 0 {
		return nil, fmt.Errorf("no encoded header found")
	}
	return decodedHeaders, nil
}

func containsEncodedWord(s string) bool {
	// if space included, split by space
	var components []string
	if strings.Contains(s, " ") {
		components = strings.Split(s, " ")
	} else {
		components = []string{s}
	}
	for _, component := range components {
		if isEncodedWord(component) {
			return true
		}
	}
	return false
}

func isEncodedWord(s string) bool {
	if !strings.HasPrefix(s, "=?") || !strings.HasSuffix(s, "?=") {
		return false
	}
	// if the number of "?" is not 4, return false
	if strings.Count(s, "?") != 4 {
		return false
	}
	// split into charset, encoding, and encoded text
	s = strings.TrimPrefix(s, "=?")
	s = strings.TrimSuffix(s, "?=")
	slice := strings.Split(s, "?")
	encoding := slice[1]
	// if encoding is not "B" or "Q", return false
	if encoding != "B" && encoding != "Q" {
		return false
	}
	return true
}

func EncodeHeader(s, charset, encoding string) (string, error) {
	var encodedBytes []byte
	var err error

	switch charset {
	case "UTF-8":
		encodedBytes = []byte(s)
	case "ISO-2022-JP":
		encodedBytes, err = japanese.ISO2022JP.NewEncoder().Bytes([]byte(s))
		if err != nil {
			return "", err
		}
	case "Shift_JIS":
		encodedBytes, err = japanese.ShiftJIS.NewEncoder().Bytes([]byte(s))
		if err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("invalid charset")
	}

	switch encoding {
	case "B":
		return mime.BEncoding.Encode(charset, string(encodedBytes)), nil
	case "Q":
		return mime.QEncoding.Encode(charset, string(encodedBytes)), nil
	default:
		return "", fmt.Errorf("invalid encoding")
	}
}