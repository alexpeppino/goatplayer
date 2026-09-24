package main

func (this *cTableKeys) Reset() {
	this.sl = nil
	// this.iSelected = 0

	//!
	if this.iOldSelected == 0 {
		this.iSelected = 0
		this.iOldSelected = UNSET
	}
	//!better
	// if this.iSelected == 0 && this.iOldSelected == 0 { // this is true only at the beginning
	// 	this.iSelected = 0
	// 	this.iOldSelected = 1
	// }

	this.iActive = UNSET
	this.iOldActive = UNSET
}

func (this *cTableKeys) Unselect() {
	this.iSelected = UNSET
}

func (this *cTableKeys) Make(size int) {
	this.sl = make([]tTableKey, size)
}

func (this *cTableKeys) Append(x *tTableKey) {
	this.sl = append(this.sl, *x)
}

func (this *cTableKeys) GetLast() *tTableKey {
	return &this.sl[len(this.sl)-1]
}

/* Selects the item #i */
func (this *cTableKeys) Select(i uint16) bool {
	if i < this.Len() {
		this.iOldSelected = this.iSelected
		this.iSelected = i
		return true
	}
// trace.Error("Keys, select, invalid index: %d", i) // t //
	return false
}

func (this *cTableKeys) Activate(i uint16) bool {
	if i < this.Len() {
		this.iOldActive = this.iActive
		this.iActive = i
		return true
	}
// trace.Error("Keys, activate, invalid index: %d", i) // t //
	return false
}

func (this *cTableKeys) GetSelected() *tTableKey {
	return this.Get(this.iSelected)
}

func (this *cTableKeys) GetActive() *tTableKey {
	return this.Get(this.iActive)
}

func (this *cTableKeys) GetIndexByFilename(filename string) uint16 {
	var i uint16
	for i = range uint16(len(this.sl)) {
		key := &this.sl[i]
		if key.v == filename {
			return i
		}
	}
	return UNSET
}

func (this *cTableKeys) GetIndexByDbIndex(dbIndex uint16) uint16 {
	var i uint16
	for i = range uint16(len(this.sl)) {
		key := &this.sl[i]
		if key.dbIndex == dbIndex {
			return i
		}
	}
	return UNSET
}

// Returns 'nil' if index is invalid.
func (this *cTableKeys) Get(i uint16) *tTableKey {
	if i >= this.Len() {
/* // b //
		var s string
		if i == UNSET {
			s = "unset"
		} else {
			s = "invalid"
		}
*/ // e //
// trace.Error("Keys, index is %s: %d/%d", s, i, this.Len()) // t //
		return nil
	}
	return &this.sl[i]
}

func (this *cTableKeys) Len() uint16 {
	return uint16(len(this.sl))
}

func (this *cTableKeys) TraceSelected(tab string) {
	if p := this.GetSelected(); p != nil {
// trace.Print("%sKeys.iSelected: %d/%d (%d) => '%s'", tab, this.iSelected, this.Len(), this.iOldSelected, p.v) // t //
	}
}

func (this *cTableKeys) TraceActive(tab string) {
	if p := this.GetActive(); p != nil {
// trace.Print("%sKeys.iActive: %d/%d (%d) => '%s'", tab, this.iActive, this.Len(), this.iOldActive, p.v) // t //
	}
}

func (this *cTableKeys) TraceAll() {

	for i := 0; i < len(this.sl); i++ {
// trace.Print("%3d (%3d) %4d %s", i, this.sl[i].iLine, this.sl[i].dbIndex, this.sl[i].v) // t //
	}
}

func (this *cTableKeys) GetActiveFiltersCount() uint16 {
	var n uint16
	for i := range this.sl {
		k := &this.sl[i]
		if k.isActiveFilter {
			n++
		}
	}
	return n
}

func (this *cTableKeys) TraceActiveFilters() {
	for i := range this.sl {
		k := &this.sl[i]
		if k.isActiveFilter {
// trace.Print("key %3d '%s' is an active filter", i, k.v) // t //
		}
	}
}

func (this *cTableLines) TraceAll() {
	for i := 0; i < len(this.sl); i++ {
// trace.Print("%2d (%2d) %s", i, this.sl[i].iKey, this.sl[i].text) // t //
	}
}

//////	t__TableLines

func (this *cTableLines) Reset() {
	this.sl = nil
	this.iOldSelected = UNSET
}

func (this *cTableLines) Unselect() {
	this.iSelected = UNSET
}

func (this *cTableLines) Make(size int) {
	this.sl = make([]tTableLine, size)
}

func (this *cTableLines) Append(x *tTableLine) {
	this.sl = append(this.sl, *x)
}

func (this *cTableLines) Select(i uint16) bool {
	if i < this.Len() {
		this.iOldSelected = this.iSelected
		this.iSelected = i
		return true
	}
// trace.Error("Lines, select, invalid index: %d", i) // t //
	return false
}

func (this *cTableLines) GetSelected() *tTableLine {
	return this.Get(this.iSelected)
}

// Returns 'nil' if the index is invalid.
func (this *cTableLines) Get(i uint16) *tTableLine {
	if i >= this.Len() {
/* // b //
		var s string
		if i == UNSET {
			s = "unset"
		} else {
			s = "invalid"
		}
*/ // e //
// trace.Error("Lines, index is %s: %d/%d", s, i, this.Len()) // t //
		return nil
	}
	return &this.sl[i]
}

func (this *cTableLines) Len() uint16 {
	return uint16(len(this.sl))
}

func (this *cTableLines) Last() *tTableLine {
	return &this.sl[len(this.sl)-1]
}

func (this *cTableLines) TraceSelected(tab string) {
	if p := this.GetSelected(); p != nil {
// trace.Print("%sLines.iSelected: %d/%d (%d) => '%s'", tab, this.iSelected, this.Len(), this.iOldSelected, p.text) // t //
	}
}

//////////////////////////////////////////////////////////////////////////////h1
//	StringArray

func (this *cStringArray) Get(i uint16) string {
	if i < uint16(len(this.sl)) {
		return this.sl[i]
	}
	// var v0 string
	// if this.Len() > 0 {
	// 	v0 = this.sl[0]
	// }
// trace.Print("---------- ERROR with array: index %d is out of range, len is %d", i, this.Len()) // t //
	Mpv.Stop()
	return this.sl[i] // OOPS
	// return PF_ERROR
}

func (this *cStringArray) Append(s string) uint16 {

	this.sl = append(this.sl, s)
	return uint16(len(this.sl))
}

func (this *cStringArray) Reset() {

	this.sl = nil
}

func (this *cStringArray) Len() uint16 {

	return uint16(len(this.sl))
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

	this.iOld = UNSET
	this.i = UNSET
}

var allIndexes []t__Index

type cTableKeys struct {
	sl            []tTableKey
	iSelected     uint16
	iActive       uint16
	iOldSelected  uint16
	iOldActive    uint16
	activeIndexes []uint16
}

type tTableLine struct {
	text           string
	normal         string
	selected       string
	active         string
	selectedActive string
	iKey           uint16
}

type cTableLines struct {
	sl           []tTableLine
	iSelected    uint16
	iOldSelected uint16
}

type cStringArray struct {
	sl []string
}

