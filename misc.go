package redis

import (
)

func quote(in []byte) []byte {
	out := make([]byte, 0)
	for _, i := range in {
		switch i {
		case '\a':
			out = append(out, `\a`...)
		case '\b':
			out = append(out, `\b`...)
		case '\f':
			out = append(out, `\f`...)
		case '\n':
			out = append(out, `\n`...)
		case '\r':
			out = append(out, `\r`...)
		case '\t':
			out = append(out, `\t`...)
		case '\v':
			out = append(out, `\v`...)
		case 0x0:
			out = append(out, `\0`...)
		case 0x20:
			out = append(out, `\s`...)
		case '\\':
			out = append(out, `\\`...)
		default:
			out = append(out, i)
		}
	}
	return out
}

func unquote(in []byte) []byte {
	out := make([]byte, 0)
	for n := 0; n < len(in); n++ {
		if in[n] == '\\' && n+1 < len(in) {
			switch in[n+1] {
			case 'a':
				out = append(out, '\a')
				n++
			case 'b':
				out = append(out, '\b')
				n++
			case 'f':
				out = append(out, '\f')
				n++
			case 'n':
				out = append(out, '\n')
				n++
			case 'r':
				out = append(out, '\r')
				n++
			case 't':
				out = append(out, '\t')
				n++
			case 'v':
				out = append(out, '\v')
				n++
			case '0':
				out = append(out, 0x0)
				n++
			case 's':
				out = append(out, 0x20)
				n++
			case '\\':
				out = append(out, '\\')
				n++
			default:
				out = append(out, in[n])
			}
		} else {
			out = append(out, in[n])
		}
	}
	return out
	return in
}
