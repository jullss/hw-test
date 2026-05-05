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
	reader   io.Reader
	total    int64
	counter  int64
	lastPerc int
}

func (pb *ProgressBar) Read(p []byte) (int, error) {
	curCount, err := pb.reader.Read(p)

	if curCount > 0 {
		pb.counter += int64(curCount)
		if pb.total > 0 {
			perc := int(float64(pb.counter) * 100 / float64(pb.total))

			if perc > pb.lastPerc {
				fmt.Printf("\rLoading: %d%%", perc)

				time.Sleep(100 * time.Millisecond)

				pb.lastPerc = perc
			}
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
		return fmt.Errorf("%w: %w", ErrNoAccessFileInfo, err)
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
		return fmt.Errorf("%w: %w", ErrCannotOpenFile, err)
	}
	defer src.Close()

	dst, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCannotCreateFile, err)
	}
	defer dst.Close()

	_, err = src.Seek(offset, 0)
	if err != nil {
		return err
	}

	totalBytes := fileSize - offset
	if limit > 0 && limit < totalBytes {
		totalBytes = limit
	}

	pb := &ProgressBar{
		reader: src,
		total:  totalBytes,
	}

	if limit == 0 {
		_, err = io.Copy(dst, pb)
	} else {
		_, err = io.CopyN(dst, pb, totalBytes)
	}

	return err
}
