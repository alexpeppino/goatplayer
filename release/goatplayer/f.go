package main

import (
	"fmt"
	AF "goatplayer/internal/audiofiles"
)

/*
There are 2 kinds of "f" functions:
1.  fName  : outer method (for an embedding type) called inside basic.Name()
2.  _fName : defined for each variable and defined in this file
*/
func _f() {
	BrowserWsp1.GenresCon._fGetBrowserString = func() string {
		return fmt.Sprintf("%s music", BrowserWsp1.GenresCon.keys.GetSelected().V)
	}
	BrowserWsp1.ArtistsCon._fGetBrowserString = func() string {
		return BrowserWsp1.ArtistsCon.keys.GetSelected().V
	}
	BrowserWsp1.AlbumsCon._fGetBrowserString = func() string {
		i := BrowserWsp1.AlbumsCon.keys.SelectedI
		return fmt.Sprintf("%s %s", BrowserWsp1.AlbumsCon.keys.Sl[i].V, BrowserWsp1.AlbumsCon.sa1.Get(i))
	}
	BrowserWsp1.YearsCon._fGetBrowserString = func() string {
		return fmt.Sprintf("%s best songs", BrowserWsp1.YearsCon.keys.GetSelected().V)
	}
	BrowserWsp1.AlbumArtistsCon._fGetBrowserString = func() string {
		return BrowserWsp1.AlbumArtistsCon.keys.GetSelected().V
	}
	BrowserWsp1.ComposersCon._fGetBrowserString = func() string {
		return BrowserWsp1.ComposersCon.keys.GetSelected().V
	}
	// BrowserWsp1.CommentsCon._fGetBrowserString = func() string {
	// 	return BrowserWsp1.CommentsCon.keys.GetSelected().V
	// }

	BrowserWsp2.GenresCon._fGetBrowserString = func() string {
		return fmt.Sprintf("%s music", BrowserWsp2.GenresCon.keys.GetSelected().V)
	}
	BrowserWsp2.ArtistsCon._fGetBrowserString = func() string {
		return BrowserWsp2.ArtistsCon.keys.GetSelected().V
	}
	BrowserWsp2.AlbumsCon._fGetBrowserString = func() string {
		i := BrowserWsp2.AlbumsCon.keys.SelectedI
		return fmt.Sprintf("%s %s", BrowserWsp2.AlbumsCon.keys.Sl[i].V, BrowserWsp2.AlbumsCon.sa1.Get(i))
	}
	BrowserWsp2.YearsCon._fGetBrowserString = func() string {
		return fmt.Sprintf("%s best songs", BrowserWsp2.YearsCon.keys.GetSelected().V)
	}
	BrowserWsp2.AlbumArtistsCon._fGetBrowserString = func() string {
		return BrowserWsp2.AlbumArtistsCon.keys.GetSelected().V
	}
	BrowserWsp2.ComposersCon._fGetBrowserString = func() string {
		return BrowserWsp2.ComposersCon.keys.GetSelected().V
	}
	// BrowserWsp2.CommentsCon._fGetBrowserString = func() string {
	// 	return BrowserWsp2.CommentsCon.keys.GetSelected().V
	// }
	BrowserWsp1.GenresCon._fGetAudiofValue = AF.GetAudiofGenre
	BrowserWsp1.ArtistsCon._fGetAudiofValue = AF.GetAudiofArtist
	BrowserWsp1.AlbumsCon._fGetAudiofValue = AF.GetAudiofAlbum
	BrowserWsp1.AlbumArtistsCon._fGetAudiofValue = AF.GetAudiofAlbumArtist
	BrowserWsp1.ComposersCon._fGetAudiofValue = AF.GetAudiofComposer
	BrowserWsp1.YearsCon._fGetAudiofValue = AF.GetAudiofYear

	BrowserWsp2.GenresCon._fGetAudiofValue = AF.GetAudiofGenre
	BrowserWsp2.ArtistsCon._fGetAudiofValue = AF.GetAudiofArtist
	BrowserWsp2.AlbumsCon._fGetAudiofValue = AF.GetAudiofAlbum
	BrowserWsp2.AlbumArtistsCon._fGetAudiofValue = AF.GetAudiofAlbumArtist
	BrowserWsp2.ComposersCon._fGetAudiofValue = AF.GetAudiofComposer
	BrowserWsp2.YearsCon._fGetAudiofValue = AF.GetAudiofYear

	// EditorWsp.PFLists._fOpen = func() {
	// 	Edit.LoadList(EditorWsp.PFLists.keys.GetSelected().V)
	// 	Edit.PrepareIndexes()
	// 	Edit.LoadAndPrepare(nil)
	// 	Edit.w.PrintWindowAnyway()
	// }

	// BrowserWsp1.View._fSetFocusToModel = func() {
	// 	echo.AddStringWithPos(54, 0, "BrowserWsp1")
	// 	echo.Flush()
	// }
	// BrowserWsp2.View._fSetFocusToModel = func() {
	// 	echo.AddStringWithPos(54, 0, "BrowserWsp2")
	// 	echo.Flush()
	// }
	// PlayerWsp.ActivePL._fSetFocusToModel = func() {
	// 	echo.AddStringWithPos(54, 0, "PlayerWsp   ")
	// 	echo.Flush()
	// }
}
