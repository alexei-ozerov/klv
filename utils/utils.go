package utils

import "github.com/muesli/reflow/wordwrap"

func WrapText(text string, width int) string {
	return wordwrap.String(text, width)
}
