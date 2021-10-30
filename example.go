package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/mitchellh/ioprogress"
)

type FakeReader struct {
	Length  int
	Step    int
	hasRead int
}

func (f *FakeReader) Read(p []byte) (int, error) {
	if f.hasRead >= f.Length {
		return 0, io.EOF
	}
	// Sleep for a while so we can see the animation.
	<-time.After(time.Second)
	f.hasRead += f.Step
	return f.Step, nil
}

func (f *FakeReader) Size() int {
	return f.Length
}

func drawDefault() {
	r := &FakeReader{Length: 100, Step: 10}
	progressR := &ioprogress.Reader{
		Reader: r,
		Size:   int64(r.Size()),
	}
	var buff bytes.Buffer
	io.Copy(&buff, progressR)
}

func drawBar() {
	r := &FakeReader{Length: 100, Step: 10}
	bar := ioprogress.DrawTextFormatBar(20)
	progressR := &ioprogress.Reader{
		Reader:   r,
		Size:     int64(r.Size()),
		DrawFunc: ioprogress.DrawTerminalf(os.Stdout, bar),
	}
	var buff bytes.Buffer
	io.Copy(&buff, progressR)
}

func drawBytes() {
	r := &FakeReader{Length: 100, Step: 10}
	progressR := &ioprogress.Reader{
		Reader:   r,
		Size:     int64(r.Size()),
		DrawFunc: ioprogress.DrawTerminalf(os.Stdout, ioprogress.DrawTextFormatBytes),
	}
	var buff bytes.Buffer
	io.Copy(&buff, progressR)
}

func main() {
	fmt.Println("default draw...")
	drawDefault()
	fmt.Println("draw bar...")
	drawBar()
	fmt.Println("draw bytes...")
	drawBytes()
}
