package main

import (
	"io"
	"os"
)

// filesHasSameContent содержат ли файлы идентичный контент в байтах?
func filesHasSameContent(path1 string, path2 string) (bool, error) {
	f1, err := os.Open(path1)
	if err != nil {
		return false, err
	}
	f2, err := os.Open(path2)
	if err != nil {
		return false, err
	}
	buf1, err := io.ReadAll(f1)
	if err != nil {
		return false, err
	}
	buf2, err := io.ReadAll(f2)
	if err != nil {
		return false, err
	}
	if len(buf1) != len(buf2) {
		return false, nil
	}
	for i := 0; i < len(buf1); i++ {
		if buf1[i] != buf2[i] {
			return false, nil
		}
	}
	return true, nil
}
