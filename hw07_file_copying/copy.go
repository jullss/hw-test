package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrNoAccessFileInfo      = errors.New("no access to file info")
	ErrCannotOpenFile        = errors.New("cannot open source file")
	ErrCannotCreateFile      = errors.New("cannot create file")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	fileInfo, err := os.Stat(fromPath)
	if err != nil {
		return ErrNoAccessFileInfo
	}

	if !fileInfo.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	fileSize := fileInfo.Size()
	if fileSize < offset {
		return ErrOffsetExceedsFileSize
	}

	src, err := os.Open(fromPath)
	if err != nil {
		return ErrCannotOpenFile
	}
	defer src.Close()

	dst, err := os.Create(toPath)
	if err != nil {
		return ErrCannotCreateFile
	}
	defer dst.Close()

	_, err = src.Seek(offset, 0)
	if err != nil {
		return err
	}

	var totalBytes int64
	if limit > 0 {
		totalBytes = limit
	} else {
		totalBytes = fileSize - offset
	}

	bar := pb.New64(totalBytes)

	bar.ShowCounters = true
	bar.ShowSpeed = true

	bar.Start()
	defer bar.Finish()

	barReader := bar.NewProxyReader(src)

	if limit == 0 {
		_, err = io.Copy(dst, barReader)
	} else {
		_, err = io.CopyN(dst, barReader, limit)
	}

	return err
}
