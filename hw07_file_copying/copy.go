package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrNoAccessFileInfo      = errors.New("no access to file info")
	ErrCannotOpenFile        = errors.New("cannot open source file")
	ErrCannotCreateFile      = errors.New("cannot create file")
)

type ProgressBar struct {
	reader  io.Reader
	total   int64
	counter int64
}

func (pb *ProgressBar) Read(p []byte) (int, error) {
	limit := 500

	if len(p) > limit {
		p = p[:limit]
	}

	curCount, err := pb.reader.Read(p)

	if curCount > 0 {
		pb.counter += int64(curCount)
		if pb.total > 0 {
			perc := float64(pb.counter) * 100 / float64(pb.total)
			fmt.Printf("\rLoading: %.2f%%", perc)

			time.Sleep(500 * time.Millisecond)
		}
	}

	if errors.Is(err, io.EOF) {
		fmt.Println()
	}

	return curCount, err
}

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

	pb := &ProgressBar{
		reader: src,
		total:  totalBytes,
	}

	if limit == 0 {
		_, err = io.Copy(dst, pb)
	} else {
		_, err = io.CopyN(dst, pb, limit)
	}

	return err
}
