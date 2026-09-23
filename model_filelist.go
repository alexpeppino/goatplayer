package main

import (
	"fmt"
	"slices"
	"strings"
)

func (this *mFilelist) SetOuter() {
	this.mBasic.outer = this
	this.mAbsFilelist.outer = this
	this.outer = this
}

func (this *mAbsFilelist) InitModel(name string, focusKey string, helpfile string) {
	this.mBasic.InitModel(name, focusKey, helpfile)
	AllMaps.allAbstractFilelists = append(AllMaps.allAbstractFilelists, this)
	this.w.isFolded = true
}

func (this *mFilelist) InitModel(name string, focusKey string, helpfile string) {
	if this.outer == nil {
		this.SetOuter()
	}
	this.mAbsFilelist.InitModel(name, focusKey, helpfile)
	AllMaps.mapFilelists[name] = this
	this.fSetFocusToModel = func() {
		MainScreen.activeWorkspace = &this.pBrowserWsp.maWorkspace
		this.UpdateDisplay()
	}
	this._fGetBrowserString = func() string {
		i := this.keys.GetSelected().dbIndex
		audiof := AF.GetAudiof(i)
		artist := audiof.artist
		title := audiof.title
		s := fmt.Sprintf("%s+%s", title, artist)
		s = strings.ReplaceAll(s, " ", "+")
		return s
	}
}

func (this *mFilelist) BringToFront_o() {
// trace.Begin_n(this, "BringToFront_o {mFilelist}") // t //
	MainScreen.pRightSection = this.pRightSection
// trace.Print("this.pSection: %s, activeWindow: %s", this.pRightSection.name, this.pRightSection.GetActiveWindow().name) // t //
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

func (this *mFilelist) UpdateDisplay() {
	Display.Part2[0] = fmt.Sprintf("Workspace:      %s", this.pWorkspace.publicName)
	Display.Part2[1] = fmt.Sprintf("Working dir:    %s", this.pBrowserWsp.Dirs.workingDir)
	Display.Part2[2] = ""
	Display.Part2[3] = ""
	Display.Part2[4] = ""
	Display.AddPartToPrint(2)
	Term.Flush()
}

func (this *mFilelist) SelectKey_o() {
// trace.Begin_n(this, "SelectKey_o {mFilelist}") // t //
	fileProp := &this.pBrowserWsp.AudioPropCon
	if fileProp.w.IsInForeground() {
		fileProp.LoadFile(this.keys.GetSelected().v)
		fileProp.PrepareView()
		fileProp.w.PrintWindowAnyway()
	}
	Term.StatusLine.SetContent(AF.sl[this.keys.GetSelected().dbIndex].title)
// trace.End() // t //
}

func (this *mAbsFilelist) PrepareAbstractFilelist() {
	this.mBasic.PrepareModel()
	AllMaps.mapBasicFilelists[this.name] = this
}

func (this *mFilelist) PrepareModel() {
	this.mAbsFilelist.PrepareAbstractFilelist()
	// this._fSelectKey = func() { this.fSelectKey() }
	this.tabCaption = this.pBrowserWsp.publicName
}

func (this *mFilelist) UseFoundItems_o() {
// trace.Print("%s use selected: %s, %s", this.GetLongName(), Find.GetSelected(), FindDialog.focusKeys.activeKey) // t //
}

// .
func (this *mAbsFilelist) AppendKey(st string, i uint16) {
	// trace.N_BeginEnd(this, "AppendKey, %s", st)
	this.keys.Append(&tTableKey{st, this.w.iLineCounter, i, false})
	// this.iKeyCounter++
}

// Multi-filters mode.
func (this *mFilelist) UseAllFilters() {}

func (this *mFilelist) PlayThisList() {
	Play.isEditor = false
	Play.UpdateDisplay()
// trace.Print("Play.isEditor: %v", Play.isEditor) // t //
	Play.PlayThatList(this)
}

func (this *mFilelist) PlayThisItem() {
// trace.Begin_n(this, "PlayThisItem") // t //
	Play.isEditor = false
	Play.UpdateDisplay()
// trace.Print("Play.isEditor: %v", Play.isEditor) // t //
	Play.PlayThatItem(this)
// trace.End() // t //
}

/* Appends a new folder line to the view. st: text of the line. */
func (this *mAbsFilelist) AppendFolderLine(s1, s2 *string) {
	// trace.N_BeginEnd(this, "AppendLine, %s", st)
	// trace.Print("AppendFolderLine %s %s", *s1, *s2)
	var text, normalLine string
	s3 := " · " + *s2
	// totalLen := RuneLen(s1) + RuneLen(&s3) + 1
	finalLen := this.w.nCols
	normalLine = styleFolder1.EchoStyle("  "+*s1, RuneLen(s1)+2)
	normalLine += styleFolder2.EchoStyle(s3, finalLen-RuneLen(s1)-4) + "  "
	this.w.lines.Append(&tTableLine{
		text,
		normalLine,
		"", "", "", UNSET})
	// this.iLineCounter++
}

func (this *mAbsFilelist) AppendLine(lineLeft, lineRight *string) {
	var iKey uint16
	iKey = this.iKeyCounter
	nCols := this.w.fGetWidth()
	// AdjustWidthLeftAndRight(lineLeft, *lineRight, nCols)
	act := "▶  " + (*lineLeft)[3:]
	this.w.lines.Append(&tTableLine{
		*lineLeft,
		styleNormal.EchoTwoStyles(*lineLeft, *lineRight, nCols, &styleNormalLight),
		styleSelected.EchoTwoStyles(*lineLeft, *lineRight, nCols, &styleSelectedLight),
		styleActive.EchoTwoStyles(act, *lineRight, nCols, &styleActive),
		styleSelActive.EchoTwoStyles(act, *lineRight, nCols, &styleSelActive),
		iKey})
}

// func (this *mPlaylist) LoadAndPrepareFromDbIndexes() {
// 	trace.Begin_n(this, "LoadAndPrepare")
// 	trace.Print("this.dbIndexes: %+v", this.dbIndexes)
// 	this.LoadAndPrepare(&this.dbIndexes)
// 	trace.End()
// }

func (this *mAbsFilelist) PrepareEmptyView() {
	this.isEmpty = true
	this.AppendKeyAndLine("The Playlist is empty.", nil)
	this.keys.Select(0)
	this.keys.iActive = UNSET
	this.w.lines.Select(0)
	this.AfterLoad()
	this.w.SetEmptyLine()
}

func MakeFunc_SortCATT(dbIndexes *tDbIndexes) f_SortSlice {
	return func(i, j int) bool {
		audiof_i := &AF.sl[dbIndexes.sl[i]]
		audiof_j := &AF.sl[dbIndexes.sl[j]]
		if audiof_i.composer != audiof_j.composer {
			return audiof_j.composer > audiof_i.composer
		} else {
			if audiof_i.album != audiof_j.album {
				return audiof_j.album > audiof_i.album
			} else {
				if audiof_i.iTrack != audiof_j.iTrack {
					return audiof_j.iTrack > audiof_i.iTrack
				} else {
					return audiof_j.title > audiof_i.title
				}
			}
		}
	}
}

/*
If not nil, dbIndexes is copied into this.dbIndexes, which is always used.
If dbIndexes is nil, old this.dbIndexes are used.
Folding is made line by line, when a new artist+album is found.
Folder lines and normal lines are added to the window; keys are added to the filelist.
*/
func (this *mAbsFilelist) LoadAndPrepare(dbIndexes *tDbIndexes) {
	var title string
	var audiof *aAudiofile
	var lineLeft, lineRight string
	var newArtist, newAlbum bool
	var fGetField func(audiof *aAudiofile) string
	w := &this.w
// trace.Begin_n(this, "LoadAndPrepare") // t //

	/// If dbIndexes is nil, old this.dbIndexes are used
	if dbIndexes != nil {
// trace.Print("dbIndex is not nil, len is %d. Copy it to this.dbIndexes...", len(dbIndexes.sl)) // t //
		this.dbIndexes.CopyFrom(dbIndexes)
	}
	if len(this.dbIndexes.sl) == 0 {
// trace.Print("this.dbIndex len == 0, PrepareEmptyView...") // t //
		this.isEmpty = true
		this.PrepareEmptyView()
		this.w.SetEmptyLine()
// trace.ReturnAdd("empty dbIndexes") // t //
		return
	} else {
// trace.Print("this.dbIndex len == %d", len(this.dbIndexes.sl)) // t //
	}

	if this.pBrowserWsp == &BrowserWsp1 {
		if Set.Bro1Folding.GetValue() == "Artist" {
// trace.Print("folding value: Artist") // t //
			fGetField = GetAudiofArtist
		} else {
// trace.Print("folding value: Composer") // t //
			fGetField = GetAudiofComposer
			this.dbIndexes.SortSlice(MakeFunc_SortCATT(&this.dbIndexes))
		}
	} else if this.pBrowserWsp == &BrowserWsp2 {
		if Set.Bro2Folding.GetValue() == "Artist" {
// trace.Print("folding value: Artist") // t //
			fGetField = GetAudiofArtist
		} else {
// trace.Print("folding value: Composer") // t //
			fGetField = GetAudiofComposer
			this.dbIndexes.SortSlice(MakeFunc_SortCATT(&this.dbIndexes))
			// this.dbIndexes.SortSlice(fSortCATT)
		}
	} else if this == &Play.mAbsFilelist { // Player
// trace.Print("player") // t //
		if Play.callingBrowser == &BrowserWsp1 {
			fGetField = AllMaps.allGetAudiof[Set.Bro1Folding.GetValue()]
		} else if Play.callingBrowser == &BrowserWsp2 {
			fGetField = AllMaps.allGetAudiof[Set.Bro2Folding.GetValue()]
		} else {
			fGetField = GetAudiofArtist
		}
	} else {
// trace.Print("??") // t //
	}

	this.isEmpty = false
	// trace.Print("dbIndexes: %+v", dbIndexes)
	// trace.Print("this.dbIndexes: %+v", this.dbIndexes)
	this.Reset()
	this.totDuration = 0
	//
	addFolderLine := func(s1, s2 string) {
		// s := fmt.Sprintf(" %s", DC_BOLD+audiof.album+DC_RESET+Pen.Normal)
		// AdjustWidth(&s, this.nCols-1)
		if s1 == "" {
			s1 = "(None)"
		}
		if s2 == "" {
			s2 = "(None)"
		}
		this.AppendFolderLine(&s1, &s2)
		w.iLineCounter++
	}
	addItemLine := func(audiof *aAudiofile) {
		// this.filenames = append(this.filenames, 	t__FoldedFilename{audiof.filename, this.line_counter})
		// artist = audiof.artist
		// album = audiof.album
		if audiof.title == UNKNOWN_TITLE {
			title = audiof.filename
		} else {
			title = audiof.title
		}
		// AdjustWidth(&artist, 27)
		// AdjustWidth(&album, 27)
		du := audiof.duration
		this.totDuration += du
		duFormat := dura.FormatSeconds(du/100, 2)
		// line = fmt.Sprintf(" %2d %6s  %s", audiof.track_i, duFormat, title)
		lineLeft = fmt.Sprintf("    %s", title)
		lineRight = fmt.Sprintf("%2d %5s  ", audiof.iTrack, duFormat)
		// line = fmt.Sprintf(" %-30s%-30s%s", artist, album, title)
		// line := fmt.Sprintf("    %s", audiof.title)
		this.AppendLine(&lineLeft, &lineRight)
		this.AppendKey(audiof.filename, audiof.dbIndex)
		this.iKeyCounter++
		w.iLineCounter++
	}
	oldArtist := ""
	oldAlbum := ""
	for _, i := range this.dbIndexes.sl {
		audiof = &AF.sl[i]
		// trace.Print("%4d  %s", i, audiof.title)
		// trace.Print("filter %s", audiof.title)
		// GetAudiofArtist()
		//!! if audiof.artist != oldArtist {
		if fGetField(audiof) != oldArtist {
			oldArtist = fGetField(audiof)
			newArtist = true
		}
		if audiof.album != oldAlbum {
			oldAlbum = audiof.album
			newAlbum = true
		}
		if newAlbum || newArtist {
			addFolderLine(oldArtist, audiof.album)
		}
		newArtist = false
		newAlbum = false
		addItemLine(audiof) //this.filenames = append(this.filenames, t__FoldedFilename{audiof.filename, i})
	}
	this.keys.TraceSelected("")
	this.w.lines.TraceSelected("")
	this.AfterLoad()
	this.PlaylistHasChanged()
	this.keys.iSelected = 0
	this.w.lines.iSelected = this.keys.GetSelected().iLine
	this.w.lines.TraceSelected("")
	this.w.SetEmptyLine()
	this.strDuration = dura.FormatSeconds(this.totDuration/100, 1)
// trace.AddToPrint("len(this.keys): %d", this.keys.Len()) // t //
// trace.AddToPrint("len(this.lines): %d", w.lines.Len()) // t //
// trace.End() // t //
}

func (this *mAbsFilelist) PlaylistHasChanged() {
	if this == &Play.mAbsFilelist {
		ui_pl.playlistHasChanged = true
		// trace.N_BeginEnd(this, "ActivePLHasChanged")
	}
}

func (this *mFilelist) LoadFolder(that *mDirs) {
	this.container = nil //&this.pBrowser.Dirs.mContainer
	dir := that.keys.GetSelected().v
// trace.BeginAdd_n(this, "LoadFolder", "[%s]", dir) // t //
	// folder := that.filenames.Get(that.selected_i)
	this.list = that.GetName() + " ▶ " + dir
	this.Reset()
	MusicScanner.Reset()
	MusicScanner.FindFiles(dir)
	// i used a slice of indexes to maintain the order of the audiofiles
	for _, filename := range MusicScanner.fileList {
		AF.selectedIndexes.sl = append(AF.selectedIndexes.sl, AF.GetIndexByFilename(filename))
	}
	slices.Sort(AF.selectedIndexes.sl)
	// trace.Print("%s, %v", dir, DB.selectedIndexes)
	this.LoadAndPrepare(&AF.selectedIndexes)
	AF.selectedIndexes.sl = nil
	// this.AfterLoad()
	// this.ActivePLHasChanged()
// trace.BeginEnd_n(this, "LoadFolder") // t //
}

func (this *mFilelist) AddSelectedToPlaylist() {
// trace.Begin_n(this, "AddSelectedToPlaylist") // t //
	// s := this.keys.GetActive().v
	if Play.isEditor {
		Play.data.AppendLine(&tPFLine{Play.data.activeListname, 0, this.keys.GetSelected().v, 0, 0, 0})
		Play.MakeIndexesFromPFList()
	} else {
		Play.dbIndexes.sl = append(Play.dbIndexes.sl, this.keys.GetSelected().dbIndex)
	}
	Play.LoadAndPrepare(nil)
	// iSl := AF.relIndexDbSl[this.keys.GetSelected().i]
	// ActivePL.indexes = append(ActivePL.indexes, iSl)
	// trace.Print("ActivePL.indexes: %v", ActivePL.indexes)
	// ActivePL.LoadAndPrepare(
	// 	func(audiof *aAudiofile) bool {
	// 		return true //audiof.dir == dir+"/"
	// 	}, &ActivePL.indexes)
// trace.End() // t //
}

// func (this *mFilelist) AddSelectedToEditor() {
// 	trace.Begin_n(this, "AddSelectedToEditor")
// 	// s := this.keys.GetActive().v

// 	Edit.data.AppendLine(&tPFLine{Edit.data.activeListname, 0, this.keys.GetSelected().v, 0, 0, 0})
// 	Edit.PrepareIndexes()
// 	Edit.LoadAndPrepare(nil)

// 	// iSl := AF.relIndexDbSl[this.keys.GetSelected().i]
// 	// ActivePL.indexes = append(ActivePL.indexes, iSl)
// 	// trace.Print("ActivePL.indexes: %v", ActivePL.indexes)
// 	// ActivePL.LoadAndPrepare(
// 	// 	func(audiof *aAudiofile) bool {
// 	// 		return true //audiof.dir == dir+"/"
// 	// 	}, &ActivePL.indexes)
// 	trace.End()
// }

func (this *mFilelist) LoadArtist() {
	artists := &this.pBrowserWsp.ArtistsCon
	this.container = artists
// trace.Begin_n(this, "LoadArtist") // t //
	artist := artists.keys.GetSelected().v
// trace.Print("---- artist: %s", artist) // t //
	this.list = artists.GetName() + " ▶ " + artist
// trace.Print("this.list: %s", this.list) // t //

	this.dbIndexes.MakeSliceByFilter(
		func(audiof *aAudiofile) bool {
			return audiof.artist == artist
		})
// trace.Print("%+v", this.dbIndexes) // t //
	this.LoadAndPrepare(nil)

	// this.LoadAndPrepare(func(audiof *aAudiofile) bool {
	// 	return audiof.artist == artist
	// }, nil)

// trace.End() // t //
}

func (this *mFilelist) LoadAlbum() {
	albums := &this.pBrowserWsp.AlbumsCon
	this.container = albums
	album := albums.keys.GetSelected().v
	artist := albums.sa1.Get(albums.keys.iSelected)
	this.list = albums.GetName() + " ▶ " + album + " / " + artist

	this.dbIndexes.MakeSliceByFilter(
		func(audiof *aAudiofile) bool {
			return audiof.artist == artist && audiof.album == album
		})
	this.LoadAndPrepare(nil)

	// this.LoadAndPrepare(func(audiof *aAudiofile) bool {
	// 	return audiof.artist == artist && audiof.album == album
	// }, nil)

// trace.Print("album: %s", album)   //@1 // t //
// trace.Print("artist: %s", artist) //@1 // t //
// trace.EndAdd("found %d files", this.keys.Len()) // t //
}

func (this *mFilelist) LoadYear(year string) {
	this.container = &this.pBrowserWsp.YearsCon
	this.list = fmt.Sprintf("Years ▶ %s", year)
	this.LoadAndPrepare(nil)
// trace.BeginEndAdd_n(this, "LoadYear", "%d", year) // t //
}

func (this *mFilelist) LoadGenre(genre string) {
	this.container = &this.pBrowserWsp.GenresCon
	this.list = "Genres ▶ " + genre
	this.LoadAndPrepare(nil)
// trace.BeginEndAdd_n(this, "LoadGenre [%s], found %d files", genre) // t //
}

func (this *mAbsFilelist) LoadPFList(pflistname string) {
	this.container = &PFLists.mContainer
// trace.Begin_n(this, "LoadPFList") // t //
	this.list = PFLists.GetName() + " ▶ " + pflistname
// trace.Print("this.list: %s", this.list) //@1 // t //
	for i := range PFListsData.sl {
		r := &PFListsData.sl[i]
		if r.listname == pflistname {
			if dbIndex := AF.GetIndexByFilename(r.filename); dbIndex != UNSET {
				AF.selectedIndexes.sl = append(AF.selectedIndexes.sl, dbIndex)
			}
		}
	}
	if len(AF.selectedIndexes.sl) == 0 {
		this.PrepareEmptyView()
// trace.ReturnAdd("playlist is empty") // t //
		return
	}
	this.LoadAndPrepare(&AF.selectedIndexes)
	AF.selectedIndexes.sl = nil
// trace.End() // t //
}

// Adds the rows to the PfListsDB.
func (this *mFilelist) SavePFList() {
	listName := Input.ReadCommand(MainScreen.LeftSection.marginBottom+1, 1, "List name: ")
// trace.BeginEndAdd_n(this, "MakePFList", "%s", listName) // t //
	PFListsData.DeleteList(listName)
	for i, key := range this.keys.sl {
		PFListsData.AppendLine(&tPFLine{listName, uint16(i), key.v, 0, 0, 0})
	}
	PFListsData.SaveDataToFile()
	PFLists.LoadListnames()
	// PfListsDB.Read()
}

type mAbsFilelist struct {
	mBasic
	pDbIndexes    *tDbIndexes // deprecated
	dbIndexes     tDbIndexes
	totDuration   uint32
	strDuration   string
	container     *mContainer
	list          string
	pRightSection *vSection
}

type mFilelist struct {
	mAbsFilelist
	outer    miFilelistOuter
	listName string
}

