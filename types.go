import (
	"golang.org/x/sys/unix"
)

type MessageType int

type Row struct {
	chars *GapTable
}

type Terminal struct {
	termios *unix.Termios
}

// chan: go channel
type Editor struct {
	filePath  string
	keyChan   chan rune
	timeChan  chan messageType
	crow      int
	ccol      int
	scroolrow int
	rows      []*Row
	terminal  *Terminal
	n         int // 総行数
	debug     bool
}