package preset

import (
	"bytes"
	"strings"

	"github.com/gobuffalo/packr"
	"github.com/sachaos/go-life/format/life106"
)

type Preset struct {
	Name  string
	Cells [][]bool
}

//go:generate packr

func LoadPresets() ([]Preset, error) {
	box := packr.NewBox("./files")

	var presets []Preset
	for _, name := range box.List() {
		byte, err := box.MustBytes(name)
		if err != nil {
			return nil, err
		}

		r := bytes.NewReader(byte)

		cells := life106.Parse(r)
		presets = append(presets, Preset{
			Name:  strings.Split(name, ".")[0],
			Cells: cells,
		})
	}
	return presets, nil
}
