package tss_test

import (
	"bytes"
	"io"
	"strings"
	"testing"
	"time"

	tss "github.com/kevinburke/tss/lib"
)

type sleepReader struct {
	count    int
	max      int
	sleepFor time.Duration
}

func (s *sleepReader) Read(p []byte) (int, error) {
	s.count++
	if s.count > s.max {
		return 0, io.EOF
	}
	if s.count == 1 {
		copy(p[:6], "hello\n")
		return 6, nil
	}
	time.Sleep(s.sleepFor)
	copy(p[:6], "hello\n")
	return 6, nil
}

func TestCopy(t *testing.T) {
	t.Parallel()
	s := &sleepReader{max: 3, sleepFor: 5 * time.Millisecond}
	buf := new(bytes.Buffer)
	n, err := tss.Copy(buf, s)
	if n != 72 {
		t.Errorf("expected n of 72, got %d", n)
	}
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	parts := strings.Split(buf.String(), "\n")
	if len(parts) != 4 {
		t.Errorf("incorrect number of parts: want 4 got %q", parts)
	}
	line1 := parts[0]
	if len(line1) != 23 {
		t.Errorf("line1 length: want %d got %d", 23, len(line1))
	}
	lineParts := strings.Fields(line1)
	if len(lineParts) != 2 {
		t.Errorf("wrong line parts")
	}
	part, err := time.ParseDuration(lineParts[0])
	if err != nil {
		t.Fatal(err)
	}
	if part > 100*time.Millisecond {
		t.Errorf("part took too long: %d", part)
	}
	lineParts = strings.Fields(parts[1])
	if len(lineParts) != 3 {
		t.Errorf("wrong number of line parts in line 2: got %d", len(lineParts))
	}
}

var scalerTests = []struct {
	in  time.Duration
	out string
}{
	{100 * time.Microsecond, "0.1ms"},
	{500 * time.Microsecond, "0.5ms"},
	{99 * time.Microsecond, "0.1ms"},
	{49 * time.Microsecond, "49.0µs"},
	{time.Millisecond, "1.0ms"},
	{56*time.Millisecond + 290*time.Microsecond, "56.3ms"},
	{56*time.Millisecond + 251*time.Microsecond, "56.3ms"},
	{56*time.Millisecond + 100*time.Microsecond, "56.1ms"},
	{0, "0.0ms"},
}

func TestTimeScaler(t *testing.T) {
	for _, tt := range scalerTests {
		v := tss.TimeScaler(tt.in)
		if v != tt.out {
			t.Errorf("timeScaler(%q): want %q, got %q", tt.in, tt.out, v)
		}
	}
}

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
