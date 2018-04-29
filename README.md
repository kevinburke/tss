# tss (1)

`tss` is like `ts` from moreutils, but prints relative durations (with
millisecond precision) by default, and can be shipped as a compiled binary.

Try it out:

```
$ (sleep 1; echo "hello"; sleep 2; echo "two sec") | tss
   995ms          hello
      3s   2.005s two sec
```

The first column is the amount of time that has elapsed since the program
started. The second column is the amount of time that has elapsed since the last
line printed.

## Installation

[Find your target operating system](https://github.com/kevinburke/tss/releases) (darwin, windows, linux) and desired bin
directory, and modify the command below as appropriate:

    curl --silent --location --output /usr/local/bin/tss https://github.com/kevinburke/tss/releases/download/0.4/tss-linux-amd64 && chmod 755 /usr/local/bin/tss

The latest version is 0.4.

If you have a Go development environment, you can also install via source code:

    go get -u github.com/kevinburke/tss

The corresponding library is available at
`github.com/kevinburke/tss/lib`. View the library documentation at
[godoc.org](https://godoc.org/github.com/kevinburke/tss/lib).

## Usage Notes

- Piping commands to `tss` may result in programs buffering their output before
flushing it to stdout file descriptor. You can avoid this by wrapping the target
program in a command like `unbuffer` (via [the expect package][expect]) or
[`stdbuf` from the coreutils package][stdbuf]. On Macs you can install with
`brew install expect` and `brew install coreutils` respectively; the stdbuf
command may be prefixed with a 'g': `gstdbuf`.

    ```
    stdbuf --output=L <mycommand> | tss
    ```

- Piping commands may also change the return code from non-zero to zero, since
Bash by default uses the return code of the last command in the pipe to decide
how to exit. This means if you are piping output to `tss` or `ts` you may
accidentally change a failing program to a passing one. Use `set -o pipefail` in
Bash scripts to ensure that Bash will return a non-zero return code if any part
of a pipe operation fails. Or add this to a Makefile:

    ```
    SHELL = /bin/bash -o pipefail
    ```


[expect]: http://expect.sourceforge.net/example/unbuffer.man.html
[stdbuf]: https://www.gnu.org/software/coreutils/manual/html_node/stdbuf-invocation.html
