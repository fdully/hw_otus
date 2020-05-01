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
		require.Errorf(t, err, "can't open temp file")
	}

	defer func() {
		if err := os.Remove(tmpfile.Name()); err != nil {
			require.Errorf(t, err, "can't remove test file")
		}
		if err := os.Remove(outFileName); err != nil {
			require.Errorf(t, err, "can't remove test file")
		}
	}() // clean up

	if _, err := tmpfile.Write(content); err != nil {
		require.Errorf(t, err, "can't write to temp file")
	}
	if err := tmpfile.Close(); err != nil {
		require.Errorf(t, err, "can't close temp file")
	}

	t.Run("plain copy", func(t *testing.T) {
		err = Copy(tmpfile.Name(), outFileName, 0, 0)
		require.Nil(t, err)
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
