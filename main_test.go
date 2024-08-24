// main_test.go
package main

import (
	"testing"
)

func TestModifyTwitterLinks(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Twitter link",
			input:    "Check out https://twitter.com/user/status/123456",
			expected: "Check out https://fxtwitter.com/user/status/123456",
		},
		{
			name:     "X link",
			input:    "Look at https://x.com/user/status/789012",
			expected: "Look at https://fixupx.com/user/status/789012",
		},
		{
			name:     "Multiple links",
			input:    "Twitter: https://twitter.com/user1/status/123 and X: https://x.com/user2/status/456",
			expected: "Twitter: https://fxtwitter.com/user1/status/123 and X: https://fixupx.com/user2/status/456",
		},
		{
			name:     "No links",
			input:    "Just a regular message",
			expected: "Just a regular message",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := modifyTwitterLinks(tc.input)
			if result != tc.expected {
				t.Errorf("modifyTwitterLinks(%q) = %q; want %q", tc.input, result, tc.expected)
			}
		})
	}
}