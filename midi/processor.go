package midi

import (
	"fmt"
	"sync"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/rtmididrv"
)

type NoteOffAction struct {
	Action  string        `yaml:"action"`
	Message []interface{} `yaml:"message"`
}

type Mapping struct {
	NoteOn   map[int][]interface{} `yaml:"noteon"`
	NoteOff  NoteOffAction         `yaml:"noteoff"`
	Controls map[int][]interface{} `yaml:"controls"`
}

type Processor struct {
	drv     *rtmididrv.Driver
	in      midi.In
	out     midi.Out
	running bool
	mu      sync.Mutex
}

func NewProcessor() *Processor {
	return &Processor{}
}

func (p *Processor) Start(inputIdx, outputIdx int, mapping Mapping) error {
	if TestMode {
		fmt.Println("[TEST] Старт маршрутизации")
		fmt.Println("[TEST] Mapping:", mapping)
		return nil
	}
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.running {
		return fmt.Errorf("already running")
	}

	drv, err := rtmididrv.New()
	if err != nil {
		return err
	}
	p.drv = drv

	ins, _ := drv.Ins()
	outs, _ := drv.Outs()

	if inputIdx < 0 || inputIdx >= len(ins) {
		return fmt.Errorf("invalid input device index")
	}
	if outputIdx < 0 || outputIdx >= len(outs) {
		return fmt.Errorf("invalid output device index")
	}

	p.in = ins[inputIdx]
	p.out = outs[outputIdx]

	if err := p.in.Open(); err != nil {
		return err
	}
	if err := p.out.Open(); err != nil {
		return err
	}

	// Слушатель MIDI сообщений
	p.in.SetListener(func(data []byte, deltaMicroseconds int64) {
		if len(data) < 3 {
			return
		}
		status := data[0] & 0xF0
		key := data[1]
		val := data[2]

		switch status {
		case 0x90: // NoteOn
			if val > 0 {
				if cmd, ok := mapping.NoteOn[int(key)]; ok {
					fmt.Printf("[MAP] NoteOn %d -> %v\n", key, cmd)
					execute(p.out, cmd, nil)
				}
			} else {
				if mapping.NoteOff.Action == "send" {
					execute(p.out, mapping.NoteOff.Message, nil)
				}
			}
		case 0x80: // NoteOff
			if mapping.NoteOff.Action == "send" {
				execute(p.out, mapping.NoteOff.Message, nil)
			}
		case 0xB0: // ControlChange
			cc := key
			if cmd, ok := mapping.Controls[int(cc)]; ok {
				execute(p.out, cmd, &val)
			}
		}
	})

	p.running = true
	return nil
}

func (p *Processor) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if !p.running {
		return
	}
	p.in.Close()
	p.out.Close()
	p.drv.Close()
	p.running = false
}

func execute(out midi.Out, message []interface{}, dynamicValue *uint8) {
	if len(message) == 0 {
		return
	}

	cmd := message[0].(string)
	switch cmd {
	case "cc":
		if len(message) < 2 {
			fmt.Println("[ERROR] Invalid CC mapping")
			return
		}
		cc := byte(message[1].(int))
		var val byte
		if len(message) >= 3 {
			val = byte(message[2].(int))
		} else if dynamicValue != nil {
			val = *dynamicValue
		}

		// Формируем CC-сообщение вручную: [status, controller, value]
		msg := []byte{0xB0, cc, val}
		out.Write(msg)
	default:
		fmt.Println("[WARN] Unknown command:", cmd)
	}
}
