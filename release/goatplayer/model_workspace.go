package main

func (this *mPlayerWsp) InitPlayerWsp(name, publicName string) {
	this.name = name
	// this.outer = this
	// this.mWorkspace.outer = this
	this.pFilelist = &this.Player.mAbsFilelist
	this.publicName = publicName
	MainScreen.mapWorkspaces[this.name] = &this.maWorkspace
	this.Player.pPlayerW = this
	this.Player.pWorkspace = &this.maWorkspace
	this.AudioPropCon.pWorkspace = &this.maWorkspace
	// this.Param.pWorkspace = &this.maWorkspace
	// this.PFListsCon.pWorkspace = &this.maWorkspace
}

// func (this *mEditorWsp) InitEditorWsp(name, publicName string) {
// 	this.name = name
// 	// this.outer = this
// 	// this.mWorkspace.outer = this
// 	this.publicName = publicName
// 	MainScreen.mapWorkspaces[this.name] = &this.maWorkspace
// 	this.Edit.pEditorWsp = this
// 	this.Edit.pWorkspace = &this.maWorkspace
// 	this.Param.pWorkspace = &this.maWorkspace
// 	this.PFLists.pWorkspace = &this.maWorkspace
// 	// this.FileProp.pBrowserWsp = &BrowserWsp1
// }

func (this *mBrowserWsp) InitBrowserWsp(name, publicName string) {
	this.name = name
	this.Dirs.workingDir = AF.musicDir
// trace.BeginSilent_n(this, "InitBrowserWsp") // t //
	this.publicName = publicName
	MainScreen.mapWorkspaces[this.name] = &this.maWorkspace
	this.wspDbIndexes.CopyFrom(&AF.allDbIndexes)

// trace.AddToPrint("dbIndexes len: %d", len(this.wspDbIndexes.sl)) // t //
	// trace.AddToPrint("this.dbIndexes:  %+v", this.dbIndexes)

	this.outer = this
	this.maWorkspace.outer = this
	this.pFilelist = &this.View.mAbsFilelist

	this.View.pBrowserWsp = this
	this.GenresCon.pBrowserWsp = this
	this.ArtistsCon.pBrowserWsp = this
	this.AlbumsCon.pBrowserWsp = this
	this.YearsCon.pBrowserWsp = this
	this.AlbumArtistsCon.pBrowserWsp = this
	this.ComposersCon.pBrowserWsp = this
	// this.CommentsCon.pBrowserWsp = this
	this.AudioPropCon.pBrowserWsp = this
	this.Dirs.pBrowserWsp = this

	this.View.pWorkspace = &this.maWorkspace
	this.GenresCon.pWorkspace = &this.maWorkspace
	this.ArtistsCon.pWorkspace = &this.maWorkspace
	this.AlbumsCon.pWorkspace = &this.maWorkspace
	this.YearsCon.pWorkspace = &this.maWorkspace
	this.AlbumArtistsCon.pWorkspace = &this.maWorkspace
	this.ComposersCon.pWorkspace = &this.maWorkspace
	// this.CommentsCon.pWorkspace = &this.maWorkspace
	this.AudioPropCon.pWorkspace = &this.maWorkspace
	this.Dirs.pWorkspace = &this.maWorkspace
// trace.End() // t //
}

func (this *maWorkspace) GetLongName() string { return this.name }

/* short=Wsp */
type maWorkspace struct {
	name          string
	pRightSection *vSection
	sections      []*vSection
	focusKeys     uFocusKeys /* focus-to-model keys */
	outer         miWorkspace
	publicName    string
	AudioPropCon  mAudioProperties
	pFilelist     *mAbsFilelist
}

type mBrowserWsp struct {
	maWorkspace
	View            mFilelist
	GenresCon       mContainer
	ArtistsCon      mContainer
	AlbumsCon       mContainer
	YearsCon        mContainer
	AlbumArtistsCon mContainer
	ComposersCon    mContainer
	// CommentsCon     mContainer
	Dirs         mDirs
	wspDbIndexes tDbIndexes
}

type mPlayerWsp struct {
	maWorkspace
	Player mPlaylist
	// PFListsCon mPFListsCon
	// Param      mParameters
}

// type mEditorWsp struct {
// 	maWorkspace
// 	Edit    mEditor
// 	PFLists mPFLists
// 	// PFListsCon mPFLists
// }

type miWorkspace interface {
// Dump_o() /*c*/
}

