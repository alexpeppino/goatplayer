package main

import (
	"fmt"
	"goatplayer/internal/K"
	"goatplayer/internal/echo"
	"goatplayer/internal/input"
	"goatplayer/internal/style"
)

func (this *vWindow) Focus() {
	this.pLayer.pSection.Focus(this)
}

func (this *vSection) Focus(that *vWindow) {
	// this.SectionKeyboardMode()
}

func IsInsideRect(re *tRect, iRow uint16, iCol uint16) bool {
	return iRow >= re.marginTop && iRow <= re.marginBottom && iCol >= re.marginLeft && iCol <= re.marginRight
}

func AddEmptyRect(rect tRect, style *style.GStyle) {
	// trace.Begin("PrintEmptyRect") // t //
	var iTRow uint16
	em := style.EchoStyle("", rect.marginRight-rect.marginLeft+1)
	for iTRow = rect.marginTop; iTRow <= rect.marginBottom; iTRow++ {
		// trace2.Print("empty line at trow %d", iTRow)
		echo.AddStringWithPos(iTRow, rect.marginLeft, K.DC_RESET+em)
	}
	// trace.EndAdd("PrintEmptyRect") // t //
}

func SetCaption(name string, focus string) string {
	if f, ok := input.MapKeypad[focus]; ok {
		return fmt.Sprintf("%s·%s", f, name) // ▸
	}
	return fmt.Sprintf("%s·%s", focus, name)
}

func (this *vDialogButton) IsClicked(iRow, iCol uint16) bool {
	// trace.Print("IsClicked")
	return iRow >= this.marginTop && iRow <= this.marginBottom && iCol >= this.marginLeft && iCol <= this.marginRight
}

func (this *vWindow) IsSingleWindow() bool {
	return !this.pLayer.IsMultiTabs()
}

type vBasicView struct {
	name string
	tRect
}

type vButton struct {
	text string
	iRow uint16
	left uint16
	len  uint16
	cmm  func()
}

type vDialogButton struct {
	tRect
	desc string
	i    uint16
}

type vButtonRect struct {
	tRect
	cmm func()
}
