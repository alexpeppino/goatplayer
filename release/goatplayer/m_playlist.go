package main

import (
	"fmt"
	AF "goatplayer/internal/audiofiles"
	"goatplayer/internal/echo"
	"goatplayer/internal/input"
	"goatplayer/internal/mpv"
	"goatplayer/internal/small/dura"
	"slices"
	"sort"
	"strings"
)

func (this *mPlaylist) AppendFromList() {
	Find.title = "Find and select the audiofile to append"
	this.AddFromList(1)
	ui_pl.playlistHasChanged = false
}

func (this *mPlaylist) InsertFromList() {
	Find.title = "Find and select the audiofile to insert"
	this.AddFromList(2)
	ui_pl.playlistHasChanged = false
}

func (this *mPlaylist) AddFromList(mode uint16) {
	// trace.Begin_n(this, "AddFromList") // t //
	var ssl []tStringsAndIndexes
	// var oldSelectedDbIndex = this.keys.GetSelected().dbIndex
	var oldActiveDbIndex = this.keys.GetActive().DbIndex
	// trace.Print("old selected index: %d  (%s)", oldSelectedDbIndex, AF.GetAudiof(oldSelectedDbIndex).title)
	// trace.Print("old active index: %d  (%s)", oldActiveDbIndex, AF.GetAudiof(oldActiveDbIndex).Title) // t //
	ssl = make([]tStringsAndIndexes, len(AF.Sl))
	for i := range AF.Sl {
		audiof := &AF.Sl[i]
		ssl[i].str = audiof.Title
		ssl[i].str2 = fmt.Sprintf("    {%s · %s}", audiof.Artist, audiof.Album)
		ssl[i].lowerStr = strings.ToLower(audiof.Title)
		ssl[i].dbIndex = AF.GetIndexByFilename(audiof.Filename)
	}
	// trace.Print("%+v", ssl)
	sort.Slice(ssl, func(i int, j int) bool {
		return ssl[i].str < ssl[j].str
	})
	if this.OpenFindDialog2(&ssl) {
		var newDbIndex = Find.keys.GetSelected().DbIndex
		// trace.Print("selected index: %d, %s", newDbIndex, AF.Sl[newDbIndex].Title) // t //
		if mode == 1 {
			this.dbIndexes.Sl = append(this.dbIndexes.Sl, newDbIndex)
		} else {
			this.dbIndexes.Sl = slices.Insert(this.dbIndexes.Sl, int(this.keys.ActiveI)+1, newDbIndex)
		}
		// trace.Print("%+v", this.dbIndexes)
		this.LoadAndPrepare(nil)
		this.SelectKeyAndLine(this.keys.GetIndexByDbIndex(oldActiveDbIndex))
		this.keys.ActiveI = this.keys.SelectedI
		this.w.PrintWindowAnyway()
		echo.Flush()
		// this.PlaySelected()
	}
	// trace.End() // t //
}

func (this *mPlaylist) UpdateDisplay() {
	this.pWorkspace = &PlayerWsp.maWorkspace
	this.pWorkspace.publicName = "Player"
	Display.Part2[0] = fmt.Sprintf("Workspace:      %s", this.pWorkspace.publicName)
	Display.Part2[1] = fmt.Sprintf("Is editor:      %v", this.isEditor)
	Display.Part2[2] = fmt.Sprintf("List name:      %v", this.data.activeListname)
	Display.Part2[3] = ""
	Display.Part2[4] = ""
	Display.AddPartToPrint(2)
}

func (this *mPlaylist) SelectKey_o() {
	// trace.Begin_n(this, "SelectKey_o {mPlaylist}") // t //
	if this.isEmpty {
		// trace.ReturnAdd("is empty") // t //
		return
	}
	if this.isEditor {
		// PlayerWsp.Param.PrepareView()
	}
	fileProp := this.FileProp
	if k := this.keys.GetSelected(); k != nil {
		fileProp.LoadFile(k.V)
	}
	fileProp.PrepareView()
	if fileProp.w.IsInForeground() {
		fileProp.w.PrintWindowAnyway()
		//! } else if PlayerWsp.Param.w.IsInForeground() {
		// 	PlayerWsp.Param.w.PrintWindowAnyway()
	}
	Term.StatusLine.SetContent(AF.Sl[this.keys.GetSelected().DbIndex].Title)
	// trace.End() // t //
}

func (this *mPlaylist) PlayThatList(that *mFilelist) {
	// trace.BeginAdd_n(this, "PlayThatList", "from %s", that.name) // t //
	this.callingBrowser = that.pBrowserWsp
	this.LoadAndPrepare(&that.dbIndexes)
	// trace.Print("that selected %d", that.keys.SelectedI) // t //
	this.SelectKeyAndLine(that.keys.SelectedI)

	this.list = "From " + that.pBrowserWsp.publicName
	this.SetStatusLine_o()
	this.w.PrintWindowAnyway()

	this.PlaySelected()
	this.GetScreen().ActivateModel(&this.mBasic)
	// this.w.BringToForeground()
	// trace.End() // t //
}

func (this *mPlaylist) PlayThatItem(that *mFilelist) {
	var dbIndexes AF.GDbIndexes
	// trace.BeginAdd_n(this, "PlayThatItem", "from %s", that.name) // t //
	this.callingBrowser = that.pBrowserWsp
	dbIndexes.AppendIndex(that.dbIndexes.Sl[that.keys.SelectedI])
	this.LoadAndPrepare(&dbIndexes)
	this.SelectKeyAndLine(0)

	this.list = "From " + that.pBrowserWsp.publicName
	this.SetStatusLine_o()
	this.w.PrintWindowAnyway()

	this.PlaySelected()
	this.GetScreen().ActivateModel(&this.mBasic)
	// this.w.BringToForeground()
	// trace.End() // t //
}

func (this *mPlaylist) ClearPlayList() {
	// trace.Begin_n(this, "ClearPlayList") // t //
	this.keys.Reset()
	this.w.lines.Reset()
	// trace.End() // t //
}

func player_WakeUp() {
	ui_pl.signal <- 1
}

func (this *mPlaylist) PlaySelected() {
	// trace.Begin_n(this, "PlaySelected") // t //
	//	Play selected file
	if !this.isEmpty {
		ui_pl.stopPlaying = false
		this.keys.ActiveI = this.keys.SelectedI
		this.w.UpdateActiveLines()
		this.w.pLayer.pSection.PrintFocusButtons(true)
		if ui_pl.playerIsWaiting {
			player_WakeUp()
		} else {
			ui_pl.playSelectedFile = true
			mpv.Stop()
		}
	}
	// trace.End() // t //
}

func (this *mPlaylist) MakeIndexesFromPFList() {
	this.dbIndexes.Sl = nil
	for i := range this.data.sl {
		line := &this.data.sl[i]
		this.dbIndexes.Sl = append(this.dbIndexes.Sl, AF.MepFilenameIndex[line.filename])
	}
}

func (this *mPlaylist) LoadList(listname string) {
	this.data.Reset()
	this.data.activeListname = listname
	for i := range PFListsData.sl {
		line := &PFListsData.sl[i]
		if line.listname == listname {
			this.data.AppendLine(line)
		}
	}
}

/* dbIndexes are already there, just prepare the pflist data */
func (this *mPlaylist) MakePFListFromIndexes() {
	// trace.Begin_n(this, "MakePFListFromIndexes") // t //
	// trace.Print("this.data.activeListname: %s", this.data.activeListname) // t //
	this.data.sl = nil
	for j, i := range this.dbIndexes.Sl {
		this.data.AppendLine(&tPFLine{this.data.activeListname, i, this.keys.Get(uint16(j)).V, 0, 0, 0})
	}
	// trace.End() // t //
}

// func (this *mPlaylist) MakeAndSaveList() {
// 	trace.Begin_n(this, "CopyPFListToPFData")
// 	if !this.isEditor {
// 		this.RenameListAskUser()
// 		this.MakePFListFromIndexes()
// 		this.list = "List: " + this.data.activeListname
// 		this.SetStatusLine_o()
// 		this.w.PrintWindowAnyway()
// 		this.isEditor = true
// 		this.UpdateDisplay()
// 	}
// 	// trace.PrintIndentedV("data:  %#v", this.data)
// 	PFListsData.DeleteList(this.data.activeListname)
// 	for i := range this.data.sl {
// 		line := &this.data.sl[i]
// 		line.iTrack = uint16(i)
// 		line.listname = this.data.activeListname
// 		trace.Print("append line:  %#v", line)
// 		PFListsData.AppendLine(line)
// 	}
// 	PFListsData.SaveDataToFile()
// 	PlayerWsp.PFListsCon.ReloadAndPrint()
// 	// trace.PrintIndentedV("PFListsData:  %#v", PFListsData)
// 	trace.End()
// }

/* Get the list's name from user */
func (this *mPlaylist) RenameListAskUser() {
	s := input.ReadLineAtPos(54, 144, 40)
	// trace.Print("'%s'", s) // t //
	this.data.activeListname = s
}

func (this *mPlaylist) NewList() {
	this.data.activeListname = PFLISTS_TMP
	this.w.statusLine = this.data.activeListname
	this.w.PrintStatusLine()
	echo.Flush()
}

func (this *mPlaylist) BringToFront_o() {
	// trace.Begin_n(this, "BringToFront_o {mPlaylist}") // t //
	MainScreen.pRightSection = this.pRightSection
	if this.pRightSection.activeLayer == nil {
		this.pRightSection.activeLayer = &this.pRightSection.layers[0]
	}
	if this.pRightSection.activeLayer.activeWindow == nil {
		this.pRightSection.activeLayer.activeWindow = this.pRightSection.activeLayer.windows[0]
	}
	// trace.Print("this.pSection: %s, activeWindow: %s", this.pRightSection.name, this.pRightSection.GetActiveWindow().name) // t //
	// trace.Print("MainScreen.PlayerSection.GetActiveWindow: %s", MainScreen.PlayerSection.GetActiveWindow().name) // t //
	///
	MainScreen.allFocusableSections = nil
	MainScreen.allFocusableSections = append(MainScreen.allFocusableSections,
		&MainScreen.LeftSection,
		MainScreen.pRightSection)
	MainScreen.UpdateFocusButtonList()
	///
	MainScreen.pRightSection.PrintFocusButtons(false)

	//!!
	// MainScreen.pRightSection.activeWindow.BringWindowToFront(false)
	MainScreen.pRightSection.GetActiveWindow().BringWindowToFront(false)

	// MainScreen.pSectionRight.UpdateFrame(false)
	// trace.End() // t //
}

// func (this *mPlaylist) BringToFront_o() {
// 	trace.Begin_n(this, "BringToFront_o {mEditor}")
// 	this.pRightSection = &MainScreen.EditorSection
// 	MainScreen.pRightSection = this.pRightSection
// 	this.pRightSection.activeLayer = &this.pRightSection.layers[0]
// 	this.pRightSection.activeLayer.activeWindow = this.pRightSection.activeLayer.windows[0]
// 	trace.Print("this.pSection: %s, activeWindow: %s", this.pRightSection.name, this.pRightSection.GetActiveWindow().name)
// 	///
// 	MainScreen.allFocusableSections = nil
// 	MainScreen.allFocusableSections = append(MainScreen.allFocusableSections,
// 		&MainScreen.LeftSection,
// 		MainScreen.pRightSection)
// 	MainScreen.UpdateFocusButtonList()
// 	///
// 	MainScreen.pRightSection.PrintFocusButtons(false)

// 	//!!
// 	// MainScreen.pRightSection.activeWindow.BringWindowToFront(false)
// 	MainScreen.pRightSection.GetActiveWindow().BringWindowToFront(false)

// 	// MainScreen.pSectionRight.UpdateFrame(false)
// 	trace.End()
// }

func (this *mPlaylist) SetOuter() {
	this.outer = this
	this.mAbsFilelist.outer = this
	this.mBasic.outer = this
}

func (this *mParameters) SetOuter() {
	this.mBasic.outer = this
	this.outer = this
}

func (this *mPlaylist) InitModel(name string, focusKey string, helpfile string) {
	if this.outer == nil {
		this.SetOuter()
	}
	this.mAbsFilelist.InitModel(name, focusKey, helpfile)
	this.isEmpty = true
	this.fSetFocusToModel = func() {
		MainScreen.activeWorkspace = &this.pPlayerW.maWorkspace
		this.UpdateDisplay()
	}
}

func (this *mParameters) UpdateViewIfInForeground() {
	if this.w.IsInForeground() {
		this.PrepareView()
		this.w.PrintWindowAnyway()
	}
}

// func (this *mPlaylist) SetStart() {
// 	i := Mpv.GetFileTime()
// 	this.data.sl[this.keys.IActive].start = i
// 	trace.Print("start: %d", i)
// 	PlayerWsp.Param.UpdateViewIfInForeground()
// }

// func (this *mPlaylist) UnsetStart() {
// 	this.data.sl[this.keys.IActive].start = 0
// 	PlayerWsp.Param.UpdateViewIfInForeground()
// }

// func (this *mPlaylist) UnsetEnd() {
// 	this.data.sl[this.keys.IActive].end = 0
// 	PlayerWsp.Param.UpdateViewIfInForeground()
// }

// func (this *mPlaylist) UnsetSpeed() {
// 	this.data.sl[this.keys.IActive].speed = 0
// 	PlayerWsp.Param.UpdateViewIfInForeground()
// }

// func (this *mPlaylist) SetEnd() {
// 	i := Mpv.GetFileTime()
// 	this.data.sl[this.keys.IActive].end = i
// 	trace.Print("end: %d", i)
// 	PlayerWsp.Param.UpdateViewIfInForeground()
// }

// func (this *mPlaylist) SetSpeed() {
// 	this.data.sl[this.keys.IActive].speed = ui_pl.Speed_i
// 	trace.Print("speed: %d", ui_pl.Speed_i)
// 	PlayerWsp.Param.UpdateViewIfInForeground()
// }

// func (this *mPlaylist) PrepareModel() {
// 	this.maFilelist.PrepareAbstractFilelist()
// 	this.tabCaption = "Editor"
// }

func (this *mPlaylist) PrepareModel() {
	this.mAbsFilelist.PrepareAbstractFilelist()
	this.tabCaption = "Player"
}

func (this *mParameters) InitModel(name string, focusKey string, helpfile string) {
	if this.outer == nil {
		this.SetOuter()
	}
	this.mBasic.InitModel(name, focusKey, helpfile)
}

func (this *mParameters) LoadList() {}

func (this *mParameters) PrepareView() {
	// trace.Begin_n(this, "PrepareView") // t //
	this.keys.Sl = nil
	this.w.lines.Sl = nil
	if len(Play.data.sl) > 0 {
		i := Play.keys.SelectedI
		line := &Play.data.sl[i]
		this.AppendKeyAndLine(fmt.Sprintf("listname:     %s", line.listname), nil)
		this.AppendKeyAndLine(fmt.Sprintf("iTrack:       %d", line.iTrack), nil)
		this.AppendKeyAndLine(fmt.Sprintf("filename:     %s", line.filename), nil)
		this.AppendKeyAndLine(fmt.Sprintf("start:        %s", dura.FormatWithCents(line.start)), nil)
		this.AppendKeyAndLine(fmt.Sprintf("end:          %s", dura.FormatWithCents(line.end)), nil)
		this.AppendKeyAndLine(fmt.Sprintf("speed:        %.2f", float32(line.speed)/100), nil)
	}
	this.w.SetEmptyLine()
	// trace.EndAdd("%d lines", this.w.lines.Len()) // t //
}

func (this *mPlaylist) CalcPrecDuration() {
	// trace.Begin_n(this, "CalcPrecDuration") // t //
	iActive := this.keys.ActiveI
	if iActive == 0 {
		this.precDuration = 0
		// trace.ReturnAdd("is the first") // t //
		return
	}
	this.precDuration = 0
	this.keys.TraceActive("")
	for i := range iActive {
		du := AF.GetAudiof(this.dbIndexes.Sl[i]).Duration
		// trace.Print("duration: %v", du)
		this.precDuration += du
	}
	// trace.Print("precDuration: %d", this.precDuration) // t //
	// trace.End() // t //

}

type mPlaylist struct {
	mAbsFilelist
	pPlayerW *mPlayerWsp
	// pEditorWsp *mEditorWsp
	FileProp       *mAudioProperties
	data           aPFLists
	isEditor       bool
	i              int
	precDuration   uint32
	callingBrowser *mBrowserWsp
	isMuted        bool
	isPaused       bool
}

// type mEditor struct {
// 	maFilelist
// }

type mParameters struct {
	mBasic
}
