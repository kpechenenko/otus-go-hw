package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRunCmd(t *testing.T) {
	testCases := []struct {
		name string
		cmd  []string
		env  Environment
		code int
	}{
		{
			"valid echo script",
			[]string{"./testdata/echo.sh"},
			nil,
			ReturnCodeOk,
		},
		{
			"non-existent echo script",
			[]string{"./testdata/oedeokd"},
			nil,
			ReturnCodeErr,
		},
		{
			"valid test script",
			[]string{"./test.sh"},
			Environment{
				"BAR":   EnvValue{Value: "bar", NeedRemove: false},
				"EMPTY": EnvValue{Value: "", NeedRemove: false},
				"FOO":   EnvValue{Value: "   foo\nwith new line", NeedRemove: false},
				"HELLO": EnvValue{Value: "\"hello\"", NeedRemove: false},
				"UNSET": EnvValue{Value: "", NeedRemove: true},
			},
			ReturnCodeOk,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			code := RunCmd(tc.cmd, tc.env)
			assert.Equal(t, tc.code, code)
		})
	}
}
