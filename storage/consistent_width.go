package storage

import (
	"io"
)

type Handle interface {
	io.Reader
	io.Seeker
	io.WriterAt
}

type ConsistentWidthHandler interface {
	GetLineWidth() int
	Read(buffer []byte) (int, int, error)
	Insert(data []byte) (lineNumber int, err error)
	Update(lineIndex int, data []byte) error
	Remove(lineIndex int) error
}

type ConsistentWidthHandlerOption interface {
	apply(*consistentWidthHandler)
}

type OptionLineIncrementWidth struct {
	value int
}

func (option OptionLineIncrementWidth) apply(cwh *consistentWidthHandler) {
	cwh.lineIncrementWidth = option.value
}

type consistentWidthHandler struct {
	handle             Handle
	initialLineCount   int64
	lineCount          int64
	lineWidth          int
	lineIncrementWidth int
	emptyLineIndexes   []int
}

func NewConsistentWidthHandler(
	handle Handle,
	options ...ConsistentWidthHandlerOption,
) (ConsistentWidthHandler, error) {
	var err error

	cwh := &consistentWidthHandler{
		handle:             handle,
		lineIncrementWidth: 500,
		emptyLineIndexes:   make([]int, 0, 10),
	}

	for _, option := range options {
		option.apply(cwh)
	}

	if err = cwh.findLineWidth(); err != nil && err == io.EOF {
		// TODO
		// this means no '\n' was found, it can mean either an empty file or
		// a file with no '\n' runes
	} else if err != nil {
		return nil, err
	}

	if err = cwh.findLineCount(); err != nil {
		return nil, err
	}

	if _, err = handle.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	return cwh, nil
}

func (cwh *consistentWidthHandler) findLineWidth() (err error) {
	var (
		bufferLen  = 1000
		handle     = cwh.handle
		i          = 0
		b          = make([]byte, bufferLen)
		foundIndex bool
	)

	if _, err = handle.Seek(0, io.SeekStart); err != nil {
		return err
	}

LOOP:
	for !foundIndex {
		var n int

		if n, err = cwh.handle.Read(b); err != nil && err != io.EOF {
			return err
		}

		for j := 0; j < n; j++ {
			if b[i] == '\n' {
				break LOOP
			}
			i++
		}

		if err != nil && err == io.EOF {
			return err
		}

		if n < bufferLen {
			return io.EOF
		}
	}

	cwh.lineWidth = i + 1

	return nil
}

func (cwh *consistentWidthHandler) findLineCount() (err error) {
	var (
		handle = cwh.handle
		size   int64
	)

	if size, err = handle.Seek(0, io.SeekEnd); err != nil {
		return err
	}

	cwh.initialLineCount = size / int64(cwh.lineWidth)
	cwh.lineCount = cwh.initialLineCount

	return nil
}

func (cwh *consistentWidthHandler) expand(needed int) error {
	return nil
}

func (cwh *consistentWidthHandler) Read(buffer []byte) (int, error) {
	return 0, nil
}

func (cwh *consistentWidthHandler) Update(lineIndex int, data []byte) error {
	return nil
}
