package ltsv

import (
	"bufio"
	"io"

	"golang.org/x/xerrors"
)

// Reader is LTSV reader.
//
// As returned by NewReader, a Reader expects input LTSV-encoded stream.
type Reader struct {
	StrictMode bool
	r          *bufio.Scanner
	lineNo     int
}

// NewReader creates a new ltsv.Reader.
func NewReader(r io.Reader, strictMode bool) *Reader {
	return &Reader{
		StrictMode: strictMode,
		r:          bufio.NewScanner(r),
	}
}

// Read reads one record from r and store to record.
// If r is EOF, io.EOF will be returned.
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

// ReadAll reads all records.
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
