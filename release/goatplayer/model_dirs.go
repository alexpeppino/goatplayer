package main

import (
	"fmt"
	"strings"
)

func (this *mDirs) SetOuter() {
	this.outer = this
	this.mBasic.outer = this
}

func (this *mDirs) InitModel(name string, sFocus string, filename string) {
	if this.outer == nil {
		this.SetOuter()
	}
	this.mBasic.InitModel(name, sFocus, filename)
	// this.mContainer.PrepareModel(name, sFocus, true)
}

func (this *mDirs) PrepareModel() {
	this.mBasic.PrepareModel()
	// this.mContainer.PrepareModel(name, sFocus, true)
}

func (this *mDirs) Load() {}

/*
Set selected as working dir for the browser.
this.pBrowserWsp.dbIndexes is updated.
*/
func (this *mDirs) SetSelectedAsWorkingDir() {
// trace.Begin_n(this, "SetSelectedAsWorkingDir") // t //
	workingDir := this.keys.GetSelected().v + "/"
	this.workingDir = workingDir
// trace.Print("selected: %s", workingDir) // t //
	var n int
	this.pBrowserWsp.wspDbIndexes.ClearSlice()
	for i := range AF.sl {
		audiof := &AF.sl[i]
		if strings.HasPrefix(audiof.filename, workingDir) {
			this.pBrowserWsp.wspDbIndexes.AppendIndex(uint16(i))
// trace.Print("%3d   %s", n, audiof.filename) // t //
			n++
		}
	}

	// this.pBrowserWsp.TraceIndexes()

	// trace.Print("slice: %+v", this.pBrowserWsp.dbIndexes.sl)
	// trace.Print("len: %d", len(this.pBrowserWsp.dbIndexes.sl))
	this.pBrowserWsp.View.LoadAndPrepare(&this.pBrowserWsp.wspDbIndexes)
	s1 := fmt.Sprintf("%s, working dir set to: %s", this.pBrowserWsp.publicName, workingDir)
	Term.StatusLine.PrintMessage(s1)
	this.pBrowserWsp.View.UpdateDisplay()
	// Term.AddStringWithPos(60, 0, styleMessage.EchoStyle(s1, Term.nCols-2))
	this.pBrowserWsp.View.w.PrintWindowAnyway()
	AF.PrepareContainers(this.pBrowserWsp, &this.pBrowserWsp.wspDbIndexes)
	this.pBrowserWsp.GenresCon.PrepareView()
	this.pBrowserWsp.GenresCon.w.PrintWindowAnyway()
	this.pBrowserWsp.ArtistsCon.PrepareView()
	this.pBrowserWsp.ArtistsCon.w.PrintWindowAnyway()
	this.pBrowserWsp.AlbumsCon.PrepareView()
	this.pBrowserWsp.AlbumsCon.w.PrintWindowAnyway()
	// this.pBrowserWsp.GenresCon.LoadFromIndexesAndPrint(&this.pBrowserWsp.dbIndexes, true)
	// this.pBrowserWsp.ArtistsCon.LoadFromIndexesAndPrint(&this.pBrowserWsp.dbIndexes, true)
	// this.pBrowserWsp.AlbumsCon.LoadFromIndexesAndPrint(&this.pBrowserWsp.dbIndexes, true)
// trace.End() // t //
}

// { Folders >> load to Filelist }
func (this *mDirs) LoadTo(that *mFilelist) {
// trace.BeginEndAdd_n(this, "LoadTo", "%s", that.name) //@1 // t //
	if this.nMusicFiles[this.keys.iSelected] == 0 {
		return
	}
	this.keys.iActive = this.keys.iSelected
	this.w.UpdateActiveLines()
	this.w.pLayer.pSection.PrintFocusButtons(true)
	// Screen.FileContainers.actLayer = this.w.layer
	that.LoadFolder(this)
	// that.PrepareLines()
	that.w.PrintWindow()
}

func (this *mDirs) PrepareView() {
// trace.BeginSilent_n(this, "PrepareView") // t //
	var tab string
	w := &this.w
	// this.lines = make([]t__TableLine, len(DirScanner.ValidDirs))
	// this.keys = make([]t__TableKey, len(DirScanner.ValidDirs))
	MusicScanner.Reset()
	var musicPath = AF.musicDir
	MusicScanner.FindValidDirs(musicPath, "", 0)
	w.lines.Make(len(MusicScanner.validDirs))
	this.keys.Make(len(MusicScanner.validDirs))
	this.nMusicFiles = make([]uint16, len(MusicScanner.validDirs))
	var str string
	this.w.isFolded = false
// trace.AddToPrint("%d lines", w.lines.Len()) // t //
	this.keys.iActive = UNSET
	var arrow string
	for i, v := range MusicScanner.validDirs {
		if i < len(MusicScanner.validDirs)-1 {
			if MusicScanner.validDirs[i+1].level > v.level {
				arrow = "▼"
			} else {
				arrow = "▶"
			}
		} else {
			arrow = "▶"
		}
		this.keys.sl[i].v = v.dirPath
		this.keys.sl[i].iLine = uint16(i)
		this.nMusicFiles[i] = v.n
		tab = strings.Repeat("   ", int(v.level))
		// if v.level > 0 {
		// 	tab = strings.Repeat("  ", int(v.level-1)) + "↳ "  // ↴↳➤▻►
		// } else {
		// 	tab = ""
		// }
		if v.dirPath == musicPath {
			str = v.dirPath
		} else {
			str = v.basename
		}
		// AdjustWidth(&str, uint16(91-len(tab)))
		// v2 := fmt.Sprintf(" %s%s  %3d", tab, str, v.n)
		// AdjustWidth(&v2, w.nCols)
		p := this.w.lines.Get(uint16(i))
		if v.n > 0 {
			width := this.w.nCols - 1
			v2 := fmt.Sprintf(" %s%s 🎵 %s", tab, arrow, str)
			nstr := fmt.Sprintf("  %d", v.n)
			AdjustWidthLeftAndRight(&v2, nstr, width)
			p.normal = styleNormal.EchoStyle(v2, width)
			p.selected = styleSelected.EchoStyle(v2, width)
		} else {
			width := this.w.nCols
			v2 := fmt.Sprintf(" %s▶    %s", tab, str)
			nstr := fmt.Sprintf("  %d", v.n)
			AdjustWidthLeftAndRight(&v2, nstr, width)
			p.normal = styleGrey.EchoStyle(v2, width)
			p.selected = styleSelected.EchoStyle(v2, width)
		}
		p.iKey = uint16(i)
	}
	//	Empty line
	this.w.SetEmptyLine()
	this.AfterLoad()
// trace.End() // t //
}

type mDirs struct {
	mBasic
	outer       miContainerOuter
	nMusicFiles []uint16
	workingDir  string
}

