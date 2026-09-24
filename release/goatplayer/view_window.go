package main

import (
	"fmt"
	"slices"
	"sort"
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
	p := this.lines.Get(iLine)
	nCols := this.nCols
	AdjustWidth(&st, nCols)
	p.normal = styleNormal.EchoStyle(st, nCols)
	p.selected = styleSelected.EchoStyle(st, nCols)
	p.active = styleActive.EchoStyle(st, nCols)
	p.selectedActive = styleSelActive.EchoStyle(st, nCols)
	p.iKey = iKey
}

func (this *vWindow) MakeActiveFilterLine(s string) string {
	nCols := this.nCols
	return styleActFilter.EchoStyle(s, nCols)
}

func (this *vWindow) MakeSelectedActiveFilterLine(s string) string {
	nCols := this.nCols
	return styleSelActFilter.EchoStyle(s, nCols)
}

func (this *vWindow) SetLine(iLine uint16, st string, iKey uint16, option uint16) {
	// trace.N_BeginEnd(this, "SetLine, %d, %s", iLine, st)
	p := this.lines.Get(iLine)
	AdjustWidth(&st, this.nCols)
	if option == 1 {
		p.normal = styleGrey.EchoStyle(st, this.nCols)
	} else {
		p.normal = styleNormal.EchoStyle(st, this.nCols)
	}
	p.text = st
	p.selected = styleSelected.EchoStyle(st, this.nCols)
	p.active = styleActive.EchoStyle(st, this.nCols)
	p.selectedActive = styleSelActive.EchoStyle(st, this.nCols)
	p.iKey = iKey
}

func GetAudiofGenre(audiof *aAudiofile) string       { return audiof.genre }
func GetAudiofArtist(audiof *aAudiofile) string      { return audiof.artist }
func GetAudiofAlbum(audiof *aAudiofile) string       { return audiof.album }
func GetAudiofAlbumArtist(audiof *aAudiofile) string { return audiof.albumArtist }
func GetAudiofComposer(audiof *aAudiofile) string    { return audiof.composer }
func GetAudiofYear(audiof *aAudiofile) string        { return audiof.year }

func (this *cFilters) Or(audiof *aAudiofile) bool {
	for i := range this.sl {
		this.b[i] = this.sl[i](audiof)
	}
	return slices.Contains(this.b, true)
}

func (this *cFilters) And(audiof *aAudiofile) bool {
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

/* Returns nil if there are no filterValues. */
func (this *mContainer) GetFieldValuesForFiltered(fGetFieldValue func(audiof *aAudiofile) string) *[]string {
	if len(this.validValues) == 0 {
		return nil
	}
	var mapA map[string]uint16
	mapA = make(map[string]uint16)
	var values []string
	// for i := range AF.sl {
	for i := range this.pBrowserWsp.wspDbIndexes.sl {
		audiof := &AF.sl[i]
		// trace.Print("audiof is %d %s", i, audiof.title)
		if this.ApplyFilterOR(audiof) {
			mapA[fGetFieldValue(audiof)]++
		}
	}
	for value := range mapA {
		values = append(values, value)
	}
	sort.Strings(values)
	return &values
}

func FilterStrings(s1, s2 *[]string) {
	var newS []string
	for _, s := range *s1 {
		if slices.Contains(*s2, s) {
// trace.Print("contains %s", s) // t //
			newS = append(newS, s)
		}
	}
	*s1 = nil
	*s1 = make([]string, len(newS))
	copy(*s1, newS)
}

func (this *mContainer) ReactivateFilters() {
// trace.Begin_n(this, "ReactivateFilters") // t //
	for i := range this.keys.sl {
		k := &this.keys.sl[i]
		if slices.Contains(this.validValues, k.v) {
			k.isActiveFilter = true
// trace.Print("%s reactivated", k.v) // t //
		}
	}
	// this.filterValues = nil
// trace.End() // t //
}

/*
   //////	Add all filters
   //	For each container in this browser, search if there are active filters
   var allFilters1 []*mContainer
   // var allFilters []*mContainer
   var containers1 []*mContainer
   var containers2 []*mContainer
   containers1 = append(containers1, &browser.Genres)

   	for _, container := range containers1 {
   		container.filterValues = nil
   		for i := range container.keys.sl {
   			k := container.keys.sl[i]
   			if k.isActiveFilter {
   				container.AppendFilterValue(k.v)
   				nFilters++
   			}
   		}
   		if len(container.filterValues) > 0 {
   			allFilters1 = append(allFilters1, container)
   			allFilters = append(allFilters, container)
   		}
   	}

// trace.Print("allFilters1 %s", LongNames(allFilters1)) // t //
// trace.Print("allFilters %s", LongNames(allFilters)) // t //

   	for _, container := range containers2 {
   		container.filterValues = nil
   		for i := range container.keys.sl {
   			k := container.keys.sl[i]
   			if k.isActiveFilter {
   				container.AppendFilterValue(k.v)
   				nFilters++
   			}
   		}
   		if len(container.filterValues) > 0 {
   			allFilters = append(allFilters, container)
   		}
   	}

// trace.Print("allFilters %s", LongNames(allFilters)) // t //

   //////	Get valid Artists and Albums

   var ArtistsValues []string
   var AlbumsValues []string

   	if len(allFilters1) > 0 {
   		ArtistsValues = *browser.Genres.GetFieldValuesForFiltered(GetAudiofArtist)
   		AlbumsValues = *browser.Genres.GetFieldValuesForFiltered(GetAudiofAlbum)
   	} else {

   		ArtistsValues = nil
   		ArtistsValues = append(ArtistsValues, browser.Artists.allValues...)
   		AlbumsValues = nil
   		AlbumsValues = append(AlbumsValues, browser.Albums.allValues...)
   	}

// trace.Print("valid artists: %#v (%d)", ArtistsValues, len(ArtistsValues)) // t //
// trace.Print("valid albums: %#v (%d)", AlbumsValues, len(AlbumsValues)) // t //
   FilterStrings(&browser.Artists.filterValues, &ArtistsValues)
   FilterStrings(&browser.Albums.filterValues, &AlbumsValues)
// trace.Print("%s Artists.filterValues, %#v", browser.name, browser.Artists.filterValues) // t //
// trace.Print("%s Albums.filterValues, %#v", browser.name, browser.Albums.filterValues) // t //

   browser.Artists.LoadFromSliceAndPrint(&ArtistsValues)
   browser.Albums.LoadFromSliceAndPrint(&AlbumsValues)
   //////	Update Artists and Albums filters

   // var newArtistsFilterValues []string
   // for _, fil := range browser.Artists.filterValues {
   // }

   //////	Apply filters to indexes; load and prepare the View

   	if len(allFilters) == 0 {
   		//	Load all audiofiles
   		browser.View.LoadAndPrepare(func(audiof *aAudiofile) bool { return true }, nil)
   		if monitoring {
   			Monitor.AddLine("no filters")
   			Monitor.UpdateFileWithContent()
   		}
   	} else {

   		if monitoring {
   			for _, container := range allFilters {
   				Monitor.AddLine("-- " + container.name)
   				for _, v := range container.filterValues {
   					Monitor.AddLine(v)
   				}
   				container.keys.TraceActiveFilters()
   			}
   			Monitor.UpdateFileWithContent()
   		}
   		//	Apply all filters
   		for i := range AF.sl {
   			audiof := &AF.sl[i]
   			// trace.Print("audiof is %d %s", i, audiof.title)
   			isValid := false
   			for _, container := range allFilters {
   				if container.ApplyFilterOr(audiof) {
   					isValid = true
   				} else {
   					isValid = false
   					break
   				}
   			}
   			if isValid {
// trace.Print("valid audiof: %s", audiof.title) // t //
   				indexes = append(indexes, uint16(i))
   			}
   		}
   		browser.View.LoadAndPrepare(func(audiof *aAudiofile) bool { return true }, &indexes)
   	}

   browser.View.w.PrintWindowAnyway()
   //////	Update containers
   // switch this {
   // case &browser.Genres:
   // 	if len(this.filterValues) == 0 {
   // 		indexes = nil
   // 		for i := range AF.sl {
   // 			indexes = append(indexes, uint16(i))
   // 		}
   // 	}
   // 	//	Artists
   // 	var mapArtists map[string]uint16
   // 	mapArtists = make(map[string]uint16)
   // 	for _, i := range indexes {
   // 		audiof := &AF.sl[i]
   // 		mapArtists[audiof.artist]++
   // 	}
   // 	// trace.Print("map artists: %v", mapArtists)
   // 	browser.Artists.Reset()
   // 	var artistList []string
   // 	for artistName := range mapArtists {
   // 		artistList = append(artistList, artistName)
   // 	}
   // 	sort.Strings(artistList)
   // 	for _, name := range artistList {
   // 		browser.Artists.AppendKey(name)
   // 	}
   // 	browser.Artists.PrepareView()
   // 	browser.Artists.w.PrintWindowAnyway()
   // 	//	Albums
   // 	var mapAlbums map[string]uint16
   // 	mapAlbums = make(map[string]uint16)
   // 	for _, i := range indexes {
   // 		audiof := &AF.sl[i]
   // 		mapAlbums[audiof.album]++
   // 	}
   // 	// trace.Print("map albums: %v", mapAlbums)
   // 	browser.Albums.Reset()
   // 	var albumList []string
   // 	for albumName := range mapAlbums {
   // 		albumList = append(albumList, albumName)
   // 	}
   // 	sort.Strings(albumList)
   // 	for _, name := range albumList {
   // 		browser.Albums.AppendKey(name)
   // 	}
   // 	browser.Albums.PrepareView()
   // 	browser.Albums.w.PrintWindowAnyway()
   // case &browser.Artists:
   // 	if len(this.filterValues) == 0 {
   // 		break
   // 	}
   // 	//	Albums
   // 	var mapAlbums map[string]uint16
   // 	mapAlbums = make(map[string]uint16)
   // 	for _, i := range indexes {
   // 		audiof := &AF.sl[i]
   // 		mapAlbums[audiof.album]++
   // 	}
   // 	// trace.Print("map albums: %v", mapAlbums)
   // 	browser.Albums.Reset()
   // 	var albumList []string
   // 	for albumName := range mapAlbums {
   // 		albumList = append(albumList, albumName)
   // 	}
   // 	sort.Strings(albumList)
   // 	for _, name := range albumList {
   // 		browser.Albums.AppendKey(name)
   // 	}
   // 	browser.Albums.PrepareView()
   // 	browser.Albums.w.PrintWindowAnyway()
   // }
// trace.N_End(this, "ActivateFilterByTRow, iTRow: %d", iTRow) // t //
*/

// Given a terminal row index, selects a table row and reprints.
func (this *vWindow) SelectLineByTRow(iTRow uint16) {
	iLine := this.CalcLineIndexByTRow(iTRow)
// trace.BeginAdd_n(this.pModel, "SelectLineByTRow", "iTRow: %d, iLine: %d", iTRow, iLine) // t //
	p := this.lines.Get(iLine)
	if p == nil || p.iKey == UNSET {
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
	if iKey == UNSET {
// trace.Error_n(this.pModel, "GetLineByKey, UNSET index, %d", iKey) // t //
		return UNSET
	}
	if p := this.pModel.keys.Get(iKey); p != nil {
		if this.isFolded {
			return p.iLine
		} else {
			return iKey
		}
	} else {
		return UNSET
	}

	// trace.N_BeginEnd(this, "GetLineByKey")
	// trace.Print("iKey: %d", iKey) //@1
	// trace.Print("this.keys[iKey].iLine: %d", this.keys[iKey].iLine) //@1
	// return this.keys[iKey].iLine
}

// Reprints the new and old active lines.
func (this *vWindow) UpdateActiveLines() {
// trace.BeginSilent_n(this, "UpdateActiveLines") // t //
	model := this.pModel
// trace.AddToPrint("active %d (%d)", model.keys.iActive, model.keys.iOldActive) // t //
	if MainScreen.LeftSection.GetActiveLayer() != Play.w.pLayer {
// trace.Return() // t //
		return
	}
	// model.keys.TraceActive("")
	this.PrintLine(this.GetLineByKey(model.keys.iOldActive))
	this.PrintLine(this.GetLineByKey(model.keys.iActive))
	Term.Flush()
	model.keys.iOldActive = model.keys.iActive
// trace.End() // t //
}

// Calculates terminal row index by line index. Returns OUT_OF_PAGE if the line is not in the page.
func (this *vWindow) CalcTRowIndex(iLine uint16) uint16 {
	var iPage = iLine / this.nRows
	if iPage != this.iPage {
		return OUT_OF_PAGE
	}
	return iLine%this.nRows + this.rowOffset + this.pLayer.pSection.marginTop
}

// .
func (this *vWindow) CalcLineIndexByTRow(iTRow uint16) uint16 {
	return this.iPage*this.nRows + iTRow - this.rowOffset - this.marginTop
}

// .
func (this *vWindow) AddLineToPrint(iTRow uint16, line *string) {
	Term.AddStringWithPos(iTRow, this.marginLeft+this.colOffset, *line)
}

// .
func (this *vWindow) PrintLine(iLine uint16) {
	// trace2.Print("PrintLine: %d", iLine)
	if iLine == UNSET {
		return
	}
	model := this.pModel
	var isActiveFilter, isActive bool
	if !this.isFolded {
		if model.keys.sl[iLine].isActiveFilter {
			isActiveFilter = true
		}
	}
	if this.isFolded || this.pModel.GetScreen().isDialog {
		if model.keys.iActive != UNSET {
			iActiveLine := this.GetLineByKey(model.keys.iActive)
			if iActiveLine == iLine {
				isActive = true
				// trace.Print("%d: isActive", iLine)
			}
		}
	}
	// trace.AddToPrint("active line: %d", iActiveLine)
	var iTRow = this.CalcTRowIndex(iLine)
	line_p := this.lines.Get(iLine)
	if iLine == this.lines.iSelected {
		if isActive {
			this.AddLineToPrint(iTRow, &line_p.selectedActive)
		} else if isActiveFilter {
			s := styleSelActFilter.EchoStyle(line_p.text, this.nCols)
			this.AddLineToPrint(iTRow, &s)
		} else {
			this.AddLineToPrint(iTRow, &line_p.selected)
		}
	} else {
		if isActive {
			this.AddLineToPrint(iTRow, &line_p.active)
		} else if isActiveFilter {
			s := styleActFilter.EchoStyle(line_p.text, this.nCols)
			this.AddLineToPrint(iTRow, &s)
		} else {
			this.AddLineToPrint(iTRow, &line_p.normal)
		}
	}
	// trace.N_BeginEndPrintAll(this, "PrintLine %d", iLine)
}

/* Remember to set the selected line */
func (this *vWindow) PrintWindowAnyway() {
	this.iOldPage = UNSET
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
		Term.AddStringWithPos(iTRow, this.marginLeft+this.colOffset, DC_RESET+this.emptyLine)
	}
// trace.End() // t //
}

func (this *vWindow) PrintPage() {
// trace.BeginEnd_n(this, "PrintPage") // t //
	var start uint16 = this.iPage * this.nRows
	var end uint16 = start + this.nRows - 1
	var iLine uint16
	Term.AddStringWithPos(this.CalcTRowIndex(start)-1, this.marginLeft+this.colOffset, DC_RESET+this.emptyLine)
	for iLine = start; iLine <= end; iLine++ {
		iTRow := this.CalcTRowIndex(iLine)
		if iLine > this.lines.Len()-1 {
			Term.AddStringWithPos(iTRow, this.marginLeft+this.colOffset, DC_RESET+this.emptyLine)
			continue
		}
		line := &this.lines.sl[iLine]
		this.AddLineToPrint(iTRow, &line.normal)
	}

	//	PrintStatusLine
	s := fmt.Sprintf(" Filepage %d/%d ", this.iPage+1, this.nPages)
	Term.AddStringWithPos(this.iStatusLineRow, this.marginLeft+this.colOffset, stylePageStatusLine.EchoStyle(s, this.nCols))
	Term.Flush()
}

func (this *vWindow) PrintCanvas() {
// trace.Begin_n(this, "PrintCanvas") // t //
	this.AddEmptyWindowToPrint()
	Term.AddSimple(this.canvasEcho)
	Term.Flush()
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
	if len(this.lines.sl) == 0 {
// trace.BeginEnd_n(this, "PrintWindow, no lines to print") // t //
		this.AddEmptyWindowToPrint()
		Term.Flush()
		return
	}
// trace.BeginSilent_n(this, "PrintWindow") // t //
	m := this.pModel
	// se := this.pLayer.pSection
	// this.layer.activeWindow = this
	// this.pLayer.pSection.PrintFocusButtons(true)
	/// Calculate page number
	var iPage uint16 = this.lines.iSelected / this.nRows
	this.nLines = this.lines.Len()
	///
// trace.AddToPrint("Page: %d OldPage: %d", iPage, this.iOldPage) // t //
// trace.AddToPrint("SelectedLine: %d, OldSelectedLine: %d", this.lines.iSelected, this.lines.iOldSelected) // t //
	// if this.oldSelected_i == this.selected_i {
	// 	this.Reprint()
	// 	return
	// }
	if iPage != this.iOldPage {
		//	Change page, reprint all Lines (including the empty ones)
		// trace.Print("-> change page")
		/// Window caption for multi window
		if !this.IsSingleWindow() {
			cap := SetCaption(m.shortName, m.focusKey)

			//!!
			// if this == this.pLayer.pSection.activeWindow {
			if this.IsActiveToSection() {
				cap = styleSelButton_SelFrame.EchoStyle(" "+cap, this.nCols+1)
				// trace.AddToPrint("Window caption for multi window, cap %s", cap)
				// Term.AddStringWithPos(this.rowOffset-3, this.marginLeft+this.colOffset, cap)
				Term.AddStringWithPos(this.rowOffset-2, this.marginLeft+this.colOffset, cap)
				Term.AddStringWithPos(this.rowOffset-1, this.marginLeft+this.colOffset, this.emptyLine)
			} else {
				cap = styleButton_SelFrame.EchoStyle(" "+cap, this.nCols+1)
				// trace.AddToPrint("Window caption for multi window, cap %s", cap)
				// Term.AddStringWithPos(this.rowOffset-3, this.marginLeft+this.colOffset, cap)
				Term.AddStringWithPos(this.rowOffset-2, this.marginLeft+this.colOffset, cap)
				Term.AddStringWithPos(this.rowOffset-1, this.marginLeft+this.colOffset, this.emptyLine)
			}
		}
		///
		this.iPage = iPage
		var start uint16 = iPage * this.nRows
		var end uint16 = start + this.nRows - 1
		var iLine uint16
		if this.IsSingleWindow() {
			if this.columnsCaption == "" {
				Term.AddStringWithPos(this.CalcTRowIndex(start)-1, this.marginLeft+this.colOffset, DC_RESET+this.emptyLine)
			} else {
				Term.AddStringWithPos(this.CalcTRowIndex(start)-1, this.marginLeft+this.colOffset, this.columnsCaption)
			}
		}
// trace.AddToPrint("Keys %d, start %d, end %d, iPage %d, nLines: %d", m.keys.Len(), start, end, iPage, this.nLines) // t //
		for iLine = start; iLine <= end; iLine++ {
			var iTRow = this.CalcTRowIndex(iLine)
			// var iRow = i % this.rows_n + this.rowOffset
			if iLine > this.nLines-1 {
				//	Print an empty line
				// log.Println("this.PrintLine > empty line")
				Term.AddStringWithPos(iTRow, this.marginLeft+this.colOffset, DC_RESET+this.emptyLine)
				continue
			}
			this.PrintLine(iLine)
		}
	} else {
// trace.AddToPrint("same page") // t //
		//	Reprint only old and new selected line
		this.PrintLine(this.lines.iOldSelected)
		this.PrintLine(this.lines.iSelected)
	}
	this.iOldPage = iPage
	this.lines.iOldSelected = this.lines.iSelected
	/// Key value long line screen bottom
	// if !this.IsSingleWindow() {
	// 	// Term.AddStringWithPos(this.marginBottom, se.marginLeft+2, styleSelTab.EchoStyle(m.keys.GetSelected().v, se.marginRight-se.marginLeft-5))
	// 	Term.AddStringWithPos(nRows, 0, styleNormal.EchoStyle("Selected: "+m.keys.GetSelected().v, nCols-22))
	// }
	///
	this.PrintStatusLine()
	// this.tableEcho = Window.contents
	Term.SetNormal()
	Term.Flush()
// trace.End() // t //
}

/* this.statusLine must be set before to call this method. */
func (this *vWindow) PrintStatusLine() {
	m := this.pModel
	m.outer.SetStatusLine_o()
	s := fmt.Sprintf("  %d/%d %s", m.keys.iSelected+1, m.keys.Len(), this.statusLine)
	Term.AddStringWithPos(this.iStatusLineRow, this.marginLeft+this.colOffset, styleWindowStatusLine.EchoStyle(s, this.nCols))
}

func (this *vWindow) SetEmptyLine() {
	var width uint16
	width = this.nCols
	this.emptyLine = ""
	this.emptyLine = styleEmptyLine.EchoStyle(this.emptyLine, width)
}

func (this *vWindow) SetEmptyLineForDialog() {
	this.emptyLine = ""
	this.emptyLine = styleEmptyLine.EchoStyle(this.emptyLine, this.nCols)
}

func (this *vWindow) SelectLine(iLine uint16) {
	m := this.pModel
// trace.BeginAdd_n(m, "SelectLine", "iLine: %d", iLine) // t //
	this.lines.Select(iLine)
	if line := this.lines.Get(iLine); line != nil {
		if this.isFolded {
			m.keys.Select(line.iKey)
		} else {
			m.keys.Select(iLine)
		}
		m.outer.SelectKey_o()
	}
// trace.End() // t //
}

/*
Appends a new line to the view. st: the text of the line. setKeyIndex: true if the key index must be set for this line.
option: nil; 1= arrow selection
*/
func (this *vWindow) AppendLine(st string, setKeyIndex bool, style *cStyle) {
	// trace.N_BeginEnd(this, "AppendLine, %s %+v", st, style)
	var iKey uint16
	if setKeyIndex {
		iKey = this.pModel.iKeyCounter
	} else {
		iKey = this.lines.Len() // it was UNSET
	}
	// AdjustWidth(&st, this.nCols)
	nCols := this.nCols
	if style == nil {
		act := " ▶ " + st[3:]
		this.lines.Append(&tTableLine{
			st,
			styleNormal.EchoStyle(st, nCols),
			styleSelected.EchoStyle(st, nCols),
			styleActive.EchoStyle(act, nCols),
			styleSelActive.EchoStyle(act, nCols),
			iKey})
	} else {
		this.lines.Append(&tTableLine{
			st,
			style.EchoStyle(st, nCols),
			styleSelected.EchoStyle(st, nCols),
			style.EchoStyle(st, nCols),
			style.EchoStyle(st, nCols),
			iKey})
	}
	// this.iLineCounter++
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
	lines          cTableLines
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

