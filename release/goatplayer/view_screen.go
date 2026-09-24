package main

import (
	"fmt"
)

func (this *vAbsDialogScreen) GetFocusMan_o() *uFocusKeys {
// trace.BeginEndAdd_n(this, "GetFocusMan_o {vDialog}", "%s", this.focusKeys.name) // t //
	return &this.focusKeys
}

func (this *vScreen) UpdateFocusButtonList() {
// trace.Begin_n(this, "UpdateFocusButtonList") // t //
	//	Focus buttons from all active sections
	this.allFocusButtons = nil
	// trace.Print("this.allFocusableSections, %v", LongNames(this.allFocusableSections))
	for _, s := range this.allFocusableSections {
		// trace.Print("focusable section %s", s.name)
		// trace.Print("focusable section buttons %v", s.focusButtons)
		for i := range s.focusButtons {
			this.allFocusButtons = append(this.allFocusButtons, &s.focusButtons[i])
		}
	}
// trace.End() // t //
}

func (this *vAbsDialogScreen) PrepareDialogButtons() {
// trace.BeginSilent_n(this, "PrepareDialogButtons") // t //
	style1 := &styleDialogButton_Frame
	// style2 := &styleFocusButton_Frame
	se := &this.Section
	var line1a, line1b string
	var maBottom, maTop, maLeft, maRight uint16
	maTop = this.Section.marginTop
	maBottom = this.Section.marginTop + 1
	vbar := "│"
	///	Draw ok-button
	// line1a = style1.EchoStyleFormatOnly("┌──┐", false)
	// line1b = style1.EchoStyleFormatOnly(vbar+"✅"+vbar, false)
	// maLeft = se.marginRight - 10
	// maRight = maLeft + 5
	// this.dialogButtons = append(this.dialogButtons, &vDialogButton{tRect{maTop, maBottom, maLeft, maRight}, "ok button", 0})
	// trace.Print("ok-button coordinates: %d %d", maLeft, maRight)
	// Term.AddStringWithPos(se.marginTop, maLeft+1, line1a)
	// Term.AddStringWithPos(se.marginTop+1, maLeft+1, line1b)
	///	Draw close-button
	line1a = style1.EchoStyleFormatOnly("┌──┐", false)
	line1b = style1.EchoStyleFormatOnly(vbar+"❌"+vbar, false)
	maLeft = se.marginRight - 6
	maRight = maLeft + 5
	this.dialogButtons = append(this.dialogButtons, &vDialogButton{tRect{maTop, maBottom, maLeft, maRight}, "close button", 1})
// trace.Print("close-button coordinates: %d %d", maLeft, maRight) // t //
	Term.AddStringWithPos(se.marginTop, maLeft+2, line1a)
	Term.AddStringWithPos(se.marginTop+1, maLeft+2, line1b)
	Term.SlurpContents(&this.closeButton)
	///	Set drag button
	maRight = maLeft - 1
	maLeft = this.Section.focusButtons[len(this.Section.focusButtons)-1].marginRight // margin right of the last focus button
	this.dialogButtons = append(this.dialogButtons, &vDialogButton{tRect{maTop, maBottom, maLeft, maRight}, "drag button", 2})
// trace.End() // t //
}

func (this *vScreen) IsFocusToModelKey(key string) bool {
// trace.BeginAdd_n(this, "IsFocusToModelKey", "key: %s", Input.GetKeyName(key)) // t //
	focus := this.outer.GetFocusMan_o()
_, ok := focus.mapKeyToModel[key]

// f, ok := focus.mapKeyToModel[key] /*c*/
	if ok {
// trace.ReturnAdd("key found for model %s", f.name) // t //
		return ok
	}
// trace.ReturnAdd("key not found") // t //
	return false
}

func (this *vMainScreen) GetFocusMan_o() *uFocusKeys {
// trace.Begin_n(this, "GetFocusMan_o {vMainScreen}") // t //
	if this.activeWorkspace == nil {
// trace.Print("activeWorkspace is unset") // t //
		this.activeWorkspace = &BrowserWsp1.maWorkspace
	}
// trace.ReturnAdd("activeWorkspace: %s, focus: %s", this.activeWorkspace.name, this.activeWorkspace.focusKeys.name) // t //
	return &this.activeWorkspace.focusKeys
}

func (this *vScreen) SeekCommandButtons() (cmmFound bool) {
// trace.Begin_n(this, "SeekCommandButtons") // t //
	if Input.inputKey != MOUSE_B1_DOWN {
		return false
	}
	activeWindow := this.GetActiveWindow()
	if activeWindow != nil {
		// trace.Print("activeWindow / model: %s / %s", activeWindow.name, activeWindow.pModel.name)
		model := this.GetActiveWindow().pModel
		///	Seek commands
		var first = model.iCmmPage * CommandPanel.nLines
		var last = (model.iCmmPage + 1) * CommandPanel.nLines
		// trace.Print("model.iCmmPage: %d,  CommandPanel.nLines: %d", model.iCmmPage, CommandPanel.nLines)
		for i := first; i < last; i++ { //:= range this.activeWindow.pModel.kCmm {
			if i < uint16(len(model.kCmm)) {
				cmm := &model.kCmm[i]
				// trace.Print("checking command %s...", cmm.desc)
				if IsInsideRect(cmm.rect, Input.iRow, Input.iCol) {
					// cmmFound = true
// trace.PrintFocus("is inside, command: %s", cmm.desc) // t //
					Term.activeCommand = cmm
					PrintPressedCommand(cmm)
					/// Execute command
					cmm.f()
// trace.ReturnAdd("command executed, return true") // t //
					return true
				}
			}
		}
	} else {
		// trace.Print("there is no activeWindow")
	}
// trace.ReturnAdd("no command found, return false") // t //
	return false
}

func (this *vScreen) SetActiveWindow(model *mBasic) {
	this.activeWindow = &model.w
	this.activeSection = this.activeWindow.pLayer.pSection
	this.activeSection.activeLayer = this.activeWindow.pLayer
	this.activeSection.activeLayer.activeWindow = this.activeWindow
}

func (this *vScreen) SeekFocusButtonsAndWindow() (model *mBasic) {
// trace.Begin_n(this, "SeekFocusButtonsAndWindow") // t //
	///	Seek focus buttons
	for i := range this.allFocusButtons {
		focusButton := this.allFocusButtons[i]
		// trace.Print("%+v", bu)
		if focusButton.IsClicked(Input.iRow, Input.iCol) {
			// if focusButton.layer == nil { // is the close-button
			// 	trace.PrintFocus("close-button clicked")
			// 	trace.N_Return(this, "MouseMode, close-button clicked")
			// 	return "ESC", false
			// }
			model = focusButton.layer.windows[0].pModel
			this.SetActiveWindow(model)
			// this.activeWindow = &model.w
// trace.PrintFocus("clicked on focus button for model: %s, layer: %s", model.name, focusButton.layer.name) // t //
			break
		}
	}
	if model == nil {
		w := this.WhichWindow()
		if w != nil {
// trace.PrintFocus("clicked on window: %s", w.name) // t //
			model = w.pModel
			this.SetActiveWindow(model)
		}
	}
	if model != nil {
		model.SetFocusToModel()
		this.SetActiveWindow(model)
		// this.activeWindow = &model.w
		// this.activeSection = this.activeWindow.pLayer.pSection
		// this.activeSection.activeLayer = this.activeWindow.pLayer
		// this.activeSection.activeLayer.activeWindow = this.activeWindow
		model.outer.CheckMouseInput_o()
	}
// trace.Return() // t //
	return model
}

func (this *vAbsDialogScreen) SeekDialogButtons() (cmmFound bool, sReturn string) {
// trace.Begin_n(this, "SeekDialogButtons") // t //
	for i := range this.dialogButtons {
		dialogButton := this.dialogButtons[i]
		if dialogButton.IsClicked(Input.iRow, Input.iCol) {
			cmmFound = true
// trace.PrintFocus("dialog button clicked: %s %d", dialogButton.desc, dialogButton.i) // t //
			switch dialogButton.i {
			case 0:
				sReturn = "ESC"
// trace.ReturnAdd("ESC") // t //
				return
			case 1:
				sReturn = "ESC"
// trace.ReturnAdd("ESC") // t //
				return
			case 2:
				// this.DragDialog()
				// // this.TraceActives()
				// trace.ReturnAdd("DragDialog")
				return
			}
			break
		}
	}
// trace.Return() // t //
	return
}

/*  */
func (this *vScreen) MouseMode_o() (string, bool) {
// trace.BeginAdd_n(this, "MouseMode_o {vScreen}", "vBasicScreen") // t //
	for {
// trace.PrintWithName("for...") // t //

		/// Click on Progress bar
		// trace.Print("row/col: %d/%d, ProgressBar top/left/right: %d %d %d", Input.iRow, Input.iCol, Term.ProgressBar_top, Term.ProgressBar_left, Term.ProgressBar_right)
		if Input.iRow == Term.ProgressBar_top && Input.iCol >= Term.ProgressBar_left && Input.iCol <= Term.ProgressBar_right {
			/// Click on Music Bar --> change audiof position
			var pos float32
			var perc uint16
			perc = 1 + 100*(Input.iCol-Term.ProgressBar_left)/(Term.ProgressBar_right-Term.ProgressBar_left)
			pos = float32(pl_ft.duration) * float32(perc) / 10000.0
// trace.Print("perc: %d, pos: %f, duration: %d", perc, pos, pl_ft.duration) // t //
			Mpv.SeekAbsolute(pos)
			goto readinput
		}

// trace.Print("check all Term buttons...........") // t //
		for i := range Term.allButtons {
			but := Term.allButtons[i]
			// trace.Print("-- button: %s", but.text)
			if but.Clicked(Input.iRow, Input.iCol) {
// trace.Print("clicked term button: %s", but.text) // t //
				but.cmm()
				Term.pPressedButton = but
				goto readinput
			}
		}

		///
		if !this.SeekCommandButtons() {
			this.SeekFocusButtonsAndWindow()
		}

	readinput:
		Input.ReadInput()
		if Input.inputType == INPUT_KEYBOARD {
			// Focus.activeKey = Input.inputKey
// trace.ReturnAdd("INPUT_KEYBOARD, active model: %s (its focus key is: %s)", this.activeWindow.pModel.name, this.activeWindow.pModel.focusKey) // t //
			return this.activeWindow.pModel.focusKey, false // false means: don't read the input at the beginning of the checkinput, since there's already one read to parse.
		}
	}
}

func (this *vAbsDialogScreen) MouseMode_o() (string, bool) {
// trace.BeginAdd_n(this, "MouseMode {vDialog}", "vBasicDialogScreen") // t //
	var cmmFound bool
	var sReturn string
	for {
		cmmFound = this.SeekCommandButtons()
		if !cmmFound {
			cmmFound, sReturn = this.SeekDialogButtons()
// trace.Print("SeekDialogButtons returned %v, '%s'", cmmFound, sReturn) // t //
			if sReturn != "" {
				return sReturn, false
			}
		}
		if !cmmFound {
			this.SeekFocusButtonsAndWindow()
		}
		Input.ReadInput()
		if Input.inputType == INPUT_KEYBOARD {
			activeWindow := this.GetActiveWindow()
			// Focus.activeKey = Input.inputKey
// trace.ReturnAdd("active model %s, focus key %s", activeWindow.name, activeWindow.pModel.focusKey) // t //
			return activeWindow.pModel.focusKey, false // false means: don't read the input at the beginning of the checkinput, since there's already one read to parse.
		}
	}
}

func (this *vAbsDialogScreen) KeyboardMode() bool {
// trace.Begin_n(this, "KeyboardMode") // t //
	var readInputAtBeginning bool = true
	focusToModel := this.focusKeys
	for {
		// Input.ReadInput()
// trace.PrintFocus("switch, key: %s", focusToModel.GetActiveKeyName()) // t //
		// this.manFocusToModel.TraceActiveKey()
		Input.DisplayInput()
		switch focusToModel.activeKey {
		case "ESC":
// trace.Print("case: ESC") // t //
			/// dialog canceled
// trace.ReturnAdd("ESC") // t //
			Term.SetCursorInvisible()
			return false
		case "ENTER":
// trace.Print("case: ENTER") // t //
			/// dialog ok
			Term.SetCursorInvisible()
// trace.ReturnAdd("ENTER") // t //
			return true
		case MOUSE_MODE:
// trace.Print("case: MOUSE_MODE") // t //
// trace.PrintFocus("KeyboardMode, call this.MouseMode") // t //
			focusToModel.activeKey, readInputAtBeginning = this.outer.MouseMode_o()
// trace.PrintFocus("KeyboardMode, this.MouseMode returned, activeKey %s, readInputAtBeginning: %v", focusToModel.activeKey, readInputAtBeginning) // t //
			// b = false
		default:
			model := focusToModel.GetModel()
			if model == nil {
				// Quit()
			}
// trace.Print("default case: set focus for model %s", model.name) // t //
			model.SetFocusToModel()
// trace.PrintFocus("KeyboardMode, call CheckInput for model %s", model.name) // t //
			focusToModel.activeKey = model.outer.CheckInput_o(readInputAtBeginning)
// trace.PrintFocus("KeyboardMode, CheckInput returned for model %s", model.name) // t //
			readInputAtBeginning = true
		}
	}
}

// .
func (this *vMainScreen) KeyboardMode() {
// trace.Begin_n(this, "KeyboardMode") // t //
	var readInputAtBeginning bool = true
	var focusToModel *uFocusKeys = &this.activeWorkspace.focusKeys
	for {
		// Input.ReadInput()
// trace.PrintFocus("KeyboardMode switch, key: %s", focusToModel.GetActiveKeyName()) // t //
// trace.Print("this.activeWorkspace: %s", this.activeWorkspace.name) // t //
		// this.manFocusToModel.TraceActiveKey()
		Input.DisplayInput()
		switch focusToModel.activeKey {
		case MOUSE_MODE:
// trace.Print("case: MOUSE_MODE") // t //
// trace.PrintFocus("KeyboardMode, call this.MouseMode") // t //
			focusToModel.activeKey, readInputAtBeginning = this.outer.MouseMode_o()
// trace.PrintFocus("KeyboardMode, this.MouseMode returned, key: %s, readInputAtBeginning: %v", focusToModel.activeKey, readInputAtBeginning) // t //
			// b = false
		default:
// trace.Print("start default case") // t //
			// trace.Print("focusToModel: %s", focusToModel.name)
			// trace.Print("focusToModel.GetActiveKeyName: %s", focusToModel.GetActiveKeyName())

			/// change of workspace, focuskeys 1-4
// trace.Print("change of workspace, focusToModel.activeKey: %s", focusToModel.activeKey) // t //
			switch focusToModel.activeKey {
			case View1.focusKey:
// trace.Print("change of workspace, BrowserWsp1") // t //
				BrowserWsp1.focusKeys.activeKey = View1.focusKey
				this.activeWorkspace = &BrowserWsp1.maWorkspace
			case View2.focusKey:
// trace.Print("change of workspace, BrowserWsp2") // t //
				BrowserWsp2.focusKeys.activeKey = View2.focusKey
				this.activeWorkspace = &BrowserWsp2.maWorkspace
			case Play.focusKey:
// trace.Print("change of workspace, PlayerWsp") // t //
				PlayerWsp.focusKeys.activeKey = Play.focusKey
				this.activeWorkspace = &PlayerWsp.maWorkspace
			// case Edit.focusKey:
			// 	trace.Print("change of workspace, EditorWsp")
			// 	EditorWsp.focusKeys.activeKey = Edit.focusKey
			// 	this.activeWorkspace = &EditorWsp.maWorkspace
			case SetM.focusKey:
// trace.Print("change of workspace, SetM") // t //
				SettingsDialog.focusKeys.activeKey = SetM.focusKey
			}

// trace.Print("after change of workspace, this.activeWorkspace: %s", this.activeWorkspace.name) // t //
			focusToModel = &this.activeWorkspace.focusKeys
// trace.Print("focusToModel: %s", focusToModel.name) // t //
// trace.Print("focusToModel.GetActiveKeyName: %s", focusToModel.GetActiveKeyName()) // t //
			m := focusToModel.GetActiveModel()
			if m == nil {
// trace.Print("model is nil, get activeWindow model") // t //
				m = this.activeWindow.pModel
			}
// trace.PrintWithName("default case: set focus for model %s", m.name) // t //
			// trace.Print("focusToModel: %s", focusToModel.name)
			// trace.Print("focusToModel.GetActiveKeyName: %s", focusToModel.GetActiveKeyName())
			m.SetFocusToModel()
// trace.PrintFocus("KeyboardMode, call CheckInput for model %s", m.name) // t //
			focusToModel.activeKey = m.outer.CheckInput_o(readInputAtBeginning)
// trace.PrintFocus("KeyboardMode, CheckInput returned key %s for model %s", focusToModel.activeKey, m.name) // t //
			// trace.Print("focusToModel: %s", focusToModel.name)
			// trace.Print("focusToModel.GetActiveKeyName: %s", focusToModel.GetActiveKeyName())
			readInputAtBeginning = true
		}
	}
}

func (this *vMainScreen) InitScreen() {
	this.SetOuter()
	this.mapWorkspaces = make(map[string]*maWorkspace)
	AllMaps.allScreens = append(AllMaps.allScreens, &this.vScreen)
}

// .
func (this *vMainScreen) PrepareScreen() {
	this.name = "MainScreen"
// trace.Begin_n(this, "PrepareScreen") // t //
	// this.manFocusToModel = &Browser1.manFocusToModel
	var width uint16 = ((Term.nCols - 90 - 25) / 3) - 1
	var height = Term.nRows - 13
	var w1 uint16 = 90
	var se *vSection
	this.content = ""
	// this.begin = TPUT_CIVIS
	// this.end = TPUT_CUP_INIT
	// contents: 	"",
	// begin:		tputScCivis,
	// end:		tputRcCnorm + "\x1b[53;41H" }	// tput rc; tput cnorm; tput cup 52 40

	/// LeftSection
	se = &this.LeftSection
	se.InitSection(3, tRect{0, height, 0, w1}, "LeftSection", &this.vScreen)
	se.layers[0].name = "LeftSection.View1"
	se.layers[0].AddWindow(&View1.w, nil)
	se.layers[1].name = "LeftSection.View2"
	se.layers[1].AddWindow(&View2.w, nil)
	se.layers[2].name = "LeftSection.Play"
	se.layers[2].AddWindow(&Play.w, nil)
	// se.layers[3].name = "LeftSection.Edit"
	// se.layers[3].AddWindow(&Edit.w, nil)
	// for i := range se.layers {
	// 	l := &se.layers[i]
	// 	trace.Print("%s - layer address:  (%d)   %p", se.name, i, l)
	// }
	se.PrepareSection()

	/// Browser1Section
	bro := &BrowserWsp1
	se = &this.Browser1Section
	se.InitSection(4, tRect{0, height, w1, w1 + 3*width}, "Browser1Section", &this.vScreen)
	se.WspTitle = BrowserWsp1.publicName
	se.layers[0].name = "Browser1Section.Filters1"
	se.layers[0].caption = "Fields"
	se.layers[0].AddWindow(&bro.GenresCon.w, &tRect{0, height, w1, w1 + width})
	se.layers[0].AddWindow(&bro.ArtistsCon.w, &tRect{0, height, w1 + width, w1 + 2*width})
	se.layers[0].AddWindow(&bro.AlbumsCon.w, &tRect{0, height, w1 + 2*width, w1 + 3*width})
	se.layers[1].name = "Browser1Section.Filters2"
	se.layers[1].caption = "Fields"
	se.layers[1].AddWindow(&bro.YearsCon.w, &tRect{0, height, w1, w1 + width})
	se.layers[1].AddWindow(&bro.AlbumArtistsCon.w, &tRect{0, height, w1 + width, w1 + 2*width})
	se.layers[1].AddWindow(&bro.ComposersCon.w, &tRect{0, height, w1 + 2*width, w1 + 3*width})
	se.layers[2].name = "Browser1Section.FileProp"
	se.layers[2].AddWindow(&bro.AudioPropCon.w, &tRect{0, height, w1, w1 + 3*width})
	se.layers[3].name = "Browser1Section.Dirs"
	se.layers[3].AddWindow(&bro.Dirs.w, &tRect{0, height, w1, w1 + 3*width})
	// se.layers[4].name = "Browser1Section.PFLists"
	// se.layers[4].AddWindow(&PFLists.w, nil)
	// se.layers[5].name = "Browser1Section.Comments"
	// se.layers[5].AddWindow(&bro.CommentsCon.w, &tRect{0, height, 90, 210})
	se.PrepareSection()

	/// Browser2Section
	bro = &BrowserWsp2
	se = &this.Browser2Section
	se.InitSection(4, tRect{0, height, w1, w1 + 3*width}, "Browser2Section", &this.vScreen)
	se.WspTitle = BrowserWsp2.publicName
	se.layers[0].name = "Browser2Section.Filters1"
	se.layers[0].caption = "Fields"
	se.layers[0].AddWindow(&bro.GenresCon.w, &tRect{0, height, w1, w1 + width})
	se.layers[0].AddWindow(&bro.ArtistsCon.w, &tRect{0, height, w1 + width, w1 + 2*width})
	se.layers[0].AddWindow(&bro.AlbumsCon.w, &tRect{0, height, w1 + 2*width, w1 + 3*width})
	se.layers[1].name = "Browser2Section.Filters2"
	se.layers[1].caption = "Fields"
	se.layers[1].AddWindow(&bro.YearsCon.w, &tRect{0, height, w1, w1 + width})
	se.layers[1].AddWindow(&bro.AlbumArtistsCon.w, &tRect{0, height, w1 + width, w1 + 2*width})
	se.layers[1].AddWindow(&bro.ComposersCon.w, &tRect{0, height, w1 + 2*width, w1 + 3*width})
	se.layers[2].name = "Browser2Section.FileProp"
	se.layers[2].AddWindow(&bro.AudioPropCon.w, &tRect{0, height, w1, w1 + 3*width})
	se.layers[3].name = "Browser2Section.Dirs"
	se.layers[3].AddWindow(&bro.Dirs.w, &tRect{0, height, w1, w1 + 3*width})
	// se.layers[4].name = "Browser2Section.PFLists"
	// se.layers[4].AddWindow(&PFLists.w, nil)
	// se.layers[5].name = "Browser2Section.Comments"
	// se.layers[5].AddWindow(&bro.CommentsCon.w, &tRect{0, height, 90, 210})
	se.PrepareSection()

	/// PlayerSection
	se = &this.PlayerSection
	se.InitSection(1, tRect{0, height, w1, w1 + 3*width}, "PlayerSection", &this.vScreen)
	se.WspTitle = "Player"
	se.layers[0].name = "PlayerSection.FileProp"
	se.layers[0].AddWindow(&PlayerWsp.AudioPropCon.w, nil)
	// se.layers[1].name = "PlayerSection.PFLists"
	// se.layers[1].AddWindow(&PlayerWsp.PFListsCon.w, nil)
	// se.layers[2].name = "PlayerSection.Param"
	// se.layers[2].AddWindow(&PlayerWsp.Param.w, nil)
	se.PrepareSection()

	this.CommandPanelSection.InitSection(1, tRect{0, height - 2, Term.nCols - 25, Term.nCols}, "CommandPanelSection", &this.vScreen)

	se = &this.DisplaySection
	se.InitSection(1, tRect{height + 2, height + 8, 0, this.Browser1Section.marginRight}, "DisplaySection", &this.vScreen)
	se.layers[0].name = "MusicPlayer.L0"
	// this.MusicPlayer.layers[0].AddWindow(&mus.w, &tRect{0, 52, w1, w1 + width})
	Display.Prepare(se.tRect)

	this.SectionCommandLine.InitSection(0, tRect{54, 54, 150, 200}, "CommandLine", &this.vScreen)

	for _, l := range this.allLayers {
		l.InitLayer()
	}
	for _, w := range this.allWindows {
		w.PrepareWindow()
	}
	this.pRightSection = &this.Browser1Section
	this.UpdateFocusButtonList()
	this.activeSection = &this.LeftSection
// trace.End() // t //
}

func (this *vFindDialog) InitScreen() {
	this.SetOuter()
}

var dialogLeft uint16 = 90
var dialogTop uint16 = 5

func (this *vFindDialog) PrepareScreen() {
	AllMaps.allScreens = append(AllMaps.allScreens, &this.vScreen)
	var top uint16 = 7
	var left uint16 = dialogLeft
	var bottom uint16 = top + 35
	var se *vSection
	// var middle uint16 = left + 30
	var right uint16 = left + 120
	this.name = "FindDialog"
// trace.Begin_n(this, "InitScreen") // t //
	// this.manFocusToModel = new(uManScreenFocusToModel)
	this.focusKeys.pScreen = &this.vScreen
	// this.manFocusToModel.InitFocus()
	this.content = ""
	// this.begin = TPUT_CIVIS
	// this.end = TPUT_CUP_INIT
	this.allFocusableSections = append(this.allFocusableSections, &this.Section)
	this.isDialog = true
	se = &this.Section
	///
	se.InitSection(1, tRect{top, bottom, left, right}, "FindSection", &this.vScreen)
	se.layers[0].name = "FindSection.Find"
	se.layers[0].AddWindow(&Find.w, &tRect{top, bottom, left, right})
	///
	this.Section.PrepareSection()
	for _, l := range this.allLayers {
		l.InitLayer()
	}
	for _, w := range this.allWindows {
		w.PrepareWindow()
	}
	Find.w.nRows--
	//
	this.Section.SetActiveWindow(&Find.w)
	this.UpdateFocusButtonList()
	this.PrepareDialogButtons()
	// this.activeSection = &this.Section
// trace.End() // t //
}

func (this *vSettingsDialog) InitScreen() {
	this.SetOuter()
}

// .
func (this *vSettingsDialog) PrepareScreen() {
	AllMaps.allScreens = append(AllMaps.allScreens, &this.vScreen)
	var top uint16 = 7
	var left uint16 = dialogLeft
	var bottom uint16 = top + 35
	var middle uint16 = left + 35
	var right uint16 = left + 120
	var se *vSection
	// this.manFocusToModel = new(uManScreenFocusToModel)
	this.name = "SettingsDialog"
// trace.Begin_n(this, "PrepareScreen") // t //
	this.focusKeys.pScreen = &this.vScreen
	// this.manFocusToModel.InitFocus()
	this.content = ""
	// this.begin = TPUT_CIVIS
	// this.end = TPUT_CUP_INIT
	this.allFocusableSections = append(this.allFocusableSections, &this.Section)
	this.isDialog = true
	se = &this.Section
	///
	se.InitSection(5, tRect{top, bottom, left, right}, "SettingsSection", &this.vScreen)

	se.layers[0].name = "SettingsSection.Help"
	se.layers[0].AddWindow(&Help.w, &tRect{top, bottom, left, middle})
	se.layers[0].AddWindow(&Help.w2, &tRect{top, bottom, middle, right})
	se.layers[0].isMonoModel = true

	se.layers[1].name = "SettingsSection.SetM"
	se.layers[1].AddWindow(&SetM.w, &tRect{top, bottom, left, right})

	se.layers[2].name = "SettingsSection.Equalizer"
	se.layers[2].AddWindow(&Equalizer.w, &tRect{top, bottom, left, middle})
	se.layers[2].AddWindow(&Equalizer.m2.w, &tRect{top, bottom, middle, right})
	se.layers[2].isMonoModel = true

	se.layers[3].name = "SettingsSection.Colors"
	se.layers[3].AddWindow(&Colors.w, &tRect{top, bottom, left, left + 35})
	se.layers[3].AddWindow(&Colors.m2.w, &tRect{top, bottom, left + 35, right})
	se.layers[3].caption = "Colors"
	se.layers[3].isMonoModel = true

	se.layers[4].name = "SettingsSection.Report"
	se.layers[4].AddWindow(&Report.w, &tRect{top, bottom, left, right})

	this.Section.PrepareSection()
	///
	for _, l := range this.allLayers {
		l.InitLayer()
	}
	for _, w := range this.allWindows {
		w.PrepareWindow()
	}
	// for _, w := range this.allWindows {
	// 	w.InitWindow()
	// }
	//
	this.Section.SetActiveWindow(&SetM.w)
	this.UpdateFocusButtonList()
	this.PrepareDialogButtons()
	// this.activeSection = &this.Section
// trace.End() // t //
}

// func (this *vHelpDialog) InitScreen() {
// 	this.SetOuter()
// }

// // .
// func (this *vHelpDialog) PrepareScreen() {
// 	AllMaps.allScreens = append(AllMaps.allScreens, &this.vScreen)
// 	var top uint16 = 7
// 	var left uint16 = dialogLeft
// 	var bottom uint16 = top + 35
// 	var middle uint16 = left + 30
// 	var right uint16 = left + 120
// 	var se *vSection
// 	// this.manFocusToModel = new(uManScreenFocusToModel)
// 	this.name = "HelpDialog"
// 	trace.Begin_n(this, "InitScreen")
// 	this.focusKeys.pScreen = &this.vScreen
// 	// this.manFocusToModel.InitFocus()
// 	this.content = ""
// 	// this.begin = TPUT_CIVIS
// 	// this.end = TPUT_CUP_INIT
// 	this.allFocusableSections = append(this.allFocusableSections, &this.Section)
// 	this.isDialog = true
// 	se = &this.Section
// 	se.InitSection(3, tRect{top, bottom, left, left + 120}, "HelpSection", &this.vScreen)
// 	se.layers[0].name = "HelpSection.Help"
// 	se.layers[0].AddWindow(&Help.w, &tRect{top, bottom, left, middle})
// 	se.layers[0].AddWindow(&Help.w2, &tRect{top, bottom, middle, right})
// 	se.layers[0].isMonoModel = true
// 	se.layers[1].name = "HelpSection.Bios"
// 	se.layers[1].AddWindow(&Bios.w, &tRect{top, bottom, left, middle})
// 	se.layers[1].AddWindow(&Bios.w2, &tRect{top, bottom, middle, right})
// 	se.layers[1].isMonoModel = true
// 	se.layers[2].name = "HelpSection.Report"
// 	se.layers[2].AddWindow(&Report.w, &tRect{top, bottom, left, right})
// 	this.Section.PrepareSection()
// 	///
// 	for _, l := range this.allLayers {
// 		l.InitLayer()
// 	}
// 	for _, w := range this.allWindows {
// 		w.PrepareWindow()
// 	}
// 	this.Section.SetActiveWindow(&Help.w)
// 	// this.activeSection = &this.Section
// 	this.PrepareDialogButtons()
// 	this.UpdateFocusButtonList()
// 	trace.End()
// }

func (this *vAbsDialogScreen) DragDialog() {
	var re tRect
	var pRect *tRect = &this.Section.tRect
	var initRow, initCol, diffRow, diffCol int16
	var initTop = pRect.marginTop
	var initBottom = pRect.marginBottom
	var initLeft = pRect.marginLeft
	var initRight = pRect.marginRight
	initRow = int16(Input.iRow)
	initCol = int16(Input.iCol)
// trace.BeginAdd_n(this, "DragDialog", "initial cursor position: %d %d %v", Input.iRow, Input.iCol, pRect) // t //
	fmt.Print(MOUSE_TRACKING_ON)
	for Input.TrackMouse() {
		diffRow = int16(Input.iRow) - initRow
		diffCol = int16(Input.iCol) - initCol
		re = tRect{initTop + uint16(diffRow), initBottom + uint16(diffRow), initLeft + uint16(diffCol), initRight + uint16(diffCol)}
		// trace.Print("Drag mode, cursor position: %d %d, diff %d %d %v", Input.iRow, Input.iCol, diffRow, diffCol, re)
		fmt.Print(styleNormal.JustTheColors() + TPUT_CLEAR)
		Term.MakeSectionFrame(&re, &styleSelButton_SelFrame, &Frames.Empty)
		Term.Flush()
		// Term.AddStringWithPos(0, 0, Input.TrackMouse())
		// Term.PrintScreen()
	}
// trace.Print("Drag mode end, cursor position: %d %d", Input.iRow, Input.iCol) // t //
	fmt.Print(MOUSE_TRACKING_OFF)
	pRect.marginTop += uint16(diffRow)
	pRect.marginBottom += uint16(diffRow)
	pRect.marginLeft += uint16(diffCol)
	pRect.marginRight += uint16(diffCol)
	this.Section.tRect = re
	if this.activeWindow != nil {
		this.activeWindow.tRect = re
	}
	this.Section.PrepareSection()
	Term.MakeSectionFrame(&re, &styleSelButton_SelFrame, &Frames.Empty)
	Term.SlurpContents(&this.Section.frameBold)

	MainScreen.LeftSection.BringSectionToFront()
	MainScreen.pRightSection.BringSectionToFront()
	MainScreen.activeSection.GetActiveWindow().BringWindowToFront(true)
	MainScreen.pRightSection.GetActiveWindow().BringWindowToFront(true)
	this.Section.BringSectionToFront()
	this.Section.GetActiveWindow().BringWindowToFront(true)
// trace.End() // t //
}

func (this *vScreen) ActivateModel(model *mBasic) {
// trace.BeginEndAdd_n(this, "ActivateModel", "%s", model.name) // t //
	this.activeWindow = &model.w
	model.SetFocusToModel()
}

// Searches in all section, in all windows of the active layer.
func (this *vScreen) WhichWindow() *vWindow {
// trace.BeginSilent_n(this, "WhichWindow") // t //
	iRow := Input.iRow
	iCol := Input.iCol
// trace.AddToPrint("iRow: %d, iCol: %d", iRow, iCol) // t //
	for _, se := range this.allFocusableSections {
		if se.GetActiveWindow() != nil {
// trace.AddToPrint("section %s, active layer: %s", se.name, se.GetActiveWindow().name) // t //
			for _, w := range se.GetActiveLayer().windows {
				if IsInsideRect(&w.tRect, iRow, iCol) {
// trace.ReturnAdd("found window %s", w.GetName()) // t //
					this.activeWindow = w
					return w
				} else {
					// trace.Print("  (window %s %v)", w.name, w.tRect)
				}
			}
		} else {
			// trace.Print("section %s, active layer: none", se.name)
		}
	}
// trace.EndAdd("no window found") // t //
	return nil
}

// .
func (this *vScreen) HighlightPoint(iRow, iCol uint16) {
	var s = styleMargins.EchoStyleFormatOnly("*", false)
	Term.AddStringWithPos(iRow, iCol, s)
	Term.Flush()
}

// .
func (this *vScreen) HighlightMargins() {
	var s string
	var le uint16
	// ast := styleMargins.EchoStyleOnlyColors("*")
	for _, se := range this.allSections {
		s = styleMargins.EchoStyleFormatOnly(fmt.Sprintf("← %d %d %s.s", se.marginTop, se.marginLeft, se.name), true)
		Term.AddStringWithPos(se.marginTop, se.marginLeft, s)
		s = fmt.Sprintf("%s.s %d %d →", se.name, se.marginBottom, se.marginRight)
		le = RuneLen(&s)
		s = styleMargins.EchoStyleFormatOnly(s, true)
		Term.AddStringWithPos(se.marginBottom, se.marginRight-le-1, s)
		if w := se.GetActiveWindow(); w != nil {
			s = styleMargins.EchoStyleFormatOnly(fmt.Sprintf("↑ %d %d %s", w.marginTop, w.marginLeft, w.name), true)
			Term.AddStringWithPos(w.marginTop+1, w.marginLeft, s)
			s = fmt.Sprintf("%s %d %d ↓", w.name, w.marginBottom, w.marginRight)
			le = RuneLen(&s)
			s = styleMargins.EchoStyleFormatOnly(s, true)
			Term.AddStringWithPos(w.marginBottom-1, w.marginRight-le-1, s)
		}
	}
	Term.Flush()
}

// Searches in all sections.
// func (this *vBasicScreen) WhichSection() *vSection {
// 	iRow := Input.iRow
// 	iCol := Input.iCol
// 	trace.AddToPrint("iRow: %d, iCol: %d", iRow, iCol)
// 	for _, se := range *this.outer.GetSections_o() {
// 		if IsInsideRect(&se.tRect, iRow, iCol) {
// 			trace.BeginEndAdd("Screen.WhichSection, found section %s", se.GetName())
// 			return se
// 		}
// 	}
// 	trace.BeginEndAdd("Screen.WhichSection", "no section found")
// 	return nil
// }

func (this *vScreen) GetName() string {
	return this.name
}

func (this *vScreen) GetLongName() string {
	return this.name
}

func OpenReadlineDialog() {
	// Term.PrepareRectFrame(&tRect{10, 20, 50, 150}, nil)
}

func (this *vAbsDialogScreen) OpenDialog(focusKey string) bool {
	Term.activeCommand = nil
	// if this == &HelpDialog.vAbsDialogScreen {
	// 	HelpDialog.isOpen = true
	// } else {
	Term.oldActiveScreen = Term.activeScreen
	Term.activeScreen = &this.vScreen
	// }
	dialogMode = true
	this.focusKeys.activeKey = focusKey
// trace.BeginAdd_n(this, "OpenDialog", "focusKey: %s, MainScreen.activeWindow: %s, Term.activeScreen: %s", Input.GetKeyName(focusKey), MainScreen.activeWindow.name, Term.activeScreen.name) // t //
	rv := this.KeyboardMode()
	// if this == &HelpDialog.vAbsDialogScreen {
	// 	HelpDialog.isOpen = false
	// } else {
	Term.activeScreen = Term.oldActiveScreen
	// }
	dialogMode = false
// trace.End() // t //
	return rv
}

// .
// func (this *vMainScreen) OpenSelectDialog() {
// 	// Input.xSelect.Reset()
// 	Input.xSelect.Reset()
// 	Input.xSelect.Prepare("select", FOCUS_UNFOCUSABLE)
// 	Input.xSelect.SetCommands(true)
// 	Input.xSelect.SetOuter()
// 	Input.xSelect.w.SetEmptyLine()
// 	for i := range Artists.keys.sl {
// 		k := &Artists.keys.sl[i]
// 		Input.xSelect.AppendKeyAndLine(k.v, nil)
// 	}
// 	Term.PrepareFrame(&Screen.DialogSelect.tRect, &styleNormal, &Frames.Double)
// 	Term.SlurpContents(&Screen.DialogSelect.frameBold)
// 	Term.AddSimple(this.DialogSelect.frameBold)
// 	Input.xSelect.keys.iSelected = 0
// 	Input.xSelect.w.lines.iSelected = 0
// 	Input.xSelect.w.PrintWindowAnyway()
// 	i := Input.xSelect.CheckInput2()
// 	trace.Print("i: %d", i)
// 	// trace.Print("Screen")
// 	// this.FileContainers.TraceSection()
// 	// this.Filelists.TraceSection()
// }

// .
func (this *vScreen) ConsoleMode() *vScreen {
	Set.consoleModeOn = true
	this.consoleModeOn = true
	Term.SttyCBreakMin1()
	Term.SttyEchoOn()
	Term.Clear()
	Input.MouseOff()
	TputCNorm()
	return this
}

// .
func (this *vScreen) TUIMode() {
	Set.consoleModeOn = false
	this.consoleModeOn = false
	Term.SttyCBreakMin1()
	Term.SttyEchoOff()
	// Term.Clear()
	Input.MouseOn()
	TputCivis()
}

type vScreen struct {
	name                 string
	content              string
	beforeContent        string
	afterContent         string
	isDialog             bool
	sections             []*vSection
	fgWindows            []*vWindow
	allWindows           []*vWindow
	allLayers            []*vLayer
	allSections          []*vSection
	allFocusButtons      []*vFocusButton
	allFocusableSections []*vSection
	outer                viScreenOuter
	activeSection        *vSection
	activeWindow         *vWindow
	consoleModeOn        bool
}

type vAbsDialogScreen struct {
	vScreen
	closeButton   string
	focusKeys     uFocusKeys /* focus-to-model keys */
	dialogButtons []*vDialogButton
	Section       vSection
}

type vMainScreen struct {
	vScreen
	///
	activeWorkspace *maWorkspace
	allWorkspaces   []*maWorkspace
	mapWorkspaces   map[string]*maWorkspace
	///
	LeftSection         vSection
	pRightSection       *vSection
	CommandPanelSection vSection
	DisplaySection      vSection
	SectionCommandLine  vSection
	/// Right sections (one of them)
	Browser1Section vSection
	Browser2Section vSection
	PlayerSection   vSection
	// EditorSection   vSection
	// pManActiveFocusToModel *uManScreenFocusToModel
}

type vFindDialog struct {
	vAbsDialogScreen
}

type vSettingsDialog struct {
	vAbsDialogScreen
}

type vHelpDialog struct {
	vAbsDialogScreen
	isOpen bool
}

type viScreenOuter interface {
	GetFocusMan_o() *uFocusKeys
	MouseMode_o() (string, bool)
	BringToFront_o()
// Dump_o() /*c*/
	GetName() string
}

func (this *vAbsDialogScreen) SetOuter() {
	this.vScreen.outer = this
	this.outer = this
}

func (this *vMainScreen) SetOuter() {
	this.vScreen.outer = this
	this.outer = this
}

// func (this *vHelpDialog) SetOuter() {
// 	this.vScreen.outer = this
// 	this.outer = this
// }

