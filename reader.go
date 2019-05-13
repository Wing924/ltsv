package ltsv

import (
	"bufio"
	"io"

	"golang.org/x/xerrors"
)

type Reader struct {
	StrictMode bool
	r          *bufio.Scanner
	lineNo     int
}

func NewReader(r io.Reader) *Reader {
	return &Reader{
		r: bufio.NewScanner(r),
	}
}

func (r *Reader) Read(record map[string]string) (map[string]string, error) {
	if !r.r.Scan() {
		err := r.r.Err()
		if err == nil {
			return nil, io.EOF
		}
		return record, xerrors.Errorf("I/O error: %w", err)
	}
	r.lineNo++
	line := r.r.Bytes()
	record, err := ParseLine(line, r.StrictMode, record)
	if err != nil {
		return nil, xerrors.Errorf("bad syntax at line %d: %w", r.lineNo, err)
	}
	return record, nil
}

func (r *Reader) ReadAll() ([]map[string]string, error) {
	var records []map[string]string
	for r.r.Scan() {
		r.lineNo++
		record, err := ParseLine(r.r.Bytes(), r.StrictMode, nil)
		if err != nil {
			return records, xerrors.Errorf("bad syntax at line %d: %w", r.lineNo, err)
		}
		records = append(records, record)
	}
	return records, nil
}
