package main

import (
	"fmt"
	"goatplayer/internal/K"
	AF "goatplayer/internal/audiofiles"
	"goatplayer/internal/echo"
	"goatplayer/internal/input"
	"goatplayer/internal/small/utils"
	"goatplayer/internal/style"
	"goatplayer/internal/table"
	"os/exec"
	"slices"
	"sort"
	"strings"
)

/* Prints this table's commands and the focus legend. */
func (this *mBasic) PrintCommands() {
	// trace.Begin_n(this, "PrintCommands") // t //
	// Focus.AddToPrint()
	echo.AddSimple(this.cmmEchoes[this.iCmmPage])
	echo.Flush()
	// trace.Print("print page %d/%d", this.iCmmPage, this.nCmmPages) // t //
	// trace.End() // t //
}

func (this *mBasic) SetFocusToModel() {
	// trace.Begin_n(this, "SetFocusToModel") // t //
	if Term.focusedModel == this {
		// trace.ReturnAdd("already focused") // t //
		return
	}
	if this.fSetFocusToModel != nil {
		// trace.Print("calling fSetFocusToModel") // t //
		this.fSetFocusToModel()
	}
	if this._fSetFocusToModel != nil {
		// trace.Print("calling _fSetFocusToModel") // t //
		this._fSetFocusToModel()
	}
	this.w.iOldPage = K.UNSET //!!
	this.outer.SelectKey_o()
	this.PrintCommands()
	CommandPanel.butPageMinus.PrintUnpressed()
	CommandPanel.butPagePlus.PrintUnpressed()
	this.w.BringWindowToFront(true)
	this.w.pLayer.pSection.BringSectionToFront()
	if this == &View1.mBasic {
	}
	Term.focusedModel = this
	if this != &Help.mBasic {
		Help.SetFilename(this.helpfile)
	}
	// trace.End() // t //
}

func (this *mBasic) SetOuter() {
	this.outer = this
}

/* Calls InitWindow */
func (this *mBasic) InitModel(name string, sFocus string, helpfile string) {
	if this.outer == nil {
		this.SetOuter()
	}
	this.name = name
	this.helpfile = helpfile
	this.shortName = name
	this.keys.SelectedI = K.NEVERSET
	if sFocus == K.FOCUS_UNFOCUSABLE {
		this.isFocusable = false
	} else {
		this.isFocusable = true
		this.focusKey = sFocus
		// Focus.items = append(Focus.items, uFocusItem{focus_s, name})
		// Focus.mapFocus[focus_s] = this.name
		// Focus.mapFocusModel[focus_s] = this
	}
	this.w.InitWindow(this)
}

/* Prepares the table, sets the margins, adds this to its Platoon and to Focus. */
func (this *mBasic) PrepareModel() {
	// trace.BeginEnd_n(this, "PrepareModel") // t //
	if this.pWorkspace != nil {
		// this.shortName = this.name
		this.name = this.pWorkspace.name + "." + this.name
		this.w.name = this.name + ".w"
	}
	AllMaps.allBasicModels = append(AllMaps.allBasicModels, this)
	AllMaps.mapBasicModels[this.name] = this
	// trace.N_BeginEnd(this, "Prepare")
	// platoon.AddWindow(this)
	// this.IOldSelectedKey = 1
	// this.keys.IActive = K.UNSET
	// this.IOldActive = K.UNSET
}

func (this *mBasic) GetScreen() *vScreen {
	if this.w.pLayer == nil {
		// trace.Error("%s, pLayer is nil", this.name) // t //
	} else if this.w.pLayer.pSection == nil {
		// trace.Error("%s, pSection is nil", this.name) // t //
	} else if this.w.pLayer.pSection.pScreen == nil {
		// trace.Error("%s, pScreen is nil", this.name) // t //
	}
	return this.w.pLayer.pSection.pScreen
}

func (this *mBasic) OpenHelpDialog() {
	if dialogMode {
		// return
	}
	// trace.Print("CMM_BASIC_SETTINGS_DIALOG") // t //
	MainScreen.BeforeDialog()
	SettingsDialog.focusKeys.activeKey = K.FOCUS_HELP
	SettingsDialog.OpenDialog(Help.focusKey)
	// trace.Print("Exiting from dialog-mode, section Settings") // t //
	//	Exited from dialog-mode
	MainScreen.AfterDialog()
	this.SetFocusToModel()
}

func (this *mBasic) OpenReportDialog() {
	if dialogMode {
		// return
	}
	// trace.Print("CMM_BASIC_SETTINGS_DIALOG") // t //
	MainScreen.BeforeDialog()
	SettingsDialog.focusKeys.activeKey = K.FOCUS_REPORT
	SettingsDialog.OpenDialog(Report.focusKey)
	// trace.Print("Exiting from dialog-mode, section Settings") // t //
	//	Exited from dialog-mode
	MainScreen.AfterDialog()
	this.SetFocusToModel()
}

func (this *mBasic) OpenSettingsDialog() {
	if dialogMode {
		// return
	}
	// trace.Print("CMM_BASIC_SETTINGS_DIALOG") // t //
	MainScreen.BeforeDialog()
	SettingsDialog.focusKeys.activeKey = K.FOCUS_SETTINGS
	SettingsDialog.OpenDialog(SetM.focusKey)
	// trace.Print("Exiting from dialog-mode, section Settings") // t //
	//	Exited from dialog-mode
	MainScreen.AfterDialog()
	this.SetFocusToModel()
}

func (this *mBasic) OpenEqualizerDialog() {
	if dialogMode {
		// return
	}
	// trace.Print("CMM_BASIC_EQUALIZER_DIALOG") // t //
	MainScreen.BeforeDialog()
	SettingsDialog.focusKeys.activeKey = K.FOCUS_EQUALIZER
	SettingsDialog.OpenDialog(Equalizer.focusKey)
	// trace.Print("Exiting from dialog-mode, Settings/Equalizer") // t //
	//	Exited from dialog-mode
	MainScreen.AfterDialog()
	this.SetFocusToModel()
}

func (this *mBasic) OpenColorsDialog() {
	if dialogMode {
		// return
	}
	// trace.Print("OpenColorsDialog") // t //
	MainScreen.BeforeDialog()
	SettingsDialog.focusKeys.activeKey = K.FOCUS_COLORS
	SettingsDialog.OpenDialog(Colors.focusKey)
	// trace.Print("Exiting from dialog-mode, Settings/Colors") // t //
	//	Exited from dialog-mode
	MainScreen.AfterDialog()
	this.SetFocusToModel()
}

func (this *mBasic) OpenFindDialog() bool {
	var ss []string
	if this == &BrowserWsp1.View.mBasic || this == &BrowserWsp2.View.mBasic {
		for i := range AF.Sl {
			audiof := &AF.Sl[i]
			ss = append(ss, audiof.Title)
		}
		slices.Sort(ss)
	} else {
		for i := range this.keys.Sl {
			k := &this.keys.Sl[i]
			ss = append(ss, k.V)
		}
	}
	Find.LoadList1(&ss, this)
	// FindScreen.activeSection = &FindScreen.Section
	FindDialog.activeWindow = &Find.w
	Find.keys.SelectedI = 0
	Find.w.lines.SelectedI = 0
	MainScreen.BeforeDialog()
	rv := FindDialog.OpenDialog(K.FOCUS_FIND)
	MainScreen.AfterDialog()
	this.SetFocusToModel()
	return rv
}

type tStringsAndIndexes struct {
	str      string
	str2     string
	lowerStr string
	dbIndex  uint16
}

func SortStringsAndIndexes(list *[]tStringsAndIndexes) {
	sort.Slice(*list, func(i int, j int) bool {
		return (*list)[i].str < (*list)[j].str
	})
}

func (this *mBasic) OpenFindDialog2(ssl *[]tStringsAndIndexes) bool {
	Find.LoadList2(ssl, this)
	// FindScreen.activeSection = &FindScreen.Section
	FindDialog.activeWindow = &Find.w
	Find.keys.SelectedI = 0
	Find.w.lines.SelectedI = 0
	MainScreen.BeforeDialog()
	rv := FindDialog.OpenDialog(K.FOCUS_FIND)
	MainScreen.AfterDialog()
	this.SetFocusToModel()
	return rv
}

func (this *vMainScreen) BeforeDialog() {
	this.LeftSection.BringSectionToBack()
	this.pRightSection.BringSectionToBack()
}

func (this *vMainScreen) AfterDialog() {
	this.activeSection.BringSectionToFront()
	this.LeftSection.GetActiveWindow().BringWindowToFront(false)
	this.pRightSection.GetActiveWindow().BringWindowToFront(false)
}

// func (this *mBasic) OpenBiosDialog() {
// 	if dialogMode {
// 		return
// 	}
// 	MainScreen.BeforeDialog()
// 	HelpDialog.OpenDialog(Bios.focusKey)
// 	MainScreen.AfterDialog()
// 	trace.Print("Exiting from dialog-mode, Settings")
// 	if Term.activeScreen == &SettingsDialog.vScreen {
// 		SettingsDialog.Section.BringSectionToFront()
// 		SettingsDialog.Section.GetActiveWindow().BringWindowToFront(true)
// 	}
// 	// filename := Bios.folderName + "/" + Artists.keys.GetSelected().V
// 	// trace.Print(filename)
// 	// if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
// 	// 	os.Create(filename)
// 	// 	go func() {
// 	// 		path := filename
// 	// 		trace.Print("open %s", path)
// 	// 		by, err := exec.Command("xdg-open", path).Output()
// 	// 		trace.Print("%v %v", by, err)
// 	// 	}()
// 	// }
// 	// Screen.ShowHelp(&Bios)
// 	// Bios.ShowFile()
// }

// .
func (this *mBasic) GetName() string {
	return this.name
}

// .
func (this *mBasic) GetLongName() string {
	return this.name
}

// Actions to be done after loading data to the table.
func (this *mBasic) AfterLoad() {
	// trace.N_BeginEnd(this, "AfterLoad")
	this.w.nPages = ((this.keys.Len() - 1) / this.w.nRows) + 1
	if this.keys.Len() > 0 {
		this.keys.SelectedI = 0
		this.w.lines.SelectedI = this.keys.GetSelected().LineI
	} else {
		// trace.Error("%s.AfterLoad: keys len is 0", this.name) // t //
	}
	// this.nKeys = uint16(len(this.keys))
	// trace.Print("%d %d %d %d", this.filenames_n, this.pages_n, this.rows_n, this.pages_n*this.rows_n)
	// trace.N_End(this, "AfterLoad")
}

// .
func (this *mBasic) GetKeyByLine(iLine uint16) uint16 {
	if iLine == K.UNSET {
		// trace.Error_n(this, "GetKeyByLine, K.UNSET index, %d", iLine) // t //
		return K.UNSET
	}
	// trace.N_BeginEnd(this, "GetKeyByLine")
	// trace.Print("iLine: %d", iLine) //@1
	// trace.Print("this.lines[iLine].IKey: %d", this.lines[iLine].IKey) //@1
	if p := this.w.lines.Get(iLine); p != nil {
		return p.IKey
	} else {
		return K.UNSET
	}
}

/* Changes all values about the selected index. Keys must be prepared before. */
func (this *mBasic) SelectKeyAndLine(iKey uint16) {
	// trace.BeginAdd_n(this, "SelectKey", "%d", iKey) // t //
	if this.isEmpty {
		// trace.ReturnAdd("is empty") // t //
		return
	}
	this.keys.Select(iKey)
	var iSelectedLine uint16
	if this.w.isFolded {
		if key := this.keys.Get(iKey); key != nil {
			iSelectedLine = key.LineI
		}
	} else {
		iSelectedLine = iKey
	}
	this.w.lines.Select(iSelectedLine)
	this.outer.SelectKey_o()
	this.keys.TraceSelected("")
	this.w.lines.TraceSelected("")
	// trace.End() // t //
}

func (this *mBasic) NextTab() {
}

// Selects the row according to a filename (key).
func (this *mBasic) SelectIndexByKey(thatKey string) {
	// trace.BeginEndAdd_n(this, "SelectLineByKey", "key: %s", thatKey) // t //
	var i uint16
	for i = 0; i < this.keys.Len(); i++ {
		key_p := this.keys.Get(i)
		if thatKey == key_p.V {
			break
		}
	}
	// trace.Print("key found at index: %d", i) // t //
	this.SelectKeyAndLine(i)
	// this.Print()
}

// .
func (this *mBasic) SelectPrev() {
	// var new_selected uint16
	var i uint16
	if this.keys.SelectedI > 0 {
		// new_selected = this.Selected_i-1
		i = this.keys.SelectedI - 1
	} else {
		i = this.keys.Len() - 1
		// new_selected = this.filenames_n - 1
	}
	// this.SelectRow(new_selected)
	this.SelectKeyAndLine(i)
	// this.ChangeSelectedIndex(i)
	this.w.PrintWindow()
}

// .
func (this *mBasic) SelectNext() {
	var i uint16
	// trace.BeginAdd_n(this, "SelectNext", "iSelected: %d", this.keys.SelectedI) // t //
	this.keys.TraceSelected("")
	if this.keys.SelectedI < this.keys.Len()-1 {
		i = this.keys.SelectedI + 1
	} else {
		i = 0
	}
	this.SelectKeyAndLine(i)
	// this.ChangeSelectedIndex(i)
	this.w.PrintWindow()
	// trace.End() // t //
}

func (this *mBasic) SelectFirst() {
	this.SelectKeyAndLine(0)
	// this.ChangeSelectedIndex(0)
	this.w.PrintWindow()
}

func (this *mBasic) SelectActive_o() {
	input.DisplayInput()
	if this.keys.ActiveI != K.UNSET {
		this.SelectKeyAndLine(this.keys.ActiveI)
		// this.ChangeSelectedIndex(this.active_i)
		this.w.PrintWindow()
	}
}

// .
func (this *mBasic) ToggleActiveFilter(iKey uint16) bool {
	k := this.keys.Get(iKey)
	k.IsActiveFilter = !k.IsActiveFilter
	return k.IsActiveFilter
	// this.keys.TraceActiveFilters()
}

func (this *mBasic) GotoPrevPage() {
	var i uint16
	if this.keys.SelectedI < this.w.nRows {
		i = 0
	} else {
		i = this.keys.SelectedI - this.w.nRows
	}
	// trace.Print("this.Selected_i %d", this.keys.SelectedI) // t //
	this.SelectKeyAndLine(uint16(i))
	// this.ChangeSelectedIndex(i)
	this.w.PrintWindow()
}

func (this *mBasic) GotoNextPage() {
	var i uint16
	var last = this.keys.Len() - 1
	i = min(this.keys.SelectedI+this.w.nRows, last)
	this.SelectKeyAndLine(uint16(i))
	// this.ChangeSelectedIndex(i)
	this.w.PrintWindow()
}

func (this *mBasic) OpenInBrowser() {
	// trace.BeginEnd_n(this, "OpenInBrowser") // t //
	if this._fGetBrowserString == nil {
		return
	}
	var s, url string
	fReplace := func() {
		s = strings.ReplaceAll(s, " ", "+")
	}
	fOpen := func(url string) {
		// trace.Print("browser search: %s", url) // t //
		go func() {
			// by, err := exec.Command("xdg-open", url).Output() /*c*/
			exec.Command("xdg-open", url).Output()
			// trace.Print("%v %v", by, err) // t //
		}()
	}
	s = this._fGetBrowserString()
	if Set.OpenInBrowserApple.Value {
		url = fmt.Sprintf("https://music.apple.com/us/search?term=%s", s)
		fOpen(url)
	}
	if Set.OpenInBrowserGoogle.Value {
		fReplace()
		url = fmt.Sprintf("https://www.google.com/search?q=%s", s)
		fOpen(url)
	}
	if Set.OpenInBrowserYoutube.Value {
		fReplace()
		url = fmt.Sprintf("https://www.youtube.com/results?search_query=%s", s)
		fOpen(url)
	}
}

/* If style is nil, do normal lines */
func (this *mBasic) AppendKeyAndLine(st string, style *style.GStyle) {
	// trace2.N_BeginEnd(this, "AppendKeyAndLine")
	this.AppendKey(st)
	st = "  " + st
	utils.AdjustWidth(&st, this.w.fGetWidth())
	this.w.AppendLine(st, false, style)
}

func (this *mBasic) AppendKeyAndLineAndDbIndex(item *tStringsAndIndexes) {
	var normal, selected string
	// trace2.N_BeginEnd(this, "AppendKeyAndLine")
	this.keys.Append(&table.TTableKey{V: item.str + item.str2, LineI: this.w.iLineCounter, DbIndex: item.dbIndex, IsActiveFilter: false})
	str := "  " + item.str
	len1 := utils.RuneLen(&str)
	len2 := this.w.nCols - len1
	if len1 >= this.w.nCols {
		normal = style.Normal.EchoStyle(str, this.w.nCols)
		selected = style.Selected.EchoStyle(str, this.w.nCols)
	} else {
		normal = style.Normal.EchoStyle(str, 0) + style.NormalLight.EchoStyle(item.str2, len2)
		selected = style.Selected.EchoStyle(str, 0) + style.SelectedLight.EchoStyle(item.str2, len2)
	}
	this.w.lines.Append(&table.TTableLine{Text: str, Normal: normal, Selected: selected, Active: "", SelectedActive: "", IKey: 0})
}

func (this *vBasicView) GetName() string     { return this.name }
func (this *vBasicView) GetLongName() string { return this.name }

// Appends a key.
func (this *mBasic) AppendKey(st string) {
	// trace.N_BeginEnd(this, "AppendKey, %s", st)
	this.keys.Append(&table.TTableKey{V: st, LineI: this.w.iLineCounter, DbIndex: 0, IsActiveFilter: false})
	// this.IKeyCounter++
}

func (this *mBasic) Reset() {
	this.keys.Reset()
	this.w.lines.Reset()
	this.tableHasChanged = true
	this.iKeyCounter = 0
	this.w.iLineCounter = 0
	this.w.iOldPage = K.UNSET
}

type mBasic struct {
	name       string
	shortName  string
	tabCaption string
	focusKey   string

	cmmEchoes []string
	nCmmPages uint16
	iCmmPage  uint16

	helpfile           string
	outer              miBasicOuter
	keys               table.GTableKeys
	iKeyCounter        uint16
	kCmm               []uCommand
	mCmm               []uCommand
	pBrowserWsp        *mBrowserWsp
	pWorkspace         *maWorkspace
	isFocusable        bool
	tableHasChanged    bool
	isEmpty            bool
	_fGetBrowserString func() string
	fSetFocusToModel   func()
	_fSetFocusToModel  func()
	w                  vWindow
}

var x__ = "················································································"
