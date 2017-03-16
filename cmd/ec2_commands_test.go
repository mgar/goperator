package cmd

import "testing"

func TestStringsSliceToStringPointersSlice(t *testing.T) {
	inputSlice := []string{
		"i-03bd8c8658f294cc5",
		"i-0692345ea5c5c6d51",
	}
	outputSlice := stringsSliceToStringPointersSlice(inputSlice)

	for item := range inputSlice {
		if inputSlice[item] != *outputSlice[item] {
			t.Errorf("Slices doesn't match: %s is not equal to %s", inputSlice[item], *outputSlice[item])
		}
	}
}
