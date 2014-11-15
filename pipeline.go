//-*- coding:utf-8; indent-tabs-mode:nil; -*-

// Copyright 2014, Wu Xi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pipeline

import (
    "bytes"
    "errors"
    "fmt"
    "os/exec"
    "syscall"
)

/* Pipeline strings together the given exec.Cmd commands in a similar fashion
   to the Unix pipeline.  Each command's standard output is connected to the
   standard input of the next command, and the output of the final command in
   the pipeline is returned, along with the collected standard error of all
   commands and the first error found (if any).

   To provide input to the pipeline, assign an io.Reader to the first's Stdin.

   Pipeline exit code:
    if all cmds start ok and exit zero, finErr=nil, otherwise finErr equals
    the value of the last (rightmost) command to exit with a non-zero status.
    like Bash set -o pipefail.

   Examples:
    1:
        stdout, err := Pipeline(&os.Stderr, exec.Command("ls", "-alh"))

    2:
        var stderr bytes.Buffer
        stdout, err := Pipeline(&stderr, exec.Command("ls", "-alh"),
            exec.Command("cat"))

    3:
        stderr, err := os.Create("./log/file.stderr")
        stdout, err := Pipeline(&stderr, exec.Command("ls", "-alh"))
*/
func Run(cmds ...*exec.Cmd) (stdout, stderr bytes.Buffer, finErr error) {
    // Require at least one command
    if len(cmds) < 1 {
        finErr = nil
        return
    }

    finErr = nil

    last := len(cmds) - 1
    for i, cmd := range cmds[:last] {
        var err error
        // Connect each command's stdin to the previous command's stdout
        if cmds[i+1].Stdin, err = cmd.StdoutPipe(); err != nil {
            finErr = err
            return
        }
        // Connect each command's stderr to a buffer
        cmd.Stderr = &stderr
    }

    // Connect the output and error for the last command
    cmds[last].Stdout, cmds[last].Stderr = &stdout, &stderr

    // Start each command
    for i, cmd := range cmds {
        if err := cmd.Start(); err != nil {
            finErr = errors.New(fmt.Sprintf("start cmd[%v] failed: %v", i, err))
            return
        }
    }

    // Wait for each command(in reverse order) to complete
    for i := last; i >= 0; i -= 1 {
        cmd := &cmds[i]
        if err := cmd.Wait(); err != nil {
            if exiterr, ok := err.(*exec.ExitError); ok {
                // The program has exited with an exit code != 0
                if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
                    finErr = errors.New(
                        fmt.Sprintf("cmd[%v] exits %v", i, status.ExitStatus()))
                    return
                } else {
                    finErr = errors.New(
                        fmt.Sprintf("get exit_code failed for cmd[%v]", i))
                    return
                }
            }
            finErr = errors.New(fmt.Sprintf("cmd[%v] failed: %v", i, err))
            return
        }
    } // for

    // Return the pipeline output and the collected standard error
    return
}
