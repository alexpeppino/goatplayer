package main

import (
	"goatplayer/internal/echo"
)

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
		echo.AddSimple(MainScreen.pRightSection.frameLight)
		echo.AddSimple(l.focusButtonsEcho)

	} else if this == MainScreen.pRightSection {
		echo.AddSimple(MainScreen.LeftSection.frameLight)
		echo.AddSimple(MainScreen.LeftSection.GetActiveLayer().focusButtonsEcho)
		echo.AddSimple(MainScreen.LeftSection.GetActiveLayer().frameLight)
	}
	echo.AddSimple(this.frameBold)
	echo.AddSimple(this.GetActiveLayer().focusButtonsActiveEcho)
	echo.AddSimple(this.GetActiveLayer().frameLight)
	this.pScreen.outer.BringToFront_o()
	echo.Flush()
	// trace.End() // t //
}

func (this *vSection) BringSectionToBack() {
	// trace.Begin_n(this, "BringSectionToBack") // t //
	if this.GetActiveWindow() == nil {
		this.SetActiveWindow(this.layers[0].windows[0])
	}
	echo.AddSimple(this.GetActiveLayer().focusButtonsEcho)
	echo.AddSimple(this.frameLight)
	echo.Flush()
	// trace.End() // t //
}
