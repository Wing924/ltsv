package ltsv

import (
	"bytes"

	"golang.org/x/xerrors"
)

var (
	ErrMissingLabel = xerrors.New("missing label")
	ErrEmptyLabel   = xerrors.New("empty label")
	ErrInvalidLabel = xerrors.New("invalid label")
	ErrInvalidValue = xerrors.New("invalid value")
)

func ParseLine(line []byte, strictMode bool, m map[string]string) (map[string]string, error) {
	if m == nil {
		m = map[string]string{}
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

		m[string(label)] = string(value)
	}
	return m, nil
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
