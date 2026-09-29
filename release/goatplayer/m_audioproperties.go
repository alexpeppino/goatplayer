package main

import (
	"fmt"
	"goatplayer/internal/K"
	AF "goatplayer/internal/audiofiles"
	"goatplayer/internal/input"
	"goatplayer/internal/small/dura"
	"goatplayer/internal/small/utils"
	"goatplayer/internal/style"
	"goatplayer/internal/table"
	"slices"

	taglib "github.com/wtolson/go-taglib"
)

func (this *mAudioProperties) SetOuter() {
	this.outer = this
	this.mContainer.outer = this
	this.mBasic.outer = this
}

func (this *mAudioProperties) InitModel(name string, sFocus string, helpfile string) {
	if this.outer == nil {
		this.SetOuter()
	}
	this.mContainer.InitModel(name, sFocus, helpfile)
}

// { FileProp >> load to Playlist }
// func (this *mAudioProperties) LoadTo(that *mFilelist) {
// 	var selTag = this.tagNames.Get(this.keys.ISelected)
// 	var tagValue string
// 	this.keys.IActive = this.keys.ISelected
// 	this.w.UpdateActiveLines()
// 	this.w.pLayer.pSection.PrintFocusButtons(true)
// 	switch selTag {
// 	case "Year":
// 		that.LoadYear(this.year)
// 		// Tags.LoadYears()
// 		// Tags.PrepareLines()
// 		// this.LoadYearTo(that)
// 		tagValue = this.year
// 	case "Genre":
// 		tagValue = this.genre
// 		that.LoadGenre(this.genre)
// 		// Tags.LoadGenres()
// 		// Tags.PrepareLines()
// 		// this.LoadGenreTo(that)
// 	case "Album":
// 		tagValue = this.album
// 		this.pBrowserWsp.AlbumsCon.SelectIndexByKey(this.album)
// 		this.pBrowserWsp.AlbumsCon.LoadTo(that)
// 		// this.LoadAlbumTo(that)
// 	}
// 	// that.PrepareLines()
// 	that.w.PrintWindow()
// 	trace.BeginEndAdd_n(this, "LoadTo", "%s, %s: %s", that.name, selTag, tagValue)
// }

func (this *mAudioProperties) EditSelected() {
	// trace.Begin_n(this, "EditSelected") // t //

	possible := []uint16{0, 1, 2, 4, 5}
	if !slices.Contains(possible, this.keys.SelectedI) {
		Term.StatusLine.PrintMessage("Can't modifiy this field")
		// trace.Return() // t //
		return
	}

	this.keys.TraceSelected("")
	sel := this.keys.SelectedI
	newValue := input.ReadLineAtPos(sel+3, this.w.marginLeft+19, 60)
	// trace.Print("input: '%s'", newValue) // t //
	audiof := AF.GetAudiofByFilename(this.pWorkspace.pFilelist.keys.GetSelected().V)
	filename := audiof.Filename
	// trace.Print("audiof: '%s' -- %s", audiof.Title, filename) // t //

	fd, err := taglib.Read(filename)
	if err != nil {
		// log.Panic(err)
	}

	switch sel {
	case 0:
		// trace.Print("SetTitle...") // t //
		fd.SetTitle(newValue)
	case 1:
		// trace.Print("SetArtist...") // t //
		fd.SetArtist(newValue)
	case 2:
		// trace.Print("SetAlbum...") // t //
		fd.SetAlbum(newValue)
	case 4:
		// trace.Print("SetGenre...") // t //
		fd.SetGenre(newValue)
	case 5:
		// trace.Print("SetYear...") // t //
		fd.SetYear(int(utils.Str2uint(newValue)))
	default:

	}
	// trace.Print("saving the new data to the audiofile...") // t //
	fd.Save()
	fd.Close()

	// trace.Print("reloading audiof from filesystem to db... (dbIndex is %d)", audiof.DbIndex) // t //
	AF.AddAudiofileFromFileSystem(int(audiof.DbIndex), filename)
	AF.DoSortAndSaveAtQuit = true
	// trace.Print("done") // t //
	this.LoadFile(this.pWorkspace.pFilelist.keys.GetSelected().V)
	this.PrepareView()
	this.w.PrintWindowAnyway()
	// trace.End() // t //
}

func (this *mAudioProperties) OpenFindWithAllValues(con *mContainer) {
	Find.LoadList1(&con.allValues, &con.mBasic)
	// FindScreen.activeSection = &FindScreen.Section
	FindDialog.activeWindow = &Find.w
	Find.keys.SelectedI = 0
	Find.w.lines.SelectedI = 0
	MainScreen.BeforeDialog()
	FindDialog.OpenDialog(K.FOCUS_FIND)
	MainScreen.AfterDialog()
}

func (this *mAudioProperties) EditSelectedFromFind() {
	// trace.Begin_n(this, "EditSelectedFromFind") // t //

	possible := []uint16{1, 2, 4}
	if !slices.Contains(possible, this.keys.SelectedI) {
		Term.StatusLine.PrintMessage("Can't modifiy this field")
		// trace.Return() // t //
		return
	}

	this.keys.TraceSelected("")
	audiof := AF.GetAudiofByFilename(this.pWorkspace.pFilelist.keys.GetSelected().V)
	filename := audiof.Filename
	// trace.Print("audiof: '%s' -- %s", audiof.Title, filename) // t //
	fd, err := taglib.Read(filename)
	if err != nil {
		// log.Panic(err)
	}
	defer fd.Close()
	switch this.keys.SelectedI {
	case 1:
		this.OpenFindWithAllValues(&this.pBrowserWsp.ArtistsCon)
		// trace.Print("SetArtist to %s...", Find.GetSelected()) // t //
		fd.SetArtist(Find.GetSelected())
	case 2:
		this.OpenFindWithAllValues(&this.pBrowserWsp.AlbumsCon)
		// trace.Print("SetAlbum to %s...", Find.GetSelected()) // t //
		fd.SetAlbum(Find.GetSelected())
	case 4:
		this.OpenFindWithAllValues(&this.pBrowserWsp.GenresCon)
		// trace.Print("SetGenre to %s...", Find.GetSelected()) // t //
		fd.SetGenre(Find.GetSelected())
	}
	// trace.Print("saving the new data to the audiofile...") // t //
	fd.Save()
	// trace.Print("reloading audiof from filesystem to db... (dbIndex is %d)", audiof.DbIndex) // t //
	AF.AddAudiofileFromFileSystem(int(audiof.DbIndex), filename)
	// trace.Print("sorting and saving the database...") // t //
	AF.DoSortAndSaveAtQuit = true
	// trace.Print("done") // t //
	this.LoadFile(this.pWorkspace.pFilelist.keys.GetSelected().V)
	this.PrepareView()
	this.w.PrintWindowAnyway()
	// trace.End() // t //
}

// Loads metadata and other properties for the given audiofile ('filename').
func (this *mAudioProperties) LoadFile(filename string) {
	// trace.BeginSilent_n(this, "LoadFile") // t //
	// trace.AddToPrint("filename: %s", filename) //@1 // t //
	// var counter uint16
	this.Reset()
	this.tagNames.Reset()
	this.tagValues.Reset()
	// GetMetadata(filename) // still needed for some tags
	audiof := AF.GetAudiofByFilename(filename)
	var tagValue, line, tagName string
	var i, j int
	append := func(tagName string) {
		// line = fmt.Sprintf("  %-15s %s", tagName, tagValue)
		// AdjustWidth(&line, this.w.nCols)
		this.tagNames.Append(tagName)
		this.tagValues.Append(tagValue)
		this.AppendKey(tagValue)
	}
	///
	tagValue = audiof.Title
	append("Title")
	///
	tagValue = audiof.Artist
	this.artist = tagValue
	append("Artist")
	///
	tagValue = audiof.Album
	this.album = tagValue
	append("Album")
	///
	i = int(audiof.TrackI)
	j = int(audiof.TrackN)
	tagName = "Track"
	if i == 0 {
		line = "(none)"
	} else if j > 0 {
		line = fmt.Sprintf("%d/%d", i, j)
	} else {
		line = fmt.Sprintf("%d", i)
	}
	utils.AdjustWidth(&line, this.w.nCols)
	this.tagNames.Append(tagName)
	this.tagValues.Append(line)
	this.AppendKey(line)
	///
	tagValue = audiof.Genre
	this.genre = tagValue
	append("Genre")
	///
	tagValue = audiof.Year
	this.year = tagValue
	append("Year")
	///
	tagValue = audiof.Composer
	this.composer = tagValue
	append("Composer")
	///
	tagValue = audiof.AlbumArtist
	this.albumArtist = tagValue
	append("AlbumArtist")
	///
	tagValue = dura.FormatSeconds(audiof.Duration/100, 1)
	append("Duration")
	///
	tagValue = audiof.ChannelsN
	this.modTime = tagValue
	append("Channels")
	///
	tagValue = audiof.SampleRate
	this.modTime = tagValue
	append("SampleRate")
	///
	tagValue = audiof.BitRate
	this.modTime = tagValue
	append("BitRate")
	///
	// tagValue = audiof.comment
	// this.comment = tagValue
	// append("Comment")
	///
	tagValue = audiof.ModTime
	this.modTime = tagValue
	append("ModTime")
	///
	tagValue = filename
	append("Filename")
	///
	tagValue = utils.Uint2str(audiof.DbIndex)
	append("DbIndex")
	///
	// trace.Print("this.genre = %s", this.genre)
	this.AfterLoad()
	// trace.End() // t //
}

func (this *mAudioProperties) PrepareView() {
	// trace.BeginSilent_n(this, "PrepareView") // t //
	var firstColLen uint16 = 18
	this.w.lines.Make(int(this.keys.Len()))
	// trace.N_Begin(this, "PrepareLines")
	for i := range this.w.lines.Sl {
		line := &this.w.lines.Sl[i]
		// w.SetLine(uint16(i), key.v, uint16(i), 0)
		s1 := "  " + this.tagNames.Get(uint16(i))
		utils.AdjustWidth(&s1, firstColLen)
		s2 := this.tagValues.Get(uint16(i))
		utils.AdjustWidth(&s2, this.w.nCols-firstColLen-2)
		s2 += "  "
		line.Normal = style.NormalLight.EchoStyle(s1, 0) + style.Normal.EchoStyle(s2, 0)
		line.Selected = style.SelectedLight.EchoStyle(s1, 0) + style.Selected.EchoStyle(s2, 0)
		line.IKey = uint16(i)
	}
	//	Empty line
	this.w.SetEmptyLine()
	// trace.EndAdd("%d lines", this.w.lines.Len()) // t //
}

type mAudioProperties struct {
	mContainer
	values      []string
	tagNames    table.GStringArray
	tagValues   table.GStringArray
	artist      string
	album       string
	genre       string
	year        string
	year_i      uint16
	composer    string
	albumArtist string
	comment     string
	modTime     string
	filename    string
}
