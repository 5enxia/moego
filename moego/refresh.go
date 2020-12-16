package moego

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
    "syscall"
)

func (e *Editor) Refresh() {
    for i := 0; i < e.terminal.height; i++ {
        e.crow = i
        e.writeRow(e.rows[e.scroolrow+i])
    }
}

// ...interface{} is variable args
// arg...
func (e *Editor) debugPrint(arg ...interface{}) {
    if e.debug {
        _, _ = fmt.Fprintln(os.Stderr, arg...)
    }
}

func (e *Editor) writeRow(r *Row) {
    var buf []byte

    for _, r := range r.chars.Runes() {
        buf = append(buf, []byte(string(r))...)
    }

    e.moveCursor(e.crow, 0)
    e.flushRow()

    // go highlight
    if filepath.Ext(e.filePath) == ".go" {
        colors := e.highlight(buf)
        e.writeWithColor(buf, colors)
    } else {
        e.write(buf)
    }
}

func (e *Editor) flushRow() {
    e.write([]byte("033[2J"))
}

func (e *Editor) highlight(buf []byte) []Color {
    colors := make([]Color, len(buf))
    for i := range colors {
        colors[i] =  DefaultColor
    }

    // ascii
    ascii := string(buf)

    // hilighting keywrods
    for key := range KeywordColor {
        index := strings.Index(ascii, string(key))
        if index != -1 {
            for i := 0; i < len(string(key)); i++ {
                colors[index+i] = KeywordColor[key]
            }
        }
    }

    // string literal
    isStringLiteral := false
    for i, b := range ascii {
        if b == '"' || isStringLiteral {
            if b == '"' {
                isStringLiteral = !isStringLiteral
            }
            colors[i] = GREEN
        }
    }

    return colors
}

func (e *Editor) writeWithColor(buf []byte, colros []Color) {
    var newBuf []byte

    for i, c := range colros {
        s := fmt.Sprintf("\033[%dm", c)
        newBuf = append(newBuf, []byte(s)...)
        newBuf = append(newBuf, buf[i])
    }
    syscall.Write(0, newBuf)
}

