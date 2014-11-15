// -*- coding:utf-8; indent-tabs-mode:nil; -*-

// Copyright 2014, Wu Xi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//
// Get error code:

package main

import (
    "fmt"
    "os"
    "os/exec"
    "os/user"
)

import "github.com/wuxicn/pipeline"

func main() {

    u, e := user.Current()
    if e != nil {
        fmt.Println("get current user failed: %v", e)
        os.Exit(255)
    }

    _, stderr, err := pipeline.Run(
        exec.Command("ls", "-alh", u.HomeDir), // list files
        exec.Command("tr", "a-z", "A-Z"),      // to upper-case
        exec.Command("./exit_non_zero.sh"),    // exit with non-zero
        exec.Command("nl"))                    // add line number

    fmt.Println("STDERR:")
    fmt.Println(stderr.String())

    if err != nil {
        e := err.(*pipeline.Error)
        fmt.Println("ERR:", e.Code, e.Err)
    }
}
