// -*- coding:utf-8; indent-tabs-mode:nil; -*-

// Copyright 2014, Wu Xi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//
// Normal example:

package main

import (
    "fmt"
    "os"
    "os/exec"
    "os/user"
)

import "github.com/wuxicn/pipeline"

func main() {

    u, err := user.Current()
    if err != nil {
        fmt.Println("get current user failed: %v", err)
        os.Exit(255)
    }

    stdout, stderr, err := pipeline.Run(
        exec.Command("ls", "-alh", u.HomeDir), // list files
        exec.Command("./echo_with_stderr.sh"), // add dummy stderr
        exec.Command("tr", "a-z", "A-Z"),      // to upper-case
        exec.Command("nl"))                    // add line number

    fmt.Println("STDOUT:")
    fmt.Println(stdout.String())
    fmt.Println("--------------")
    fmt.Println("STDERR:")
    fmt.Println(stderr.String())
    fmt.Println("--------------")
    fmt.Println("ERR:", err)
}
