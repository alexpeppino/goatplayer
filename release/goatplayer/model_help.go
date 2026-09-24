package main

import (
	"fmt"
	"os"
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
	this.w.emptyLine = styleNormal.EchoStyle(this.w.emptyLine, this.w.fGetWidth())
	this.w2.emptyLine = ""
	this.w2.emptyLine = styleNormal.EchoStyle(this.w2.emptyLine, this.w2.fGetWidth())
// trace.BeginEndAdd_n(this, "InitModel", "folder name: %s", this.folderName) // t //
}

func (this *mHelp) SetFilename(s string) {
	this.filename = s
}

func (this *mHelp) PrepareModel() {
// trace.BeginSilent_n(this, "PrepareModel") // t //
	this.mBasic.PrepareModel()
	this.keys.iSelected = 0
	this.keys.iActive = UNSET
	this.w.lines.iSelected = 0
// trace.EndAdd("len %d", len(this.keys.sl)) // t //
}

/* Prepares the list of help files. */
func (this *mHelp) PrepareView() {
// trace.BeginSilent_n(this, "PrepareView") // t //
	this.Reset()
	entries, _ := os.ReadDir(this.folderName)
// trace.AddToPrint("found %d entries", len(entries)) // t //
	for _, entry := range entries {
		this.AppendKeyAndLine(entry.Name(), &styleNormal)
	}
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
	var style *cStyle
	var filepath string
// trace.Begin_n(this, "LoadFile") // t //
	if this.filename == "" {
		this.filename = this.keys.GetSelected().v
// trace.Print("no filename set, use selected: %s", this.filename) // t //
	} else {
// trace.Print("filename set: %s", this.filename) // t //
		key := this.keys.GetIndexByFilename(this.filename)
		if key == UNSET {
			this.keys.iSelected = 0
			this.w.lines.iSelected = 0
			Term.StatusLine.PrintMessage(fmt.Sprintf("Help file '%s' not found", this.filename))
			this.filename = this.keys.GetSelected().v
		} else {
			this.keys.iSelected = key
			this.w.lines.iSelected = key
		}
	}
	this.w2.lines.Reset()
	filepath = this.folderName + "/" + this.filename
// trace.PrintGreen("path: %s", filepath) // t //

	bytes, _ := os.ReadFile(filepath)
	// bytes, err := exec.Command("fold", "-sw", "80", filepath).Output()
	// trace.PrintGreen("%v", err)
	s := string(bytes)
	lines := strings.Split(s, "\n")
	for _, line := range lines {
		// trace.PrintGreen("-%s-", line)
		style = &styleNormal
		if len(line) >= 3 {
			switch line[0:3] {
			case "<b>":
				style = &stylePageBold
				line = line[3:]
			case "<i>":
				style = &stylePageItalic
				line = line[3:]
			case "<c>":
				style = &stylePageColored
				line = line[3:]
			case "<u>":
				// trace.Print("------------------ underline")
				style = &stylePageUnderline
				line = line[3:]
			}
		}
		this.w2.AppendLine("  "+line, false, style)
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

