package display

import (
	"errors"
	"log"
	"strings"
)

type (
	LCD interface {
		Open() error
		Write(line Line, text string) error
		Enable(yes bool) error
		Listen(l func(btn int, released bool) bool)
		Close() error
	}
	Line int
)

var (
	ErrClosed            = errors.New("display closed")
	ErrDisplayNotWorking = errors.New("display not working")
	ErrMsgSizeMismatch   = errors.New("msg size mismatch")

	filledSquare = string([]byte{0xff})
)

const (
	LineOne    Line = 0
	LineTwo    Line = 1
	DefaultTTy      = "/dev/ttyS1"
	c16             = 16
)

func Find() *LCD {
	lcd, err := NewQnapLCD("")
	if err == nil {
		return &lcd
	} else {
		log.Println(err)
	}
	return nil
}

func prepareTxt(txt string) string {
	l := len(txt)
	if l > c16 {
		txt = txt[0:c16]
	} else if l < c16 {
		txt += strings.Repeat(" ", c16-l)
	}
	return txt
}

func Progress(perc int) string {
	chars := percentOf(c16, 100, perc)
	return strings.Repeat(filledSquare, chars) + strings.Repeat("-", c16-chars)
}

func percentOf(maxVal, maxPercent, currentPercent int) int {
	return (maxVal * currentPercent) / maxPercent
}
