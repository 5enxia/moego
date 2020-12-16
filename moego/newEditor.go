package moego

import (
    "os"
    "io/ioutil"
	"golang.org/x/sys/unix"
)

func NewEditor(filepath string, debug bool) *Editor {
	terminal := newTerminal(0)

    // ファイルの存在チェック
    // ファイルを読み混み＆表示
    if existsFile(filepath) {
        e := loadFile(filepath)
        e.terminal = terminal
        return e
    }

    rows := makeRows()
    return &Editor {
        crow: 0,
        ccol: 0,
        scroolrow: 0,
        rows: rows,
        filePath: filepath,
        keyChan: make(chan rune),
        timeChan: make(chan MessageType),
        terminal: terminal,
        n: 1,
        debug: debug,
    }
}

func newTerminal(fd int) *Terminal { // fd ファイルディスクリプタ
	termios := initTermios(fd)
    width, height := getWindowSize(fd)

    terminal := &Terminal{
        termios: termios,
        width: width,
        height: height - 2, // save 2 lines for status & message bar
    }

    return terminal
}

func initTermios(fd int) *unix.Termios {
	// TIOCGETA:端末の状態をtermios構造体に格納し取得する
    // termios関数群は非同期通信ポートを制御するための汎用ターミナルインターフェース
    // ref: https://linuxjm.osdn.jp/html/LDP_man-pages/man3/termios.3.html
	termios, err := unix.IoctlGetTermios(fd, unix.TIOCGETA)
	if err != nil {
		panic(err)
	}

    termios.Iflag &^= unix.IGNBRK | unix.BRKINT | unix.PARMRK | unix.ISTRIP | unix.INLCR | unix.IGNCR | unix.ICRNL | unix.IXON
	termios.Oflag &^= unix.OPOST
	termios.Lflag &^= unix.ECHO | unix.ECHONL | unix.ICANON | unix.ISIG | unix.IEXTEN
	termios.Cflag &^= unix.CSIZE | unix.PARENB
	termios.Cflag |= unix.CS8
	termios.Cc[unix.VMIN] = 1
	termios.Cc[unix.VTIME] = 0

    if err := unix.IoctlSetTermios(fd, unix.TIOCGETA, termios); err != nil {
        panic(err)
    }

    return termios
}

func getWindowSize(fd int) (int, int) {
    ws, err := unix.IoctlGetWinsize(fd, unix.TIOCGWINSZ)
    if err != nil {
        panic(err)
    }

    return int(ws.Col), int(ws.Row)
}

func existsFile(filepath string) bool {
    _, err := os.Stat(filepath)
    return err == nil
}

func loadFile(filepath string) *Editor {
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

func makeRows() []*Row {
    rows := make([]*Row, 1024)
    for i := range rows {
        rows[i] = &Row {
            chars: NewGapTable(128), // 要勉強
        }
    }
    return rows
}