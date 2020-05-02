package main //nolint:golint,stylecheck

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath string, toPath string, offset, limit int64) error {
	// Place your code here
	from, err := os.OpenFile(fromPath, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer from.Close()

	to, err := os.OpenFile(toPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer to.Close()

	fromInfo, err := os.Stat(fromPath)
	if err != nil {
		return err
	}
	if !fromInfo.Mode().IsRegular() {
		return ErrUnsupportedFile
	}
	toInfo, err := os.Stat(toPath)
	if err != nil {
		return err
	}
	if !toInfo.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	var fromSize = fromInfo.Size()

	if offset > fromSize {
		return ErrOffsetExceedsFileSize
	}

	if offset > 0 {
		_, err := from.Seek(offset, io.SeekStart)
		if err != nil {
			return err
		}
	}

	var bytesToCopy int64
	if limit == 0 || limit >= fromSize || limit >= fromSize-offset {
		bytesToCopy = fromSize - offset
	} else {
		bytesToCopy = limit
	}

	// start new bar
	bar := pb.Full.Start64(bytesToCopy)
	defer bar.Finish()
	// create proxy reader
	barWriter := bar.NewProxyWriter(to)

	_, err = io.CopyN(barWriter, from, bytesToCopy)
	if err != io.EOF && err != nil {
		return err
	}

	return nil
}
