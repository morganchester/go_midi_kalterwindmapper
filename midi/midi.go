package midi

import (
	"gitlab.com/gomidi/rtmididrv"
)

func GetDevices() (inputs []string, outputs []string) {
	drv, err := rtmididrv.New()
	if err != nil {
		return []string{}, []string{}
	}
	defer drv.Close()

	ins, _ := drv.Ins()
	outs, _ := drv.Outs()

	for _, in := range ins {
		inputs = append(inputs, in.String())
	}
	for _, out := range outs {
		outputs = append(outputs, out.String())
	}

	return inputs, outputs
}
