package main

import (
	"goatplayer/internal/K"
	"goatplayer/internal/input"
)

func (this *mBasic) CheckInput_o(b bool) string {
	// trace.Begin_n(this, "CheckInput_o {mBasic}") // t //
	screen := this.w.pLayer.pSection.pScreen
	for {
		if b {
			input.ReadInput()
		}
		b = true
		if input.InputType == K.INPUT_MOUSE {
			//	Switch to mouse-mode
			// trace.ReturnAdd("mouse mode") // t //
			return K.MOUSE_MODE
		}
		// trace.PrintWithName("switch inputkey...") // t //
		switch input.InputKey {
		case K.KEYB_TAB:
			// trace.ReturnAdd("tab key found") // t //
			// dump.TraceModelParents(this) // t //
			return this.GetScreen().outer.GetFocusMan_o().NextItem()
		case K.KEYB_SHIFT_TAB:
			// trace.ReturnAdd("shift+tab key found") // t //
			// dump.TraceModelParents(this) // t //
			return this.GetScreen().outer.GetFocusMan_o().PrevItem()
		}
		// trace.PrintWithName("seek focus keys...") // t //
		if screen.IsFocusToModelKey(input.InputKey) {
			// trace.ReturnAdd("found focus key %s", input.InputKey) // t //
			return input.InputKey
		}
		// trace.PrintWithName("seek commands...") // t //
		for i := range this.kCmm {
			cmm := &this.kCmm[i]
			if input.InputKey == cmm.key {
				// trace.Print("found command: %s", cmm.desc) // t //
				cmm.f()
				break
			}
		}
	}
}

func (this *mFind) CheckInput_o(doReadInputAtBeginning bool) string {
	// trace.Begin_n(this, "CheckInput OUTER") // t //
	var isCommandFound bool
	this.PrintTitle(this.title)
	// focusMan := this.w.pLayer.pSection.pScreen.focusMan
	for {
		if doReadInputAtBeginning {
			input.ReadInput()
		}
		doReadInputAtBeginning = true
		if input.InputType == K.INPUT_MOUSE {
			//	Switch to mouse-mode
			// trace.ReturnAdd("mouse mode") // t //
			return K.MOUSE_MODE
		}
		switch input.InputKey {
		case K.KEYB_ESC:
			// trace.ReturnAdd("ESC") // t //
			return "ESC"
		case K.KEYB_ENTER:
			// trace.ReturnAdd("ENTER") // t //
			return "ENTER"
		}
		// if focusMan.IsKey(input.InputKey) {
		// 	trace.N_Return(this, "CheckInput, is key %s", input.InputKey)
		// 	return input.InputKey
		// }
		isCommandFound = false
		for i := range this.kCmm {
			cmm := &this.kCmm[i]
			if input.InputKey == cmm.key {
				cmm.f()
				isCommandFound = true
				break
			}
		}
		if !isCommandFound {
			switch this.mode {
			case 1:
				this.UpdatePattern1()
			case 2:
				this.UpdatePattern2()
			}
		}
	}
}

func (this *mBasic) CheckInputBasic() {
	// trace.N_BeginEnd(this, "CheckInput")
	// trace.Print("input %d %s %d %d", input.InputType, input.InputKey, input.X, input.Y)
	switch input.InputType {
	case K.INPUT_KEYBOARD:
		for i := range this.kCmm {
			cmm := &this.kCmm[i]
			if input.InputKey == cmm.key {
				cmm.f()
				break
			}
		}
	}
}

func (this *mBasic) CheckMouseInput_o() {
	// trace.Begin_n(this, "CheckMouseInput") // t //
	for i := range this.mCmm {
		mcmm := &this.mCmm[i]
		if input.InputKey == mcmm.key {
			mcmm.f()
			break
		}
	}
	// trace.End() // t //
}

// func (this *mPlaylist) CheckMouseInput() {
// 	trace.N_Begin(this, "CheckMouseInput")
// 	if this.mBasic.CheckMouseInputBasic() {
// 		trace.N_Return(this, "CheckMouseInput")
// 		return
// 	}

// 	trace.N_End(this, "CheckMouseInput")
// }

// func (this *mColors) CheckMouseInput() {
// }

// func (this *mColors) CheckInput(b bool) string {
// 	// trace.N_BeginEnd(this, "CheckInput")
// 	for {
// 		if b {
// 			input.ReadInput()
// 		}
// 		b = true
// 		if input.InputType == INPUT_MOUSE {
// 			//	Switch to mouse-mode
// 			return K.MOUSE_MODE
// 		}
// 		this.mBasic.CheckInputBasic()
// 		switch input.InputType {
// 		case INPUT_KEYBOARD:
// 			if Focus.IsKey(input.InputKey) {
// 				return input.InputKey
// 			}
// 			switch input.InputKey {
// 			}
// 		}
// 	}
// }

func (this *mDialog) CheckInput_o(bool) string { return "" }

func (this *mDialog) CheckInput2() uint16 {
	// trace.N_BeginEnd(this, "CheckInput")
	for {
		input.ReadInput()
		if input.InputType == K.INPUT_MOUSE {
			//	Switch to mouse-mode
			return K.UNSET
		}
		this.mBasic.CheckInputBasic()
		switch input.InputType {
		case K.INPUT_KEYBOARD:
			switch input.InputKey {
			case "\n":
				return this.keys.SelectedI
			case K.KEYB_ESC:
				return K.UNSET
			}
		}
	}
}

// func (this *mHelp) CheckMouseInput() {
// 	this.mBasic.CheckMouseInputBasic()
// }

func (this *mHelp) CheckInput_o(b bool) string {
	// trace.Begin_n(this, "CheckInput_o {mHelp}") // t //
	for {
		if b {
			input.ReadInput()
		}
		b = true
		if input.InputType == K.INPUT_MOUSE {
			//	Switch to mouse-mode
			// trace.Return() // t //
			return K.MOUSE_MODE
		}
		// this.mBasic.CheckInputBasic()
		if this.w.pLayer.pSection.pScreen.IsFocusToModelKey(input.InputKey) {
			// trace.Return() // t //
			return input.InputKey
		}
		switch input.InputKey {
		case K.KEYB_ESC:
			// trace.Return() // t //
			return "ESC"
		}
		for i := range this.kCmm {
			cmm := &this.kCmm[i]
			if input.InputKey == cmm.key {
				cmm.f()
				break
			}
		}
	}
}

// func (this *mSettings) CheckMouseInput() {
// 	this.mBasic.CheckMouseInputBasic()
// }

// b: read input at start
func (this *mSettings) CheckInput_o(b bool) string {
	// trace.Begin_n(this, "CheckInput") // t //
	for {
		if b {
			input.ReadInput()
		}
		b = true
		if input.InputType == K.INPUT_MOUSE {
			//	Switch to mouse-mode
			// trace.Return() // t //
			return K.MOUSE_MODE
		}
		if this.w.pLayer.pSection.pScreen.IsFocusToModelKey(input.InputKey) {
			// trace.Return() // t //
			return input.InputKey
		}
		switch input.InputKey {
		case K.KEYB_ESC:
			// trace.Return() // t //
			return "ESC"
		}
		for i := range this.kCmm {
			cmm := &this.kCmm[i]
			if input.InputKey == cmm.key {
				cmm.f()
				break
			}
		}
	}
}

func (this *mColorsPicker) CheckInput_o(b bool) string {
	return Colors.CheckInput_o(b)
}

func (this *mColors) CheckInput_o(b bool) string {
	// trace.Begin_n(this, "CheckInput") // t //
	for {
		if b {
			input.ReadInput()
		}
		b = true
		if input.InputType == K.INPUT_MOUSE {
			//	Switch to mouse-mode
			// trace.Return() // t //
			return K.MOUSE_MODE
		}
		if this.w.pLayer.pSection.pScreen.IsFocusToModelKey(input.InputKey) {
			// trace.Return() // t //
			return input.InputKey
		}
		switch input.InputKey {
		case K.KEYB_ESC:
			// trace.Return() // t //
			return "ESC"
		}
		for i := range this.kCmm {
			cmm := &this.kCmm[i]
			if input.InputKey == cmm.key {
				cmm.f()
				break
			}
		}
	}
}

func (this *mReport) CheckInput_o(b bool) string {
	// trace.Begin_n(this, "CheckInput_o {mReport}") // t //
	for {
		if b {
			input.ReadInput()
		}
		b = true
		if input.InputType == K.INPUT_MOUSE {
			//	Switch to mouse-mode
			// trace.Return() // t //
			return K.MOUSE_MODE
		}
		if this.w.pLayer.pSection.pScreen.IsFocusToModelKey(input.InputKey) {
			// trace.Return() // t //
			return input.InputKey
		}
		switch input.InputKey {
		case K.KEYB_ESC:
			// trace.Return() // t //
			return "ESC"
		}
		for i := range this.kCmm {
			cmm := &this.kCmm[i]
			if input.InputKey == cmm.key {
				cmm.f()
				break
			}
		}
	}
}

// func (this *mEqualizer) CheckMouseInput() {
// 	this.mBasic.CheckMouseInputBasic()
// }

// b: read input at start
func (this *mEqualizer) CheckInput_o(b bool) string {
	// trace.Begin_n(this, "CheckInput_o {mEqualizer}") // t //
	screen := this.w.pLayer.pSection.pScreen
	for {
		if b {
			input.ReadInput()
		}
		b = true
		if input.InputType == K.INPUT_MOUSE {
			//	Switch to mouse-mode
			// trace.ReturnAdd("K.MOUSE_MODE") // t //
			return K.MOUSE_MODE
		}
		if screen.IsFocusToModelKey(input.InputKey) {
			// trace.Return() // t //
			return input.InputKey
		}
		switch input.InputKey {
		case K.KEYB_ESC:
			// trace.ReturnAdd("ESC") // t //
			return "ESC"
		}
		for i := range this.kCmm {
			cmm := &this.kCmm[i]
			if input.InputKey == cmm.key {
				// trace.Print("found command: %s", cmm.desc) // t //
				cmm.f()
				break
			}
		}
	}
}

func (this *mEqualizer2) CheckInput_o(b bool) string {
	return Equalizer.CheckInput_o(b)
}

//////	mFileProperties

// func (this *mFileProperties) CheckMouseInput() {
// 	if this.mBasic.CheckMouseInputBasic() {
// 		return
// 	}
// }

// func (this *mFileProperties) CheckInput(b bool) string {
// 	// trace.N_BeginEnd(this, "CheckInput")
// 	for {
// 		if b {
// 			input.ReadInput()
// 		}
// 		b = true
// 		if input.InputType == INPUT_MOUSE {
// 			//	Switch to mouse-mode
// 			return K.MOUSE_MODE
// 		}
// 		this.mBasic.CheckInputBasic()
// 		switch input.InputType {
// 		case INPUT_KEYBOARD:
// 			if Focus.IsKey(input.InputKey) {
// 				return input.InputKey
// 			}
// 			switch input.InputKey {
// 			// case this.cmm[CMM_CONTAINER_LOAD_TO_OPENPL].key:
// 			// 	this.LoadTo(PLView1)
// 			// case this.cmm[CMM_CONTAINER_LOAD_TO_ACTIVEPL].key:
// 			// 	this.LoadTo(PLView2)
// 			// case this.cmm[CMM_FILEPROP_EDIT_SELECTED].key:
// 			// 	this.EditSelected()
// 			}
// 		case INPUT_MOUSE:
// 		}
// 	}
// }

//////	mContainer

// func (this *mContainer) CheckMouseInput() {
// 	if this.mBasic.CheckMouseInputBasic() {
// 		return
// 	}
// 	switch input.InputKey {
// 	case K.MOUSE_B1_DOWN:
// 		this.w.SelectLineByTRow(input.RowI)
// 	case K.MOUSE_B3_DOWN:
// 		iCol := input.ColI
// 		if iCol <= this.w.marginLeft || iCol >= this.w.marginRight {
// 		} else {
// 			this.ActivateFilterByTRow(input.RowI)
// 		}
// 	}
// }

// func (this *mContainer) CheckInput(b bool) string {
// 	// trace.N_BeginEnd(this, "CheckInput")
// 	focusMan := this.w.pLayer.pSection.pScreen.focusMan
// 	for {
// 		if b {
// 			input.ReadInput()
// 		}
// 		b = true
// 		if input.InputType == INPUT_MOUSE {
// 			//	Switch to mouse-mode
// 			return K.MOUSE_MODE
// 		}
// 		if focusMan.IsKey(input.InputKey) {
// 			return input.InputKey
// 		}
// 		for i := range this.cmm {
// 			cmm := &this.cmm[i]
// 			if input.InputKey == cmm.key {
// 				cmm.f()
// 				break
// 			}
// 		}
// 	}
// }

//////

// func (this *mReport) CheckMouseInput() {
// 	if this.mBasic.CheckMouseInputBasic() {
// 		return
// 	}

// }

// func (this *mReport) CheckInput(b bool) string {
// 	// trace.N_BeginEnd(this, "CheckInput")
// 	for {
// 		if b {
// 			input.ReadInput()
// 		}
// 		b = true
// 		if input.InputType == INPUT_MOUSE {
// 			//	Switch to mouse-mode
// 			return K.MOUSE_MODE
// 		}
// 		this.mBasic.CheckInputBasic()
// 		switch input.InputType {
// 		case INPUT_KEYBOARD:
// 			if Focus.IsKey(input.InputKey) {
// 				return input.InputKey
// 			}
// 		}
// 	}
// }

//////

// func (this *mSettings) CheckInput2() {
// 	// trace.N_BeginEnd(this, "CheckInput")
// 	for {
// 		input.ReadInput()
// 		this.mBasic.CheckInputBasic()
// 		switch input.InputType {
// 		case INPUT_MOUSE:
// 			if Screen.WhichSection() == &Screen.HelpSection {
// 				this.CheckMouseInput()
// 			} else {
// 				trace.Print("CheckInput2 out of window")
// 			}
// 		case INPUT_KEYBOARD:
// 			if Focus.IsKey(input.InputKey) {
// 				return
// 			}
// 			switch input.InputKey {
// 			case KEYPAD_ARROW_RIGHT:
// 				this.ModifySelected(true)
// 			case K.KEYB_ESC:
// 				return
// 			case KEYPAD_ARROW_LEFT:
// 				this.ModifySelected(false)
// 			}
// 		}
// 	}
// }

//////////////////////////////////////h2

// func (this *mDirs) CheckMouseInput() {
// 	if this.mBasic.CheckMouseInputBasic() {
// 		return
// 	}
// }

// func (this *mDirs) CheckInput(b bool) string {
// 	trace.N_BeginEnd(this, "CheckInput")
// 	focusMan := this.w.pLayer.pSection.pScreen.focusMan
// 	for {
// 		if b {
// 			input.ReadInput()
// 		}
// 		b = true
// 		if input.InputType == INPUT_MOUSE {
// 			//	Switch to mouse-mode
// 			return K.MOUSE_MODE
// 		}
// 		if focusMan.IsKey(input.InputKey) {
// 			return input.InputKey
// 		}
// 		for i := range this.cmm {
// 			cmm := &this.cmm[i]
// 			if input.InputKey == cmm.key {
// 				cmm.f()
// 				break
// 			}
// 		}
// 	}
// }
