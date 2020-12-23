package moego

import (
    "golang.org/x/sys/unix"
    "fmt"
    "os"
    "strings"
    "io/ioutil"
)

func (r *Row) len() int        { return r.chars.Len() }

// emacs like command
func (e *Editor) interpretKey() {
    for {
        // recive channel
        r := <-e.keyChan

        switch r {
        case ControlA:
            e.setRowCol(e.crow, 0)

        case ControlB, ArrowLeft:
            e.back()

        case ControlC:
            e.exit()
            return

        case ControlE:
            e.setRowCol(e.crow, e.numberOfRunesInRow())

        case ControlF, ArrowRight:
            e.next()

        case ControlH, BackSpace:
            e.backspace()

        case ControlN, ArrowDown:
            e.setRowCol(e.crow+1, e.ccol)

        case Tab:
            for i := 0; i < 4; i++ {
                e.insertRune(e.currentRow(), e.ccol, rune(' '))
            }
            e.setColPosition(e.ccol + 4)

        case Enter:
            e.newLine()

        case ControlS:
            saveFile(e.filePath, e.rows)
            e.WriteHelpMenu("Saved!")
            e.timeChan <- RESET_MESSAGE

        case ControlP:
            e.setRowCol(e.crow-1, e.ccol)

        // DEBUG
        case ControlV:
            e.debugDetailPrint(e)

        default:
            e.insertRune(e.currentRow(), e.ccol, r)
            e.setColPosition(e.ccol +1)
        }
    }
}

// BACK
func (e *Editor) back() {
    if e.ccol == 0 {
        if e.crow > 0 {
            e.setRowCol(e.crow-1, e.rows[e.crow + e.scroolrow-1].visibleLen())
        }
    } else {
        e.setRowCol(e.crow, e.ccol-1)
    }
}

// EXIT
func (e *Editor) exit() { e.restoreTerminal(0) }
func (e *Editor) restoreTerminal(fd int) {
    if err := unix.IoctlSetTermios(fd, unix.TIOCSETA, e.terminal.termios); err != nil {
        panic(err)
    }
}

func (e *Editor) numberOfRunesInRow() int { return e.currentRow().chars.Len() }

// NEXT
func (e *Editor) next() {
    if e.ccol >= e.currentRow().visibleLen() {
        if e.crow+1 < e.n {
            e.setRowCol(e.crow+1, 0)
        }
    } else {
        e.setRowCol(e.crow, e.ccol+1)
    }
}

// BACKSPACE
func (e *Editor) backspace() {
    row := e.currentRow()

    if e.ccol == 0 {
        if e.crow + e.scroolrow > 0 {
            prevRowPos := e.crow + e.scroolrow - 1
            prevRow := e.rows[prevRowPos]

            // Update the previous row
            newRunes := append([]rune{}, prevRow.chars.Runes()[:prevRow.len()-1]...)
            newRunes = append(newRunes, row.chars.Runes()...)
            e.replaceRune(prevRowPos, newRunes)

            // Delete the current row
            currentRowPos := e.crow + e.scroolrow
            e.deleteRow(currentRowPos)
            e.setRowCol(e.crow - 1, prevRow.len() - 1)
        }
    } else {
        e.deleteRune(row, e.ccol -1)
    }

    e.debugRowRunes()
}

// REPLACE RUNE
func (e *Editor) replaceRune(row int, newRune []rune) {
    gt := NewGapTable(128)

    for _, r := range newRune {
        gt.AppendRune(r)
    }

    r := &Row{
        chars: gt,
    }

    e.rows[row] = r

    prevRowPos := e.crow
    e.crow = row - e.scroolrow
    e.updateRowRunes(r)
    e.crow = prevRowPos
}

// UPDATE RUNE
func (e *Editor) updateRowRunes(row *Row) {
    if e.crow < e.terminal.height {
        e.DebugPrint("DEBUG: row's view updated at", e.crow + e.scroolrow, "for", row.chars.Runes())
        e.writeRow(row)
    }
}

func (e *Editor) debugDetailPrint(a ...interface{}) {
	if e.debug {
		_, _ = fmt.Fprintf(os.Stderr, "%+v\n", a...)
	}
}

func saveFile(filepath string, rows []*Row) {
    sb := strings.Builder{}

    for _, r := range rows {
        if r.len() >= 1 {
            for _, ch := range r.chars.Runes() {
                sb.WriteRune(ch)
            }
        }
    }

    _ = ioutil.WriteFile(filepath, []byte(sb.String()), 0644)
}

func (e *Editor) insertRow(row int, runes []rune) {
    gt := NewGapTable(128)

    for _, r := range runes {
        gt.AppendRune(r)
    }

    r := &Row {
        chars:gt,
    }

    // https://github.com/golang/go/wiki/SliceTricks
	e.rows = append(e.rows[:row], append([]*Row{ r }, e.rows[row:]...)...)
	e.n += 1
    e.reallocBufferIfNeed()

    prevRowPos := e.crow
    e.Refresh()
    e.crow = prevRowPos
}

func (e *Editor) reallocBufferIfNeed () {
    if e.n == len(e.rows) {
        newCap := cap(e.rows) * 2
        newRows := make([]*Row, newCap)
        copy(newRows, e.rows)
        e.DebugPrint("DEBUG: realloc occurred")
    }
}

func (e *Editor) newLine() {
    // newLine = insert line
    currentLineRowPos := e.crow + e.scroolrow
    currentLineRow := e.rows[currentLineRowPos]

    newLineRowPos := e.crow + e.scroolrow
    newRowRunes := append([]rune{}, currentLineRow.chars.Runes()[e.ccol:]...)
    e.insertRow(newLineRowPos, newRowRunes)

    currentRowNewRunes := append([]rune{}, currentLineRow.chars.Runes()[:e.ccol]...)
    currentRowNewRunes = append(currentRowNewRunes, '\n')
    e.replaceRune(e.crow + e.scroolrow, currentRowNewRunes)

    e.setRowCol(e.crow + 1, 0)
    e.debugRowRunes()
}

func (e *Editor) deleteRow(row int) {
    e.rows = append(e.rows[:row], e.rows[row+1:]...)
    e.n -= 1

    prevRowPos := e.crow
    e.Refresh()
    e.crow = prevRowPos
}

func (e *Editor) deleteRune(row *Row, col int) {
    row.deleteAt(col)
    e.updateRowRunes(row)
    e.setRowCol(e.crow, e.ccol - 1)
}

func (r *Row) deleteAt (col int) {
    if col >= r.len() {
        return
    }

    r.chars.DeleteAt(col)
}

func (r *Row) insertAt (colPosition int, newRune rune) {
    if colPosition > r.len() {
        colPosition = r.len()
    }

    r.chars.InsertAt(colPosition, newRune)
}

func (e *Editor) insertRune (row *Row, col int, newRune rune) {
    row.insertAt(col, newRune)
    e.updateRowRunes(row)
}

func (e *Editor) debugRowRunes() {
    if e.debug {
        for i := 0; i < e.n; i++ {
            _, _ = fmt.Fprintln(os.Stderr, i, ":", e.rows[i].chars.Runes())
        }
    }
}

