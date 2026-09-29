package main

import (
	"fmt"
	"goatplayer/internal/K"
	"time"
)

func (this *mReport) SetOuter() {
	this.mBasic.outer = this
	this.outer = this
}

func (this *mReport) AddReportLine(format string, a ...any) {
	t := time.Now()
	s := t.Format("15:04:05 ")
	this.keys.ActiveI = K.UNSET
	this.AppendKeyAndLine(s+fmt.Sprintf(format, a...), nil)
	if this.w.IsInForeground() {
		this.keys.SelectedI = 0
		this.w.lines.SelectedI = 0
		this.w.PrintWindowAnyway()
	}
}

func (this *mReport) AddErrorLine(s string) {
	this.keys.ActiveI = K.UNSET
	this.AppendKeyAndLine(s, nil)
	if this.w.IsInForeground() {
		this.keys.SelectedI = 0
		this.w.lines.SelectedI = 0
		this.w.PrintWindowAnyway()
	}
}

func (this *mReport) PrepareView() {
	this.Reset()
	this.keys.SelectedI = 0
	this.w.lines.SelectedI = 0
	this.w.SetEmptyLine()
	// trace.BeginEndAdd_n(this, "PrepareView", "%d lines", this.w.lines.Len()) //@1 // t //
}

type mReport struct {
	mBasic
}
