package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"path"
	"strings"
	"unicode"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// NeedSetToEnv нужно ли присваивать значение в окружение?
// У переменных, которые нужно присвоить, NeedRemove в читающей функции записывается в false.
// Сравнивать только по v.Value != "" не совсем корректно, т.к. в value может быть пустая строка.
func (v EnvValue) NeedSetToEnv() bool {
	return v.Value != "" || !v.NeedRemove
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	env := make(Environment)
	for _, file := range entries {
		if file.IsDir() || strings.Contains(file.Name(), "=") {
			continue
		}
		fInfo, err := file.Info()
		if err != nil {
			return nil, err
		}
		if fInfo.Size() == 0 {
			env[file.Name()] = EnvValue{NeedRemove: true}
		} else {
			val, err := readEnvVal(path.Join(dir, file.Name()))
			if err != nil {
				return nil, err
			}
			env[file.Name()] = EnvValue{Value: val}
		}
	}
	return env, nil
}

// readEnvVal вычитать из файла значение переменной окружения. Удаляет лидирующие и завершающие пробелы.
// Заменяет терминальные нули.
func readEnvVal(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	r := bufio.NewReader(f)
	s, _, err := r.ReadLine()
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}
	s = bytes.TrimRightFunc(s, unicode.IsSpace)
	s = bytes.ReplaceAll(s, []byte("\x00"), []byte("\n"))
	return string(s), nil
}
