package main

import "fmt"

func (this *tDisplaySection) AddLineToPrint(iPart, iLine uint16, s string) *vTerm {
	switch iPart {
	case 1:
		if iLine == 4 {
			var strPaused, strMuted string
			if Play.isPaused {
				strPaused = "[Paused]"
			}
			if Play.isMuted {
				strMuted = "[Muted]"
			}
			var strRepeat = Set.RepeatMode.GetValue()
			this.Part1[4] = fmt.Sprintf("Volume:  %-3s   [%3s]   %8s   %7s", ui_pl.Volume, strRepeat, strPaused, strMuted)
		} else {
			this.Part1[iLine] = s
		}
		Term.AddStringWithPos(this.marginTop+1+iLine, this.part1_left, styleDisplay.EchoStyle("  "+this.Part1[iLine], this.part1_width))
	case 2:
		this.Part2[iLine] = s
		Term.AddStringWithPos(this.marginTop+1+iLine, this.part2_left, styleDisplay.EchoStyle("  "+this.Part2[iLine], this.part2_width))
	case 3:
		this.Part3[iLine] = s
		Term.AddStringWithPos(this.marginTop+1+iLine, this.part3_left, styleDisplay.EchoStyle("  "+this.Part3[iLine], this.part3_width))
	}
	return &Term
}

func (this *tDisplaySection) AddPartToPrint(iPart uint16) {
// trace.Begin_n(this, "Print") // t //
	var i uint16
// trace.Print("rect:  %+v", this.tRect) // t //
	Term.AddStringWithPos(this.marginTop, this.marginLeft, this.frame)
	switch iPart {
	case 1:
		Term.AddStringWithPos(this.marginTop+1+i, this.part1_left, styleDisplayMessage.EchoStyle("  "+this.Part1[i], this.part1_width))
		for i = 1; i < 5; i++ {
			Term.AddStringWithPos(this.marginTop+1+i, this.part1_left, styleDisplay.EchoStyle("  "+this.Part1[i], this.part1_width))
		}
		go func() {
			ti2++
			utils.Sleep(5)
			ti2--
			if ti2 == 0 {
				Term.AddStringWithPos(this.marginTop+1, this.part1_left, styleDisplay.EchoStyle("  "+this.Part1[0], this.part1_width))
				// Term.AddStringWithPos(top+1, left, styleDisplay.EchoStyle(s1, 78))
				Term.Flush()
			}
		}()
		for i = range 5 {
			Term.AddStringWithPos(this.marginTop+1+i, this.part3_left, styleDisplay.EchoStyle("  ▏  "+this.Part3[i], this.part3_width))
		}
	case 2:
		for i = range 5 {
			Term.AddStringWithPos(this.marginTop+1+i, this.part2_left, styleDisplay.EchoStyle("  ▏  "+this.Part2[i], this.part2_width))
		}
		for i = range 5 {
			Term.AddStringWithPos(this.marginTop+1+i, this.part3_left, styleDisplay.EchoStyle("  ▏  "+this.Part3[i], this.part3_width))
		}
	case 3:
		for i = range 5 {
			Term.AddStringWithPos(this.marginTop+1+i, this.part3_left, styleDisplay.EchoStyle("  ▏  "+this.Part3[i], this.part3_width))
		}
	}
// trace.End() // t //
}

func (this *tDisplaySection) AddFrameToPrint() {
	Term.AddStringWithPos(this.marginTop, this.marginLeft, this.frame)
}

func (this *tDisplaySection) Prepare(rect tRect) {
// trace.Begin_n(this, "Prepare") // t //
	this.tRect = rect
// trace.Print("rect:  %+v", this.tRect) // t //
	this.frame = Term.MakeFrame_Rect(this.tRect, &styleDisplay)
// trace.Print("frame:  %#v", this.frame) // t //
	this.part1_left = 1
	this.part1_width = 80
	this.part2_left = 81
	this.part2_width = 60
	this.part3_left = 141
	this.part3_width = this.marginRight - this.part3_left
// trace.End() // t //
}

type tDisplaySection struct {
	vBasicView
	frame       string
	Part1       [5]string
	Part2       [5]string
	Part3       [5]string
	part1_left  uint16
	part1_width uint16
	part2_left  uint16
	part2_width uint16
	part3_left  uint16
	part3_width uint16
}

