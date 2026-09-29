package main

import (
	"fmt"
	"goatplayer/internal/K"
	AF "goatplayer/internal/audiofiles"
	"goatplayer/internal/echo"
	"goatplayer/internal/mpv"
	"goatplayer/internal/small/dura"
	"goatplayer/internal/small/utils"
	"goatplayer/internal/style"
	"os/exec"
	"time"
)

type vMusicDisplay struct {
	vBasicView
	caption    string
	frameLight string
}

var MusicDisplay vMusicDisplay

func (this *vMusicDisplay) Print() {
	iActive := Play.keys.GetActive()
	if iActive == nil {
		Play.keys.TraceActive("")
		return
	}
	st := iActive.V
	audiof := AF.GetAudiofByFilename(st)
	Display.Part1[0] = "Title:   " + audiof.Title
	Display.Part1[1] = "Artist:  " + audiof.Artist
	Display.Part1[2] = "Album:   " + audiof.Album
	Display.AddPartToPrint(1)
	echo.Flush()
	Term.SetTermTitle(fmt.Sprintf("GoatPlayer  ::  %s  ::  %s", audiof.Title, audiof.Artist))
}

var ti2 uint16

const MPV_SOCKET string = "--input-ipc-server=/tmp/mpvsocket"
const MPV_NO_TERM string = "--no-terminal"
const MPV_NO_AUDIO_DISPLAY string = "--no-audio-display"

func PlayFile(filename string) {
	var start, end, speed, pitch string
	if ui_pl.startPoint != 0 {
		start = fmt.Sprintf("--start=%.2f", float32(ui_pl.startPoint)/100.0)
		ui_pl.startPoint = 0
	}
	if Play.isEditor {
		line := &Play.data.sl[Play.keys.ActiveI]
		if line.start != 0 {
			start = fmt.Sprintf("--start=%.2f", float32(line.start)/100)
		}
		if line.end != 0 && line.end > line.start {
			end = fmt.Sprintf("--end=%.2f", float32(line.end)/100)
		}
		if line.speed != 0 {
			speed = fmt.Sprintf("--speed=%.2f", float32(line.speed)/100)
		}
	}
	if speed == "" {
		speed = "--speed=" + ui_pl.Speed
	}
	if Set.PlaySecs.Pos != 0 {
		end = "--end=" + Set.PlaySecs.GetValue()
	}
	if !Set.AudioPitchCorrection.Value {
		pitch = "--audio-pitch-correction=no"
	}

	var parameters []string
	parameters = append(parameters,
		filename,
		// "--speed="+ui_pl.Speed,
		"--volume="+ui_pl.Volume,
		start,
		end,
		speed,
		Equalizer.GetEqusetString(K.UNSET, true),
		MPV_NO_TERM,
		pitch,
		// "--idle=yes",
		MPV_NO_AUDIO_DISPLAY,
		MPV_SOCKET)
	// for i := range parameters { /*c*/
	// trace.Player("-- %2d  %s", i, parameters[i]) // t //
	// } /*c*/

	exec.Command("mpv", parameters...).Output()
	// _, err := exec.Command("mpv",
	// 	name,
	// 	"--speed="+ui_pl.Speed,
	// 	"--volume="+ui_pl.Volume,
	// 	start, end,
	// 	Equalizer.GetEqusetString(K.UNSET, true),
	// 	MPV_NO_TERM,
	// 	pitch,
	// 	// "--idle=yes",
	// 	MPV_NO_AUDIO_DISPLAY,
	// 	MPV_SOCKET).Output()

	//!
	// if err != nil {
	// 	log.Print(err)
	// 	dump.TraceAndReportError("PlayFile/exec.Command mpv", false)
	// }

	// output := string(out[:])
	// fmt.Println("output: ", output)
}

// func GetFileDuration(filename string) uint32 {
// 	mutex1.Lock()
// 	defer mutex1.Unlock()
// 	return mapFilenameDuration[filename]
// }

// Returns a hund
func PrintFileDuration(filename string) uint32 {
	u1 := AF.GetAudiofByFilename(filename).Duration
	// form := dura.FormatSeconds(u1/100, 3)
	// s := "█  Time:    "
	// trace.Print("file duration: -%s- %f %d -%s-", ss, ff, u1, form)
	// top := MainScreen.DisplaySection.marginTop
	// left := MainScreen.DisplaySection.marginLeft
	// echo.AddStringWithPos(top+4, left, styleMusicPlayer.EchoStyle(s, 78))
	// echo.Flush()
	return u1
}

// Player co-routine
func PlayerTime_go(sig chan uint16) {
	// pl_ft.filetimeIsWaiting = true
	for {
		<-sig
		// trace.Player("filetime start") // t //
		pl_ft.filetimeIsWaiting = false
		pl_ft.stop = false
		var percent uint32 = 0
		// var oldPercent uint32 = 1
		for {
			///	Stop playing ?
			if pl_ft.stop {
				break
			}
			audiofTime := mpv.GetFileTime()
			if pl_ft.duration != 0 {
				percent = 100 * audiofTime / pl_ft.duration
			} else {
				percent = 0
			}
			if percent > 100 {
				percent = 100
			}
			// trace.Player("filetime %d %d", ti, x)
			if !Set.consoleModeOn {
				var str1, str2 string
				var perc1, perc2 uint32
				var firstCol1 uint16 = 0
				var firstCol2 uint16
				var segLen uint16

				Term.SaveCursor()
				echo.AddSimple(K.DC_RESET)
				// echo.AddStringWithPos(59, 10, styleMusicPlayer.EchoStyle(dura.FormatSeconds(ti/100, 3), 8))
				// top := MainScreen.SectionMusicPlayer.marginTop
				// left := MainScreen.SectionMusicPlayer.marginLeft

				str1, perc1 = dura.FormatPartialAndTotalBoth(audiofTime/100, pl_ft.duration/100, "Time: ")
				str2, perc2 = dura.FormatPartialAndTotalBoth((Play.precDuration+audiofTime)/100, Play.totDuration/100, "List: ")

				var totLen = Display.marginRight + 1
				segLen = (totLen - 2) / 2
				var empLen = totLen - 2*segLen
				perc1 = perc1*uint32(segLen)/100 + 1
				perc2 = perc2*uint32(segLen)/100 + 1
				firstCol2 = firstCol1 + segLen + empLen

				utils.AdjustWidth(&str1, segLen)
				utils.AdjustWidth(&str2, segLen)
				str1a := style.ProgressBarInv.EchoStyle(str1[:perc1], 0)
				str1b := style.ProgressBar.EchoStyle(str1[perc1:], 0)
				Term.ProgressBar_top = Term.nRows - 3
				Term.ProgressBar_left = firstCol1
				Term.ProgressBar_right = firstCol1 + segLen
				echo.AddStringWithPos(Term.nRows-3, firstCol1, str1a+str1b)
				str2a := style.ProgressBarInv.EchoStyle(str2[:perc2], 0)
				str2b := style.ProgressBar.EchoStyle(str2[perc2:], 0)
				echo.AddStringWithPos(Term.nRows-3, firstCol2, str2a+str2b)
				// ▁

				// 	var str1 string
				// 	var perc uint32
				// 	str1, _ = dura.FormatPartialAndTotal(audiofTime/100, pl_ft.duration/100)
				// 	echo.AddStringWithPos(Term.nRows-2, 0, styleMusicPlayer.EchoStyle(str1, 30))
				// 	str1, perc = dura.FormatPartialAndTotal((Play.precDuration+audiofTime)/100, Play.totDuration/100)
				// 	echo.AddStringWithPos(Term.nRows-2, 100, styleMusicPlayer.EchoStyle(str1, 30))
				// 	if percent != oldPercent {
				// 		echo.AddSimple(styleTimeBar.JustTheColors())
				// 		echo.AddStringWithPos(Term.nRows-2, 30, strings.Repeat(" ", 70))
				// 		echo.AddStringWithPos(Term.nRows-2, 30, strings.Repeat(Set.MusicBarRune.GetValue(), int(percent*70/100)))
				// 		echo.AddStringWithPos(Term.nRows-2, 130, strings.Repeat(" ", 70))
				// 		echo.AddStringWithPos(Term.nRows-2, 130, strings.Repeat(Set.MusicBarRune.GetValue(), int(perc*70/100)))
				// 		// echo.AddStringWithPos(FilePlayer.marginTop+1, 80, strings.Repeat("░", 100))
				// 		// echo.AddStringWithPos(FilePlayer.marginTop+1, 80, strings.Repeat("▒", int(x)))
				// 		oldPercent = percent
				// 	}

				echo.Flush()
				if ui_pl.isVerbose {
					// trace.Player("x: %d, ti: %d, pl_ft.duration: %d", percent, audiofTime, pl_ft.duration) // t //
				}
				Term.RestoreCursor()
			}
			time.Sleep(time.Duration(Set.SleepTime.GetUintValue()) * time.Millisecond)
		}
		// trace.Player("filetime pause") // t //
		pl_ft.filetimeIsWaiting = true
	}
}

// Goroutine.
func Player_go(sig chan uint16) {
	for {
		<-sig
		// trace.Player("start") // t //
		ui_pl.playerIsWaiting = false
		ui_pl.playlistHasChanged = false
		if pl_ft.filetimeIsWaiting {
			//	Start fileTime
			pl_ft.signal <- 1
		}
		for {
			//	Stop playing ?
			if ui_pl.stopPlaying {
				break
			}

			/// Play the audio file
			// trace.Player("before playing %d => %s", Play.keys.ActiveI, Play.keys.GetActive().V) // t //
			MusicDisplay.Print()
			Play.CalcPrecDuration()
			pl_ft.duration = PrintFileDuration(Play.keys.GetActive().V)
			PlayFile(Play.keys.GetActive().V)

			///{ Audio file has finished playing }
			/// Repeat 1
			if Set.RepeatMode.Pos == K.SET_POS_REPEAT_1 {
				// trace.Player("repeat audiofile") // t //
				continue
			}
			/// Stop playing ?
			if ui_pl.stopPlaying {
				Play.keys.ActiveI = K.UNSET
				Play.w.UpdateActiveLines()
				break
			}
			if ui_pl.playSelectedFile {
				// trace.Player("play selected/active file: %d", Play.keys.ActiveI) // t //
				ui_pl.playSelectedFile = false
				ui_pl.playlistHasChanged = false
			} else if ui_pl.playlistHasChanged {
				/// Playlist has changed
				// ActivePL.iActive++
				if Play.keys.ActiveI >= Play.keys.Len() {
					Play.keys.ActiveI = 0
				}
				// trace.Player("playlist has changed  active_i= %d", Play.keys.ActiveI) // t //
				ui_pl.playlistHasChanged = false
			} else {
				///	Playlist has not changed
				Play.keys.ActiveI++
				if Play.keys.ActiveI == Play.keys.Len() {
					///	Reached the end of list
					if ui_pl.RepeatList {
						///	Repeat the list
						Play.keys.ActiveI = 0
					} else {
						///	Stop playing
						Play.keys.ActiveI = K.UNSET
						break
					}
				}
				// trace.Player("playlist has not changed, iActive set to %d", Play.keys.ActiveI) // t //
			}
			/// Reprint the active line
			if Set.PlayerSelectActive.Value {
				Play.SelectKeyAndLine(Play.keys.ActiveI)
				if Play.w.IsInForeground() {
					Play.w.PrintWindowAnyway()
				}
			}
			Play.w.UpdateActiveLines()
		}
		/// Pause playing
		pl_ft.stop = true
		// trace.Player("pause") // t //
		ui_pl.playerIsWaiting = true
	}
}
