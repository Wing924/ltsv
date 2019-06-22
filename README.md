# ltsv

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Build Status](https://travis-ci.com/Wing924/ltsv.svg?branch=master)](https://travis-ci.com/Wing924/ltsv)
[![Go Report Card](https://goreportcard.com/badge/github.com/Wing924/ltsv)](https://goreportcard.com/report/github.com/Wing924/ltsv)
[![codecov](https://codecov.io/gh/Wing924/ltsv/branch/master/graph/badge.svg)](https://codecov.io/gh/Wing924/ltsv)
[![GoDoc](https://godoc.org/github.com/Wing924/ltsv?status.svg)](https://godoc.org/github.com/Wing924/ltsv)

High performance LTSV (Labeled Tab Separeted Value) parser for Go.

About LTSV: http://ltsv.org/

	Labeled Tab-separated Values (LTSV) format is a variant of 
	Tab-separated Values (TSV). Each record in a LTSV file is represented 
	as a single line. Each field is separated by TAB and has a label and
	 a value. The label and the value have been separated by ':'. With 
	the LTSV format, you can parse each line by spliting with TAB (like 
	original TSV format) easily, and extend any fields with unique labels 
	in no particular order.

## Installation

```bash
go get github.com/Wing924/ltsv
```

## Examples

```go
package main

import (
	"fmt"
	"github.com/Wing924/ltsv"
)

func main() {
	line := []byte("foo:123\tbar:456")
    record, err := ltsv.ParseLineAsMap(line, nil)
    if err != nil {
        panic(err)
    }
    fmt.Printf("%#v", record) // map[string]string{"foo":"123", "bar":"456"}
}
```

## Benchmarks

Benchmark against 
* [Songmu/go-ltsv](https://github.com/Songmu/go-ltsv): 635% faster
* [ymotongpoo/goltsv](https://github.com/ymotongpoo/goltsv): 365% faster
* [najeira/ltsv](https://github.com/najeira/ltsv): 782% faster

Source code: [bench/line_test.go](https://github.com/Wing924/ltsv/blob/master/bench/line_test.go).

### Result

```
$ go test -bench . -benchmem
goos: darwin
goarch: amd64
pkg: github.com/Wing924/ltsv/bench
Benchmark_line_Wing924_ltsv-4          	 2000000	       626 ns/op	     224 B/op	      17 allocs/op
Benchmark_line_Wing924_ltsv_strict-4   	 2000000	       788 ns/op	     224 B/op	      17 allocs/op
Benchmark_line_Songmu_goltsv-4         	  300000	      3975 ns/op	    1841 B/op	      32 allocs/op
Benchmark_line_ymotongpoo_goltsv-4     	  500000	      2286 ns/op	    5793 B/op	      17 allocs/op
Benchmark_line_najeira_ltsv-4          	  300000	      4896 ns/op	    5529 B/op	      26 allocs/op
PASS
ok  	github.com/Wing924/ltsv/bench	8.245s
```
