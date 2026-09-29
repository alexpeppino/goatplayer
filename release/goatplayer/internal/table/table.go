package table

import (
	"goatplayer/internal/K"
	"goatplayer/internal/mpv"
)

func (this *GTableKeys) Reset() {
	this.Sl = nil
	// this.iSelected = 0

	//!
	if this.OldSelectedI == 0 {
		this.SelectedI = 0
		this.OldSelectedI = K.UNSET
	}
	//!better
	// if this.iSelected == 0 && this.IOldSelected == 0 { // this is true only at the beginning
	// 	this.iSelected = 0
	// 	this.IOldSelected = 1
	// }

	this.ActiveI = K.UNSET
	this.OldActiveI = K.UNSET
}

func (this *GTableKeys) Unselect() {
	this.SelectedI = K.UNSET
}

func (this *GTableKeys) Make(size int) {
	this.Sl = make([]TTableKey, size)
}

func (this *GTableKeys) Append(x *TTableKey) {
	this.Sl = append(this.Sl, *x)
}

func (this *GTableKeys) GetLast() *TTableKey {
	return &this.Sl[len(this.Sl)-1]
}

/* Selects the item #i */
func (this *GTableKeys) Select(i uint16) bool {
	if i < this.Len() {
		this.OldSelectedI = this.SelectedI
		this.SelectedI = i
		return true
	}
// trace.Error("Keys, select, invalid index: %d", i) // t //
	return false
}

func (this *GTableKeys) Activate(i uint16) bool {
	if i < this.Len() {
		this.OldActiveI = this.ActiveI
		this.ActiveI = i
		return true
	}
// trace.Error("Keys, activate, invalid index: %d", i) // t //
	return false
}

func (this *GTableKeys) GetSelected() *TTableKey {
	return this.Get(this.SelectedI)
}

func (this *GTableKeys) GetActive() *TTableKey {
	return this.Get(this.ActiveI)
}

func (this *GTableKeys) GetIndexByFilename(filename string) uint16 {
	var i uint16
	for i = range uint16(len(this.Sl)) {
		key := &this.Sl[i]
		if key.V == filename {
			return i
		}
	}
	return K.UNSET
}

func (this *GTableKeys) GetIndexByDbIndex(dbIndex uint16) uint16 {
	var i uint16
	for i = range uint16(len(this.Sl)) {
		key := &this.Sl[i]
		if key.DbIndex == dbIndex {
			return i
		}
	}
	return K.UNSET
}

// Returns 'nil' if index is invalid.
func (this *GTableKeys) Get(i uint16) *TTableKey {
	if i >= this.Len() {
/* // b //
		var s string
		if i == K.UNSET {
			s = "unset"
		} else {
			s = "invalid"
		}
*/ // e //
// trace.Error("Keys, index is %s: %d/%d", s, i, this.Len()) // t //
		return nil
	}
	return &this.Sl[i]
}

func (this *GTableKeys) Len() uint16 {
	return uint16(len(this.Sl))
}

func (this *GTableKeys) TraceSelected(tab string) {
	if p := this.GetSelected(); p != nil {
// trace.Print("%skeys.ISelected: %d/%d (%d) => '%s'", tab, this.SelectedI, this.Len(), this.OldSelectedI, p.V) // t //
	}
}

func (this *GTableKeys) TraceActive(tab string) {
	if p := this.GetActive(); p != nil {
// trace.Print("%skeys.IActive: %d/%d (%d) => '%s'", tab, this.ActiveI, this.Len(), this.OldActiveI, p.V) // t //
	}
}

func (this *GTableKeys) TraceAll() {

	for i := 0; i < len(this.Sl); i++ {
// trace.Print("%3d (%3d) %4d %s", i, this.Sl[i].LineI, this.Sl[i].DbIndex, this.Sl[i].V) // t //
	}
}

func (this *GTableKeys) GetActiveFiltersCount() uint16 {
	var n uint16
	for i := range this.Sl {
		k := &this.Sl[i]
		if k.IsActiveFilter {
			n++
		}
	}
	return n
}

func (this *GTableKeys) TraceActiveFilters() {
	for i := range this.Sl {
		k := &this.Sl[i]
		if k.IsActiveFilter {
// trace.Print("key %3d '%s' is an active filter", i, k.V) // t //
		}
	}
}

func (this *GTableLines) TraceAll() {
	for i := 0; i < len(this.Sl); i++ {
// trace.Print("%2d (%2d) %s", i, this.Sl[i].IKey, this.Sl[i].Text) // t //
	}
}

//////	t__TableLines

func (this *GTableLines) Reset() {
	this.Sl = nil
	this.OldSelectedI = K.UNSET
}

func (this *GTableLines) Unselect() {
	this.SelectedI = K.UNSET
}

func (this *GTableLines) Make(size int) {
	this.Sl = make([]TTableLine, size)
}

func (this *GTableLines) Append(x *TTableLine) {
	this.Sl = append(this.Sl, *x)
}

func (this *GTableLines) Select(i uint16) bool {
	if i < this.Len() {
		this.OldSelectedI = this.SelectedI
		this.SelectedI = i
		return true
	}
// trace.Error("Lines, select, invalid index: %d", i) // t //
	return false
}

func (this *GTableLines) GetSelected() *TTableLine {
	return this.Get(this.SelectedI)
}

// Returns 'nil' if the index is invalid.
func (this *GTableLines) Get(i uint16) *TTableLine {
	if i >= this.Len() {
/* // b //
		var s string
		if i == K.UNSET {
			s = "unset"
		} else {
			s = "invalid"
		}
*/ // e //
// trace.Error("Lines, index is %s: %d/%d", s, i, this.Len()) // t //
		return nil
	}
	return &this.Sl[i]
}

func (this *GTableLines) Len() uint16 {
	return uint16(len(this.Sl))
}

func (this *GTableLines) Last() *TTableLine {
	return &this.Sl[len(this.Sl)-1]
}

func (this *GTableLines) TraceSelected(tab string) {
	if p := this.GetSelected(); p != nil {
// trace.Print("%slines.ISelected: %d/%d (%d) => '%s'", tab, this.SelectedI, this.Len(), this.OldSelectedI, p.Text) // t //
	}
}

//////////////////////////////////////////////////////////////////////////////h1
//	StringArray

func (this *GStringArray) Get(i uint16) string {
	if i < uint16(len(this.Sl)) {
		return this.Sl[i]
	}
	// var v0 string
	// if this.Len() > 0 {
	// 	v0 = this.sl[0]
	// }
// trace.Print("---------- ERROR with array: index %d is out of range, len is %d", i, this.Len()) // t //
	mpv.Stop()
	return this.Sl[i] // OOPS
	// return PF_ERROR
}

func (this *GStringArray) Append(s string) uint16 {

	this.Sl = append(this.Sl, s)
	return uint16(len(this.Sl))
}

func (this *GStringArray) Reset() {

	this.Sl = nil
}

func (this *GStringArray) Len() uint16 {

	return uint16(len(this.Sl))
}

//////

type t__Flag struct {
	b bool
}

func (this *t__Flag) SetFlag() {

	this.b = true
}

func (this *t__Flag) UnsetFlag() bool {

	b2 := this.b
	this.b = false
	return b2
}

var allFlags []t__Flag

type t__Index struct {
	i    uint16
	iOld uint16
}

func (this *t__Index) SetIndex(i uint16) {

	this.iOld = this.i
	this.i = i
}

func (this *t__Index) UnsetIndex(i uint16) {

	this.iOld = K.UNSET
	this.i = K.UNSET
}

var allIndexes []t__Index

type GTableKeys struct {
	Sl            []TTableKey
	SelectedI     uint16
	ActiveI       uint16
	OldSelectedI  uint16
	OldActiveI    uint16
	ActiveIndexes []uint16
}

type TTableLine struct {
	Text           string
	Normal         string
	Selected       string
	Active         string
	SelectedActive string
	IKey           uint16
}

type GTableLines struct {
	Sl           []TTableLine
	SelectedI    uint16
	OldSelectedI uint16
}

type TTableKey struct {
	V              string
	LineI          uint16
	DbIndex        uint16
	IsActiveFilter bool
}

type GStringArray struct {
	Sl []string
}

// const (
// 	K.UNSET       uint16 = 65535
// 	NEVERSET    uint16 = 65534
// 	OUT_OF_PAGE uint16 = 65535
// 	UNDEF              = "UNDEFINED"
// 	K.UNSET_S            = "K.UNSET"
// 	NONE               = "NONE"
// )

