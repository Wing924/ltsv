package ltsv

import (
	"bufio"
	"io"
)

type Reader struct {
	r *bufio.Scanner
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
			err = io.EOF
		}
		return record, err
	}
	line := r.r.Bytes()
	return ParseLine(line, record)
}

func (r *Reader) ReadAll() ([]map[string]string, error) {
	var records []map[string]string
	for r.r.Scan() {
		record, err := ParseLine(r.r.Bytes(), nil)
		if err != nil {
			return records, err
		}
		records = append(records, record)
	}
	return records, nil
}
