package ltsv

import (
	"testing"

	"golang.org/x/xerrors"

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
		{"NL", "a:1\n", map[string]string{}, ErrInvalidValue},
		{"leading tab", "\ta:1", map[string]string{"a": "1"}, nil},
		{"no label", "a\tb", map[string]string{}, ErrMissingLabel},
		{"missing label", ":a", map[string]string{}, ErrEmptyLabel},
		{"bad label", "a\rb:1", map[string]string{}, ErrInvalidLabel},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m, err := ParseLine([]byte(test.line), true, nil)
			if test.err != nil {
				assert.Error(t, err)
				assert.True(t, xerrors.Is(err, test.err))
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
			label, value, err := parseField([]byte(test.line), true)
			if test.err != nil {
				assert.Error(t, err)
				assert.Truef(t, xerrors.Is(err, test.err), "%+v", err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.label, string(label))
			assert.Equal(t, test.value, string(value))
		})
	}
}

func BenchmarkParseLine(b *testing.B) {
	line := []byte("host:127.0.0.1\tident:-\tuser:frank\ttime:[10/Oct/2000:13:55:36 -0700]\treq:GET /apache_pb.gif HTTP/1.0\tstatus:200\tsize:2326\treferer:http://www.example.com/start.html\tua:Mozilla/4.08 [en] (Win98; I ;Nav)")
	m := make(map[string]string, 11)
	for i := 0; i < b.N; i++ {
		ParseLine(line, false, m)
	}
}
