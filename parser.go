package ltsv

import (
	"bytes"

	"golang.org/x/xerrors"
)

type (
	// Field is a struct to hold label-value pair.
	Field struct {
		Label string
		Value string
	}

	// Parser is for parsing LTSV-encoded format.
	Parser struct {
		// FieldDelimiter is the delimiter of fields. It defaults to '\t'.
		FieldDelimiter byte
		// ValueDelimiter is the delimiter of label-value pairs. It defaults to ':'.
		ValueDelimiter byte
		// StrictMode is a flag to check if labels and values are valid.
		// If strictMode is false,
		// the parser just split fields with `FieldDelimiter`
		// and split label and value with `ValueDelimiter` without checking if they are valid.
		// The valid label is `/[0-9A-Za-z_.-]+/`.
		// The valid value is `/[^\b\t\r\n]*/`.
		StrictMode bool
	}
)

// DefaultParser is the default parser
var DefaultParser = Parser{
	FieldDelimiter: '\t',
	ValueDelimiter: ':',
	StrictMode:     true,
}

var (
	// ErrMissingLabel is an error to describe label is missing (ex. 'my_value')
	ErrMissingLabel = xerrors.New("missing label")
	// ErrEmptyLabel is an error to describe label is empty (ex. ':my_value')
	ErrEmptyLabel = xerrors.New("empty label")
	// ErrInvalidLabel is an error to describe label contains invalid char (ex. 'my\tlabel:my_value')
	ErrInvalidLabel = xerrors.New("invalid label")
	// ErrInvalidValue is an error to describe value contains invalid char (ex. 'my_label:my_value\n')
	ErrInvalidValue = xerrors.New("invalid value")
	// Break is an error for break loop
	Break = xerrors.New("break")
)

// ParseField parse LTSV-encoded field and return the label and value.
// The result share same memory with inputted field.
func (p Parser) ParseField(field []byte) (label []byte, value []byte, err error) {
	idx := bytes.IndexByte(field, p.ValueDelimiter)
	if idx > 0 {
		label = field[0:idx]
		value = field[idx+1:]
		if p.StrictMode {
			if err = validateLabel(label); err != nil {
				return nil, nil, xerrors.Errorf("bad field label syntax %q: %w", string(field), err)
			}
			if err = validateValue(value); err != nil {
				return nil, nil, xerrors.Errorf("bad field value syntax %q: %w", string(field), err)
			}
		}
	} else {
		switch idx {
		case -1:
			err = xerrors.Errorf("bad field syntax %q: %w", string(field), ErrMissingLabel)
		case 0:
			err = xerrors.Errorf("bad field syntax %q: %w", string(field), ErrEmptyLabel)
		}
	}
	return
}

// ParseLine parse one line of LTSV-encoded data and call callback.
// The callback function will be called for each field.
func (p Parser) ParseLine(line []byte, callback func(label []byte, value []byte) error) error {
	oriLine := line
	for len(line) > 0 {
		idx := bytes.IndexByte(line, p.FieldDelimiter)
		var field []byte
		if idx == -1 {
			field = line
			line = nil
		} else {
			field = line[0:idx]
			line = line[idx+1:]
		}
		if len(field) == 0 {
			continue
		}
		label, value, err := p.ParseField(field)
		if err != nil {
			return xerrors.Errorf("bad line syntax %q: %w", string(oriLine), err)
		}

		if err = callback(label, value); err != nil {
			if err == Break {
				break
			}
			return xerrors.Errorf("ParseLine callback error: %w", err)
		}
	}
	return nil
}

// ParseLineAsMap parse one line of LTSV-encoded data and return the map[string]string.
// For reducing memory allocation, you can pass a map to record to reuse the given map.
func (p Parser) ParseLineAsMap(line []byte, record map[string]string) (map[string]string, error) {
	if record == nil {
		record = map[string]string{}
	}
	err := p.ParseLine(line, func(label []byte, value []byte) error {
		record[string(label)] = string(value)
		return nil
	})
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return record, nil
}

// ParseLineAsSlice parse one line of LTSV-encoded data and return the []Field.
// For reducing memory allocation, you can pass a slice to record to reuse the given slice.
func (p Parser) ParseLineAsSlice(line []byte, record []Field) ([]Field, error) {
	record = record[:0]
	err := p.ParseLine(line, func(label []byte, value []byte) error {
		record = append(record, Field{string(label), string(value)})
		return nil
	})
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return record, nil
}

func validateLabel(label []byte) error {
	for _, c := range label {
		if !isValidKey(c) {
			return xerrors.Errorf("invalid char %q used in label %q: %w", c, string(label), ErrInvalidLabel)
		}
	}
	return nil
}

func validateValue(value []byte) error {
	for _, c := range value {
		if !isValidValue(c) {
			return xerrors.Errorf("invalid char %q used in value %q: %w", c, string(value), ErrInvalidValue)
		}
	}
	return nil
}

func isValidKey(ch byte) bool { // [0-9A-Za-z_.-]
	switch ch {
	case '_', '.', '-',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
		return true
	}
	return false
}

func isValidValue(ch byte) bool {
	// %x01-08 / %x0B / %x0C / %x0E-FF
	switch ch {
	case '\b', '\t', '\r', '\n':
		return false
	}
	return true
}
