package bench

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	Songmu "github.com/Songmu/go-ltsv"
	Wing924 "github.com/Wing924/ltsv"
	najeira "github.com/najeira/ltsv"
	ymotongpoo "github.com/ymotongpoo/goltsv"
)

var line = []byte("host:127.0.0.1	ident:-	user:frank	time:[10/Oct/2000:13:55:36 -0700]	req:GET /apache_pb3.gif HTTP/1.0	status:200	size:2326	referer:http://www.example.com/start.html	ua:Mozilla/4.08 [en] (Win98; I ;Nav)")
var value = map[string]string{
	"host":    "127.0.0.1",
	"ident":   "-",
	"user":    "frank",
	"time":    "[10/Oct/2000:13:55:36 -0700]",
	"req":     "GET /apache_pb3.gif HTTP/1.0",
	"status":  "200",
	"size":    "2326",
	"referer": "http://www.example.com/start.html",
	"ua":      "Mozilla/4.08 [en] (Win98; I ;Nav)",
}

// Wing924/ltsv
func Test_line_Wing924_ltsv(t *testing.T) {
	parser := Wing924.DefaultParser
	parser.StrictMode = false
	m, err := parser.ParseLineAsMap(line, nil)
	assert.NoError(t, err)
	assert.EqualValues(t, value, m)
}

func Test_line_Wing924_ltsv_strict(t *testing.T) {
	parser := Wing924.DefaultParser
	m, err := parser.ParseLineAsMap(line, nil)
	assert.NoError(t, err)
	assert.EqualValues(t, value, m)
}

func Benchmark_line_Wing924_ltsv(b *testing.B) {
	parser := Wing924.DefaultParser
	parser.StrictMode = false
	m := make(map[string]string, 17)
	for i := 0; i < b.N; i++ {
		parser.ParseLineAsMap(line, m)
	}
}

func Benchmark_line_Wing924_ltsv_strict(b *testing.B) {
	parser := Wing924.DefaultParser
	m := make(map[string]string, 17)
	for i := 0; i < b.N; i++ {
		parser.ParseLineAsMap(line, m)
	}
}

// Songmu/go-ltsv
func Test_line_Songmu_goltsv(t *testing.T) {
	m := make(map[string]string, 17)
	err := Songmu.Unmarshal(line, &m)
	assert.NoError(t, err)
	assert.EqualValues(t, value, m)
}

func Benchmark_line_Songmu_goltsv(b *testing.B) {
	m := make(map[string]string, 17)
	for i := 0; i < b.N; i++ {
		Songmu.Unmarshal(line, &m)
	}
}

// ymotongpoo/goltsv
func Test_line_ymotongpoo_goltsv(t *testing.T) {
	buf := bytes.NewBuffer(line)
	m, err := ymotongpoo.NewReader(buf).Read()
	assert.NoError(t, err)
	assert.EqualValues(t, value, m)
}

func Benchmark_line_ymotongpoo_goltsv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := bytes.NewBuffer(line)
		ymotongpoo.NewReader(buf).Read()
	}
}

// najeira/ltsv
func Test_line_najeira_ltsv(t *testing.T) {
	buf := bytes.NewBuffer(line)
	m, err := najeira.NewReader(buf).Read()
	assert.NoError(t, err)
	assert.EqualValues(t, value, m)
}

func Benchmark_line_najeira_ltsv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := bytes.NewBuffer(line)
		najeira.NewReader(buf).Read()
	}
}
