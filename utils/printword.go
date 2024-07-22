package utils

import (
	"fmt"
	"strings"
)

// printWord prints a word using predefined lines from the content file
func PrintWord(word string, contentLines []string) (st string) {
	linesOfSlice := make([]string, 9)
	for _, v := range word {
		for i := 1; i <= 9; i++ {
			linesOfSlice[i-1] += contentLines[int(v-32)*9+i]
		}
	}
	st += fmt.Sprint(strings.Join(linesOfSlice, "\n"))
	return
}
