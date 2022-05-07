package write_line

import (
	"fmt"
	"io"
	"strings"
	"tiktaalik-honeypot-go/src/write-line/colors"
	"time"
)

const Welcome = "嗨嗨害,来了嗷~"

func Write(w io.Writer, str string) {
	chars := strings.Split(str, "")

	for _, ch := range chars {
		fmt.Fprint(w, ch)
		time.Sleep(30 * time.Millisecond)
	}
}

// ColorWrite 显示鬼畜颜色特效
func ColorWrite(w io.Writer, str, color string) {
	Write(w, color+str+colors.Reset)
}

func PrintEnd(w io.Writer, ends int) {
	for i := 0; i < ends; i++ {
		fmt.Fprint(w, "\n")
	}
}
