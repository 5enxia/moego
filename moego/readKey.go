package moego

import (
    "syscall"
    "unicode/utf8"
)

func (e *Editor) ReadKeys() {
    buf := make([]byte, 64)

    // 無限ループ
    for {
        if n, err := syscall.Read(0, buf); err == nil {
            b := buf[:n]
            for {
                r, n := e.parseKey(b)

                if n == 0 {
                    break
                }

                e.keyChan <- r
                b = b[n:]
            }
        }
    }
}

func (e *Editor) parseKey(buf []byte) (rune, int) {
    // parse escape sequence
    if len(buf) == 3 {
        if buf[0] == byte(27) && buf[1] == '[' {
            switch buf[2] {
            case 'A':
                return ArrowUp, 3
            case 'B':
                return ArrowDown, 3
            case 'C':
                return ArrowRight, 3
            case 'D':
                return ArrowLeft, 3
             default:
                 return DummyKey, 0
            }
        }
    }

    // parse bytes as utf-8
    return utf8.DecodeRune(buf)
}
