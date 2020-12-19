package moego

import (
    "syscall"
)

func (e *Editor) {
    buf := make([]byte, 64)

    // 無限ループ
    for {
        if n, err := syscall.Read(0, buf); err == nil {
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
