package main

import (
	"fmt"
	"goatplayer/internal/K"
	"goatplayer/internal/help"
	"goatplayer/internal/style"
	"strings"
)

func (this *mHelp) SetOuter() {
	this.mBasic.outer = this
	this.outer = this
}

func (this *mHelp) InitModel(name string, sFocus string, folderName string) {
	if this.outer == nil {
		this.SetOuter()
	}
	this.mBasic.InitModel(name, sFocus, "")
	this.w2.InitWindow(&this.mBasic)
	this.w2.isPage = true
	this.folderName = folderName
	this.w.emptyLine = ""
	this.w.emptyLine = style.Normal.EchoStyle(this.w.emptyLine, this.w.fGetWidth())
	this.w2.emptyLine = ""
	this.w2.emptyLine = style.Normal.EchoStyle(this.w2.emptyLine, this.w2.fGetWidth())
	// trace.BeginEndAdd_n(this, "InitModel", "folder name: %s", this.folderName) // t //
}

func (this *mHelp) SetFilename(filename string) {
	this.filename = filename
}

func (this *mHelp) PrepareModel() {
	// trace.BeginSilent_n(this, "PrepareModel") // t //
	this.mBasic.PrepareModel()
	this.keys.SelectedI = 0
	this.keys.ActiveI = K.UNSET
	this.w.lines.SelectedI = 0
	// trace.EndAdd("len %d", len(this.keys.Sl)) // t //
}

/* Prepares the list of help files. */
func (this *mHelp) PrepareView() {
	// trace.BeginSilent_n(this, "PrepareView") // t //
	this.Reset()

	////!old
	// entries, _ := os.ReadDir(this.folderName)
	// trace.AddToPrint("found %d entries", len(entries))
	// for _, entry := range entries {
	// 	this.AppendKeyAndLine(entry.Name(), &style.Normal)
	// }
	////!new
	for _, entry := range help.AllTitles {
		this.AppendKeyAndLine(entry, &style.Normal)
	}
	////!

	this.w.SetEmptyLineForDialog()
	this.SelectKeyAndLine(0)
	this.LoadFile(false)
	// trace.End() // t //
}

func (this *mHelp) SelectKey_o() {
	// trace.Begin_n(this, "SelectKey_o {mHelp}") // t //
	this.LoadFile(true)
	// trace.End() // t //
}

/* Loads the content of a helpfile into the page. */
func (this *mHelp) LoadFile(doPrint bool) {
	var style1 *style.GStyle
	// var filepath string
	// trace.Begin_n(this, "LoadFile") // t //
	if this.filename == "" {
		this.filename = this.keys.GetSelected().V
		// trace.Print("no filename set, use selected: %s", this.filename) // t //
	} else {
		// trace.Print("filename set: %s", this.filename) // t //
		key := this.keys.GetIndexByFilename(this.filename)
		if key == K.UNSET {
			this.keys.SelectedI = 0
			this.w.lines.SelectedI = 0
			Term.StatusLine.PrintMessage(fmt.Sprintf("Help file '%s' not found", this.filename))
			this.filename = this.keys.GetSelected().V
		} else {
			this.keys.SelectedI = key
			this.w.lines.SelectedI = key
		}
	}
	this.w2.lines.Reset()
	// filepath = this.folderName + "/" + this.filename
	// trace.PrintGreen("path: %s", filepath)

	////!old
	// bytes, _ := os.ReadFile(filepath)
	// s := string(bytes)
	////new
	s := help.MepTitles[this.filename]
	////!

	lines := strings.Split(*s, "\n")
	for _, line := range lines {
		// trace.PrintGreen("-%s-", line)
		style1 = &style.Normal
		if len(line) >= 3 {
			switch line[0:3] {
			case "<b>":
				style1 = &style.PageBold
				line = line[3:]
			case "<i>":
				style1 = &style.PageItalic
				line = line[3:]
			case "<c>":
				style1 = &style.PageColored
				line = line[3:]
			case "<u>":
				// trace.Print("------------------ underline")
				style1 = &style.PageUnderline
				line = line[3:]
			}
		}
		this.w2.AppendLine("  "+line, false, style1)
	}
	// trace.Print("w2 lines len: %d", this.w2.lines.Len()) // t //
	this.w2.iPage = 0
	this.w2.nPages = ((this.w2.lines.Len() - 1) / this.w2.nRows) + 1
	if doPrint {
		this.w2.PrintPage()
	}
	this.filename = ""
	// trace.End() // t //
}

type mHelp struct {
	mBasic
	w2         vWindow
	folderName string
	filename   string
}
