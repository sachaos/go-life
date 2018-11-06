package format

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/sachaos/go-life/format/life106"
	"github.com/sachaos/go-life/format/rle"
)

const (
	Life105 = iota
	Life106
	RLE
)

func Parse(r io.Reader) ([][]bool, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	formatReader := bytes.NewReader(b)
	contentReader := bytes.NewReader(b)
	switch DetectFormat(formatReader) {
	case Life106:
		return life106.Parse(contentReader), nil
	case Life105:
		return nil, fmt.Errorf("Life1.05 is not implemented")
	default:
		return rle.Parse(contentReader), nil
	}
}

func DetectFormat(r io.Reader) int {
	br := bufio.NewReader(r)

	firstLine, _, _ := br.ReadLine()
	firstLineStr := string(firstLine)

	switch firstLineStr {
	case "#Life 1.05":
		return Life105
	case "#Life 1.06":
		return Life106
	default:
		return RLE
	}
}
