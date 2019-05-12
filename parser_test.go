package ltsv

import (
	"bufio"
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseLine(t *testing.T) {
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
		{"NL", "a:1\n", map[string]string{"a": "1"}, nil},
		{"leading tab", "\ta:1\n", map[string]string{"a": "1"}, nil},
		{"bad syntax", "a\tb", map[string]string{}, ErrInvalidLabel},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m, err := ParseLine([]byte(test.line), nil)
			if test.err != nil {
				assert.Equal(t, test.err, err)
				return
			}
			require.NoError(t, err)
			assert.EqualValues(t, test.want, m)
		})
	}
}

func TestParseField(t *testing.T) {
	tests := []struct {
		name  string
		line  string
		label string
		value string
		rest  string
		err   error
	}{
		{"simple", "abc:123", "abc", "123", "", nil},
		{"1 char label", "a:1", "a", "1", "", nil},
		{"empty_value", "abc:", "abc", "", "", nil},

		{"simple with tab", "abc:123\ta:1", "abc", "123", "\ta:1", nil},
		{"1 char label with tab", "a:123\ta:1", "a", "123", "\ta:1", nil},
		{"empty_value with tab", "abc:\ta:1", "abc", "", "\ta:1", nil},

		{"empty", "", "", "", "", ErrMissingLabel},
		{"bad_syntax", "abc", "", "", "", ErrMissingLabel},
		{"key has tab", "a\tc:123", "", "", "", ErrInvalidLabel},
		{"empty_label", ":123", "", ":123", "", ErrMissingLabel},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rest, label, value, err := parseField([]byte(test.line))
			if test.err != nil {
				assert.Equal(t, test.err, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.label, label)
			assert.Equal(t, test.value, value)
			assert.Equal(t, test.rest, string(rest))
		})
	}
}

func TestParseLabel(t *testing.T) {
	tests := []struct {
		name string
		line string
		want string
		rest string
		err  error
	}{
		{"simple", "abc:123", "abc", ":123", nil},
		{"1 char label", "a:123", "a", ":123", nil},
		{"empty_value", "abc:", "abc", ":", nil},
		{"empty", "", "", "", ErrMissingLabel},
		{"bad_syntax", "abc", "", "", ErrMissingLabel},
		{"key has tab", "a\tc:123", "", "", ErrInvalidLabel},
		{"empty_label", ":123", "", ":123", ErrMissingLabel},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rest, label, err := parseLabel([]byte(test.line))
			if test.err != nil {
				assert.Equal(t, test.err, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.want, label)
			assert.Equal(t, test.rest, string(rest))
		})
	}
}

func TestParseValue(t *testing.T) {
	tests := []struct {
		name string
		line string
		want string
		rest string
		err  error
	}{
		{"empty", "", "", "", nil},
		{"1 char", "a", "a", "", nil},
		{"multi-chars", "abc", "abc", "", nil},
		{"empty with rest", "\tdef", "", "\tdef", nil},
		{"1 char with rest", "a\tdef", "a", "\tdef", nil},
		{"multi-chars with rest", "abc\tdef", "abc", "\tdef", nil},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rest, label, err := parseValue([]byte(test.line))
			if test.err != nil {
				assert.Equal(t, test.err, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.want, label)
			assert.Equal(t, test.rest, string(rest))
		})
	}
}

func BenchmarkParseLine(b *testing.B) {
	line := []byte("host:127.0.0.1\tident:-\tuser:frank\ttime:[10/Oct/2000:13:55:36 -0700]\treq:GET /apache_pb.gif HTTP/1.0\tstatus:200\tsize:2326\treferer:http://www.example.com/start.html\tua:Mozilla/4.08 [en] (Win98; I ;Nav)")
	m := make(map[string]string, 11)
	for i := 0; i < b.N; i++ {
		ParseLine(line, m)
	}
}

func BenchmarkReadLine(b *testing.B) {
	const str = `aaaaaaaaaaaa
bbbbbbbbbbbbbbb
cccccccccccccccccccc
dddddddddddddd
eeeeeeeeeeeeeee
fffffff`
	buf := bytes.NewBufferString(str)
	b.Run("bufio.Scanner", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			buf.Reset()
			scanner := bufio.NewScanner(buf)
			for scanner.Scan() {
				scanner.Text()
			}
		}
	})
	b.Run("bufio.Reader", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			buf.Reset()
			reader := bufio.NewReader(buf)
			for {
				_, err := reader.ReadBytes('\n')
				if err != nil {
					if err == io.EOF {
						break
					}
					panic(err)
				}
			}
		}
	})
}
