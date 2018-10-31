package rle

import (
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	f, err := os.Open("./test/glider.rle")
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

func TestParseReplicator(t *testing.T) {
	f, err := os.Open("./test/replicator.rle")
	if err != nil {
		t.Fatal(err.Error())
	}
	b := Parse(f)
	expected := [][]bool{
		{false, false, true, true, true},
		{false, true, false, false, true},
		{true, false, false, false, true},
		{true, false, false, true, false},
		{true, true, true, false, false},
	}

	for i, row := range expected {
		for j, expectedState := range row {
			if b[i][j] != expectedState {
				t.Errorf("Expected %v, But actually %v", expected, b)
			}
		}
	}
}
