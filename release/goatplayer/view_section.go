package main

import (
	"strings"
)

/* Prints the caption line for this section. */
func (this *vSection) PrintFocusButtons(isActive bool) {
// trace.BeginAdd_n(this, "PrintFocusButtons", "%v", isActive) // t //
	if this.GetActiveWindow() == nil {
// trace.ReturnAdd("activeWindow: nil") // t //
		return
	}
	if isActive {
		Term.AddSimple(this.GetActiveLayer().focusButtonsActiveEcho)
	} else {
		Term.AddSimple(this.GetActiveLayer().focusButtonsEcho)
	}
	Term.AddSimple(this.GetActiveLayer().frameLight)
	// Term.AddSimple(this.closeButton)
	Term.Flush()
// trace.End() // t //
}

func (this *vSection) PrepareFocusButtons(actLayer *vLayer) {
// trace.BeginSilent_n(this, "PrepareFocusButtons") // t //
	var ss, caption string
	var line1a, line1b, seg1a, seg1b string
	var line2a, line2b, seg2a, seg2b string
	var style1, style2 *cStyle
	var pos uint16
	var maBottom, maTop, maLeft, maRight uint16
	for i := range len(this.layers) {
		//	For each layer
		l := &this.layers[i]
		m := l.windows[0].pModel
		if l.IsMultiTabs() {
			if !this.pScreen.isDialog {
				//	First and last windows
				mLast := l.windows[len(l.windows)-1].pModel
				caption = SetCaption(l.caption, m.focusKey+"/"+mLast.focusKey)
			}
		} else {
			//!! if error here, check main.go initAll
			if m.tabCaption != "" {
				ss = m.tabCaption
			} else {
				ss = m.shortName
			}
			caption = SetCaption(ss, m.focusKey)
		}
		captionLen := RuneLen(&caption)
		s4 := strings.Repeat("─", int(captionLen)+2)
		if this.pScreen.isDialog {
			if l == actLayer {
				// s1 = styleSelTab.EchoStyleFormatOnly(caption)
				style1 = &styleDialogSelButton_Frame
				style2 = &styleDialogSelButton_Frame
			} else {
				style1 = &styleDialogButton_Frame
				style2 = &styleDialogButton_Frame
			}
		} else {
			if l == actLayer {
				// s1 = styleSelTab.EchoStyleFormatOnly(caption)
				style1 = &styleSelButton_SelFrame
				style2 = &styleSelButton_Frame
			} else {
				style1 = &styleButton_SelFrame
				style2 = &styleButton_Frame
			}
		}
		seg1a = style1.EchoStyleFormatOnly("┌"+s4+"┐", false)
		seg2a = style2.EchoStyleFormatOnly("┌"+s4+"┐", false)
		vbar1 := style1.EchoStyleFormatOnly("│", false)
		vbar2 := style2.EchoStyleFormatOnly("│", false)
		seg1b = vbar1 + style1.EchoStyleFormatOnly(caption, true) + vbar1
		seg2b = vbar2 + style2.EchoStyleFormatOnly(caption, true) + vbar2
		// s2 += s1
		line1a += seg1a
		line1b += seg1b
		line2a += seg2a
		line2b += seg2b
		if actLayer == &actLayer.pSection.layers[0] { // do this only when the first layer is called (only once)
			/// Prepare the focus-button
			maLeft = this.marginLeft + 2 + pos
			maRight = captionLen + 4 + this.marginLeft + 2 + pos - 1
			maTop = this.marginTop
			maBottom = this.marginTop + 1
			this.focusButtons = append(this.focusButtons, vFocusButton{tRect{maTop, maBottom, maLeft, maRight}, l, 0})
// trace.Print("focus button coordinates for layer %s: %d %d", l.name, maLeft, maRight) // t //
			pos += captionLen + 4
		}
	}
	// s2 = fmt.Sprintf("%s%s  %6d %6d %6d", s2, Pen.Normal, this.selectedTable.selected_i+1, this.selectedTable.filenames_n, this.activeTable.active_i+1)
	// this.caption = s2
	// Term.AddStringWithPos(this.marginTop, this.marginLeft+2, s2)

	/// When section is active
	Term.AddStringWithPos(this.marginTop, this.marginLeft+1, line1a)
	Term.AddStringWithPos(this.marginTop+1, this.marginLeft+1, line1b)
	if actLayer.IsMultiTabs() {
		s1 := styleSelButton_SelFrame.EchoStyle(" ", this.marginRight-this.marginLeft)
		Term.AddStringWithPos(this.marginTop+2, this.marginLeft+1, s1)
	}
	if this.WspTitle != "" {
		var x = styleButton_SelFrame.EchoStyleFormatOnly(this.WspTitle, true)
		Term.AddStringWithPos(this.marginTop, this.marginRight-RuneLen(&this.WspTitle)-1, x)
	}
	Term.SlurpContents(&actLayer.focusButtonsActiveEcho)

	/// When section is not active
	Term.AddStringWithPos(this.marginTop, this.marginLeft+1, line2a)
	Term.AddStringWithPos(this.marginTop+1, this.marginLeft+1, line2b)
	if actLayer.IsMultiTabs() {
		s1 := styleButton_Frame.EchoStyle(" ", this.marginRight-this.marginLeft)
		Term.AddStringWithPos(this.marginTop+2, this.marginLeft+1, s1)
	}
	if this.WspTitle != "" {
		var x = styleButton_Frame.EchoStyleFormatOnly(this.WspTitle, true)
		Term.AddStringWithPos(this.marginTop, this.marginRight-RuneLen(&this.WspTitle)-1, x)
	}
	Term.SlurpContents(&actLayer.focusButtonsEcho)

	// trace.Print("%q", actLayer.sectionFocusButtonsActive)
	// trace.Print("%q", actLayer.sectionFocusButtons)
// trace.End() // t //
}

func (this *vSection) PrepareSection() {
	this.tobePrepared = true
	this.focusButtons = nil
	this.PrepareFramesAndButtons()
	// this.pScreen.dialogButtons = nil
	// slurps all caption echoes
}

func (this *vSection) PrepareFramesAndButtons() {
	if this.pScreen.isDialog {
// trace.Print("isDialog") // t //
		Term.MakeSectionFrame(&this.tRect, &styleDialogButton_Frame, &Frames.Empty).SlurpContents(&this.frameLight)
		Term.MakeSectionFrame(&this.tRect, &styleDialogButton_Frame, &Frames.Empty).SlurpContents(&this.frameBold)
	} else {
// trace.Print("is not Dialog") // t //
		Term.MakeSectionFrame(&this.tRect, &styleButton_Frame, &Frames.Empty).SlurpContents(&this.frameLight)
		Term.MakeSectionFrame(&this.tRect, &styleButton_SelFrame, &Frames.Empty).SlurpContents(&this.frameBold)
	}
	for i := range this.layers {
		l := &this.layers[i]
		this.PrepareFocusButtons(l)
	}
}

func (this *vSection) IsThisSection() bool {
	iRow := Input.iRow
	iCol := Input.iCol
	return IsInsideRect(&this.tRect, iRow, iCol)
}

func (this *vSection) InitSection(nLayers uint16, rect tRect, name string, screen *vScreen) {
	this.name = name
// trace.BeginSilent_n(this, "InitSection") // t //
	this.pScreen = screen
	screen.allSections = append(screen.allSections, this)
	AllMaps.allSections = append(AllMaps.allSections, this)

	if nLayers > 0 {
		this.layers = make([]vLayer, nLayers)
	}
	for i := range nLayers {
		l := &this.layers[i]
		l.pSection = this
		screen.allLayers = append(screen.allLayers, l)
	}
	// this.focusToModel.mapKeyToModel = make(map[string]*mBasic)
	// this.focusToModel.mapKeyToName = make(map[string]string)
	// this.focusToModel.mapNameToModel = make(map[string]*mBasic)

	this.tRect = rect

	// if this.pScreen.isDialog {
	// 	trace.Print("isDialog")
	// 	Term.PrepareSectionFrame(&this.tRect, &styleDialogButton_Frame, &Frames.Empty).SlurpContents(&this.frameLight)
	// 	Term.PrepareSectionFrame(&this.tRect, &styleDialogButton_Frame, &Frames.Empty).SlurpContents(&this.frameBold)
	// } else {
	// 	trace.Print("is not Dialog")
	// 	Term.PrepareSectionFrame(&this.tRect, &styleButton_Frame, &Frames.Empty).SlurpContents(&this.frameLight)
	// 	Term.PrepareSectionFrame(&this.tRect, &styleButton_SelFrame, &Frames.Empty).SlurpContents(&this.frameBold)
	// }
// trace.End() // t //
}

// func (this *vSection) AddLayer(name string) *vLayer {
// 	trace.BeginEndAdd_n(this, "AddLayer", name)
// 	var layer vLayer
// 	var l *vLayer
// 	this.layers = append(this.layers, layer)
// 	l = &this.layers[len(this.layers)-1]
// 	l.pSection = this
// 	l.name = name
// 	this.pScreen.allLayers = append(this.pScreen.allLayers, l)
// 	trace.Print("new layer address: %p", l)
// 	return l
// }

func (this *vSection) DumpSection() {
// trace.Begin_n(this, "DumpScreen {vSection}") // t //

	for i := range len(this.layers) {
		this.layers[i].DumpLayer()
	}
// trace.End() // t //
}

type vSection struct {
	vBasicView
	WspTitle    string
	owner       *vLayer
	pScreen     *vScreen
	layers      []vLayer
	l1          vLayer
	l2          vLayer
	frameLight  string
	frameBold   string
	activeLayer *vLayer
	// activeWindow        *vWindow
	caption      string
	focusButtons []vFocusButton /* focus-to-layer buttons */
	// focusToModel        uSectionFocusToModel
	tobePrepared bool // when reprinting, call PrepareView for this section
}

