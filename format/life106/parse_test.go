package life106

import (
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	f, err := os.Open("./test/glider.life")
	if err != nil {
		t.Fatal(err.Error())
	}
	b := Parse(f)
	expected := [][]bool{
		{false, true, false},
		{false, false, true},
		{true, true, true},
	}

	for i, row := range expected {
		for j, expectedState := range row {
			if b[i][j] != expectedState {
				t.Errorf("Expected %v, But actually %v", expected, b)
			}
		}
	}
}
