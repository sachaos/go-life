package format

import (
	"os"
	"testing"
)

func TestDetectFormat(t *testing.T) {
	testCases := []struct {
		Name           string
		Filename       string
		ExpectedFormat int
	}{
		{
			Name:           "life106",
			Filename:       "./life106/test/glider.life",
			ExpectedFormat: Life106,
		},
		{
			Name:           "rle",
			Filename:       "./rle/test/glider.rle",
			ExpectedFormat: RLE,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			f, err := os.Open(testCase.Filename)
			if err != nil {
				t.Fatal(err)
			}

			detectedFormat := DetectFormat(f)
			if detectedFormat != testCase.ExpectedFormat {
				t.Errorf("ExpectedFormat: %v, but actually: %v", testCase.ExpectedFormat, detectedFormat)
			}
		})
	}
}
