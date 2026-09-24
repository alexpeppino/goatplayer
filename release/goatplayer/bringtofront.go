package main

// func (this *vSection) IsActiveToScreen() bool {
// return this == this.pScreen.activeSection
// }

/*  */
func (this *vLayer) BringLayerToFront() {
// trace.Begin_n(this, "BringLayerToFront") // t //
	// this.pSection.activeLayer = this
	for i := range len(this.windows) {
		this.windows[i].PrintWindowAnyway()
	}
// trace.End() // t //
}

func (this *vSection) BringSectionToFront() {
// trace.Begin_n(this, "BringSectionToFront") // t //
	if this.GetActiveWindow() == nil {
// trace.Print("activeWindow: nil, set to layers[0].windows[0]") // t //
		this.SetActiveWindow(this.layers[0].windows[0])
	}
	if this == &MainScreen.LeftSection {
		l := MainScreen.pRightSection.GetActiveLayer()
		Term.AddSimple(MainScreen.pRightSection.frameLight)
		Term.AddSimple(l.focusButtonsEcho)

	} else if this == MainScreen.pRightSection {
		Term.AddSimple(MainScreen.LeftSection.frameLight)
		Term.AddSimple(MainScreen.LeftSection.GetActiveLayer().focusButtonsEcho)
		Term.AddSimple(MainScreen.LeftSection.GetActiveLayer().frameLight)
	}
	Term.AddSimple(this.frameBold)
	Term.AddSimple(this.GetActiveLayer().focusButtonsActiveEcho)
	Term.AddSimple(this.GetActiveLayer().frameLight)
	this.pScreen.outer.BringToFront_o()
	Term.Flush()
// trace.End() // t //
}

func (this *vSection) BringSectionToBack() {
// trace.Begin_n(this, "BringSectionToBack") // t //
	if this.GetActiveWindow() == nil {
		this.SetActiveWindow(this.layers[0].windows[0])
	}
	Term.AddSimple(this.GetActiveLayer().focusButtonsEcho)
	Term.AddSimple(this.frameLight)
	Term.Flush()
// trace.End() // t //
}

