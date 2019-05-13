package ltsv

import (
	"bytes"

	"golang.org/x/xerrors"
)

var (
	// ErrMissingLabel is an error to describe label is missing (ex. 'my_value')
	ErrMissingLabel = xerrors.New("missing label")
	// ErrEmptyLabel is an error to describe label is empty (ex. ':my_value')
	ErrEmptyLabel = xerrors.New("empty label")
	// ErrInvalidLabel is an error to describe label contains invalid char (ex. 'my\tlabel:my_value')
	ErrInvalidLabel = xerrors.New("invalid label")
	// ErrInvalidValue is an error to describe value contains invalid char (ex. 'my_label:my_value\n')
	ErrInvalidValue = xerrors.New("invalid value")
)

// ParseLine parse LTSV-encoded data and return the result.
// If strictMode is false, it just split fields with '\t' and split label and value with ':' without checking if format is valid.
// For reducing memory allocation, you can pass a map to record to reuse the given map.
func ParseLine(line []byte, strictMode bool, record map[string]string) (map[string]string, error) {
	if record == nil {
		record = map[string]string{}
	}
	oriLine := line
	for len(line) > 0 {
		idx := bytes.IndexByte(line, '\t')
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
		label, value, err := parseField(field, strictMode)
		if err != nil {
			return nil, xerrors.Errorf("bad line syntax %q: %w", string(oriLine), err)
		}

		record[string(label)] = string(value)
	}
	return record, nil
}

func parseField(field []byte, strictMode bool) (label []byte, value []byte, err error) {
	idx := bytes.IndexByte(field, ':')
	if idx > 0 {
		label = field[0:idx]
		value = field[idx+1:]
		if strictMode {
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
