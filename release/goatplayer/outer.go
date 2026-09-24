package main

import "fmt"

type miBasicOuter interface {
	GetName() string
	CheckInput_o(bool) string
	CheckMouseInput_o()
	SetStatusLine_o()
	SelectActive_o()
	UseFoundItems_o()
	SelectKey_o()
	BringToFront_o()
}

func (this *mBasic) SetStatusLine_o() {}
func (this *mBasic) UseFoundItems_o() {}
func (this *mBasic) SelectKey_o() {
// trace.Begin_n(this, "SelectKey_o {mBasic}") // t //
// trace.End() // t //
}
func (this *mBasic) BringToFront_o() {}

type miContainerOuter interface {
	miBasicOuter
	LoadTo(*mFilelist)
}

type miFilelistOuter interface {
	miBasicOuter
}

func (this *vMainScreen) BringToFront_o() {
// trace.BeginEnd_n(this, "BringToFront_o {vMainScreen}") // t //

}

func (this *vAbsDialogScreen) BringToFront_o() {
// trace.BeginEnd_n(this, "BringToFront_o {vDialog}") // t //
	Term.AddSimple(this.closeButton)
}

// type miDirsOuter interface {}

//////	Set outer

func (this *mDialog) SetOuter() {
	this.mBasic.outer = this
	this.outer = this
}

func (this *mFind) SetOuter() {
	this.mBasic.outer = this
	this.outer = this
}

func (this *mSettings) SetOuter() {
	this.mBasic.outer = this
	this.outer = this
}

func (this *mColorsPicker) SetOuter() {
	this.mBasic.outer = this
	this.outer = this
}

// func (this *mDirs) SetOuter(outer miDirsOuter) {}

//////	Simple Outer methods

func (this *mContainer) SetStatusLine_o() {
	count := this.keys.GetActiveFiltersCount()
	if count == 0 {
		this.w.statusLine = "{all}"
	} else {
		this.w.statusLine = fmt.Sprintf("{%d}", count)
	}
}

func (this *mDirs) SetStatusLine_o()   { this.w.statusLine = this.keys.GetSelected().v }
func (this *mDialog) SetStatusLine_o() { this.w.statusLine = "Select a line and press Enter/Esc" }
func (this *mFind) SetStatusLine_o()   { this.w.statusLine = this.modelName }

func (this *mFilelist) SetStatusLine_o() {
	this.w.statusLine = fmt.Sprintf("[%s] %s", this.strDuration, this.list)
}

func (this *mPlaylist) SetStatusLine_o() {
	this.w.statusLine = fmt.Sprintf("[%s] %s", this.strDuration, this.list)
}

// func (this *mBasic) CheckInput(bool) string {
// 	return ""
// }

// func (this *mBasic) CheckMouseInput() {}

