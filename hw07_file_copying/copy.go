package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrNegativeOffset        = errors.New("offset is negative")
	ErrNegativeLimit         = errors.New("limit is negative")
)

// Copy скопировать limit байтов из файла fromPath в файл toPath, начиная с offset байт.
func Copy(fromPath, toPath string, offset, limit int64) error {
	if offset < 0 {
		return ErrNegativeOffset
	}
	if limit < 0 {
		return ErrNegativeLimit
	}
	f, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer f.Close()
	fInfo, err := f.Stat()
	if err != nil {
		return err
	}
	if offset > fInfo.Size() {
		return ErrOffsetExceedsFileSize
	}
	if !fInfo.Mode().IsRegular() {
		return ErrUnsupportedFile
	}
	// продвинуться внутри файла на offset байт, если необходимо
	if offset > 0 {
		_, err = f.Seek(offset, io.SeekStart)
		if err != nil {
			return err
		}
	}
	newF, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer newF.Close()
	// вычитать весь файл
	if limit == 0 || limit > fInfo.Size()-offset {
		limit = fInfo.Size() - offset
	}
	// прогресс бар в консоль
	bar := pb.Full.Start64(limit)
	proxyReader := bar.NewProxyReader(f)
	_, err = io.CopyN(newF, proxyReader, limit)
	if err != nil {
		return err
	}
	bar.Finish()
	return nil
}
