package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	// Place your code here
	content := []byte("temporary file's content")
	outFileName := "temp-out.txt"

	tmpfile, err := ioutil.TempFile("", "tempfile")
	if err != nil {
		require.Fail(t, "can't open temp file", err)
	}

	defer func() {
		if err := os.Remove(tmpfile.Name()); err != nil {
			require.Fail(t, "can't remove test file", err)
		}
		if err := os.Remove(outFileName); err != nil {
			require.Fail(t, "can't remove test file", err)
		}
	}() // clean up

	if _, err := tmpfile.Write(content); err != nil {
		require.Fail(t, "can't write to temp file", err)
	}
	if err := tmpfile.Close(); err != nil {
		require.Fail(t, "can't close temp file", err)
	}

	t.Run("plain copy", func(t *testing.T) {
		err = Copy(tmpfile.Name(), outFileName, 0, 0)
		require.Nil(t, err)

		copyResult, err := ioutil.ReadFile(outFileName)
		if err != nil {
			require.Fail(t, "can't open copied file", err)
		}
		require.Equal(t, content, copyResult)
	})

	t.Run("offset error", func(t *testing.T) {
		err = Copy(tmpfile.Name(), outFileName, 50, 0)
		require.Equal(t, ErrOffsetExceedsFileSize, err)
	})

	t.Run("not regular file", func(t *testing.T) {
		err = Copy(tmpfile.Name(), "/dev/urandom", 0, 0)
		require.Equal(t, ErrUnsupportedFile, err)
	})

}
