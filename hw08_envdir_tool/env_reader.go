package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	result := Environment{}
	envFiles, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, envFile := range envFiles {
		fileName := envFile.Name()
		if strings.Contains(fileName, "=") {
			continue
		}
		envValue, err := handleEnvFile(dir, fileName)
		if err != nil {
			return nil, err
		}
		result[fileName] = *envValue
	}

	return result, nil
}

func handleEnvFile(dir string, fileName string) (*EnvValue, error) {
	envFile, err := os.Open(fmt.Sprintf("%s/%s", dir, fileName))
	defer func() {
		_ = envFile.Close()
	}()

	if err != nil {
		return nil, err
	}

	stat, err := envFile.Stat()
	if err != nil {
		return nil, err
	}
	if stat.Size() == 0 {
		envValue := EnvValue{
			Value:      "",
			NeedRemove: true,
		}
		return &envValue, nil
	}

	reader := bufio.NewReader(envFile)
	val, _, err := reader.ReadLine()
	if err != nil {
		return nil, err
	}
	val = bytes.ReplaceAll(val, []byte{0x00}, []byte("\n"))
	valStr := string(val)
	valStr = strings.TrimRight(valStr, " \t")

	envValue := EnvValue{
		Value:      valStr,
		NeedRemove: false,
	}

	return &envValue, nil
}
