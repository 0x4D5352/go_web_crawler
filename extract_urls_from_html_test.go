package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestGetURLsFromHTML(t *testing.T) {
	tests := []struct {
		name          string
		inputURL      string
		inputBody     string
		expected      []string
		errorContains string
	}{
		{
			name:     "absolute and relative URLs",
			inputURL: "https://blog.boot.dev",
			inputBody: `
<html>
	<body>
		<a href="/path/one">
			<span>Boot.dev</span>
		</a>
		<a href="https://other.com/path/one">
			<span>Boot.dev</span>
		</a>
	</body>
</html>
`,
			expected: []string{"https://blog.boot.dev/path/one", "https://other.com/path/one"},
		},
		{
			name:     "test nested div anchors",
			inputURL: "https://blog.boot.dev",
			inputBody: `
<html>
	<body>
		<a href="/path/one">
			<span>Boot.dev</span>
		</a>
		<a href="https://other.com/path/one">
			<span>Boot.dev</span>
		</a>
		<div>
			<a href="/path/two">
				<span>Boot.dev</span>
			</a>
			<div>
				<a href="/path/three">
					<span>Boot.dev</span>
				</a>
				<div>
					<a href="https://google.com">
						<span>Google It</span>
					</a>
				</div>
			</div>
	</body>
</html>
`,
			expected: []string{"https://blog.boot.dev/path/one", "https://other.com/path/one", "https://blog.boot.dev/path/two", "https://blog.boot.dev/path/three", "https://google.com"},
		},
		{
			name:     "test empty slice",
			inputURL: "https://blog.boot.dev",
			inputBody: `
<html>
	<body>
		<p>Hello!</p>
		<div>
			<p>Welcome to my <em>blog</em>.</p>
			<p>I hope you enjoy!</p>
		</div>
	</body>
</html>
`,
			expected: []string{},
		},
		{
			name:          "handle missing body",
			inputURL:      "https://blog.boot.dev",
			inputBody:     "",
			expected:      []string{},
			errorContains: "missing body",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := getURLsFromHTML(tc.inputBody, tc.inputURL)
			if err != nil && !strings.Contains(err.Error(), tc.errorContains) {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			} else if err != nil && tc.errorContains == "" {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			} else if err == nil && tc.errorContains != "" {
				t.Errorf("Test %v - '%s' FAIL: expected error containing '%v', got none.", i, tc.name, tc.errorContains)
				return
			}

			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Test %v - %s FAIL: expected URL: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}
