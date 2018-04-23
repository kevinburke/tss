package tss

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"
)

func Copy(w io.Writer, r io.Reader) (written int64, err error) {
	return CopyTime(w, r, time.Now())
}

func CopyTime(w io.Writer, r io.Reader, start time.Time) (written int64, err error) {
	bs := bufio.NewScanner(r)
	n := int64(0)
	var lastLine time.Time
	var buf bytes.Buffer
	for bs.Scan() {
		gotLine := time.Now()
		sinceLastLine := gotLine.Sub(lastLine).Round(time.Millisecond)
		sinceStart := gotLine.Sub(start).Round(time.Millisecond)
		fmt.Fprintf(&buf, "%8s ", sinceStart.String())
		if lastLine.IsZero() {
			buf.WriteString(strings.Repeat(" ", 9))
		} else {
			fmt.Fprintf(&buf, "%8s ", sinceLastLine.String())
		}
		buf.Write(bs.Bytes())
		buf.WriteByte('\n')
		wn, err := w.Write(buf.Bytes())
		n += int64(wn)
		if err != nil {
			return n, err
		}
		buf.Reset()
		lastLine = gotLine
	}
	return n, bs.Err()
}
