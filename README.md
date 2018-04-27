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

    curl --silent --location --output=/usr/local/bin/tss https://github.com/kevinburke/tss/releases/download/0.3/tss-linux-amd64 && chmod 755 /usr/local/bin/tss

The latest version is 0.3.

If you have a Go development environment, you can also install via source code:

    go get -u github.com/kevinburke/tss

The corresponding library is available at
`github.com/kevinburke/tss/lib`. View the library documentation at
[godoc.org](https://godoc.org/github.com/kevinburke/tss/lib).
