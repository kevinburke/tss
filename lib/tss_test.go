package tss_test

import (
	"bytes"
	"testing"

	tss "github.com/kevinburke/tss/lib"
)

func BenchmarkCopy(b *testing.B) {
	bs := bytes.Repeat([]byte{'a'}, 512+1)
	for i := 0; i < len(bs); i += 40 {
		bs[i] = '\n'
	}
	rd := bytes.NewReader(bs)
	buf := new(bytes.Buffer)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tss.Copy(buf, rd)
		rd.Reset(bs)
	}
}
