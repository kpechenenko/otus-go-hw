package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilesHasSameContent(t *testing.T) {
	testCases := []struct {
		name  string
		path1 string
		path2 string
		same  bool
		err   error
	}{
		{
			name:  "1 file does not exist",
			path1: "testdata/efjewofweofkewp;",
			path2: "testdata/input.txt",
			same:  false,
			err:   os.ErrNotExist,
		},
		{
			name:  "2 file does not exist",
			path1: "testdata/input.txt",
			path2: "testdata/dlwpdewpdlewpdpew",
			same:  false,
			err:   os.ErrNotExist,
		},
		{
			name:  "same file",
			path1: "testdata/input.txt",
			path2: "testdata/input.txt",
			same:  true,
			err:   nil,
		},
		{
			name:  "another same file",
			path1: "testdata/out_offset0_limit0.txt",
			path2: "testdata/out_offset0_limit0.txt",
			same:  true,
			err:   nil,
		},
		{
			name:  "different files",
			path1: "testdata/out_offset0_limit1000.txt",
			path2: "testdata/input.txt",
			same:  false,
			err:   nil,
		},
		{
			name:  "other different files",
			path1: "testdata/out_offset100_limit1000.txt",
			path2: "testdata/out_offset0_limit10000.txt",
			same:  false,
			err:   nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ok, err := filesHasSameContent(tc.path1, tc.path2)
			if tc.err != nil {
				assert.ErrorIs(t, err, tc.err)
			}
			assert.Equal(t, ok, tc.same)
		})
	}
}
