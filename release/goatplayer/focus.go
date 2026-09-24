package main

func (this *uFocusKeys) GetActiveKeyName() string {
	return Input.GetKeyName(this.activeKey)
}

func (this *uFocusKeys) GetActiveModel() *mBasic {
// trace.BeginSilent_n(this, "GetActiveModel") // t //
	if model, ok := this.mapKeyToModel[this.activeKey]; ok {
// trace.ReturnAdd("found model: %s, activeKey: %s", model.name, this.activeKey) // t //
		return model
	}
// trace.ReturnAdd("no model found, activeKey: %s", this.activeKey) // t //
	return nil
}

func SetWorkspaceAndDialogFocusKeys() {
// trace.Begin("SetWorkspaceAndDialogFocusKeys") // t //
	BrowserWsp1.sections = append(BrowserWsp1.sections, &MainScreen.LeftSection, &MainScreen.Browser1Section)
	BrowserWsp2.sections = append(BrowserWsp2.sections, &MainScreen.LeftSection, &MainScreen.Browser2Section)
	PlayerWsp.sections = append(PlayerWsp.sections, &MainScreen.LeftSection, &MainScreen.PlayerSection)
	// EditorWsp.sections = append(EditorWsp.sections, &MainScreen.LeftSection, &MainScreen.EditorSection)
	SettingsDialog.sections = append(SettingsDialog.sections, &SettingsDialog.Section)
	// HelpDialog.sections = append(HelpDialog.sections, &HelpDialog.Section)
	FindDialog.sections = append(FindDialog.sections, &FindDialog.Section)

	MainScreen.allWorkspaces = append(MainScreen.allWorkspaces,
		&BrowserWsp1.maWorkspace,
		&BrowserWsp2.maWorkspace,
		&PlayerWsp.maWorkspace)
	// &EditorWsp.maWorkspace)

	AllMaps.allDialogs = append(AllMaps.allDialogs,
		&SettingsDialog.vAbsDialogScreen,
		// &HelpDialog.vAbsDialogScreen,
		&FindDialog.vAbsDialogScreen)

	for _, workspace := range MainScreen.allWorkspaces {
		// trace.Print("workspace %s", workspace.name)
		focus := &workspace.focusKeys
		focus.name = workspace.name + ".focusToModel"
		focus.mapKeyToModel = make(map[string]*mBasic)
		focus.mapKeyToName = make(map[string]string)
		focus.mapNameToModel = make(map[string]*mBasic)
		for _, section := range workspace.sections {
			// trace.Print("  section %s", section.name)
			for _, layer := range section.layers {
				// trace.Print("    layer %s", layer.name)
				for _, window := range layer.windows {
					model := window.pModel
					// trace.Print("      window %s, model %s", window.name, model.name)
					focus.items = append(focus.items, uFocusKey{model.focusKey, model.name, model})
					if model.pWorkspace == workspace {
						focus.itemsTab = append(focus.itemsTab, model.focusKey)
					}
					focus.mapKeyToName[model.focusKey] = model.name
					focus.mapKeyToModel[model.focusKey] = model
					focus.mapNameToModel[model.name] = model
				}
			}
		}
	}

	for _, dialog := range AllMaps.allDialogs {
		// trace.Print("workspace %s", workspace.name)
		focus := &dialog.focusKeys
		focus.name = dialog.name + ".focusToModel"
		focus.mapKeyToModel = make(map[string]*mBasic)
		focus.mapKeyToName = make(map[string]string)
		focus.mapNameToModel = make(map[string]*mBasic)
		for _, section := range dialog.sections {
			// trace.Print("  section %s", section.name)
			for _, layer := range section.layers {
				// trace.Print("    layer %s", layer.name)
				for _, window := range layer.windows {
					model := window.pModel
					// trace.Print("      window %s, model %s", window.name, model.name)
					focus.items = append(focus.items, uFocusKey{model.focusKey, model.name, model})
					focus.itemsTab = append(focus.itemsTab, model.focusKey)
					focus.mapKeyToName[model.focusKey] = model.name
					focus.mapKeyToModel[model.focusKey] = model
					focus.mapNameToModel[model.name] = model
				}
			}
		}
	}

// dump.TraceAllFocuses() // t //
// trace.End() // t //
}

func (this *uFocusKeys) GetName() string     { return this.name }
func (this *uFocusKeys) GetLongName() string { return this.name }

// Give the focus key, get a pointer to the model or to ActivePL if the key does not exist.
func (this *uFocusKeys) GetModel() *mBasic {
// trace.Begin_n(this, "GetModel") // t //
	model, ok := this.mapKeyToModel[this.activeKey]
	if ok {
		// trace.PrintFocus("%s.focusMan.IsKey, found key %s/%s", se.name, Input.GetKeyDesc(key), f)
// trace.ReturnAdd("found model: %s", model.name) // t //
		return model
	}
// this.TraceItems() /*c*/
// trace.ReturnAdd("no model found") // t //
	return nil
	// for _, se := range *this.outer.GetSections_o() {
	// 	focusToModel := se.focusToModel
	// 	model, ok := focusToModel.mapKeyToModel[key]
	// 	if ok {
	// 		// trace.PrintFocus("%s.focusMan.IsKey, found key %s/%s", se.name, Input.GetKeyDesc(key), f)
	// 		return model
	// 	}
	// }
	// return nil
}

// func (this *uManScreenFocusToModel) TraceActiveKey() {
// 	if s, ok := this.mapKeyToName[this.activeKey]; ok {
// 		trace.PrintFocus("focus TraceActiveKey: key %s found: %s", Input.GetKeyDesc(this.activeKey), s)
// 	} else {
// 		trace.PrintFocus("focus TraceActiveKey: key %s not found", Input.GetKeyDesc(this.activeKey))
// 	}
// }

// func (this *uSectionFocusToModel) InitFocus() {
// 	this.mapKeyToName = make(map[string]string)
// 	this.mapKeyToModel = make(map[string]*mBasic)
// }

func (this *uFocusKeys) NextItem() string {
// trace.Begin("Focus.NextItem") // t //
// trace.Print("name %s", this.name) // t //
// trace.Print("activeKey: %s", this.activeKey) // t //
// trace.Print("mapKeyToName: %+v", this.mapKeyToName) // t //
	for i, key := range this.itemsTab {
		// trace.Print("%s  %s", item.key, this.activeKey)
		if key == this.activeKey {
			if i == len(this.itemsTab)-1 {
				this.activeKey = this.itemsTab[0]
			} else {
				this.activeKey = this.itemsTab[i+1]
			}
// trace.Print("new activeKey: %s", this.activeKey) // t //
			break
		}
	}
// trace.Return() // t //
	return this.activeKey
}

func (this *uFocusKeys) PrevItem() string {
// trace.Begin("Focus.PrevItem") // t //
// trace.Print("name %s", this.name) // t //
// trace.Print("activeKey: %s", this.activeKey) // t //
// trace.Print("mapKeyToName: %+v", this.mapKeyToName) // t //
	for i, key := range this.itemsTab {
		// trace.Print("%s  %s", item.key, this.activeKey)
		if key == this.activeKey {
			if i == 0 {
				this.activeKey = this.itemsTab[len(this.itemsTab)-1]
			} else {
				this.activeKey = this.itemsTab[i-1]
			}
// trace.Print("new activeKey: %s", this.activeKey) // t //
			break
		}
	}
// trace.Return() // t //
	return this.activeKey
}

// Adds the focus legend to the Screen. The first time it is called, prepares the echo string.
func (this *uFocusKeys) AddToPrint() {
	// trace.Begin("Focus.AddToPrint")
	if this.echo == "" {
		iRow := 1
		Term.AddStringWithPos(uint16(iRow), MainScreen.Browser1Section.marginRight+2, Pen.Normal+DC_UNDERLINE+"Focus"+DC_RESET)
		iRow++
		for i, v := range this.items {
			s := Pen.StatusLine + v.key + Pen.Normal + " " + v.name
			Term.AddStringWithPos(uint16(iRow+i), MainScreen.Browser1Section.marginRight+2, s)
			// if this.activeKey == v.key {
			// 	Window.Add(uint16(iRow+i), FileContainers.marginRight+2, DC_BOLD+s+DC_RESET)
			// } else {
			// 	Window.Add(uint16(iRow+i), FileContainers.marginRight+2, s)
			// }
		}
		Term.SlurpContents(&this.echo)
	}
	Term.AddSimple(this.echo)
	// trace.End("Focus.AddToPrint")
}

func (this *vFocusButton) IsClicked(iRow, iCol uint16) bool {
	// trace.Print("IsClicked")
	return iRow >= this.marginTop && iRow <= this.marginBottom && iCol >= this.marginLeft && iCol <= this.marginRight
}

/* Focus to model */
type uFocusKeys struct {
	name           string
	items          []uFocusKey
	itemsTab       []string /* used for the TAB sequence */
	activeKey      string
	echo           string
	pScreen        *vScreen
	mapKeyToName   map[string]string
	mapKeyToModel  map[string]*mBasic
	mapNameToModel map[string]*mBasic
}

// type uSectionFocusToModel struct {
// 	items          []uFocusToModelItem
// 	mapKeyToName   map[string]string
// 	mapKeyToModel  map[string]*mBasic
// 	mapNameToModel map[string]*mBasic
// 	echo           string
// 	pScreen        *vScreen
// }

type uFocusKey struct {
	key   string
	name  string
	model *mBasic
}

/* Focus to layer */
type vFocusButton struct {
	tRect
	layer *vLayer
	i     uint16
}

