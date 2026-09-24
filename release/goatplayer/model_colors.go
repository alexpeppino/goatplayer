package main

import (
	"fmt"
	"slices"
)

// func (this *mColors) SetLine(tableLine *tTableLine, st string, iKey uint16, style *cStyle) {
// 	tableLine.text = st
// 	tableLine.normal = styleBlackWhite.EchoStyle("  ", 2) + style.EchoStyle(st, nCols-2)
// 	tableLine.selected = styleBlackWhite.EchoStyle(" ▶", 2) + style.EchoStyle(st, nCols-2)
// 	tableLine.active = styleBlackWhite.EchoStyle("  ", 2) + style.EchoStyle(st, nCols-2)
// 	tableLine.selectedActive = styleBlackWhite.EchoStyle(" ▶", 2) + style.EchoStyle(st, nCols-2)
// 	tableLine.iKey = iKey
// }

// func (this *mColors) AppendLine(style *cStyle) {
// 	trace2.Begin_n(this, "AppendLine")
// 	width := this.w.nCols - 3
// 	this.w.lines.Append(&tTableLine{
// 		style.name,
// 		styleBlackWhite.EchoStyle("   ", 0) + style.EchoStyle(style.name, width),
// 		styleBlackWhite.EchoStyle(" ▶ ", 0) + style.EchoStyle(style.name, width),
// 		"", "", this.iKeyCounter})
// 	trace2.End()
// }

// func (this *mColors) AppendKeyAndLine(style *cStyle) {
// 	trace2.Begin_n(this, "AppendKeyAndLine")
// 	this.AppendKey(style.name)
// 	this.AppendLine(style)
// 	trace2.End()
// }

func (this *mColors) ModifyAndPrintLine(iLine uint16, style *cStyle) {
// trace.BeginAdd_n(this, "ModifyAndPrintLine", "iLine: ", iLine) // t //
	line := &this.w.lines.sl[iLine]
	width := this.w.nCols - 2
	line.normal = styleBlackWhite.EchoStyle("  ", 0) + style.EchoStyle(" "+style.name, width)
	line.selected = styleBlackWhite.EchoStyle("▶ ", 0) + style.EchoStyle(" "+style.name, width)
	this.w.PrintLine(iLine)
	Term.Flush()
	if style == &styleProgressBar {
// trace.Print("progress bar") // t //
		styleProgressBarInv = &cStyle{
			styleProgressBar.backColor,
			styleProgressBar.textColor,
			0, 0, ALIGN_LEFT, FONT_NORMAL, "", nil, nil}
	}
// trace.End() // t //
}

func (this *mColors) PrepareView() {
// trace.Begin_n(this, "PrepareView") // t //
	this.Reset()
	width := this.w.nCols - 2
	for _, style := range AllMaps.allConfigurableStyles {
		this.AppendKey(style.name)
		this.w.lines.Append(&tTableLine{
			style.name,
			styleBlackWhite.EchoStyle("  ", 0) + style.EchoStyle(" "+style.name, width),
			styleBlackWhite.EchoStyle("▶ ", 0) + style.EchoStyle(" "+style.name, width),
			"", "", this.iKeyCounter})
	}
	styleProgressBarInv = &cStyle{
		styleProgressBar.backColor,
		styleProgressBar.textColor,
		0, 0, ALIGN_LEFT, FONT_NORMAL, "", nil, nil}
	this.w.SetEmptyLineForDialog()
	if this.m2.w.canvasEcho == "" {
		this.PreparePalette()
	}
	this.SelectKeyAndLine(this.keys.iSelected)
// trace.End() // t //
}

func (this *mColors) SetOuter() {
	this.outer = this
	this.mBasic.outer = this
}

func (this *mColors) InitModel(name string, focusKey string, helpfile string) {
	if this.outer == nil {
		this.SetOuter()
	}
	this.mBasic.InitModel(name, focusKey, helpfile)
}

func (this *mColors) PrepareModel() {
	this.mBasic.PrepareModel()
}

func (this *mColors) SelectKey_o() {
// trace.Begin_n(this, "SelectKey_o") // t //
	var s1, s2 string
	iSel := this.keys.iSelected
	style := AllMaps.allConfigurableStyles[iSel]
	if style.textColor == UNSET {
		s1 = fmt.Sprintf("Text color is inherited from %s.", style.otherStyleTextColor.name)
	} else {
		s1 = fmt.Sprintf("Text color is %03d.", style.textColor)
	}
	if style.backColor == UNSET {
		s2 = fmt.Sprintf("Background color is inherited from %s.", style.otherStyleBackColor.name)
	} else {
		s2 = fmt.Sprintf("Background color is %03d.", style.backColor)
	}
	s3 := "Pick a text color and a background color."
	Term.AddStringWithPos(this.w.marginTop+20, this.w.marginLeft+38, styleNormal.EchoStyle(s3, 60))
	Term.AddStringWithPos(this.w.marginTop+22, this.w.marginLeft+38, styleNormal.EchoStyle(s1, 60))
	Term.AddStringWithPos(this.w.marginTop+23, this.w.marginLeft+38, styleNormal.EchoStyle(s2, 60))
	Term.Flush()
// trace.End() // t //
}

func (this *mColors) BringToFront_o() {
// trace.Begin_n(this, "BringToFront_o") // t //
	Term.AddSimple(this.m2.w.canvasEcho)
	Term.Flush()
	this.SelectKey_o()
// trace.End() // t //
}

func (this *mColors) ModifyFgColor(iColor uint16) {
	iSel := this.keys.iSelected
	style := AllMaps.allConfigurableStyles[iSel]
	if style.textColor == UNSET {
		Term.StatusLine.PrintMessage("Can't modify this color")
		return
	}
	style.textColor = iColor
	AllMaps.mapStyleSettings[style.name].TextColor = iColor

	// this.ModifyAndPrintLine(iSel, style)
	this.keys.iOldSelected = 9999
	this.PrepareView()
	this.w.PrintWindowAnyway()
	this.SelectKey_o()
}

func (this *mColors) ModifyBgColor(iColor uint16) {
	iSel := this.keys.iSelected
	style := AllMaps.allConfigurableStyles[iSel]
	if style.backColor == UNSET {
		Term.StatusLine.PrintMessage("Can't modify this color")
		return
	}
	style.backColor = iColor
	AllMaps.mapStyleSettings[style.name].BackColor = iColor

	// this.ModifyAndPrintLine(iSel, style)
	this.keys.iOldSelected = 9999
	this.PrepareView()
	this.w.PrintWindowAnyway()
	this.SelectKey_o()
}

func (this *mColors) ResetColors() {
// trace.BeginEnd_n(this, "ResetColors") // t //
	for _, style := range AllMaps.allConfigurableStyles {
// trace.Print("%3d %3d", style.textColor, style.defTextColor) // t //
		style.textColor = style.defTextColor
		style.backColor = style.defBackColor
		AllMaps.mapStyleSettings[style.name].TextColor = style.textColor
		AllMaps.mapStyleSettings[style.name].BackColor = style.backColor
	}
	this.Reset()
	this.PrepareView()
	this.SelectKeyAndLine(0)
	this.w.PrintWindowAnyway()
	this.BringToFront_o()
}

func (this *mColors) PreparePalette() {
	var iRow uint16 = this.w.marginTop + 2
	var iCol uint16 = this.w.marginLeft + 38
	var iColor uint16
	var st string
	var i, j uint16
	var white = []uint16{0, 8, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 88, 89, 232, 233, 234, 235, 236, 237, 238, 239, 240, 241, 242, 243, 244, 245, 246}
	// var iFgColor = "15"
	// var iBgColor = "16"
	var w2 *vWindow = &this.m2.w
	AddEmptyRect(tRect{w2.marginTop + 2, w2.marginBottom - 2, w2.marginLeft + 1, w2.marginRight - 1}, &styleNormal)
	this.iFgFirstRow = iRow + 1
	this.iFgFirstCol = iCol
	this.iBgFirstRow = iRow + 19
	this.iBgFirstCol = iCol
	//	All fg colors
	var iText uint16
	for i = range 16 {
		//	for each line
		st = ""
		for j = range 16 {
			//	for each color
			iColor = i*16 + j
			if slices.Contains(white, iColor) {
				iText = 15
			} else {
				iText = 16
			}
			// if Set.ShowNumbers.Value {
			st += fmt.Sprintf("\x1b[38;5;%d;48;5;%dm %03d ", iText, iColor, iColor)
			// } else {
			// 	st += fmt.Sprintf("\x1b[48;5;%dm     ", iColor)
			// }
		}
		Term.AddStringWithPos(uint16(i+iRow+1), iCol, st)
	}
	Term.SlurpContents(&this.m2.w.canvasEcho)
}

func (this *mColorsPicker) SelectKey_o() {}

type mColorsPicker struct {
	mBasic
}

type mColors struct {
	mBasic
	iFgFirstRow uint16
	iFgFirstCol uint16
	iBgFirstRow uint16
	iBgFirstCol uint16
	m2          mColorsPicker
}

