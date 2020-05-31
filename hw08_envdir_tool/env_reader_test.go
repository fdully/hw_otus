package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	// Place your code here
	testFilesDir := "./testdata/env"

	t.Run("read files from testdata dir", func(t *testing.T) {
		var expected = Environment{"BAR": "bar", "FOO": "   foo\nwith new line", "HELLO": "\"hello\"", "UNSET": ""}

		env, err := ReadDir(testFilesDir)
		if err != nil {
			t.Fatalf("can't read dir %v\n", err)
		}

		require.Equal(t, expected, env)
	})

	t.Run("ignoring files with '=' ';' in name", func(t *testing.T) {
		f, err := ioutil.TempFile(testFilesDir, "file=test.txt")
		if err != nil {
			t.Fatalf("can't create tempfile in testdata dir: %s\n", err.Error())
		}
		var expected = Environment{"BAR": "bar", "FOO": "   foo\nwith new line", "HELLO": "\"hello\"", "UNSET": ""}
		defer os.Remove(f.Name())

		list, err := ioutil.ReadDir(testFilesDir)
		if err != nil {
			t.Fatalf("can't read testdata dir %s\n", err.Error())
		}

		env, err := ReadDir(testFilesDir)
		if err != nil {
			t.Fatalf("can't read dir %v\n", err)
		}

		require.Equal(t, expected, env)
		require.Equal(t, 5, len(list))
	})

}
