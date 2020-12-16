package main

import (
	"flag"
	"fmt"
	"os"
    "./moego"
)

func main() {
	// コマンドライン引数の検証
	flag.Parse()
	if flag.NArg() < 1 || flag.NArg() > 3 {
		fmt.Println("Usage: moego <filename>")
		os.Exit(1)
	}

	filepath := flag.Arg(0)

	// デバッグフラグの検証
	debug := flag.NArg() == 2 && flag.Arg(1) == "--debug"

	e := moego.NewEditor(filepath, debug) // エディタを初期化
	e.InitTerminal() // ターミナルを初期化
	 e.Refresh() // 画面をリフレッシュ
	// e.SetRowCol(0, 0) // バーの位置を原点へ

	// go e.ReadKey() // マルチスレッドでキーボード入力を読む
	// go e.PollTimeEvent() // polling

	// e.InterpretKey()
}
