package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]string

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	// Place your code here
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var env = make(Environment)

	for _, v := range entries {
		if !v.Mode().IsRegular() {
			continue
		}
		// ignore some files
		if strings.ContainsAny(v.Name(), "=;") {
			continue
		}

		// make full path to file
		fileName, err := filepath.Abs(filepath.Join(dir, v.Name()))
		if err != nil {
			return nil, err
		}

		value, err := getValue(fileName)
		if err != nil {
			return nil, err
		}

		env[v.Name()] = value
	}

	return env, nil
}

func getValue(fileName string) (string, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0644)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	if scanner.Err() != nil {
		return "", err
	}

	line := scanner.Text()
	if line == "" {
		return "", nil
	}

	// if 0x00 exist in line, then replace all with new line
	if bytes.Contains([]byte(line), []byte{0}) {
		line = string(bytes.ReplaceAll([]byte(line), []byte{0}, []byte{10}))
	}

	return strings.TrimRight(line, " \t"), nil
}
