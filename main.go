package main

import (
	"flag"
	"fmt"
	"os"
)

type Editor struct {
	filePath string
	keyChan  chan rune
}

func main() {
	// コマンドライン引数の検証
	flag.Parse()
	if flag.NArg() < 1 || flag.NArg() > 3 {
		fmt.Println("Usage: moego <filename>")
		os.Exit(1)
	}

	// デバッグフラグの検証
	debug := flag.NArg() == 2 && flag.Arg(1) == "--debug"
}
