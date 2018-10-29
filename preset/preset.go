package preset

import (
	"bytes"
	"sort"
	"strings"

	"github.com/gobuffalo/packr"
	"github.com/sachaos/go-life/format/life106"
)

type Preset struct {
	Name  string
	Cells [][]bool
}

func (p *Preset) Size() (int, int) {
	height := len(p.Cells)
	if height == 0 {
		return 0, 0
	}

	width := len(p.Cells[0])
	return width, height
}

//go:generate packr

func LoadPresets() ([]Preset, error) {
	box := packr.NewBox("./files")

	var names []string
	for _, name := range box.List() {
		names = append(names, name)
	}

	sort.Strings(names)

	var presets []Preset
	for _, name := range names {
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
