package tss

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"
)

// TimeScaler returns a format string for the given Duration where all of the
// decimals will line up in the same column (fourth from the end).
func TimeScaler(d time.Duration) string {
	switch {
	case d == 0:
		return "0.0ms"
	case d >= time.Second:
		return fmt.Sprintf("%.2fs", float64(d.Nanoseconds())/1e9)
	case d >= 50*time.Microsecond:
		return fmt.Sprintf("%.1fms", float64(d.Nanoseconds())/1e9*1000)
	case d >= time.Microsecond:
		return fmt.Sprintf("%.1fµs", float64(d.Nanoseconds())/1e9*1000*1000)
	default:
		return fmt.Sprintf("%.1fns", float64(d.Nanoseconds()))
	}
}

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
		sinceLastLine := gotLine.Sub(lastLine).Round(100 * time.Microsecond)
		sinceStart := gotLine.Sub(start).Round(100 * time.Microsecond)
		fmt.Fprintf(&buf, "%8s ", TimeScaler(sinceStart))
		if lastLine.IsZero() {
			buf.WriteString(strings.Repeat(" ", 9))
		} else {
			fmt.Fprintf(&buf, "%8s ", TimeScaler(sinceLastLine))
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
