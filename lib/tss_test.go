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
	if s.count == 2 {
		copy(p[:3], "hel")
		return 3, nil
	}
	if s.count == 3 {
		copy(p[:3], "lo\n")
		return 3, nil
	}
	if s.count == 4 {
		copy(p[:15], "hello\nhello\nhel")
		return 15, nil
	}
	if s.count == 5 {
		copy(p[:3], "lo\n")
		return 3, nil
	}
	copy(p[:6], "hello\n")
	return 6, nil
}

func TestWriter(t *testing.T) {
	t.Parallel()
	max := 6
	s := &sleepReader{max: max, sleepFor: 2 * time.Millisecond}
	buf := new(bytes.Buffer)
	w := tss.NewWriter(buf, time.Time{})
	n, err := io.Copy(w, s)
	if err != nil {
		t.Fatal(err)
	}
	if int(n) != len("hello\n")*max {
		t.Errorf("expected n of 36, got %d:\n%s", n, buf)
	}
}

func TestCopy(t *testing.T) {
	t.Parallel()
	max := 6
	s := &sleepReader{max: max, sleepFor: 2 * time.Millisecond}
	buf := new(bytes.Buffer)
	n, err := tss.Copy(buf, s)
	want := len("hello\n") * 6
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if int(n) != want {
		t.Errorf("expected n of %d, got %d", want, n)
	}
	parts := strings.Split(buf.String(), "\n")
	if len(parts) != 7 {
		t.Errorf("incorrect number of parts: want 6 got %d:\n%q", len(parts), parts)
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
		tss.CopyTime(buf, rd, time.Now().Add(-50*time.Millisecond))
		buf.Reset()
		rd.Reset(bs)
	}
}

func BenchmarkWriter(b *testing.B) {
	bs := bytes.Repeat([]byte{'a'}, 2<<12)
	for i := 0; i < len(bs); i += 50 {
		bs[i] = '\n'
	}
	b.SetBytes(int64(len(bs)))
	rd := bytes.NewReader(bs)
	buf := new(bytes.Buffer)
	w := tss.NewWriter(buf, time.Now().Add(-50*time.Millisecond))
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		io.Copy(w, rd)
		buf.Reset()
		rd.Reset(bs)
	}
}
