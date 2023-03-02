package main

import "codexcts/lib/console/tcell_console_wrapper"

var (
	cw *tcell_console_wrapper.ConsoleWrapper
	io *tcellRenderer
)

func onInit() {
	cw = &tcell_console_wrapper.ConsoleWrapper{}
	cw.Init()
	io = &tcellRenderer{}
	io.updateBounds()
}

func onClose() {
	cw.Close()
}

func readKey() string {
	return cw.ReadKey()
}
