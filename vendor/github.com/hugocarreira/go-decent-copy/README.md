[![GoDoc](https://godoc.org/github.com/hugocarreira/go-decent-copy?status.svg)](https://godoc.org/github.com/hugocarreira/go-decent-copy) [![Build Status](https://travis-ci.org/hugocarreira/go-decent-copy.svg?branch=master)](https://travis-ci.org/hugocarreira/go-decent-copy) [![Go Report Card](https://goreportcard.com/badge/github.com/hugocarreira/go-decent-copy)](https://goreportcard.com/report/github.com/hugocarreira/go-decent-copy) [![codecov](https://codecov.io/gh/hugocarreira/go-decent-copy/branch/master/graph/badge.svg)](https://codecov.io/gh/hugocarreira/go-decent-copy)

# Go-Decent-Copy

go-decent-copy provides a copy file for humans

## Usage

```go
package main

import "github.com/hugocarreira/go-decent-copy"

func main() {
    execPath, _ := os.Getwd()

    fileOrigin := "/testCopy.txt"
    fileDestiny := fmt.Sprintf(execPath + "/folder/" + fileOrigin)
    err := Copy(fileOrigin, fileDestiny)
    if err != nil {
        fmt.Printf("Error in copy file : %#v ", err.Error())
    }
}
```

## License

This project is licensed under the [MIT License](LICENSE).
