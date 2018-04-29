package tss

import (
	"bytes"
	"io"
	"testing"
	"time"
)

func BenchmarkCopy(b *testing.B) {
	forceNonZeroTestVal = 5 * time.Millisecond
	defer func() {
		forceNonZeroTestVal = 0
	}()
	bs := bytes.Repeat([]byte{'a'}, 2<<12+1)
	for i := 0; i < len(bs); i += 50 {
		bs[i] = '\n'
	}
	b.SetBytes(int64(len(bs)))
	rd := bytes.NewReader(bs)
	buf := new(bytes.Buffer)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		CopyTime(buf, rd, time.Now().Add(-50*time.Millisecond))
		buf.Reset()
		rd.Reset(bs)
	}
}

func BenchmarkWriter(b *testing.B) {
	forceNonZeroTestVal = 5 * time.Millisecond
	defer func() {
		forceNonZeroTestVal = 0
	}()
	bs := bytes.Repeat([]byte{'a'}, 2<<12)
	for i := 0; i < len(bs); i += 50 {
		bs[i] = '\n'
	}
	b.SetBytes(int64(len(bs)))
	rd := bytes.NewReader(bs)
	buf := new(bytes.Buffer)
	w := NewWriter(buf, time.Now().Add(-50*time.Millisecond))
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		io.Copy(w, rd)
		buf.Reset()
		rd.Reset(bs)
	}
}
