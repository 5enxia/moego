package moego

import (
    "fmt"
    "syscall"
)

func (e *Editor) initTerminal() {
    e.flush()
    e.writeHelpMenu(HELP_MESSAGE)
    e.writeStatusBar()
    e.moveCursor(e.crow, e.ccol)
}

func (e *Editor) flush() {
    e.write(([]byte("\033[2J")))
}

func (e *Editor) writeHelpMenu (message string){
    prevRow, prevCol := e.crow, e.ccol

    for i, ch := range message {
        e.moveCursor(e.terminal.height+1, i)
        e.write([]byte(string(ch)))
    }

    // write blanks
    for i := len(message); i < e.terminal.width; i++ {
        e.moveCursor(e.terminal.height+1, i)
        e.write([]byte{' '})
    }

    e.moveCursor(prevRow, prevCol)
}


func (e *Editor) moveCursor(row, col int) {
    s := fmt.Sprintf("\033[%d;%dH", row+1, col+1)
    e.write([]byte(s))
}

func (e *Editor) write(b []byte) {
    syscall.Write(0, b)
}

func (e *Editor) writeStatusBar() {
    // set status-bar background color
    e.setBackgroundColor(BG_CYAN)
    defer e.setBackgroundColor(BLACK)

    // show current text file under status-bar
    for i, ch := range e.filePath {
        e.moveCursor(e.terminal.height+1, i)
        e.write([]byte(string(ch)))
    }

    // write blanks
    for i := len(e.filePath); i < e.terminal.width; i++ {
        e.moveCursor(e.terminal.height, i)
        e.write([]byte{' '})
    }
}

func (e *Editor) setBackgroundColor(color Color) {
    s := fmt.Sprintf("\033[%dm", color)
    e.write([]byte(s))
}
