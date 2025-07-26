package ui

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"kalterwind-mapper/midi"
)

func RunApp() {
	a := app.New()
	w := a.NewWindow("Kalterwind Mapper (Go Edition)")
	w.Resize(fyne.NewSize(500, 400))

	// Списки устройств
	inputSelect := widget.NewSelect([]string{}, func(s string) {})
	outputSelect := widget.NewSelect([]string{}, func(s string) {})
	presetSelect := widget.NewSelect([]string{}, func(s string) {})

	// Загрузка пресетов
	presets := loadPresets("./presets")
	if len(presets) == 0 {
		presetSelect.Options = []string{"Нет доступных пресетов"}
		presetSelect.Disable()
	} else {
		presetSelect.Options = presets
	}

	// Кнопка обновления списка устройств
	refreshBtn := widget.NewButton("Обновить устройства", func() {
		ins, outs := midi.GetDevices()
		if len(ins) == 0 {
			inputSelect.Options = []string{"Нет входных устройств"}
			inputSelect.Disable()
		} else {
			inputSelect.Options = ins
			inputSelect.Enable()
		}

		if len(outs) == 0 {
			outputSelect.Options = []string{"Нет выходных устройств"}
			outputSelect.Disable()
		} else {
			outputSelect.Options = outs
			outputSelect.Enable()
		}

		inputSelect.Refresh()
		outputSelect.Refresh()
	})

	// Кнопка запуска
	startBtn := widget.NewButton("Запустить", func() {
		selPreset := presetSelect.Selected
		selIn := inputSelect.Selected
		selOut := outputSelect.Selected

		fmt.Println("Выбранный пресет:", selPreset)
		fmt.Println("Вход:", selIn)
		fmt.Println("Выход:", selOut)

		if strings.Contains(selIn, "Нет") || strings.Contains(selOut, "Нет") {
			fmt.Println("Невозможно запустить: нет устройств")
			return
		}
	})

	content := container.NewVBox(
		widget.NewLabel("Входное устройство:"),
		inputSelect,
		widget.NewLabel("Выходное устройство:"),
		outputSelect,
		widget.NewLabel("Пресет:"),
		presetSelect,
		refreshBtn,
		startBtn,
	)

	w.SetContent(content)
	w.ShowAndRun()
}

func loadPresets(path string) []string {
	var presets []string
	_ = filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".yaml") {
			presets = append(presets, d.Name())
		}
		return nil
	})
	return presets
}
