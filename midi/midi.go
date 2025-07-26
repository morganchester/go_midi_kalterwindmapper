package midi

import (
	"fmt"
	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/rtmididrv"
)

var TestMode = false // True для ночных тестов

var (
	VirtualIn  midi.In
	VirtualOut midi.Out
)

// CreateVirtualPorts создает виртуальные MIDI-вход и выход
func CreateVirtualPorts() error {
	drv, err := rtmididrv.New()
	if err != nil {
		return fmt.Errorf("ошибка драйвера: %w", err)
	}

	in, err := drv.OpenVirtualIn("KalterwindMapper In")
	if err != nil {
		return fmt.Errorf("не удалось создать виртуальный вход: %w", err)
	}

	out, err := drv.OpenVirtualOut("KalterwindMapper Out")
	if err != nil {
		return fmt.Errorf("не удалось создать виртуальный выход: %w", err)
	}

	in.SetListener(func(data []byte, delta int64) {
		fmt.Printf("[Virtual IN] % X\n", data)
		// НЕ отправляем напрямую в out, чтобы не создать петлю
	})

	VirtualIn = in
	VirtualOut = out

	fmt.Println("[VIRTUAL] Подняты порты: KalterwindMapper In / Out")
	return nil
}

// GetDevices возвращает списки имён входных и выходных устройств
func GetDevices() ([]string, []string) {
	if TestMode {
		return []string{"Virtual In"}, []string{"Virtual Out"}
	}
	drv, err := rtmididrv.New()
	if err != nil {
		return []string{}, []string{}
	}
	defer drv.Close()

	ins, _ := drv.Ins()
	outs, _ := drv.Outs()

	inNames := []string{}
	outNames := []string{}

	for _, i := range ins {
		inNames = append(inNames, i.String())
	}
	for _, o := range outs {
		outNames = append(outNames, o.String())
	}

	return inNames, outNames
}

// FindDeviceIndexes ищет индексы устройств по именам
func FindDeviceIndexes(inputName, outputName string) (int, int) {
	if TestMode {
		return 0, 0 // всегда виртуальные устройства
	}
	drv, err := rtmididrv.New()
	if err != nil {
		return -1, -1
	}
	defer drv.Close()

	ins, _ := drv.Ins()
	outs, _ := drv.Outs()

	inIdx := -1
	outIdx := -1

	for idx, i := range ins {
		if i.String() == inputName {
			inIdx = idx
			break
		}
	}

	for idx, o := range outs {
		if o.String() == outputName {
			outIdx = idx
			break
		}
	}

	return inIdx, outIdx
}
