package main

import (
	"goatplayer/internal/echo"
	"goatplayer/internal/input"
	"goatplayer/internal/small/utils"
	"goatplayer/internal/style"
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
		echo.AddSimple(this.GetActiveLayer().focusButtonsActiveEcho)
	} else {
		echo.AddSimple(this.GetActiveLayer().focusButtonsEcho)
	}
	echo.AddSimple(this.GetActiveLayer().frameLight)
	// echo.AddSimple(this.closeButton)
	echo.Flush()
	// trace.End() // t //
}

func (this *vSection) PrepareFocusButtons(actLayer *vLayer) {
	// trace.BeginSilent_n(this, "PrepareFocusButtons") // t //
	var ss, caption string
	var line1a, line1b, seg1a, seg1b string
	var line2a, line2b, seg2a, seg2b string
	var style1, style2 *style.GStyle
	var pos uint16
	var maBottom, maTop, maLeft, maRight uint16
	for i := range len(this.layers) {
		//	For each layer
		l := &this.layers[i]
		model := l.windows[0].pModel
		if l.IsMultiTabs() {
			if !this.pScreen.isDialog {
				//	First and last windows
				mLast := l.windows[len(l.windows)-1].pModel
				caption = SetCaption(l.caption, model.focusKey+"/"+mLast.focusKey)
			}
		} else {
			//!! if error here, check main.go initAll
			if model.tabCaption != "" {
				ss = model.tabCaption
			} else {
				ss = model.shortName
			}
			caption = SetCaption(ss, model.focusKey)
		}
		captionLen := utils.RuneLen(&caption)
		s4 := strings.Repeat("─", int(captionLen)+2)
		if this.pScreen.isDialog {
			if l == actLayer {
				// s1 = styleSelTab.EchoStyleFormatOnly(caption)
				style1 = &style.DialogSelButton_Frame
				style2 = &style.DialogSelButton_Frame
			} else {
				style1 = &style.DialogButton_Frame
				style2 = &style.DialogButton_Frame
			}
		} else {
			if l == actLayer {
				// s1 = styleSelTab.EchoStyleFormatOnly(caption)
				style1 = &style.SelButton_SelFrame
				style2 = &style.SelButton_Frame
			} else {
				style1 = &style.Button_SelFrame
				style2 = &style.Button_Frame
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
	// s2 = fmt.Sprintf("%s%s  %6d %6d %6d", s2, Pen.Normal, this.SelectedTable.Selected_i+1, this.SelectedTable.filenames_n, this.activeTable.active_i+1)
	// this.caption = s2
	// echo.AddStringWithPos(this.marginTop, this.marginLeft+2, s2)

	/// When section is active
	echo.AddStringWithPos(this.marginTop, this.marginLeft+1, line1a)
	echo.AddStringWithPos(this.marginTop+1, this.marginLeft+1, line1b)
	if actLayer.IsMultiTabs() {
		s1 := style.SelButton_SelFrame.EchoStyle(" ", this.marginRight-this.marginLeft)
		echo.AddStringWithPos(this.marginTop+2, this.marginLeft+1, s1)
	}
	if this.WspTitle != "" {
		var x = style.Button_SelFrame.EchoStyleFormatOnly(this.WspTitle, true)
		echo.AddStringWithPos(this.marginTop, this.marginRight-utils.RuneLen(&this.WspTitle)-1, x)
	}
	echo.SlurpContents(&actLayer.focusButtonsActiveEcho)

	/// When section is not active
	echo.AddStringWithPos(this.marginTop, this.marginLeft+1, line2a)
	echo.AddStringWithPos(this.marginTop+1, this.marginLeft+1, line2b)
	if actLayer.IsMultiTabs() {
		s1 := style.Button_Frame.EchoStyle(" ", this.marginRight-this.marginLeft)
		echo.AddStringWithPos(this.marginTop+2, this.marginLeft+1, s1)
	}
	if this.WspTitle != "" {
		var x = style.Button_Frame.EchoStyleFormatOnly(this.WspTitle, true)
		echo.AddStringWithPos(this.marginTop, this.marginRight-utils.RuneLen(&this.WspTitle)-1, x)
	}
	echo.SlurpContents(&actLayer.focusButtonsEcho)

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
		Term.MakeSectionFrame(&this.tRect, &style.DialogButton_Frame, &Frames.Empty)
		echo.SlurpContents(&this.frameLight)
		Term.MakeSectionFrame(&this.tRect, &style.DialogButton_Frame, &Frames.Empty)
		echo.SlurpContents(&this.frameBold)
	} else {
		// trace.Print("is not Dialog") // t //
		Term.MakeSectionFrame(&this.tRect, &style.Button_Frame, &Frames.Empty)
		echo.SlurpContents(&this.frameLight)
		Term.MakeSectionFrame(&this.tRect, &style.Button_SelFrame, &Frames.Empty)
		echo.SlurpContents(&this.frameBold)
	}
	for i := range this.layers {
		l := &this.layers[i]
		this.PrepareFocusButtons(l)
	}
}

func (this *vSection) IsThisSection() bool {
	iRow := input.RowI
	iCol := input.ColI
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
