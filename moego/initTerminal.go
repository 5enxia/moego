package moego

import (
    "fmt"
    "syscall"
)

func (e *Editor) InitTerminal() {
    e.flush()
    e.WriteHelpMenu(HELP_MESSAGE)
    e.writeStatusBar()
    e.MoveCursor(e.crow, e.ccol)
}

func (e *Editor) flush() {
    e.write(([]byte("\033[2J")))
}

func (e *Editor) WriteHelpMenu (message string){
    prevRow, prevCol := e.crow, e.ccol

    for i, ch := range message {
        e.MoveCursor(e.terminal.height+1, i)
        e.write([]byte(string(ch)))
    }

    // write blanks
    for i := len(message); i < e.terminal.width; i++ {
        e.MoveCursor(e.terminal.height+1, i)
        e.write([]byte{' '})
    }

    e.MoveCursor(prevRow, prevCol)
}


func (e *Editor) MoveCursor(row, col int) {
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
        e.MoveCursor(e.terminal.height+1, i)
        e.write([]byte(string(ch)))
    }

    // write blanks
    for i := len(e.filePath); i < e.terminal.width; i++ {
        e.MoveCursor(e.terminal.height, i)
        e.write([]byte{' '})
    }
}

func (e *Editor) setBackgroundColor(color Color) {
    s := fmt.Sprintf("\033[%dm", color)
    e.write([]byte(s))
}
