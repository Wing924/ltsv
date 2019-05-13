# ltsv

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Build Status](https://travis-ci.com/Wing924/ltsv.svg?branch=master)](https://travis-ci.com/Wing924/ltsv)
[![Go Report Card](https://goreportcard.com/badge/github.com/Wing924/ltsv)](https://goreportcard.com/report/github.com/Wing924/ltsv)
[![codecov](https://codecov.io/gh/Wing924/ltsv/branch/master/graph/badge.svg)](https://codecov.io/gh/Wing924/ltsv)
[![GoDoc](https://godoc.org/github.com/Wing924/ltsv?status.svg)](https://godoc.org/github.com/Wing924/ltsv)

High performance LTSV (Labeled Tab Separeted Value) reader for Go.

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

## Benchmarks

```bash
$ cd bench && go test -bench . -benchmem
goos: darwin
goarch: amd64
pkg: github.com/Wing924/ltsv/bench
Benchmark_line_Wing924_ltsv-4          	 2000000	       588 ns/op	     224 B/op	      17 allocs/op
Benchmark_line_Wing924_ltsv_strict-4   	 2000000	       800 ns/op	     224 B/op	      17 allocs/op
Benchmark_line_Songmu_goltsv-4         	  300000	      3906 ns/op	    1841 B/op	      32 allocs/op
Benchmark_line_ymotongpoo_goltsv-4     	 1000000	      2324 ns/op	    5793 B/op	      17 allocs/op
Benchmark_line_najeira_ltsv-4          	  300000	      4933 ns/op	    5529 B/op	      26 allocs/op
PASS
ok  	github.com/Wing924/ltsv/bench	9.315s
```
