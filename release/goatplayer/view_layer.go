package main

// Adds a table (tab) to this platoon. If rect is nil, use the the section rect.
func (this *vLayer) AddWindow(window *vWindow, rect *tRect) {
// trace.BeginSilent_n(this, "AddWindow") // t //
// trace.AddToPrint("%s, rect: %+v, ", window.name, rect) // t //
	screen := this.pSection.pScreen
	// model := w.pModel
	this.windows = append(this.windows, window)
	window.pLayer = this
	if rect != nil {
		window.tRect = *rect
	} else {
		window.tRect = this.pSection.tRect
	}
	// this.pSection.pScreen.focusMan.InitFocus()
	screen.allWindows = append(screen.allWindows, window)
	/// Focus
	// trace.PrintV(w, screen)
	// screen.manFocusToModel.items = append(screen.manFocusToModel.items, uFocusToModelItem{model.focusKey, model.name})
	// screen.manFocusToModel.mapKeyToName[model.focusKey] = model.name
	// screen.manFocusToModel.mapKeyToModel[model.focusKey] = model
	// trace.N_BeginEnd(this, "AddTable %s, len: %d", table.tableName, len(this.tables))
// trace.End() // t //
}

// Get this platoon's name.
func (this *vLayer) GetName() string     { return this.name }
func (this *vLayer) GetLongName() string { return this.name }

func (this *vLayer) InitLayer() {
	if this.IsMultiWindow() {
		Term.MakeLayerFrame(this, &styleButton_Frame, &Frames.Empty).SlurpContents(&this.frameLight)
		Term.MakeLayerFrame(this, &styleSelButton_SelFrame, &Frames.Empty).SlurpContents(&this.frameBold)
	}
}

// func (this *vLayer) ActivateLayer() {
// 	trace.N_Begin(this, "ActivateLayer")
// 	// trace.Print("setActive %s (%s)", plt.name, activePlt.name)
// 	this.pSection.activeLayer = this
// 	var actSection = this.pSection
// 	var actScreen = this.pSection.pScreen
// 	if actSection == &MainScreen.SectionLeft {
// 		Term.AddSimple(actSection.frameBold)
// 		Term.AddSimple(actSection.activeLayer.sectionFocusButtonsActive)
// 		Term.AddSimple(actSection.activeLayer.frameBold)
// 		Term.AddSimple(MainScreen.pSectionRight.frameLight)
// 		Term.AddSimple(MainScreen.pSectionRight.activeLayer.sectionFocusButtons)
// 	} else if actSection == MainScreen.pSectionRight {
// 		Term.AddSimple(actSection.frameBold)
// 		Term.AddSimple(actSection.activeLayer.sectionFocusButtonsActive)
// 		Term.AddSimple(actSection.activeLayer.frameBold)
// 		Term.AddSimple(MainScreen.SectionLeft.frameLight)
// 		Term.AddSimple(MainScreen.SectionLeft.activeLayer.sectionFocusButtons)
// 	}
// 	if actScreen.isDialog {
// 		// Term.AddSimple(activeBrowser.View.pSection.frameLight)
// 		// Term.AddSimple(activeBrowser.View.pSection.activeLayer.sectionFocusButtons)
// 		Term.AddSimple(MainScreen.SectionLeft.frameLight)
// 		Term.AddSimple(MainScreen.SectionLeft.activeLayer.sectionFocusButtons)
// 		Term.AddSimple(actSection.frameBold)
// 		Term.AddSimple(actSection.activeLayer.sectionFocusButtonsActive)
// 		Term.AddSimple(this.frameLight)
// 		Term.AddSimple(this.pSection.closeButton)
// 		// Term.AddSimple(actSection.frameBold)
// 		// Term.AddSimple(actSection.activeLayer.sectionFocusButtonsActive)
// 		Term.AddSimple(actSection.activeLayer.frameBold)
// 	}
// 	Term.Flush()
// 	// this.pSection.PrintFocusButtons()
// 	trace.N_End(this, "ActivateLayer")
// }

// func (this *vSection) UpdateFrame(active bool) {
// 	trace.N_Begin(this, "UpdateFrame: %v", active)
// 	this.PrintFocusButtons(active)
// 	trace.N_End(this, "UpdateFrame")
// }

func (this *vLayer) IsMultiWindow() bool {
	return len(this.windows) != 1
}

func (this *vLayer) IsMultiTabs() bool {
	if len(this.windows) == 1 {
		// trace.N_BeginEnd(this, "IsMultiTabs false")
		return false
	}
	if this.isMonoModel {
		// trace.N_BeginEnd(this, "IsMultiTabs false")
		return false
	}
	// trace.N_BeginEnd(this, "IsMultiTabs true")
	return true
}

func (this *vLayer) DumpLayer() {
// trace.Begin_n(this, "DumpScreen {vLayer}") // t //
// trace.Print("this.name                  %v", this.name) // t //
// trace.Print("this.caption               %v", this.caption) // t //
// trace.Print("this.isMonoModel           %v", this.isMonoModel) // t //
// trace.Print("this.focusButtonsEcho      %#v", this.focusButtonsEcho) // t //
// trace.Print("this.focusButtonsActiveEcho  %#v", this.focusButtonsActiveEcho) // t //

	for i := range len(this.windows) {
		this.windows[i].DumpWindow()
	}
// trace.End() // t //
}

type vLayer struct {
	name                   string
	caption                string
	pSection               *vSection
	subSections            []*vSection
	windows                []*vWindow
	w1                     vWindow
	activeWindow           *vWindow
	frameLight             string
	frameBold              string
	isMonoModel            bool /* must be true if layer is multiwindow but not multitabs */
	focusButtonsEcho       string
	focusButtonsActiveEcho string
}

