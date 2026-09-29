package main

import (
	"fmt"
	"goatplayer/internal/K"
	AF "goatplayer/internal/audiofiles"
	"goatplayer/internal/echo"
	"goatplayer/internal/small/utils"
	"goatplayer/internal/style"
	"goatplayer/internal/table"
	"slices"
)

//////	Window

func (this *vWindow) GetName() string { return this.name }

func (this *vWindow) InitWindow(model *mBasic) {
	this.name = model.name + ".w"
	this.pModel = model
	this.colOffset = 1
	this.rowOffset = 3
	this.pageOffset = 1
	this.iOldPage = 1
	this.fGetWidth = func() uint16 {
		// trace.Print("%s  %d", this.name, this.nCols)
		return this.nCols
	}
}

func (this *vWindow) PrepareWindow() {
	this.nRows = this.marginBottom - this.marginTop + 1 - this.rowOffset - 3
	this.iStatusLineRow = this.marginBottom - 1
	this.nCols = this.marginRight - this.marginLeft - 1
	if !this.IsSingleWindow() {
		this.nRows -= 2
		this.rowOffset += 2
	}
	this.SetEmptyLine()
	// trace.BeginEndAdd_n(this, "PrepareWindow", "this.fGetWidth: %d", this.fGetWidth())
}

func (this *vWindow) BringWindowToFront(activate bool) {
	// trace.BeginAdd_n(this, "BringWindowToFront", "activate: %v", activate) // t //
	if this.isCanvas {
		return
	}
	this.SetActiveToSection()
	// this.pLayer.pSection.activeWindow = this
	// this.pLayer.pSection.pScreen.activeWindow = this
	if activate {
		this.SetActiveToScreen()
		this.pModel.outer.BringToFront_o()
	}
	this.pLayer.BringLayerToFront()
	// trace.Print("MainScreen.pRightSection.name: %s", MainScreen.pRightSection.name) // t //
	// trace.Print("MainScreen.PlayerSection.GetActiveWindow: %s", MainScreen.PlayerSection.GetActiveWindow().name) // t //
	// trace.End() // t //
}

/*
Used when the slice of lines already exists. iLine: index of the line. st: text of the line. iKey: index of the related key. options: 1= use grey for the normal line.
*/
func (this *vWindow) SetLineForDialog(iLine uint16, st string, iKey uint16, option uint16) {
	// trace.N_BeginEnd(this, "SetLine, %d, %s", iLine, st)
	line := this.lines.Get(iLine)
	nCols := this.nCols
	utils.AdjustWidth(&st, nCols)
	line.Normal = style.Normal.EchoStyle(st, nCols)
	line.Selected = style.Selected.EchoStyle(st, nCols)
	line.Active = style.Active.EchoStyle(st, nCols)
	line.SelectedActive = style.SelActive.EchoStyle(st, nCols)
	line.IKey = iKey
}

func (this *vWindow) MakeActiveFilterLine(s string) string {
	nCols := this.nCols
	return style.ActFilter.EchoStyle(s, nCols)
}

func (this *vWindow) MakeSelectedActiveFilterLine(s string) string {
	nCols := this.nCols
	return style.SelActFilter.EchoStyle(s, nCols)
}

func (this *vWindow) SetLine(iLine uint16, st string, iKey uint16, option uint16) {
	// trace.N_BeginEnd(this, "SetLine, %d, %s", iLine, st)
	line := this.lines.Get(iLine)
	utils.AdjustWidth(&st, this.nCols)
	if option == 1 {
		line.Normal = style.Grey.EchoStyle(st, this.nCols)
	} else {
		line.Normal = style.Normal.EchoStyle(st, this.nCols)
	}
	line.Text = st
	line.Selected = style.Selected.EchoStyle(st, this.nCols)
	line.Active = style.Active.EchoStyle(st, this.nCols)
	line.SelectedActive = style.SelActive.EchoStyle(st, this.nCols)
	line.IKey = iKey
}

func (this *cFilters) Or(audiof *AF.TAudiofile) bool {
	for i := range this.sl {
		this.b[i] = this.sl[i](audiof)
	}
	return slices.Contains(this.b, true)
}

func (this *cFilters) And(audiof *AF.TAudiofile) bool {
	for i := range this.sl {
		this.b[i] = this.sl[i](audiof)
	}
	return !slices.Contains(this.b, false)
}

func (this *cLogical) AddItem(b bool) {
	this.b = append(this.b, b)
}

func (this *cLogical) Or() bool {
	return slices.Contains(this.b, true)
}

func (this *cLogical) And() bool {
	return !slices.Contains(this.b, false)
}

// Given a terminal row index, selects a table row and reprints.
func (this *vWindow) SelectLineByTRow(iTRow uint16) {
	iLine := this.CalcLineIndexByTRow(iTRow)
	// trace.BeginAdd_n(this.pModel, "SelectLineByTRow", "iTRow: %d, iLine: %d", iTRow, iLine) // t //
	p := this.lines.Get(iLine)
	if p == nil || p.IKey == K.UNSET {
		// trace.N_Error(this, "SelectLineByTRow, invalid index %d", iRow)
		// trace.ReturnAdd("invalid index") // t //
		return
	}
	this.SelectLine(iLine)
	this.pModel.keys.TraceSelected("")
	this.lines.TraceSelected("")
	this.PrintWindow()
	// trace.End() // t //
}

// Give key index, get line index.
func (this *vWindow) GetLineByKey(iKey uint16) uint16 {
	if iKey == K.UNSET {
		// trace.Error_n(this.pModel, "GetLineByKey, K.UNSET index, %d", iKey) // t //
		return K.UNSET
	}
	if key := this.pModel.keys.Get(iKey); key != nil {
		if this.isFolded {
			return key.LineI
		} else {
			return iKey
		}
	} else {
		return K.UNSET
	}

	// trace.N_BeginEnd(this, "GetLineByKey")
	// trace.Print("iKey: %d", iKey) //@1
	// trace.Print("this.keys[iKey].ILine: %d", this.keys[iKey].ILine) //@1
	// return this.keys[iKey].ILine
}

// Reprints the new and old active lines.
func (this *vWindow) UpdateActiveLines() {
	// trace.BeginSilent_n(this, "UpdateActiveLines") // t //
	model := this.pModel
	// trace.AddToPrint("active %d (%d)", model.keys.ActiveI, model.keys.OldActiveI) // t //
	if MainScreen.LeftSection.GetActiveLayer() != Play.w.pLayer {
		// trace.Return() // t //
		return
	}
	// model.keys.TraceActive("")
	this.PrintLine(this.GetLineByKey(model.keys.OldActiveI))
	this.PrintLine(this.GetLineByKey(model.keys.ActiveI))
	echo.Flush()
	model.keys.OldActiveI = model.keys.ActiveI
	// trace.End() // t //
}

// Calculates terminal row index by line index. Returns OUT_OF_PAGE if the line is not in the page.
func (this *vWindow) CalcTRowIndex(iLine uint16) uint16 {
	var iPage = iLine / this.nRows
	if iPage != this.iPage {
		return K.OUT_OF_PAGE
	}
	return iLine%this.nRows + this.rowOffset + this.pLayer.pSection.marginTop
}

// .
func (this *vWindow) CalcLineIndexByTRow(iTRow uint16) uint16 {
	return this.iPage*this.nRows + iTRow - this.rowOffset - this.marginTop
}

// .
func (this *vWindow) AddLineToPrint(iTRow uint16, line *string) {
	echo.AddStringWithPos(iTRow, this.marginLeft+this.colOffset, *line)
}

// .
func (this *vWindow) PrintLine(iLine uint16) {
	// trace2.Print("PrintLine: %d", iLine)
	if iLine == K.UNSET {
		return
	}
	model := this.pModel
	var isActiveFilter, isActive bool
	if !this.isFolded {
		if model.keys.Sl[iLine].IsActiveFilter {
			isActiveFilter = true
		}
	}
	if this.isFolded || this.pModel.GetScreen().isDialog {
		if model.keys.ActiveI != K.UNSET {
			iActiveLine := this.GetLineByKey(model.keys.ActiveI)
			if iActiveLine == iLine {
				isActive = true
				// trace.Print("%d: isActive", iLine)
			}
		}
	}
	// trace.AddToPrint("active line: %d", iActiveLine)
	var iTRow = this.CalcTRowIndex(iLine)
	line_p := this.lines.Get(iLine)
	if iLine == this.lines.SelectedI {
		if isActive {
			this.AddLineToPrint(iTRow, &line_p.SelectedActive)
		} else if isActiveFilter {
			s := style.SelActFilter.EchoStyle(line_p.Text, this.nCols)
			this.AddLineToPrint(iTRow, &s)
		} else {
			this.AddLineToPrint(iTRow, &line_p.Selected)
		}
	} else {
		if isActive {
			this.AddLineToPrint(iTRow, &line_p.Active)
		} else if isActiveFilter {
			s := style.ActFilter.EchoStyle(line_p.Text, this.nCols)
			this.AddLineToPrint(iTRow, &s)
		} else {
			this.AddLineToPrint(iTRow, &line_p.Normal)
		}
	}
	// trace.N_BeginEndPrintAll(this, "PrintLine %d", iLine)
}

/* Remember to set the selected line */
func (this *vWindow) PrintWindowAnyway() {
	this.iOldPage = K.UNSET
	this.PrintWindow()
}

func (this *vWindow) AddEmptyWindowToPrint() {
	// trace.Begin_n(this, "AddEmptyWindowToPrint") // t //
	this.iPage = 0
	var end uint16 = this.nRows - 1
	// trace.Print("this.nRows: %d", this.nRows) // t //
	var iLine uint16
	for iLine = 0; iLine <= end; iLine++ {
		var iTRow = this.CalcTRowIndex(iLine)
		// trace2.Print("empty line at line %d trow %d, nrows %d", iLine, iTRow, this.nRows)
		echo.AddStringWithPos(iTRow, this.marginLeft+this.colOffset, K.DC_RESET+this.emptyLine)
	}
	// trace.End() // t //
}

func (this *vWindow) PrintPage() {
	// trace.BeginEnd_n(this, "PrintPage") // t //
	var start uint16 = this.iPage * this.nRows
	var end uint16 = start + this.nRows - 1
	var iLine uint16
	echo.AddStringWithPos(this.CalcTRowIndex(start)-1, this.marginLeft+this.colOffset, K.DC_RESET+this.emptyLine)
	for iLine = start; iLine <= end; iLine++ {
		iTRow := this.CalcTRowIndex(iLine)
		if iLine > this.lines.Len()-1 {
			echo.AddStringWithPos(iTRow, this.marginLeft+this.colOffset, K.DC_RESET+this.emptyLine)
			continue
		}
		line := &this.lines.Sl[iLine]
		this.AddLineToPrint(iTRow, &line.Normal)
	}

	//	PrintStatusLine
	s := fmt.Sprintf(" Filepage %d/%d ", this.iPage+1, this.nPages)
	echo.AddStringWithPos(this.iStatusLineRow, this.marginLeft+this.colOffset, style.PageStatusLine.EchoStyle(s, this.nCols))
	echo.Flush()
}

func (this *vWindow) PrintCanvas() {
	// trace.Begin_n(this, "PrintCanvas") // t //
	this.AddEmptyWindowToPrint()
	echo.AddSimple(this.canvasEcho)
	echo.Flush()
	// trace.End() // t //
}

/*
This is the generic method called by all window's subtypes.
There are 2 ways to reprint the window:
1. reprint all the lines
2. reprint only the modified lines
Prints the table. If the page has changed, prints the whole page, otherwise prints only the modified rows.

Remember to set the selected line
*/
func (this *vWindow) PrintWindow() {
	// trace.N_Begin(this, "Print")
	if this.isCanvas {
		this.PrintCanvas()
		return
	}
	if this.isPage {
		return
	}
	if len(this.lines.Sl) == 0 {
		// trace.BeginEnd_n(this, "PrintWindow, no lines to print") // t //
		this.AddEmptyWindowToPrint()
		echo.Flush()
		return
	}
	// trace.BeginSilent_n(this, "PrintWindow") // t //
	model := this.pModel
	// se := this.pLayer.pSection
	// this.layer.activeWindow = this
	// this.pLayer.pSection.PrintFocusButtons(true)
	/// Calculate page number
	var iPage uint16 = this.lines.SelectedI / this.nRows
	this.nLines = this.lines.Len()
	///
	// trace.AddToPrint("Page: %d OldPage: %d", iPage, this.iOldPage) // t //
	// trace.AddToPrint("SelectedLine: %d, OldSelectedLine: %d", this.lines.SelectedI, this.lines.OldSelectedI) // t //
	// if this.oldSelected_i == this.Selected_i {
	// 	this.Reprint()
	// 	return
	// }
	if iPage != this.iOldPage {
		//	Change page, reprint all Lines (including the empty ones)
		// trace.Print("-> change page")
		/// Window caption for multi window
		if !this.IsSingleWindow() {
			cap := SetCaption(model.shortName, model.focusKey)

			//!!
			// if this == this.pLayer.pSection.activeWindow {
			if this.IsActiveToSection() {
				cap = style.SelButton_SelFrame.EchoStyle(" "+cap, this.nCols+1)
				// trace.AddToPrint("Window caption for multi window, cap %s", cap)
				// echo.AddStringWithPos(this.rowOffset-3, this.marginLeft+this.colOffset, cap)
				echo.AddStringWithPos(this.rowOffset-2, this.marginLeft+this.colOffset, cap)
				echo.AddStringWithPos(this.rowOffset-1, this.marginLeft+this.colOffset, this.emptyLine)
			} else {
				cap = style.Button_SelFrame.EchoStyle(" "+cap, this.nCols+1)
				// trace.AddToPrint("Window caption for multi window, cap %s", cap)
				// echo.AddStringWithPos(this.rowOffset-3, this.marginLeft+this.colOffset, cap)
				echo.AddStringWithPos(this.rowOffset-2, this.marginLeft+this.colOffset, cap)
				echo.AddStringWithPos(this.rowOffset-1, this.marginLeft+this.colOffset, this.emptyLine)
			}
		}
		///
		this.iPage = iPage
		var start uint16 = iPage * this.nRows
		var end uint16 = start + this.nRows - 1
		var iLine uint16
		if this.IsSingleWindow() {
			if this.columnsCaption == "" {
				echo.AddStringWithPos(this.CalcTRowIndex(start)-1, this.marginLeft+this.colOffset, K.DC_RESET+this.emptyLine)
			} else {
				echo.AddStringWithPos(this.CalcTRowIndex(start)-1, this.marginLeft+this.colOffset, this.columnsCaption)
			}
		}
		// trace.AddToPrint("Keys %d, start %d, end %d, iPage %d, nLines: %d", model.keys.Len(), start, end, iPage, this.nLines) // t //
		for iLine = start; iLine <= end; iLine++ {
			var iTRow = this.CalcTRowIndex(iLine)
			// var iRow = i % this.rows_n + this.rowOffset
			if iLine > this.nLines-1 {
				//	Print an empty line
				// log.Println("this.PrintLine > empty line")
				echo.AddStringWithPos(iTRow, this.marginLeft+this.colOffset, K.DC_RESET+this.emptyLine)
				continue
			}
			this.PrintLine(iLine)
		}
	} else {
		// trace.AddToPrint("same page") // t //
		//	Reprint only old and new selected line
		this.PrintLine(this.lines.OldSelectedI)
		this.PrintLine(this.lines.SelectedI)
	}
	this.iOldPage = iPage
	this.lines.OldSelectedI = this.lines.SelectedI
	/// Key value long line screen bottom
	// if !this.IsSingleWindow() {
	// 	// echo.AddStringWithPos(this.marginBottom, se.marginLeft+2, styleSelTab.EchoStyle(m.keys.GetSelected().V, se.marginRight-se.marginLeft-5))
	// 	echo.AddStringWithPos(nRows, 0, style.Normal.EchoStyle("Selected: "+m.keys.GetSelected().V, nCols-22))
	// }
	///
	this.PrintStatusLine()
	// this.tableEcho = Window.contents
	Term.SetNormal()
	echo.Flush()
	// trace.End() // t //
}

/* this.statusLine must be set before to call this method. */
func (this *vWindow) PrintStatusLine() {
	model := this.pModel
	model.outer.SetStatusLine_o()
	s := fmt.Sprintf("  %d/%d %s", model.keys.SelectedI+1, model.keys.Len(), this.statusLine)
	echo.AddStringWithPos(this.iStatusLineRow, this.marginLeft+this.colOffset, style.WindowStatusLine.EchoStyle(s, this.nCols))
}

func (this *vWindow) SetEmptyLine() {
	var width uint16
	width = this.nCols
	this.emptyLine = ""
	this.emptyLine = style.EmptyLine.EchoStyle(this.emptyLine, width)
}

func (this *vWindow) SetEmptyLineForDialog() {
	this.emptyLine = ""
	this.emptyLine = style.EmptyLine.EchoStyle(this.emptyLine, this.nCols)
}

func (this *vWindow) SelectLine(iLine uint16) {
	model := this.pModel
	// trace.BeginAdd_n(model, "SelectLine", "iLine: %d", iLine) // t //
	this.lines.Select(iLine)
	if line := this.lines.Get(iLine); line != nil {
		if this.isFolded {
			model.keys.Select(line.IKey)
		} else {
			model.keys.Select(iLine)
		}
		model.outer.SelectKey_o()
	}
	// trace.End() // t //
}

/*
Appends a new line to the view. st: the text of the line. setKeyIndex: true if the key index must be set for this line.
option: nil; 1= arrow selection
*/
func (this *vWindow) AppendLine(st string, setKeyIndex bool, style1 *style.GStyle) {
	// trace.N_BeginEnd(this, "AppendLine, %s %+v", st, style)
	var iKey uint16
	if setKeyIndex {
		iKey = this.pModel.iKeyCounter
	} else {
		iKey = this.lines.Len() // it was K.UNSET
	}
	// AdjustWidth(&st, this.nCols)
	nCols := this.nCols
	if style1 == nil {
		act := " ▶ " + st[3:]
		this.lines.Append(&table.TTableLine{
			Text:           st,
			Normal:         style.Normal.EchoStyle(st, nCols),
			Selected:       style.Selected.EchoStyle(st, nCols),
			Active:         style.Active.EchoStyle(act, nCols),
			SelectedActive: style.SelActive.EchoStyle(act, nCols),
			IKey:           iKey})
	} else {
		this.lines.Append(&table.TTableLine{
			Text:           st,
			Normal:         style1.EchoStyle(st, nCols),
			Selected:       style.Selected.EchoStyle(st, nCols),
			Active:         style1.EchoStyle(st, nCols),
			SelectedActive: style1.EchoStyle(st, nCols),
			IKey:           iKey})
	}
	// this.ILineCounter++
}

// func (this *v__Layer) Trace() {

// 	trace.Print("  %s", this.name)
// 	for _, table := range this.windows {
// 		table.pModel.Trace()
// 	}
// }

func (this *vWindow) IsInForeground() bool {
	if this.pLayer.IsActiveToSection() && Term.activeScreen == this.pModel.GetScreen() {
		// trace.BeginEnd_n(this, "IsInForeground true") // t //
		return true
	}
	// trace.BeginEnd_n(this, "IsInForeground false") // t //
	return false
}

func (this *vWindow) DumpWindow() {
	// trace.Begin_n(this, "DumpScreen {vWindow}") // t //
	// trace.Print("this.name                  %s", this.name) // t //
	// trace.Print("this.tRect                 %+v", this.tRect) // t //
	// trace.End() // t //
}

// A printable view object.
type vWindow struct {
	vBasicView
	// focus              string
	pLayer         *vLayer
	pModel         *mBasic
	lines          table.GTableLines
	emptyLine      string
	iPage          uint16
	iStatusLineRow uint16
	iLineCounter   uint16
	nRows          uint16 // number of effective rows for printing the lines of each page
	nCols          uint16
	nLines         uint16
	nPages         uint16
	rowOffset      uint16
	pageOffset     uint16
	colOffset      uint16
	tableName      string
	statusLine     string
	list           string
	isFolded       bool
	isCanvas       bool
	isPage         bool
	canvasEcho     string
	iOldPage       uint16
	columnsCaption string
	fGetWidth      func() uint16
}

type vColorWindow struct {
	vWindow
	wPalette vWindow
}
