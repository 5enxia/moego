package moego

func (e *Editor) SetRowCol(row int, col int) {
    if row > e.n &&  col > e.currentRow().visibleLen() {
        return
    }

    e.setRowPosition(row)
    e.setColPosition(col)
}

func (e *Editor) setRowPosition(row int) {
    if row >= e.n {
        row = e.n - 1
    }

    if row < 0 {
        if e.scroolrow > 0 {
            e.scroolrow -= 1
            e.Refresh()
        }

        row = 0
    }

    if row >= e.terminal.height {
        if row + e.scroolrow <= e.n {
            e.scroolrow += 1
        }
        row = e.terminal.height-1
        e.Refresh()
    }

    e.crow = row
    e.MoveCursor(row, e.ccol)
}

func (e *Editor) setColPosition(col int) {
    if col < 0 {
        col = 0
    }

    if col >= e.currentRow().visibleLen() {
        col = e.currentRow().visibleLen()
    }

    if col >= e.terminal.width {
        col = e.terminal.width - 1
    }

    e.ccol = col
    e.MoveCursor(e.crow, e.ccol)
}

func (e *Editor) currentRow() *Row {
    return e.rows[e.crow + e.scroolrow]
}

func (r *Row) visibleLen() int {
    return r.chars.VisibleLen()
}

