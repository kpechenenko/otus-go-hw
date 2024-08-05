package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestReadDir(t *testing.T) {
	testCases := []struct {
		name    string
		dir     string
		wantEnt Environment
		err     error
	}{
		{
			"valid env",
			"./testdata/env",
			Environment{
				"BAR":   EnvValue{Value: "bar", NeedRemove: false},
				"EMPTY": EnvValue{Value: "", NeedRemove: false},
				"FOO":   EnvValue{Value: "   foo\nwith new line", NeedRemove: false},
				"HELLO": EnvValue{Value: "\"hello\"", NeedRemove: false},
				"UNSET": EnvValue{Value: "", NeedRemove: true},
			},
			nil,
		},
		{
			"non-existent dir",
			"./testdata/efoewf",
			nil,
			os.ErrNotExist,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := ReadDir(tc.dir)
			if tc.err != nil {
				assert.ErrorIs(t, err, tc.err)
			}
			assert.Equal(t, tc.wantEnt, actual)
		})
	}
}

func TestEnvValue_NeedSetToEnv(t *testing.T) {
	testCases := []struct {
		name  string
		value EnvValue
		want  bool
	}{
		{
			"valid to set",
			EnvValue{Value: "", NeedRemove: false},
			true,
		},
		{
			"valid to set non empty",
			EnvValue{Value: "efefe", NeedRemove: true},
			true,
		},
		{
			"valid to unset",
			EnvValue{Value: "", NeedRemove: true},
			false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.value.NeedSetToEnv(); got != tc.want {
				t.Fail()
			}
		})
	}
}
