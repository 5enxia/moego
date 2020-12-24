package moego

import (
    "strings"
    "io/ioutil"
)

// LOAD FILE
func LoadFile(filepath string) *Editor {
    e := &Editor{
		crow:      0,
		ccol:      0,
		scroolrow: 0,
		filePath:  filepath,
		keyChan:   make(chan rune),
		timeChan:  make(chan MessageType),
		n:         1,
    }

    rows := makeRows()

    bytes, err := ioutil.ReadFile(filepath)
    if err != nil {
        panic(err)
    }

    // GapTableについては要勉強
    gt := NewGapTable(128)

    for _, b := range bytes {
		// Treat TAB as 4 spaces.
		if b == Tab {
			gt.AppendRune(rune(0x20))
			gt.AppendRune(rune(0x20))
			gt.AppendRune(rune(0x20))
			gt.AppendRune(rune(0x20))
			continue
		}

		// ASCII-only
		gt.AppendRune(rune(b))

		if b == '\n' {
			rows[e.n-1] = &Row{chars: gt}
			e.n += 1
			gt = NewGapTable(128)
		}
	}

	rows[e.n-1] = &Row{chars: gt}
	e.rows = rows

	return e
}

// SAVE FILE
func SaveFile(filepath string, rows []*Row) {
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
