package main

type cAllMaps struct {
	allSections           []*vSection
	allScreens            []*vScreen
	allDialogs            []*vAbsDialogScreen
	allBasicModels        []*mBasic
	mapBasicModels        map[string]*mBasic
	allAbstractFilelists  []*mAbsFilelist
	mapBasicFilelists     map[string]*mAbsFilelist
	mapFilelists          map[string]*mFilelist
	allContainers         []*mContainer
	mapContainers         map[string]*mContainer
	allBrowsers           []*mBrowserWsp
	allModels_in          []inAllModels
	mapAllModels_in       map[string]inAllModels
	allViews_in           []inAllViews
	allScreens_in         []inAllScreens
	mapAllScreens_in      map[string]inAllScreens
	mapStyle              map[string]*cStyle // map: style name > style reference
	allConfigurableStyles []*cStyle
	allStyleSettings      []cStyleSettings
	mapStyleSettings      map[string]*cStyleSettings
	allGetAudiof          map[string]func(*aAudiofile) string
}

var AllMaps cAllMaps

type inAllModels interface {
	PrepareModel()
	PrepareView()
	SetCommands()
	GetName() string
// DumpModel() /*c*/
}

type inAllViews interface {
	GetName() string
}

type inAllScreens interface {
	InitScreen()
	PrepareScreen()
	GetName() string
// TraceScreen() /*c*/
// DumpScreen()  /*c*/
}

// func (this *cAllMaps) InitAllModels() {
// 	for _, model := range this.allModels {
// 		model.Init()
// 	}
// }

func (this *cAllMaps) Init() {
	this.mapBasicFilelists = make(map[string]*mAbsFilelist)
	this.mapBasicModels = make(map[string]*mBasic)
	this.mapContainers = make(map[string]*mContainer)
	this.mapFilelists = make(map[string]*mFilelist)
	this.allGetAudiof = make(map[string]func(*aAudiofile) string)
	this.InitAllInterfaces()

	this.allGetAudiof["Artist"] = GetAudiofArtist
	this.allGetAudiof["Album"] = GetAudiofAlbum
	this.allGetAudiof["Genre"] = GetAudiofGenre
	this.allGetAudiof["AlbumArtist"] = GetAudiofAlbumArtist
	this.allGetAudiof["Composer"] = GetAudiofComposer
	this.allGetAudiof["Year"] = GetAudiofYear
}

// type iModel interface {
// 	Init()
// }

////// interfaces

func (this *mBasic) PrepareView() {}

func (cAllMaps) InitAllInterfaces() {
	AllMaps.allModels_in = append(AllMaps.allModels_in,
		/// BrowserWsp1
		&BrowserWsp1.View,
		&BrowserWsp1.GenresCon,
		&BrowserWsp1.ArtistsCon,
		&BrowserWsp1.AlbumsCon,
		&BrowserWsp1.YearsCon,
		&BrowserWsp1.AlbumArtistsCon,
		&BrowserWsp1.ComposersCon,
		// &BrowserWsp1.CommentsCon,
		&BrowserWsp1.AudioPropCon,
		&BrowserWsp1.Dirs,
		/// BrowserWsp2
		&BrowserWsp2.View,
		&BrowserWsp2.GenresCon,
		&BrowserWsp2.ArtistsCon,
		&BrowserWsp2.AlbumsCon,
		&BrowserWsp2.YearsCon,
		&BrowserWsp2.AlbumArtistsCon,
		&BrowserWsp2.ComposersCon,
		// &BrowserWsp2.CommentsCon,
		&BrowserWsp2.AudioPropCon,
		&BrowserWsp2.Dirs,
		/// PlayerWsp
		&PlayerWsp.Player,
		&PlayerWsp.AudioPropCon,
		// &PlayerWsp.Param,
		// &PlayerWsp.PFListsCon,

		/// EditorWsp
		// &EditorWsp.Edit,
		// &EditorWsp.Param,
		// &EditorWsp.PFLists,
		// &EditorWsp.PFListsCon,
		/// SettingsDialog
		&SetM,
		&Equalizer,
		&Equalizer.m2,
		&Colors,
		&Colors.m2,
		/// HelpDialog
		&Help,
		// &Bios,
		&Report,
		/// FindDialog
		&Find)
	///
	// &PFLists)

	AllMaps.allViews_in = append(AllMaps.allViews_in,
		/// Screens
		&MainScreen,
		&SettingsDialog,
		// &HelpDialog,
		&FindDialog,
		/// Sections
		&MainScreen.Browser1Section,
		&MainScreen.Browser2Section,
		&MainScreen.PlayerSection,
		// &MainScreen.EditorSection,
		&SettingsDialog.Section,
		// &HelpDialog.Section,
		&FindDialog.Section)

	AllMaps.allScreens_in = append(AllMaps.allScreens_in,
		&MainScreen,
		&SettingsDialog,
		// &HelpDialog,
		&FindDialog)

	AllMaps.mapAllModels_in = make(map[string]inAllModels)
	AllMaps.mapAllScreens_in = make(map[string]inAllScreens)

	for _, in := range AllMaps.allModels_in {
		AllMaps.mapAllModels_in[in.GetName()] = in
	}

	for _, in := range AllMaps.allScreens_in {
		AllMaps.mapAllScreens_in[in.GetName()] = in
	}
}

