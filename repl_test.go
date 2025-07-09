package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "Go Is tHe bEsT",
			expected: []string{"go", "is", "the", "best"},
		},
		{
			input:    "ComPiled   Languages Go  C   c++ 	Rust",
			expected: []string{"compiled", "languages", "go", "c", "c++", "rust"},
		},
		{
			input:    " 10 20 30 apples	",
			expected: []string{"10", "20", "30", "apples"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)

		if len(actual) != len(c.expected) {
			t.Errorf("output length doesn't match expected length")
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]

			if word != expectedWord {
				t.Errorf("Output doesn't match expectation")
			}
		}
	}
}
