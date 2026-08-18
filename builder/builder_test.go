package builder

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestBuilder_WriteString(t *testing.T) {
	b := &Builder{}
	n, err := b.WriteString("hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 5 {
		t.Errorf("WriteString returned %d, want 5", n)
	}
	if got := b.String(); got != "hello" {
		t.Errorf("String() = %q, want %q", got, "hello")
	}
}

func TestBuilder_WriteByte(t *testing.T) {
	b := &Builder{}
	err := b.WriteByte('a')
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := b.String(); got != "a" {
		t.Errorf("String() = %q, want %q", got, "a")
	}
}

func TestBuilder_WriteRune(t *testing.T) {
	b := &Builder{}
	n, err := b.WriteRune('⌘')
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 3 { // UTF-8 length of '⌘'
		t.Errorf("WriteRune returned %d, want 3", n)
	}
	if got := b.String(); got != "⌘" {
		t.Errorf("String() = %q, want %q", got, "⌘")
	}
}

func TestBuilder_Write(t *testing.T) {
	b := &Builder{}
	n, err := b.Write([]byte("abc"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := b.String(); got != "abc" {
		t.Errorf("String() = %q, want %q", got, "abc")
	}
	if n != 3 {
		t.Errorf("n = %d, want %d", n, 3)
	}
}

func TestBuilder_GrowNegative(t *testing.T) {
	b := &Builder{}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Grow(-1) did not panic")
		}
	}()
	b.Grow(-1)
}

func TestBuilder_Reset(t *testing.T) {
	b := &Builder{}
	b.WriteString("hello")
	b.Reset()
	if b.Len() != 0 {
		t.Errorf("Len after Reset = %d, want 0", b.Len())
	}
	if b.String() != "" {
		t.Errorf("String after Reset = %q, want empty", b.String())
	}

	b.WriteString("world")
	if got := b.String(); got != "world" {
		t.Errorf("after Reset and rewrite = %q, want %q", got, "world")
	}
}

func TestBuilder_Grow(t *testing.T) {
	b := &Builder{}
	b.Grow(100)
	if cap(b.buf) < 100 {
		t.Errorf("cap after Grow(100) = %d, want at least 100", cap(b.buf))
	}

	b.Grow(300)
	if cap(b.buf) < 300 {
		t.Errorf("cap after Grow(100) = %d, want at least 100", cap(b.buf))
	}

	// Проверяем, что после Grow можно писать
	b.WriteString("test")
	if got := b.String(); got != "test" {
		t.Errorf("after Grow and WriteString = %q, want %q", got, "test")
	}
}

func TestBuilder_Copy(t *testing.T) {
	b := &Builder{}
	b.WriteString("hello")
	var dst bytes.Buffer
	n, err := b.Copy(&dst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 5 {
		t.Errorf("Copy returned %d, want 5", n)
	}
	if got := dst.String(); got != "hello" {
		t.Errorf("dst.String() = %q, want %q", got, "hello")
	}
}

func TestBuilder_WriteFormat(t *testing.T) {
	tests := []struct {
		name   string
		format string
		args   []any
		want   string
	}{
		{"string", "%s", []any{"world"}, "world"},
		{"int", "%d", []any{int(123)}, "123"},
		{"int8", "%d", []any{int8(123)}, "123"},
		{"int16", "%d", []any{int16(123)}, "123"},
		{"int32", "%d", []any{int32(123)}, "123"},
		{"int64", "%d", []any{int64(123)}, "123"},
		{"uint", "%d", []any{uint(123)}, "123"},
		{"uint8", "%d", []any{uint8(123)}, "123"},
		{"uint16", "%d", []any{uint16(123)}, "123"},
		{"uint32", "%d", []any{uint32(123)}, "123"},
		{"uint64", "%d", []any{uint64(123)}, "123"},
		{"float32", "%f", []any{float32(1.23)}, "1.230"},
		{"float64", "%f", []any{1.23}, "1.230"},
		{"bool", "%v %v", []any{true, false}, "true false"},
		{"percent", "%%", []any{}, "%"},
		{"mixed", "%s %d %f", []any{"test", 42, 3.14}, "test 42 3.140"},
		{"missing arg", "%s %d", []any{"only"}, "only %!MISSING"},
		{"unsupported", "%v", []any{struct{ A int }{42}}, "%!UNSUPPORTED"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Builder{}
			b.WriteFormat(tt.format, tt.args...)
			if got := b.String(); got != tt.want {
				t.Errorf("WriteFormat(%q, %v) = %q, want %q", tt.format, tt.args, got, tt.want)
			}
		})
	}
}

func TestBuilder_Bytes(t *testing.T) {
	b := &Builder{}
	b.WriteString("hello")
	got := b.Bytes()
	if string(got) != "hello" {
		t.Errorf("Bytes() = %q, want %q", got, "hello")
	}
}

func TestBuilder_New(t *testing.T) {
	b := New(10)
	if b.Len() != 0 {
		t.Errorf("Len after New = %d, want %d", b.Len(), 0)
	}
	if b.Cap() < 0 {
		t.Errorf("Cap after New = %d, want %d", b.Cap(), 10)
	}
}

func TestBuilder_LenCap(t *testing.T) {
	b := &Builder{}
	if b.Len() != 0 {
		t.Errorf("Len before write = %d, want 0", b.Len())
	}
	if b.Cap() < 0 {
		t.Errorf("Cap before write = %d, want >=0", b.Cap())
	}
	b.WriteString("abc")
	if b.Len() != 3 {
		t.Errorf("Len after write = %d, want 3", b.Len())
	}
	if b.Cap() < 3 {
		t.Errorf("Cap after write = %d, want at least 3", b.Cap())
	}
}

func BenchmarkStringsBuilder_WriteString(b *testing.B) {
	bb := &strings.Builder{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bb.WriteString("hello")
	}
}

func BenchmarkTuiBuilder_WriteString(b *testing.B) {
	bb := &Builder{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bb.WriteString("hello")
	}
}

func BenchmarkTuiBuilder_WriteFormat(b *testing.B) {
	bb := &Builder{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bb.WriteFormat("%s %d %f", "test", 42, 1.23)
	}
}

func BenchmarkStringsBuilder_WriteFmtSprintf(b *testing.B) {
	bb := &strings.Builder{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bb.WriteString(fmt.Sprintf("%s %d %f", "test", 42, 1.23))
	}
}

func makeData() []byte {
	b := make([]byte, 1024)
	for i := range b {
		b[i] = byte(i)
	}
	return b
}

func BenchmarkStringsBuilderStringToBytes(b *testing.B) {
	bb := &strings.Builder{}
	bb.Write(makeData())
	var d []byte
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d = []byte(bb.String())
	}
	b.StopTimer()
	d[0] = 1
}

func BenchmarkTuiBuilderBytes(b *testing.B) {
	bb := &Builder{}
	bb.Write(makeData())
	var d []byte
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d = bb.Bytes()
	}
	b.StopTimer()
	d[0] = 1
}

func BenchmarkStringsBuilderStringToBytesToWriter(b *testing.B) {
	wr := io.Discard
	bb := &strings.Builder{}
	bb.Write(makeData())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wr.Write([]byte(bb.String()))
	}
}

func BenchmarkTuiBuilderCopy(b *testing.B) {
	wr := io.Discard
	bb := &Builder{}
	bb.Write(makeData())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bb.Copy(wr)
	}
}
