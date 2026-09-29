package main

import (
	"fmt"
	"goatplayer/internal/K"
	AF "goatplayer/internal/audiofiles"
	"goatplayer/internal/musicscanner"
	"goatplayer/internal/small/utils"
	"goatplayer/internal/style"
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
	workingDir := this.keys.GetSelected().V + "/"
	this.workingDir = workingDir
	// trace.Print("selected: %s", workingDir) // t //
	var n int
	this.pBrowserWsp.wspDbIndexes.ClearSlice()
	for i := range AF.Sl {
		audiof := &AF.Sl[i]
		if strings.HasPrefix(audiof.Filename, workingDir) {
			this.pBrowserWsp.wspDbIndexes.AppendIndex(uint16(i))
			// trace.Print("%3d   %s", n, audiof.Filename) // t //
			n++
		}
	}

	// this.pBrowserWsp.TraceIndexes()

	// trace.Print("slice: %+v", this.pBrowserWsp.dbIndexes.Sl)
	// trace.Print("len: %d", len(this.pBrowserWsp.dbIndexes.Sl))
	this.pBrowserWsp.View.LoadAndPrepare(&this.pBrowserWsp.wspDbIndexes)
	s1 := fmt.Sprintf("%s, working dir set to: %s", this.pBrowserWsp.publicName, workingDir)
	Term.StatusLine.PrintMessage(s1)
	this.pBrowserWsp.View.UpdateDisplay()
	// echo.AddStringWithPos(60, 0, styleMessage.EchoStyle(s1, Term.nCols-2))
	this.pBrowserWsp.View.w.PrintWindowAnyway()
	this.pBrowserWsp.PrepareContainers(&this.pBrowserWsp.wspDbIndexes)
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
	if this.nMusicFiles[this.keys.SelectedI] == 0 {
		return
	}
	this.keys.ActiveI = this.keys.SelectedI
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
	musicscanner.Reset()
	var musicPath = AF.MusicDir
	musicscanner.FindValidDirs(musicPath, "", 0)
	w.lines.Make(len(musicscanner.ValidDirs))
	this.keys.Make(len(musicscanner.ValidDirs))
	this.nMusicFiles = make([]uint16, len(musicscanner.ValidDirs))
	var str string
	this.w.isFolded = false
	// trace.AddToPrint("%d lines", w.lines.Len()) // t //
	this.keys.ActiveI = K.UNSET
	var arrow string
	for i, v := range musicscanner.ValidDirs {
		if i < len(musicscanner.ValidDirs)-1 {
			if musicscanner.ValidDirs[i+1].Level > v.Level {
				arrow = "▼"
			} else {
				arrow = "▶"
			}
		} else {
			arrow = "▶"
		}
		this.keys.Sl[i].V = v.DirPath
		this.keys.Sl[i].LineI = uint16(i)
		this.nMusicFiles[i] = v.N
		tab = strings.Repeat("   ", int(v.Level))
		// if v.level > 0 {
		// 	tab = strings.Repeat("  ", int(v.level-1)) + "↳ "  // ↴↳➤▻►
		// } else {
		// 	tab = ""
		// }
		if v.DirPath == musicPath {
			str = v.DirPath
		} else {
			str = v.Basename
		}
		// AdjustWidth(&str, uint16(91-len(tab)))
		// v2 := fmt.Sprintf(" %s%s  %3d", tab, str, v.n)
		// AdjustWidth(&v2, w.nCols)
		p := this.w.lines.Get(uint16(i))
		if v.N > 0 {
			width := this.w.nCols - 1
			v2 := fmt.Sprintf(" %s%s 🎵 %s", tab, arrow, str)
			nstr := fmt.Sprintf("  %d", v.N)
			utils.AdjustWidthLeftAndRight(&v2, nstr, width)
			p.Normal = style.Normal.EchoStyle(v2, width)
			p.Selected = style.Selected.EchoStyle(v2, width)
		} else {
			width := this.w.nCols
			v2 := fmt.Sprintf(" %s▶    %s", tab, str)
			nstr := fmt.Sprintf("  %d", v.N)
			utils.AdjustWidthLeftAndRight(&v2, nstr, width)
			p.Normal = style.Grey.EchoStyle(v2, width)
			p.Selected = style.Selected.EchoStyle(v2, width)
		}
		p.IKey = uint16(i)
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
