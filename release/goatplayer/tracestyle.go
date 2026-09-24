package main

import "fmt"

type cTraceStyle struct {
	textColor int
	// backColor int
	fontStyle string
}

func (this *cTraceStyle) JustTheColor() string {
	return fmt.Sprintf("\x1b[38;5;%dm", this.textColor)
}

func (this *cTraceStyle) Bold(s string) string {
	return DC_BOLD + s + DC_RESET + traceStyleNormal.JustTheColor()
}

func (this *cTraceStyle) Echo(s string) string {
	s = fmt.Sprintf("\x1b[38;5;%dm", this.textColor) +
		this.fontStyle + s
	if true { //this.font != FONT_NORMAL {
		s += FONT_RESET
	}
	return s
}

