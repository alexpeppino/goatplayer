package main

import (
	"strings"
)

func (this *mFind) GetSelected() string {
	return this.keys.GetSelected().v
}

func (this *mFind) PrintTitle(s string) {
	top := FindDialog.Section.focusButtons[0].marginTop + 1
	left := FindDialog.Section.focusButtons[0].marginRight + 2
	Term.AddStringWithPos(top, left, styleDialogSelButton_Frame.EchoStyle(s, 0))
	Term.Flush()
}

func (this *mFind) ApplyFilter1() {
	lowerPattern := strings.ToLower(this.pattern)
	this.Reset()
	// trace.Print("-- pattern: %s", lowerPattern)
	for i, item := range this.lowerList {
		if strings.Contains(item, lowerPattern) {
			s := (*(this.list1))[i]
			this.AppendKeyAndLine(s, nil)
			// trace.Print("found item: %s", item)
		}
	}
}

func (this *mFind) ApplyFilter2() {
	lowerPattern := strings.ToLower(this.pattern)
	this.Reset()
	// trace.Print("-- pattern: %s", lowerPattern)
	for i := range *this.list2 {
		item := &(*this.list2)[i]
		if strings.Contains(item.lowerStr, lowerPattern) {
			this.AppendKeyAndLineAndDbIndex(item)
			// trace.Print("found item: %s", item)
		}
	}
}

func (this *mFind) UpdatePattern1() {
	k := Input.inputKey
	le := len(this.pattern)
	if k[0:1] == KEYB_ESC {
		/// Skip escape sequences
// trace.Print("esc") // t //
		return
	}
	if k == KEYB_BACKSPACE {
		if le > 0 {
			this.pattern = this.pattern[:le-1]
		}
	} else {
		this.pattern += k
	}
	// trace.Print("input key: %s", this.pattern)
	if this.pattern == "" {
		this.Reset()
		for _, item := range *this.list1 {
			this.AppendKeyAndLine(item, nil)
		}
	} else {
		this.ApplyFilter1()
	}
	this.keys.iActive = UNSET
	this.keys.iSelected = 0
	this.w.lines.iSelected = 0
	this.w.BringWindowToFront(true)
	s := styleFindPattern.EchoStyle(" "+this.pattern, this.w.nCols)
	Term.SetCursorNormal()
	Term.AddStringWithPos(this.w.iStatusLineRow-1, this.w.marginLeft+1, s).Flush()
	Term.SetCursorPos(this.w.iStatusLineRow-1, this.w.marginLeft+2+uint16(len(this.pattern)))
}

func (this *mFind) UpdatePattern2() {
	k := Input.inputKey
	le := len(this.pattern)
	if k[0:1] == KEYB_ESC {
		/// Skip escape sequences
// trace.Print("esc") // t //
		return
	}
	if k == KEYB_BACKSPACE {
		if le > 0 {
			this.pattern = this.pattern[:le-1]
		}
	} else {
		this.pattern += k
	}
	// trace.Print("input key: %s", this.pattern)
	if this.pattern == "" {
		this.Reset()
		for i := range *this.list2 {
			item := &(*this.list2)[i]
			this.AppendKeyAndLineAndDbIndex(item)
		}
	} else {
		this.ApplyFilter2()
	}
	this.keys.iActive = UNSET
	this.keys.iSelected = 0
	this.w.lines.iSelected = 0
	this.w.BringWindowToFront(true)
	s := styleFindPattern.EchoStyle(" "+this.pattern, this.w.nCols)
	Term.SetCursorNormal()
	Term.AddStringWithPos(this.w.iStatusLineRow-1, this.w.marginLeft+1, s).Flush()
	Term.SetCursorPos(this.w.iStatusLineRow-1, this.w.marginLeft+2+uint16(len(this.pattern)))
}

func (this *mFind) LoadList1(list *[]string, model *mBasic) {
	this.modelName = model.GetName()
	this.mode = 1
	// this.w.pLayer.pSection = this.w.pLayer.pSection
	this.GetScreen().activeWindow = &this.w
	this.Reset()
	this.sSelected = nil
	this.list1 = list
	this.pattern = ""
	this.lowerList = make([]string, len(*list))
	for i, s := range *list {
		this.AppendKeyAndLine(s, nil)
		this.lowerList[i] = strings.ToLower(s)
	}
	this.w.SetEmptyLine()
}

func (this *mFind) LoadList2(list *[]tStringsAndIndexes, model *mBasic) {
	this.modelName = model.GetName()
	this.mode = 2
	// this.w.pLayer.pSection = this.w.pLayer.pSection
	this.GetScreen().activeWindow = &this.w
	this.Reset()
	this.sSelected = nil
	this.list2 = list
	this.pattern = ""
	// this.lowerList = make([]string, len(*list))
	for i := range *list {
		item := &(*list)[i]
		this.AppendKeyAndLineAndDbIndex(item)
	}
	this.w.SetEmptyLine()
}

func (this *mFind) SelectKey_o() {
// trace.Begin_n(this, "SelectKey_o {mFind}") // t //
	Term.StatusLine.SetContent(this.keys.GetSelected().v)
// trace.End() // t //
}

type mFind struct {
	mBasic
	list1     *[]string
	lowerList []string
	list2     *[]tStringsAndIndexes
	// lowerList []string
	pattern   string
	sSelected []string
	modelName string
	mode      uint16
	title     string
}

