package ltsv

func ParseLine(line []byte, m map[string]string) (map[string]string, error) {
	if m == nil {
		m = map[string]string{}
	}
	line = skipTabs(line)
	for len(line) > 0 {
		rest, label, value, err := parseField(line)
		line = skipTabs(rest)
		if err != nil {
			return nil, err
		}
		m[label] = value
	}
	return m, nil
}

func skipTabs(line []byte) []byte {
	for len(line) > 0 {
		switch line[0] {
		case '\t':
			line = line[1:]
		case '\r', '\n':
			return nil
		default:
			return line
		}
	}
	return nil
}

func parseField(line []byte) (rest []byte, label string, value string, err error) {
	line, label, err = parseLabel(line)
	if err != nil {
		return
	}
	rest, value, err = parseValue(line[1:])
	return
}

func parseLabel(line []byte) (rest []byte, label string, err error) {
	for i, c := range line {
		if isValidKey(c) {
			continue
		}

		if c != ':' {
			err = ErrInvalidLabel
			return
		}
		if i > 0 {
			label = string(line[0:i])
			rest = line[i:]
			return
		}
		err = ErrMissingLabel
		return
	}
	err = ErrMissingLabel
	return
}

func parseValue(line []byte) (rest []byte, value string, err error) {
	for i, c := range line {
		if isValidValue(c) {
			continue
		}
		value = string(line[0:i])
		rest = line[i:]
		return
	}
	value = string(line)
	return
}

//
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
