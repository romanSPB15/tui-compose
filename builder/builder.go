package builder

import (
	"io"
	"strconv"
	"unicode/utf8"
	"unsafe"
)

type noCopy struct{}

func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}

type Builder struct {
	_   noCopy
	buf []byte
}

// String возвращает строку, разделяющую память с внутренним буфером.
// Результат валиден до следующей записи в Builder.
func (b *Builder) String() string {
	return unsafe.String(unsafe.SliceData(b.buf), len(b.buf))
}

// String возвращает копию содержимого буфера в string.
func (b *Builder) StringCopy() string {
	return string(b.buf)
}

// Len возвращает длину буфера.
func (b *Builder) Len() int { return len(b.buf) }

// Cap возвращает ёмкость буфера.
func (b *Builder) Cap() int { return cap(b.buf) }

// Reset сбрасывает длину буфера.
func (b *Builder) Reset() {
	b.buf = b.buf[:0]
}

// Grow выделяет N байт.
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

// Write реализует io.Writer.
func (b *Builder) Write(p []byte) (int, error) {
	b.buf = append(b.buf, p...)
	return len(p), nil
}

// WriteByte дописывает байт в Builder.
func (b *Builder) WriteByte(c byte) error {
	b.buf = append(b.buf, c)
	return nil
}

// WriteByte дописывает руну в Builder.
func (b *Builder) WriteRune(r rune) (int, error) {
	n := len(b.buf)
	b.buf = utf8.AppendRune(b.buf, r)
	return len(b.buf) - n, nil
}

// WriteByte дописывает строку в Builder.
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

func (b *Builder) WriteInt(v int) {
	b.buf = strconv.AppendInt(b.buf, int64(v), 10)
}

func (b *Builder) WriteUint(v uint) {
	b.buf = strconv.AppendUint(b.buf, uint64(v), 10)
}

// WriteFormat записывает форматированную строку по подмножеству правил
// fmt.Printf: %s, %d, %f (точность 3), %v (string/int/float/bool), %%.
// Width, precision, флаги и прочие глаголы не поддерживаются.
// При нехватке аргументов пишет "%!MISSING" и прекращает разбор.
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

// Bytes возращает буфер Builder.
func (b *Builder) Bytes() []byte {
	return b.buf
}

// Copy копирует буфер в io.Writer.
func (b *Builder) Copy(dst io.Writer) (int, error) {
	total := 0
	for total < len(b.buf) {
		n, err := dst.Write(b.buf[total:])
		total += n
		if err != nil {
			return total, err
		}
		if n == 0 {
			return total, io.ErrShortWrite
		}
	}
	return total, nil
}

// New создаёт Builder с указанной ёмкостью.
func New(size int) *Builder {
	return &Builder{buf: make([]byte, 0, size)}
}
