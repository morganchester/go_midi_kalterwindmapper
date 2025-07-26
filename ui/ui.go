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
	InputDevice  string `yaml:"inputDevice"`
	OutputDevice string `yaml:"outputDevice"`
	Octaves      int    `yaml:"octaves"`
	Description  string `yaml:"description"`
	Mapping      struct {
		NoteOn  map[int][]interface{} `yaml:"noteon"`
		NoteOff struct {
			Action  string        `yaml:"action"`
			Message []interface{} `yaml:"message"`
		} `yaml:"noteoff"`
	} `yaml:"mapping"`
	Controls map[int][]interface{} `yaml:"controls"`
}

func RunApp() {
	a := app.New()
	w := a.NewWindow("Kalterwind Mapper (Go Edition)")
	w.Resize(fyne.NewSize(600, 500))

	// UI элементы
	inputSelect := widget.NewSelect([]string{}, nil)
	outputSelect := widget.NewSelect([]string{}, nil)
	presetSelect := widget.NewSelect([]string{}, nil)
	infoLabel := widget.NewLabel("Информация о пресете появится здесь")

	// Загрузка YAML пресетов
	presets, mappings := loadPresets("./presets")
	if len(presets) == 0 {
		presetSelect.Options = []string{"Нет доступных пресетов"}
		presetSelect.Disable()
	} else {
		presetSelect.Options = presets
	}

	// Выбор пресета → показать информацию
	presetSelect.OnChanged = func(s string) {
		if s == "" || s == "Нет доступных пресетов" {
			infoLabel.SetText("Пресет не выбран")
			return
		}
		m := mappings[s]
		info := fmt.Sprintf(
			"Описание: %s\nОктавы: %d\nInput: %s\nOutput: %s",
			m.Description, m.Octaves, m.InputDevice, m.OutputDevice,
		)
		infoLabel.SetText(info)
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

	// Кнопка запуска маршрутизации
	startBtn := widget.NewButton("Запустить", func() {
		selPreset := presetSelect.Selected
		selIn := inputSelect.Selected
		selOut := outputSelect.Selected

		if strings.Contains(selIn, "Нет") || strings.Contains(selOut, "Нет") {
			infoLabel.SetText("Невозможно запустить: нет устройств")
			return
		}
		if selPreset == "" || selPreset == "Нет доступных пресетов" {
			infoLabel.SetText("Выберите пресет")
			return
		}

		inIdx, outIdx := midi.FindDeviceIndexes(selIn, selOut)
		if inIdx == -1 || outIdx == -1 {
			infoLabel.SetText("Устройства не найдены")
			return
		}

		// Конвертация в структуру из midi/processor.go
		m := mappings[selPreset]
		mapping := midi.Mapping{
			NoteOn:   m.Mapping.NoteOn,
			NoteOff:  midi.NoteOffAction{Action: m.Mapping.NoteOff.Action, Message: m.Mapping.NoteOff.Message},
			Controls: m.Controls,
		}

		processor := midi.NewProcessor()
		err := processor.Start(inIdx, outIdx, mapping)
		if err != nil {
			infoLabel.SetText(fmt.Sprintf("Ошибка: %v", err))
		} else {
			infoLabel.SetText("Маршрутизация запущена!")
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

// Загрузка YAML-файлов из ./presets
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
