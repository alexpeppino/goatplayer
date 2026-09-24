package main

import (
	"fmt"
	"time"
)

func (this *mReport) SetOuter() {
	this.mBasic.outer = this
	this.outer = this
}

func (this *mReport) AddReportLine(format string, a ...any) {
	t := time.Now()
	s := t.Format("15:04:05 ")
	this.keys.iActive = UNSET
	this.AppendKeyAndLine(s+fmt.Sprintf(format, a...), nil)
	if this.w.IsInForeground() {
		this.keys.iSelected = 0
		this.w.lines.iSelected = 0
		this.w.PrintWindowAnyway()
	}
}

func (this *mReport) AddErrorLine(s string) {
	this.keys.iActive = UNSET
	this.AppendKeyAndLine(s, nil)
	if this.w.IsInForeground() {
		this.keys.iSelected = 0
		this.w.lines.iSelected = 0
		this.w.PrintWindowAnyway()
	}
}

func (this *mReport) PrepareView() {
	this.Reset()
	this.keys.iSelected = 0
	this.w.lines.iSelected = 0
	this.w.SetEmptyLine()
// trace.BeginEndAdd_n(this, "PrepareView", "%d lines", this.w.lines.Len()) //@1 // t //
}

type mReport struct {
	mBasic
}

