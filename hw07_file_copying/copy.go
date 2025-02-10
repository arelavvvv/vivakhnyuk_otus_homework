package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	//nolint:depguard
	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrSameFile              = errors.New("source and destination are the same file")
	ErrUnknownFileSize       = errors.New("cannot use progress bar with file of unknown size (e.g., /dev/urandom)")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	absFromPath, err := filepath.Abs(fromPath)
	if err != nil {
		return fmt.Errorf("ошибка при получении абсолютного пути для исходного файла: %w", err)
	}

	absToPath, err := filepath.Abs(toPath)
	if err != nil {
		return fmt.Errorf("ошибка при получении абсолютного пути для целевого файла: %w", err)
	}

	if absFromPath == absToPath {
		return ErrSameFile
	}

	srcFile, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("ошибка при открытии исходного файла: %w", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("ошибка при создании целевого файла: %w", err)
	}
	defer dstFile.Close()

	srcInfo, err := srcFile.Stat()
	if !srcInfo.Mode().IsRegular() { // e.g. /dev/urandom
		return ErrUnknownFileSize
	}

	if offset > srcInfo.Size() {
		return fmt.Errorf("offset больше, чем размер файла: %d > %d", offset, srcInfo.Size())
	}

	var bytesToCopy int64
	remainingBytes := srcInfo.Size() - offset

	if limit == 0 {
		bytesToCopy = remainingBytes
	} else {
		if limit > remainingBytes {
			bytesToCopy = remainingBytes
		} else {
			bytesToCopy = limit
		}
	}

	_, err = srcFile.Seek(offset, io.SeekStart)
	if err != nil {
		return fmt.Errorf("ошибка при перемещении указателя файла: %w", err)
	}

	bar := pb.Full.Start64(bytesToCopy)
	barReader := bar.NewProxyReader(srcFile)

	_, err = io.CopyN(dstFile, barReader, bytesToCopy)
	if err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("ошибка при копировании данных: %w", err)
	}

	bar.Finish()

	return nil
}
