package utils

import "testing"

func TestDecodeHeaders(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected map[string]string
	}{
		{
			name:     "Simple",
			input:		"../test_files/simple.eml",
			expected: map[string]string{
				"From": "John Doe （ジョン　ドゥー） <john@example.com>",
				"To": "ジェーン・ドゥー <jane@example.co.jp>",
				"Subject": "Re: ご飯に行きませんか？",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			decoded, err := DecodeHeaders(tc.input)
			if err != nil {
				t.Fatalf("failed to decode: %v", err)
			}
			if len(decoded) != len(tc.expected) {
				t.Fatalf("expected %d headers, got %d", len(tc.expected), len(decoded))
			}
			for key, value := range decoded {
				if tc.expected[key] != value {
					t.Fatalf("expected %s to be %s, got %s", key, tc.expected[key], value)
				}
			}
		})
	}
}