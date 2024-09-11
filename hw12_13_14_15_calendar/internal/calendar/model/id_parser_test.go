package model

import "testing"

func TestParseEventID(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		ok    bool
	}{
		{
			"invalid empty string",
			"",
			false,
		},
		{
			"invalid input",
			"efjowefjoewf",
			false,
		},
		{
			"valid input",
			"bc94b5a3-182b-48b1-8edd-63b5c32dcfdc",
			true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ParseEventID(tc.input)
			if tc.ok && err != nil {
				t.Fatalf("unexpected err: %s", err)
			}
			if !tc.ok && err == nil {
				t.Fatalf("expected to return an error, but nothing")
			}
		})
	}
}
