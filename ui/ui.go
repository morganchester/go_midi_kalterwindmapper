package ui

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"gopkg.in/yaml.v3"
	"kalterwind-mapper/midi"
)

type Mapping struct {
	InputDevice  string      `yaml:"inputDevice"`
	OutputDevice string      `yaml:"outputDevice"`
	Octaves      interface{} `yaml:"octaves"`
	Description  string      `yaml:"description"`
	Mapping      struct {
		NoteOn  map[int]int `yaml:"noteon"`
		NoteOff string      `yaml:"noteoff"`
	} `yaml:"mapping"`
	Controls map[int]string `yaml:"controls"`
}

func RunApp() {
	a := app.New()
	w := a.NewWindow("Kalterwind Mapper (Go Edition)")
	w.Resize(fyne.NewSize(550, 500))

	// UI элементы
	inputSelect := widget.NewSelect([]string{}, nil)
	outputSelect := widget.NewSelect([]string{}, nil)
	presetSelect := widget.NewSelect([]string{}, nil)

	infoLabel := widget.NewLabel("Информация о пресете появится здесь")

	// Загрузка пресетов
	presets, mappings := loadPresets("./presets")
	if len(presets) == 0 {
		presetSelect.Options = []string{"Нет доступных пресетов"}
		presetSelect.Disable()
	} else {
		presetSelect.Options = presets
	}

	// Выбор пресета → показать инфо
	presetSelect.OnChanged = func(s string) {
		if s == "" || s == "Нет доступных пресетов" {
			infoLabel.SetText("Пресет не выбран")
			return
		}
		m := mappings[s]
		info := fmt.Sprintf(
			"Описание: %s\nОктавы: %v\nInput: %s\nOutput: %s",
			m.Description, m.Octaves, m.InputDevice, m.OutputDevice,
		)
		infoLabel.SetText(info)
	}

	// Кнопка обновления устройств
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

	// Контейнер
	content := container.NewVBox(
		widget.NewLabel("Входное устройство:"),
		inputSelect,
		widget.NewLabel("Выходное устройство:"),
		outputSelect,
		widget.NewLabel("Пресет:"),
		presetSelect,
		infoLabel,
		refreshBtn,
		startBtn,
	)

	w.SetContent(content)
	w.ShowAndRun()
}

func loadPresets(path string) ([]string, map[string]Mapping) {
	var presets []string
	mappings := make(map[string]Mapping)

	_ = filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".yaml") {
			presets = append(presets, d.Name())

			data, err := os.ReadFile(p)
			if err == nil {
				var m Mapping
				if err := yaml.Unmarshal(data, &m); err == nil {
					mappings[d.Name()] = m
				}
			}
		}
		return nil
	})
	return presets, mappings
}
