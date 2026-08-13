package builder

import (
	"io"
	"strconv"
	"unicode/utf8"
	"unsafe"
)

type Builder struct {
	buf []byte
}

func (b *Builder) String() string {
	return unsafe.String(unsafe.SliceData(b.buf), len(b.buf))
}

func (b *Builder) Len() int { return len(b.buf) }

func (b *Builder) Cap() int { return cap(b.buf) }

func (b *Builder) Reset() {
	b.buf = b.buf[:0]
}

func (b *Builder) Grow(n int) {
	if n < 0 {
		panic("builder.Grow: negative count")
	}
	if cap(b.buf)-len(b.buf) < n {
		newCap := 2*cap(b.buf) + n
		if newCap < len(b.buf)+n {
			newCap = len(b.buf) + n
		}
		newBuf := make([]byte, len(b.buf), newCap)
		copy(newBuf, b.buf)
		b.buf = newBuf
	}
}

func (b *Builder) Write(p []byte) (int, error) {
	b.buf = append(b.buf, p...)
	return len(p), nil
}

func (b *Builder) WriteByte(c byte) error {
	b.buf = append(b.buf, c)
	return nil
}

func (b *Builder) WriteRune(r rune) (int, error) {
	n := len(b.buf)
	b.buf = utf8.AppendRune(b.buf, r)
	return len(b.buf) - n, nil
}

func (b *Builder) WriteString(s string) (int, error) {
	b.buf = append(b.buf, s...)
	return len(s), nil
}

func writeInt(b *Builder, v int64) {
	b.buf = strconv.AppendInt(b.buf, v, 10)
}

func writeUint(b *Builder, v uint64) {
	b.buf = strconv.AppendUint(b.buf, v, 10)
}

func writeFloat64(b *Builder, v float64) {
	b.buf = strconv.AppendFloat(b.buf, v, 'f', 3, 64)
}

func writeFloat32(b *Builder, v float32) {
	b.buf = strconv.AppendFloat(b.buf, float64(v), 'f', 3, 32)
}

func (b *Builder) WriteFormat(s string, args ...any) {
	arg := -1
	estimate := len(s) + len(args)*10
	if b.Cap() < estimate {
		b.Grow(estimate)
	}

	for i := 0; i < len(s); i++ {
		r := s[i]
		if r != '%' {
			b.WriteByte(r)
			continue
		}
		if i == len(s)-1 {
			continue
		}
		second := s[i+1]
		if second == '%' {
			b.WriteByte('%')

			i++
			continue
		}
		i++
		arg++
		if arg >= len(args) {
			b.WriteString("%!MISSING")
			break
		}

		v := second == 'v'
		val := args[arg]

		if str, ok := val.(string); ok && (second == 's' || v) {
			b.WriteString(str)
			continue
		}

		if second == 'f' || v {
			switch x := val.(type) {
			case float64:
				writeFloat64(b, x)
				continue
			case float32:
				writeFloat32(b, x)
				continue
			}
		}

		if second == 'd' || v {
			switch x := val.(type) {
			case int:
				writeInt(b, int64(x))
			case int8:
				writeInt(b, int64(x))
			case int16:
				writeInt(b, int64(x))
			case int32:
				writeInt(b, int64(x))
			case int64:
				writeInt(b, x)
			case uint:
				writeUint(b, uint64(x))
			case uint8:
				writeUint(b, uint64(x))
			case uint16:
				writeUint(b, uint64(x))
			case uint32:
				writeUint(b, uint64(x))
			case uint64:
				writeUint(b, x)
			default:
				goto cnt
			}
			continue
		}
	cnt:

		if bv, ok := val.(bool); ok {
			if bv {
				b.WriteString("true")
			} else {
				b.WriteString("false")
			}
			continue
		}

		b.WriteString("%!UNSUPPORTED")
	}
}

func (b *Builder) Bytes() []byte {
	return b.buf
}

func (b *Builder) Copy(dst io.Writer) (int, error) {
	return dst.Write(b.buf)
}

func New(size int) *Builder {
	return &Builder{buf: make([]byte, 0, size)}
}
