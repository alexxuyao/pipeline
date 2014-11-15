// -*- coding:utf-8; indent-tabs-mode:nil; -*-

// Copyright 2014, Wu Xi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//
// With user input:
// usage: type any text and Ctrl-D

package main

import (
    "fmt"
    "os"
    "os/exec"
)

import "github.com/wuxicn/pipeline"

func main() {

    cmd0 := exec.Command("cat")
    cmd0.Stdin = os.Stdin
    stdout, _, err := pipeline.Run(cmd0, exec.Command("tr", "a-z", "A-Z"))

    fmt.Println("STDOUT:")
    fmt.Println(stdout.String())

    if err != nil {
        e := err.(*pipeline.Error)
        fmt.Println("ERR:", e.Code, e.Err)
    }
}
