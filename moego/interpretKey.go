package moego

import (
    "golang.org/x/sys/unix"
    "fmt"
    "os"
)

func (r *Row) len() int        { return r.chars.Len() }

// emacs like command
func (e *Editor) InterpretKey() {
    for {
        // recive channel
        r := <-e.keyChan

        switch r {
        case ControlA:
            e.SetRowCol(e.crow, 0)

        case ControlB, ArrowLeft:
            e.back()

        case ControlC:
            e.exit()
            return

        case ControlE:
            e.SetRowCol(e.crow, e.numberOfRunesInRow())

        case ControlF, ArrowRight:
            e.next()

        case ControlH, BackSpace:
            e.backspace()

        case ControlN, ArrowDown:
            e.SetRowCol(e.crow+1, e.ccol)

        case Tab:
            for i := 0; i < 4; i++ {
                e.insertRune(e.currentRow(), e.ccol, rune(' '))
            }
            e.setColPos(e.ccol + 4)

        case Enter:
            e.newLine()

        case ControlS:
            saveFile(e.filePath, e.rows)
            e.writeHelpMenu("Saved!")
            e.timeChan <- RESET_MESSAGE

        case ControlP:
            e.SetRowCol(e.crow-1, e.ccol)

        // DEBUG
        case ControlV:
            e.debugDetailPrint(e)

        default:
            e.insertRune(e.currentRow(), e.ccol, r)
            e.setColPos(e.ccol +1)
        }
    }
}

// BACK
func (e *Editor) back() {
    if e.ccol == 0 {
        if e.crow > 0 {
            e.SetRowCol(e.crow-1, e.rows[e.crow + e.scroolrow-1].visibleLen())
        }
    } else {
        e.SetRowCol(e.crow, e.ccol-1)
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
            e.SetRowCol(e.crow+1, 0)
        }
    } else {
        e.SetRowCol(e.crow, e.ccol+1)
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
            e.SetRowCol(e.crow - 1, prevRow.len() - 1)
        }
    } else {
        e.deleteRune(row, e.ccol -1)
    }

    e.debugRowRunes()
}

// DELETE ROW
func (e *Editor) deleteRow(row int) {
    e.rows = append(e.rows[:row], e.rows[row+1:]...)
    e.n -= 1

    prevRowPos := e.crow
    e.Refresh()
    e.crow = prevRowPos
}

// DELETE RUNE
func (e *Editor) deleteRune(row *Row, col int) {
    row.deleteAt(col)
    e.updateRowRunes(row)
    e.SetRowCol(e.crow, e.ccol - 1)
}

func (r *Row) deleteAt (col int) {
    if col >= r.len() {
        return
    }

    r.chars.DeleteAt(col)
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
        e.debugPrint("DEBUG: row's view updated at", e.crow + e.scroolrow, "for", row.chars.Runes())
        e.writeRow(row)
    }
}

// DEBUG PRINT
func (e *Editor) debugPrint(a ...interface{}) {
	if e.debug {
		_, _ = fmt.Fprintln(os.Stderr, a...)
	}
}

// INSERT RUNE
func (e *Editor) insertRune (row *Row, col int, newRune rune) {
    row.insertAt(col, newRune)
    e.updateRowRunes(row)
}

// INSERT AT
func (r *Row) insertAt (colPos int, newRune rune) {
    if colPos > r.len() {
        colPos = r.len()
    }

    r.chars.InsertAt(colPos, newRune)
}

// NEWLINE
func (e *Editor) newLine() {
    // newLine => insert line
    currentLineRowPos := e.crow + e.scroolrow
    currentLineRow := e.rows[currentLineRowPos]

    newLineRowPos := currentLineRowPos + 1

    newRowRunes := append([]rune{}, currentLineRow.chars.Runes()[e.ccol:]...)
    e.insertRow(newLineRowPos, newRowRunes)

    // update current row
    currentRowNewRunes := append([]rune{}, currentLineRow.chars.Runes()[:e.ccol]...)
    currentRowNewRunes = append(currentRowNewRunes, '\n')
    e.replaceRune(e.crow + e.scroolrow, currentRowNewRunes)

    e.SetRowCol(e.crow + 1, 0)
    e.debugRowRunes()
}

// INSERT ROW
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

// DEBUG ROW RUNES
func (e *Editor) debugRowRunes() {
    if e.debug {
        for i := 0; i < e.n; i++ {
            _, _ = fmt.Fprintln(os.Stderr, i, ":", e.rows[i].chars.Runes())
        }
    }
}


func (e *Editor) reallocBufferIfNeed () {
    if e.n == len(e.rows) {
        newCap := cap(e.rows) * 2
        newRows := make([]*Row, newCap)
        copy(newRows, e.rows)
        e.DebugPrint("DEBUG: realloc occurred")
    }
}

// DEBUG DETAIL PRINT
func (e *Editor) debugDetailPrint(a ...interface{}) {
	if e.debug {
		_, _ = fmt.Fprintf(os.Stderr, "%+v\n", a...)
	}
}
