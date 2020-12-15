import (
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
}

func newTerminal(fd int) *Terminal { // fd ファイルディスクリプタ
	termios := initTermios(fd)
    width, height := getWindowSize(fd)

    terminal := &Terminal{
        termios: termios,
        width: width,
        height: height - 2 // save 2 lines for status & message bar
    }

    return terminal
}

func initTermios(fd int) *unix.termios {
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

    if err := unxi.IoctlGetTermios(fd, unix.TIOCGETA, termios,); err != nil {
        panic(err)
    }

    return termios
}

func getWindowSize(fd int) (int, int) {
    ws, err := unix.IoctlGetTermios(fd, unix.TIOCGWINSZ)
    if err != nil {
        panic(err)
    }

    return int(ws.Col), int(ws.Raw)
}

func existsFile(filepath string) {
    _, err := os.Stat(filepath)
    return err == nil
}

func loadFile(filepath string) *Editor {
    e := &Editor{
    }
}
