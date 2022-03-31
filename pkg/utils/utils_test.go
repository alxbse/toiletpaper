package utils

import (
	"testing"
)

func TestSnakeCasePath(t *testing.T) {
	tt := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "test simple path",
			path:     "simple",
			expected: "simple",
		},
		{
			name:     "test nested path",
			path:     "nested/path",
			expected: "nestedPath",
		},
		{
			name:     "test multiple nested path",
			path:     "multiple/nested/path",
			expected: "multipleNestedPath",
		},
		{
			name:     "test path with extension",
			path:     "test.go",
			expected: "testGo",
		},
		{
			name:     "test path with dashes",
			path:     "test-file",
			expected: "testFile",
		},
		{
			name:     "test combination of most things",
			path:     "test/combination/Of-most.things",
			expected: "testCombinationOfMostThings",
		},
		{
			name:     "test repeating punctuation",
			path:     "test//repeating--punctuation",
			expected: "testRepeatingPunctuation",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			actual := SnakeCasePath(tc.path)
			if actual != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, actual)
			}
		})
	}
}
