# Builder
[![build status](https://secure.travis-ci.org/demaggus83/go-builder.svg?branch=master)](http://travis-ci.org/demaggus83/go-builder) 
[![GoDoc](https://godoc.org/github.com/demaggus83/go-builder?status.png)](http://godoc.org/github.com/demaggus83/go-builder) 
[![Go Report Card](https://goreportcard.com/badge/github.com/demaggus83/go-builder)](https://goreportcard.com/report/github.com/demaggus83/go-builder) 
[![Sourcegraph Badge](https://sourcegraph.com/github.com/demaggus83/go-builder/-/badge.svg)](https://sourcegraph.com/github.com/demaggus83/go-builder?badge)

    Notice: cgo i686 builds are untested

## About

Builder is a simple CLI tool I use to cross compile my go applications. \
It supports cgo, produces static linked binaries with stripped debugging information to reduce the file size, generate an upx version of the binary to reduce the size even more and inject some build meta data into the "main" package. 

## Usage

```
./builder help

./builder _example/cmd/example/main.go [...flags]
```

## Build Meta Data

Builder will inject some meta data into the following vars in the main package.

Example: \
```./builder _example/cmd/example/main.go -git 2190aca43 -version 0.1.0```

```go
package main

import (
	"fmt"
)

var (
    // time.Now() at build start formatted with time.RFC1123
	BuilderBuilt   string 
    // runtime.Version()
	BuilderGo      string
    // upx --version if available
	BuilderUPX     string

    // cli arg "-git" 
	BuilderGit     string
    // cli arg "-version"
	BuilderVersion string 
)

func main() {
	fmt.Println(BuilderBuilt)   // Tue, 23 Jul 2019 08:15:24 UTC
	fmt.Println(BuilderGo)      // go1.12.7
    fmt.Println(BuilderGit)     // 2190aca43
	fmt.Println(BuilderUPX)     // 3.95
	fmt.Println(BuilderVersion) // 0.1.0
}
```

### Changelog

#### 0.1.0
+ init
