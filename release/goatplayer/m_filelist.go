package main

import (
	"fmt"
	"goatplayer/internal/K"
	AF "goatplayer/internal/audiofiles"
	"goatplayer/internal/echo"
	"goatplayer/internal/input"
	"goatplayer/internal/musicscanner"
	"goatplayer/internal/small/dura"
	"goatplayer/internal/small/utils"
	"goatplayer/internal/style"
	"goatplayer/internal/table"
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
		i := this.keys.GetSelected().DbIndex
		audiof := AF.GetAudiof(i)
		artist := audiof.Artist
		title := audiof.Title
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
	echo.Flush()
}

func (this *mFilelist) SelectKey_o() {
	// trace.Begin_n(this, "SelectKey_o {mFilelist}") // t //
	fileProp := &this.pBrowserWsp.AudioPropCon
	if fileProp.w.IsInForeground() {
		fileProp.LoadFile(this.keys.GetSelected().V)
		fileProp.PrepareView()
		fileProp.w.PrintWindowAnyway()
	}
	Term.StatusLine.SetContent(AF.Sl[this.keys.GetSelected().DbIndex].Title)
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
	this.keys.Append(&table.TTableKey{V: st, LineI: this.w.iLineCounter, DbIndex: i, IsActiveFilter: false})
	// this.IKeyCounter++
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
	normalLine = style.Folder1.EchoStyle("  "+*s1, utils.RuneLen(s1)+2)
	normalLine += style.Folder2.EchoStyle(s3, finalLen-utils.RuneLen(s1)-4) + "  "
	this.w.lines.Append(&table.TTableLine{Text: text, Normal: normalLine, Selected: "", Active: "", SelectedActive: "", IKey: K.UNSET})
	// this.ILineCounter++
}

func (this *mAbsFilelist) AppendLine(lineLeft, lineRight *string) {
	var iKey uint16
	iKey = this.iKeyCounter
	nCols := this.w.fGetWidth()
	// AdjustWidthLeftAndRight(lineLeft, *lineRight, nCols)
	act := "▶  " + (*lineLeft)[3:]
	this.w.lines.Append(&table.TTableLine{
		Text:           *lineLeft,
		Normal:         style.Normal.EchoTwoStyles(*lineLeft, *lineRight, nCols, &style.NormalLight),
		Selected:       style.Selected.EchoTwoStyles(*lineLeft, *lineRight, nCols, &style.SelectedLight),
		Active:         style.Active.EchoTwoStyles(act, *lineRight, nCols, &style.Active),
		SelectedActive: style.SelActive.EchoTwoStyles(act, *lineRight, nCols, &style.SelActive),
		IKey:           iKey})
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
	this.keys.ActiveI = K.UNSET
	this.w.lines.Select(0)
	this.AfterLoad()
	this.w.SetEmptyLine()
}

/*
If not nil, dbIndexes is copied into this.dbIndexes, which is always used.
If dbIndexes is nil, old this.dbIndexes are used.
Folding is made line by line, when a new artist+album is found.
Folder lines and normal lines are added to the window; keys are added to the filelist.
*/
func (this *mAbsFilelist) LoadAndPrepare(dbIndexes *AF.GDbIndexes) {
	var title string
	var audiof *AF.TAudiofile
	var lineLeft, lineRight string
	var newArtist, newAlbum bool
	var fGetField func(audiof *AF.TAudiofile) string
	w := &this.w
	// trace.Begin_n(this, "LoadAndPrepare") // t //

	/// If dbIndexes is nil, old this.dbIndexes are used
	if dbIndexes != nil {
		// trace.Print("dbIndex is not nil, len is %d. Copy it to this.dbIndexes...", len(dbIndexes.Sl)) // t //
		this.dbIndexes.CopyFrom(dbIndexes)
	}
	if len(this.dbIndexes.Sl) == 0 {
		// trace.Print("this.dbIndex len == 0, PrepareEmptyView...") // t //
		this.isEmpty = true
		this.PrepareEmptyView()
		this.w.SetEmptyLine()
		// trace.ReturnAdd("empty dbIndexes") // t //
		return
	} else {
		// trace.Print("this.dbIndex len == %d", len(this.dbIndexes.Sl)) // t //
	}

	if this.pBrowserWsp == &BrowserWsp1 {
		if Set.Bro1Folding.GetValue() == "Artist" {
			// trace.Print("folding value: Artist") // t //
			fGetField = AF.GetAudiofArtist
		} else {
			// trace.Print("folding value: Composer") // t //
			fGetField = AF.GetAudiofComposer
			this.dbIndexes.SortSlice(AF.MakeFunc_SortCATT(&this.dbIndexes))
		}
	} else if this.pBrowserWsp == &BrowserWsp2 {
		if Set.Bro2Folding.GetValue() == "Artist" {
			// trace.Print("folding value: Artist") // t //
			fGetField = AF.GetAudiofArtist
		} else {
			// trace.Print("folding value: Composer") // t //
			fGetField = AF.GetAudiofComposer
			this.dbIndexes.SortSlice(AF.MakeFunc_SortCATT(&this.dbIndexes))
			// this.dbIndexes.SortSlice(fSortCATT)
		}
	} else if this == &Play.mAbsFilelist { // Player
		// trace.Print("player") // t //
		if Play.callingBrowser == &BrowserWsp1 {
			fGetField = AllMaps.allGetAudiof[Set.Bro1Folding.GetValue()]
		} else if Play.callingBrowser == &BrowserWsp2 {
			fGetField = AllMaps.allGetAudiof[Set.Bro2Folding.GetValue()]
		} else {
			fGetField = AF.GetAudiofArtist
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
		// s := fmt.Sprintf(" %s", DC_BOLD+audiof.Album+DC_RESET+Pen.Normal)
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
	addItemLine := func(audiof *AF.TAudiofile) {
		// this.filenames = append(this.filenames, 	t__FoldedFilename{audiof.Filename, this.line_counter})
		// artist = audiof.Artist
		// album = audiof.Album
		if audiof.Title == K.UNKNOWN_TITLE {
			title = audiof.Filename
		} else {
			title = audiof.Title
		}
		// AdjustWidth(&artist, 27)
		// AdjustWidth(&album, 27)
		du := audiof.Duration
		this.totDuration += du
		duFormat := dura.FormatSeconds(du/100, 2)
		// line = fmt.Sprintf(" %2d %6s  %s", audiof.track_i, duFormat, title)
		lineLeft = fmt.Sprintf("    %s", title)
		lineRight = fmt.Sprintf("%2d %5s  ", audiof.TrackI, duFormat)
		// line = fmt.Sprintf(" %-30s%-30s%s", artist, album, title)
		// line := fmt.Sprintf("    %s", audiof.Title)
		this.AppendLine(&lineLeft, &lineRight)
		this.AppendKey(audiof.Filename, audiof.DbIndex)
		this.iKeyCounter++
		w.iLineCounter++
	}
	oldArtist := ""
	oldAlbum := ""
	for _, i := range this.dbIndexes.Sl {
		audiof = &AF.Sl[i]
		// trace.Print("%4d  %s", i, audiof.Title)
		// trace.Print("filter %s", audiof.Title)
		// GetAudiofArtist()
		//!! if audiof.Artist != oldArtist {
		if fGetField(audiof) != oldArtist {
			oldArtist = fGetField(audiof)
			newArtist = true
		}
		if audiof.Album != oldAlbum {
			oldAlbum = audiof.Album
			newAlbum = true
		}
		if newAlbum || newArtist {
			addFolderLine(oldArtist, audiof.Album)
		}
		newArtist = false
		newAlbum = false
		addItemLine(audiof) //this.filenames = append(this.filenames, t__FoldedFilename{audiof.Filename, i})
	}
	this.keys.TraceSelected("")
	this.w.lines.TraceSelected("")
	this.AfterLoad()
	this.PlaylistHasChanged()
	this.keys.SelectedI = 0
	this.w.lines.SelectedI = this.keys.GetSelected().LineI
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
	dir := that.keys.GetSelected().V
	// trace.BeginAdd_n(this, "LoadFolder", "[%s]", dir) // t //
	// folder := that.filenames.Get(that.Selected_i)
	this.list = that.GetName() + " ▶ " + dir
	this.Reset()
	musicscanner.Reset()
	musicscanner.FindFiles(dir)
	// i used a slice of indexes to maintain the order of the audiofiles
	for _, filename := range musicscanner.FileList {
		AF.SelectedIndexes.Sl = append(AF.SelectedIndexes.Sl, AF.GetIndexByFilename(filename))
	}
	slices.Sort(AF.SelectedIndexes.Sl)
	// trace.Print("%s, %v", dir, DB.SelectedIndexes)
	this.LoadAndPrepare(&AF.SelectedIndexes)
	AF.SelectedIndexes.Sl = nil
	// this.AfterLoad()
	// this.ActivePLHasChanged()
	// trace.BeginEnd_n(this, "LoadFolder") // t //
}

func (this *mFilelist) AddSelectedToPlaylist() {
	// trace.Begin_n(this, "AddSelectedToPlaylist") // t //
	// s := this.keys.GetActive().v
	if Play.isEditor {
		Play.data.AppendLine(&tPFLine{Play.data.activeListname, 0, this.keys.GetSelected().V, 0, 0, 0})
		Play.MakeIndexesFromPFList()
	} else {
		Play.dbIndexes.Sl = append(Play.dbIndexes.Sl, this.keys.GetSelected().DbIndex)
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

// 	Edit.data.AppendLine(&tPFLine{Edit.data.activeListname, 0, this.keys.GetSelected().V, 0, 0, 0})
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
	artist := artists.keys.GetSelected().V
	// trace.Print("---- artist: %s", artist) // t //
	this.list = artists.GetName() + " ▶ " + artist
	// trace.Print("this.list: %s", this.list) // t //

	this.dbIndexes.MakeSliceByFilter(
		func(audiof *AF.TAudiofile) bool {
			return audiof.Artist == artist
		})
	// trace.Print("%+v", this.dbIndexes) // t //
	this.LoadAndPrepare(nil)

	// this.LoadAndPrepare(func(audiof *aAudiofile) bool {
	// 	return audiof.Artist == artist
	// }, nil)

	// trace.End() // t //
}

func (this *mFilelist) LoadAlbum() {
	albums := &this.pBrowserWsp.AlbumsCon
	this.container = albums
	album := albums.keys.GetSelected().V
	artist := albums.sa1.Get(albums.keys.SelectedI)
	this.list = albums.GetName() + " ▶ " + album + " / " + artist

	this.dbIndexes.MakeSliceByFilter(
		func(audiof *AF.TAudiofile) bool {
			return audiof.Artist == artist && audiof.Album == album
		})
	this.LoadAndPrepare(nil)

	// this.LoadAndPrepare(func(audiof *aAudiofile) bool {
	// 	return audiof.Artist == artist && audiof.Album == album
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
	// this.container = &PFLists.mContainer
	// trace.Begin_n(this, "LoadPFList") // t //
	// this.list = PFLists.GetName() + " ▶ " + pflistname
	// trace.Print("this.list: %s", this.list) //@1 // t //
	for i := range PFListsData.sl {
		r := &PFListsData.sl[i]
		if r.listname == pflistname {
			if dbIndex := AF.GetIndexByFilename(r.filename); dbIndex != K.UNSET {
				AF.SelectedIndexes.Sl = append(AF.SelectedIndexes.Sl, dbIndex)
			}
		}
	}
	if len(AF.SelectedIndexes.Sl) == 0 {
		this.PrepareEmptyView()
		// trace.ReturnAdd("playlist is empty") // t //
		return
	}
	this.LoadAndPrepare(&AF.SelectedIndexes)
	AF.SelectedIndexes.Sl = nil
	// trace.End() // t //
}

// Adds the rows to the PfListsDB.
func (this *mFilelist) SavePFList() {
	listName := input.ReadCommand(MainScreen.LeftSection.marginBottom+1, 1, "List name: ")
	// trace.BeginEndAdd_n(this, "MakePFList", "%s", listName) // t //
	PFListsData.DeleteList(listName)
	for i, key := range this.keys.Sl {
		PFListsData.AppendLine(&tPFLine{listName, uint16(i), key.V, 0, 0, 0})
	}
	PFListsData.SaveDataToFile()
	// PFLists.LoadListnames()
	// PfListsDB.Read()
}

type mAbsFilelist struct {
	mBasic
	pDbIndexes    AF.GDbIndexes // deprecated
	dbIndexes     AF.GDbIndexes
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
