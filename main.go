package main

import (
	"fmt"
	"kalterwind-mapper/midi"
	"kalterwind-mapper/ui"
)

func main() {
	err := midi.CreateVirtualPorts()
	if err != nil {
		fmt.Println("Ошибка создания виртуальных портов:", err)
	}

	ui.RunApp()
}
