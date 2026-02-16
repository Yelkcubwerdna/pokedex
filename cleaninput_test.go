package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "   hello world   ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "",
			expected: []string{},
		},
		{
			input:    "beats  in my   head ",
			expected: []string{"beats", "in", "my", "head"},
		},
		// add more cases here
	}

	for _, c := range cases {
		actual := cleanInput(c.input)

		if len(actual) != len(c.expected) {
			t.Errorf("Slice lengths do not match\n")
			t.Errorf("Expected: %d\n", len(c.expected))
			t.Errorf("     %v", c.expected)
			t.Errorf("Actual: %d\n", len(actual))
			t.Errorf("     %v", actual)
			t.FailNow()
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]

			if word != expectedWord {
				t.Errorf("Word mismatch\n")
				t.Errorf("Expected: %s", expectedWord)
				t.Errorf("Actual: %s", word)
				t.Fail()
			}
		}
	}
}
