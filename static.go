package ltsv

// ParseField parse LTSV-encoded field and return the label and value.
// The result share same memory with inputted field
func ParseField(field []byte) (label []byte, value []byte, err error) {
	return DefaultParser.ParseField(field)
}

// ParseLine parse one line of LTSV-encoded data and call callback.
// The callback function will be called for each field.
func ParseLine(line []byte, callback func(label []byte, value []byte)) error {
	return DefaultParser.ParseLine(line, callback)
}

// ParseLineAsMap parse one line of LTSV-encoded data and return the map[string]string.
// For reducing memory allocation, you can pass a map to record to reuse the given map.
func ParseLineAsMap(line []byte, record map[string]string) (map[string]string, error) {
	return DefaultParser.ParseLineAsMap(line, record)
}

// ParseLineAsSlice parse one line of LTSV-encoded data and return the []Field.
// For reducing memory allocation, you can pass a slice to record to reuse the given slice.
func ParseLineAsSlice(line []byte, record []Field) ([]Field, error) {
	return DefaultParser.ParseLineAsSlice(line, record)
}
