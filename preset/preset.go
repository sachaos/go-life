package preset

import (
	"bytes"
	"sort"
	"strings"

	"github.com/gobuffalo/packr"
	"github.com/sachaos/go-life/format/life106"
	"github.com/sachaos/go-life/format/rle"
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
	for _, filename := range names {
		byte, err := box.Find(filename)
		if err != nil {
			return nil, err
		}

		filenameSlice := strings.Split(filename, ".")
		name := filenameSlice[0]
		format := filenameSlice[1]

		r := bytes.NewReader(byte)

		var cells [][]bool
		if format == "rle" {
			cells = rle.Parse(r)
		} else {
			cells = life106.Parse(r)
		}

		presets = append(presets, Preset{
			Name:  strings.Split(name, ".")[0],
			Cells: cells,
		})
	}
	return presets, nil
}
