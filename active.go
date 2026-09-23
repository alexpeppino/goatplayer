package main

/// Window

func (this *vWindow) SetActiveToLayer() {
	this.pLayer.activeWindow = this
}

func (this *vWindow) SetActiveToSection() {
	this.SetActiveToLayer()
	this.pLayer.SetActiveToSection()
}

func (this *vScreen) GetActiveWindow() *vWindow {
	return this.activeSection.activeLayer.activeWindow
}

func (this *vWindow) SetActiveToScreen() {
	this.SetActiveToSection()
	this.pLayer.pSection.SetActiveToScreen()
}

func (this *vWindow) IsActiveToLayer() bool {
	return this.pLayer.activeWindow == this
}

func (this *vWindow) IsActiveToSection() bool {
	return this.IsActiveToLayer() && this.pLayer.IsActiveToSection()
}

/// Layer

func (this *vLayer) SetActiveToSection() {
	this.pSection.activeLayer = this
}

func (this *vLayer) IsActiveToSection() bool {
	return this.pSection.activeLayer == this
}

func (this *vLayer) GetActiveWindow() *vWindow {
	if this.activeWindow == nil {
		this.activeWindow = this.windows[0]
	}
	return this.activeWindow
}

/// Section

func (this *vSection) SetActiveToScreen() {
	this.pScreen.activeSection = this
}

func (this *vSection) IsActiveToScreen() bool {
	return this.pScreen.activeSection == this
}

func (this *vSection) GetActiveLayer() *vLayer {
	return this.activeLayer
}

func (this *vSection) GetActiveWindow() *vWindow {
	if this.activeLayer == nil {
		this.activeLayer = &this.layers[0]
	}
	return this.activeLayer.GetActiveWindow()
}

func (this *vSection) SetActiveWindow(that *vWindow) {
	if this.activeLayer == nil {
		this.activeLayer = that.pLayer
	}
	this.activeLayer.activeWindow = that
}

///

// Get the name of the active table (tab).
func (this *vLayer) GetActiveName() string {
	if this.activeWindow != nil {
		return this.activeWindow.pModel.name
	}
	return NONE
}

// func (this *vLayer) GetActiveWindow() *vWindow {
// 	if this.activeWindow != nil {
// 		return this.activeWindow
// 	} else {
// 		// ErrorDialog.Open("...")
// 		trace.Error("GetActiveWindow, null pointer, fixed")
// 		return this.windows[0]
// 	}
// }

// func (this *vLayer) SetActiveWindow(w *vWindow) {
// 	this.activeWindow = w
// 	// Monitor.Update...
// }

// func (this *vSection) GetActiveLayer() *vLayer {
// 	if this.activeWindow == nil {
// 		log.Print("this.activeWindow == nil")
// 		utils.TraceAndReportError("GetActiveLayer", true)
// 		// Term.PrintError("this.activeWindow == nil")
// 		return &this.layers[0]
// 	}
// 	if this.activeWindow.pLayer == nil {
// 		Term.PrintError("this.activeWindow.pLayer == nil")
// 		return &this.layers[0]
// 	}
// 	return this.activeWindow.pLayer
// }
