package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopy(t *testing.T) {
	testCases := []struct {
		name          string
		fromPath      string
		toPath        string
		validCopyPath string
		offset        int64
		limit         int64
		err           error
	}{
		{
			name:          "unsupported file",
			fromPath:      "/dev/urandom",
			toPath:        "/tmp/urandom.txt",
			validCopyPath: "testdata/out_offset0_limit0.txt",
			offset:        0,
			limit:         0,
			err:           ErrUnsupportedFile,
		},
		{
			name:          "file from does not exist",
			fromPath:      "testdata/ofkwe0fwepfkew",
			toPath:        "/tmp/out_0_0.txt",
			validCopyPath: "testdata/out_offset0_limit0.txt",
			offset:        0,
			limit:         0,
			err:           os.ErrNotExist,
		},
		{
			name:          "negative offset",
			fromPath:      "testdata/input.txt",
			toPath:        "/tmp/out_0_0.txt",
			validCopyPath: "testdata/out_offset0_limit0.txt",
			offset:        -1,
			limit:         0,
			err:           ErrNegativeOffset,
		},
		{
			name:          "negative limit",
			fromPath:      "testdata/input.txt",
			toPath:        "/tmp/out_0_0.txt",
			validCopyPath: "testdata/out_offset0_limit0.txt",
			offset:        0,
			limit:         -1,
			err:           ErrNegativeLimit,
		},
		{
			name:          "offset exceeds file size",
			fromPath:      "testdata/input.txt",
			toPath:        "/tmp/out_0_0.txt",
			validCopyPath: "testdata/out_offset0_limit0.txt",
			offset:        1_000_000_000,
			limit:         0,
			err:           ErrOffsetExceedsFileSize,
		},
		{
			name:          "valid files offset 0 limit 0",
			fromPath:      "testdata/input.txt",
			toPath:        "/tmp/out_0_0.txt",
			validCopyPath: "testdata/out_offset0_limit0.txt",
			offset:        0,
			limit:         0,
			err:           nil,
		},
		{
			name:          "valid files offset 0 limit 10",
			fromPath:      "testdata/input.txt",
			toPath:        "/tmp/out_0_10.txt",
			validCopyPath: "testdata/out_offset0_limit10.txt",
			offset:        0,
			limit:         10,
			err:           nil,
		},
		{
			name:          "valid files offset 0 limit 1000",
			fromPath:      "testdata/input.txt",
			toPath:        "/tmp/out_0_1000.txt",
			validCopyPath: "testdata/out_offset0_limit1000.txt",
			offset:        0,
			limit:         1000,
			err:           nil,
		},
		{
			name:          "valid files offset 0 limit 10000",
			fromPath:      "testdata/input.txt",
			toPath:        "/tmp/out_0_10000.txt",
			validCopyPath: "testdata/out_offset0_limit10000.txt",
			offset:        0,
			limit:         10000,
			err:           nil,
		},
		{
			name:          "valid files offset 100 limit 1000",
			fromPath:      "testdata/input.txt",
			toPath:        "/tmp/out_100_1000.txt",
			validCopyPath: "testdata/out_offset100_limit1000.txt",
			offset:        100,
			limit:         1000,
			err:           nil,
		},
		{
			name:          "valid files offset 6000 limit 1000",
			fromPath:      "testdata/input.txt",
			toPath:        "/tmp/out_6000_1000.txt",
			validCopyPath: "testdata/out_offset6000_limit1000.txt",
			offset:        6000,
			limit:         1000,
			err:           nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer os.Remove(tc.toPath)
			assert.FileExists(t, tc.validCopyPath)
			err := Copy(tc.fromPath, tc.toPath, tc.offset, tc.limit)
			if tc.err != nil {
				assert.ErrorIs(t, err, tc.err)
				return
			}
			assert.FileExists(t, tc.toPath)
			ok, err := filesHasSameContent(tc.toPath, tc.validCopyPath)
			assert.NoError(t, err)
			assert.True(t, ok)
		})
	}
}
