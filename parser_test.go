package ltsv

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseLine(t *testing.T) {
	line := []byte("foo:123\tbar:456")
	ParseLine(line, func(label []byte, value []byte) error {
		val := string(value)
		switch string(label) {
		case "foo":
			assert.Equal(t, "123", val)
		case "bar":
			assert.Equal(t, "456", val)
		default:
			t.Errorf("unknown label: %s", string(label))
		}
		return nil
	})
}

func TestParseLine_break(t *testing.T) {
	line := []byte("foo:123\tbar:456")
	counter := 0
	err := ParseLine(line, func(label []byte, value []byte) error {
		counter++
		return Break
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, counter)
}

func TestParseLine_error(t *testing.T) {
	line := []byte("foo:123\tbar:456")
	counter := 0
	customErr := errors.New("custom error")
	err := ParseLine(line, func(label []byte, value []byte) error {
		counter++
		return customErr
	})
	assert.Error(t, err)
	assert.True(t, errors.Is(err, customErr))
	assert.Equal(t, 1, counter)
}

func TestParseLineAsMap(t *testing.T) {
	tests := []struct {
		name string
		line string
		want map[string]string
		err  error
	}{
		{"empty", "", map[string]string{}, nil},
		{"single", "a:1", map[string]string{"a": "1"}, nil},
		{"simple", "a:1\tb:2", map[string]string{"a": "1", "b": "2"}, nil},
		{"extra tab", "a:1\t\tb:2\t\t", map[string]string{"a": "1", "b": "2"}, nil},
		{"NL", "a:1\n", map[string]string{}, ErrInvalidValue},
		{"leading tab", "\ta:1", map[string]string{"a": "1"}, nil},
		{"no label", "a\tb", map[string]string{}, ErrMissingLabel},
		{"missing label", ":a", map[string]string{}, ErrEmptyLabel},
		{"bad label", "a\rb:1", map[string]string{}, ErrInvalidLabel},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m, err := ParseLineAsMap([]byte(test.line), nil)
			if test.err != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, test.err)
				return
			}
			require.NoError(t, err)
			assert.EqualValues(t, test.want, m)
		})
	}
}

func ExampleParseLineAsMap() {
	line := []byte("foo:123\tbar:456")
	record, err := ParseLineAsMap(line, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v", record) // map[string]string{"foo":"123", "bar":"456"}
}

func TestParseLineAsSlice(t *testing.T) {
	tests := []struct {
		name string
		line string
		want []Field
		err  error
	}{
		{"empty", "", nil, nil},
		{"single", "a:1", []Field{{"a", "1"}}, nil},
		{"simple", "a:1\tb:2", []Field{{"a", "1"}, {"b", "2"}}, nil},
		{"extra tab", "a:1\t\tb:2\t\t", []Field{{"a", "1"}, {"b", "2"}}, nil},
		{"NL", "a:1\n", nil, ErrInvalidValue},
		//{"leading tab", "\ta:1", map[string]string{"a": "1"}, nil},
		{"no label", "a\tb", nil, ErrMissingLabel},
		{"missing label", ":a", nil, ErrEmptyLabel},
		{"bad label", "a\rb:1", nil, ErrInvalidLabel},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m, err := ParseLineAsSlice([]byte(test.line), nil)
			if test.err != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, test.err)
				return
			}
			require.NoError(t, err)
			assert.EqualValues(t, test.want, m)
		})
	}
}

func ExampleParseLineAsSlice() {
	line := []byte("foo:123\tbar:456")
	record, err := ParseLineAsSlice(line, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", record) // [{Label:foo Value:123} {Label:bar Value:456}]
}

func TestParseField(t *testing.T) {
	tests := []struct {
		name  string
		line  string
		label string
		value string
		err   error
	}{
		{"simple", "abc:123", "abc", "123", nil},
		{"1 char label", "a:1", "a", "1", nil},
		{"empty_value", "abc:", "abc", "", nil},

		{"bad_syntax", "abc", "", "", ErrMissingLabel},
		{"key has tab", "a\tc:123", "", "", ErrInvalidLabel},
		{"empty_label", ":123", "", "", ErrEmptyLabel},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			label, value, err := ParseField([]byte(test.line))
			if test.err != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, test.err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.label, string(label))
			assert.Equal(t, test.value, string(value))
		})
	}
}

func BenchmarkParseLineAsMap(b *testing.B) {
	parser := DefaultParser
	parser.StrictMode = false
	line := []byte("host:127.0.0.1\tident:-\tuser:frank\ttime:[10/Oct/2000:13:55:36 -0700]\treq:GET /apache_pb.gif HTTP/1.0\tstatus:200\tsize:2326\treferer:http://www.example.com/start.html\tua:Mozilla/4.08 [en] (Win98; I ;Nav)")
	m := make(map[string]string, 11)
	for i := 0; i < b.N; i++ {
		parser.ParseLineAsMap(line, m)
	}
}

func BenchmarkParseLine(b *testing.B) {
	parser := DefaultParser
	parser.StrictMode = false
	line := []byte("host:127.0.0.1\tident:-\tuser:frank\ttime:[10/Oct/2000:13:55:36 -0700]\treq:GET /apache_pb.gif HTTP/1.0\tstatus:200\tsize:2326\treferer:http://www.example.com/start.html\tua:Mozilla/4.08 [en] (Win98; I ;Nav)")
	for i := 0; i < b.N; i++ {
		parser.ParseLine(line, func(label []byte, value []byte) error {
			return nil
		})
	}
}
