package ltsv

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReader_Read(t *testing.T) {
	f := newTestData()
	defer f.Close()

	r := NewReader(f, true)

	m, err := r.Read(nil)
	assert.NoError(t, err)
	assert.EqualValues(t, expectRecord(0), m)

	m, err = r.Read(m)
	assert.NoError(t, err)
	assert.EqualValues(t, expectRecord(1), m)

	m, err = r.Read(m)
	assert.NoError(t, err)
	assert.EqualValues(t, expectRecord(2), m)

	_, err = r.Read(m)
	assert.Equal(t, io.EOF, err)
}

func TestReader_ReadAll(t *testing.T) {
	f := newTestData()
	defer f.Close()

	r := NewReader(f, true)

	expected := []map[string]string{
		expectRecord(0),
		expectRecord(1),
		expectRecord(2),
	}
	actual, err := r.ReadAll()
	assert.NoError(t, err)
	assert.EqualValues(t, expected, actual)
}

func newTestData() *os.File {
	f, err := os.Open("testdata/access_log.txt")
	if err != nil {
		panic(err)
	}
	return f
}

func expectRecord(row int) map[string]string {
	// host:127.0.0.1	ident:-	user:frank	time:[10/Oct/2000:13:55:36 -0700]	req:GET /apache_pb1.gif HTTP/1.0	status:200	size:2326	referer:http://www.example.com/start.html	ua:Mozilla/4.08 [en] (Win98; I ;Nav)
	return map[string]string{
		"host":    "127.0.0.1",
		"ident":   "-",
		"user":    "frank",
		"time":    "[10/Oct/2000:13:55:36 -0700]",
		"req":     fmt.Sprint("GET /apache_pb", row+1, ".gif HTTP/1.0"),
		"status":  "200",
		"size":    "2326",
		"referer": "http://www.example.com/start.html",
		"ua":      "Mozilla/4.08 [en] (Win98; I ;Nav)",
	}
}
