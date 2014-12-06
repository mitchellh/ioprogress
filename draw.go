package ioprogress

import (
	"fmt"
	"io"
	"os"
)

// DrawFunc is the callback type for drawing progress.
type DrawFunc func(int64, int64) error

// DrawTextFormatFunc is a callback used by DrawFuncs that draw text in
// order to format the text into some more human friendly format.
type DrawTextFormatFunc func(int64, int64) string

var defaultDrawFunc DrawFunc

func init() {
	defaultDrawFunc = DrawTerminal(os.Stdout)
}

// DrawTerminal returns a DrawFunc that draws a progress bar to an io.Writer
// that is assumed to be a terminal (and therefore respects carriage returns).
func DrawTerminal(w io.Writer) DrawFunc {
	return DrawTerminalf(w, func(progress, total int64) string {
		return fmt.Sprintf("%d/%d", progress, total)
	})
}

// DrawTerminalf returns a DrawFunc that draws a progress bar to an io.Writer
// that is formatted with the given formatting function.
func DrawTerminalf(w io.Writer, f DrawTextFormatFunc) DrawFunc {
	return func(progress, total int64) error {
		if progress == -1 || total == -1 {
			_, err := fmt.Fprintf(w, "\n")
			return err
		}

		_, err := fmt.Fprint(w, f(progress, total)+"\r")
		return err
	}
}

var byteUnits = []string{"B", "KB", "MB", "GB", "TB", "PB"}

// DrawTextFormatBytes is a DrawTextFormatFunc that formats the progress
// and total into human-friendly byte formats.
func DrawTextFormatBytes(progress, total int64) string {
	return fmt.Sprintf("%s/%s", byteUnitStr(progress), byteUnitStr(total))
}

func byteUnitStr(n int64) string {
	var unit string
	size := float64(n)
	for i := 1; i < len(byteUnits); i++ {
		mult := float64(i * 1000)
		if size < mult {
			unit = byteUnits[i-1]
			break
		}

		size = size / mult
	}

	return fmt.Sprintf("%.3g %s", size, unit)
}
