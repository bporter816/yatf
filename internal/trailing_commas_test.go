package internal

import "testing"

func TestListTrailingCommas(t *testing.T) {
	tests := []TestCase{
		NewUnchangedTestCase(`
resource "test" "test" {
  list = [1, 2, 3]
}`),
		NewUnchangedTestCase(`
resource "test" "test" {
  list = [
    1,
    2,
  ]
}`),
		{
			input: `
resource "test" "test" {
  list = [
    1,
    2
  ]
}`,
			expected: `
resource "test" "test" {
  list = [
    1,
    2,
  ]
}`,
		},
		{
			input: `
resource "test" "test" {
  map = {
    a = 1,
    b = 2
  }
}`,
			expected: `
resource "test" "test" {
  map = {
    a = 1
    b = 2
  }
}`,
		},
		{
			input: `
resource "test" "test" {
  map = {
    a = 1,
    l = [
      1,
      [
        1
      ],
      {
        a = 1,
      },
      3,
      "v"
    ]
    b = 2,
  }
}`,
			expected: `
resource "test" "test" {
  map = {
    a = 1
    l = [
      1,
      [
        1,
      ],
      {
        a = 1
      },
      3,
      "v",
    ]
    b = 2
  }
}`,
		},
	}

	for i, tc := range tests {
		got := LintTrailingCommasString(tc.input)
		if got != tc.expected {
			t.Fatalf("test case %v failed", i)
		}
	}
}
