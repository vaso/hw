package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require" //nolint:all
)

func TestCopy(t *testing.T) {
	t.Run("negative offset", func(t *testing.T) {
		err := Copy("testdata/input.txt", "someoutput.txt", -1, 0)
		require.Error(t, err)
	})
	t.Run("negative limit", func(t *testing.T) {
		err := Copy("testdata/input.txt", "someoutput.txt", 0, -1)
		require.Error(t, err)
	})

	t.Run("offset 0 limit 0", func(t *testing.T) {
		testCopiedFile(t, 0, 0)
	})

	t.Run("offset 0 limit 10", func(t *testing.T) {
		testCopiedFile(t, 0, 10)
	})

	t.Run("offset 100 limit 1000", func(t *testing.T) {
		testCopiedFile(t, 100, 1000)
	})
}

func testCopiedFile(t *testing.T, offset int64, limit int64) {
	t.Helper()
	f, err := os.CreateTemp("./", "file.*.tmp")
	if err != nil {
		log.Fatal(err)
	}
	outputFileName := f.Name()

	err = Copy("testdata/input.txt", outputFileName, offset, limit)
	require.NoError(t, err)
	require.FileExists(t, outputFileName)

	expectedFilePath := fmt.Sprintf("testdata/out_offset%d_limit%d.txt", offset, limit)
	testEqualFile(t, expectedFilePath, outputFileName)

	_ = os.Remove(outputFileName)
	_ = f.Close()
}

func testEqualFile(t *testing.T, expectedFileName string, actualFileName string) {
	t.Helper()
	expectedFile, err := os.Open(expectedFileName)
	require.NoError(t, err)
	defer func(expectedFile *os.File) {
		_ = expectedFile.Close()
	}(expectedFile)

	actualFile, err := os.Open(actualFileName)
	require.NoError(t, err)
	defer func(expectedFile *os.File) {
		_ = expectedFile.Close()
	}(expectedFile)

	expectedContent, err := io.ReadAll(expectedFile)
	require.NoError(t, err)

	actualContent, err := io.ReadAll(actualFile)
	require.NoError(t, err)

	require.Equal(t, expectedContent, actualContent)
}
