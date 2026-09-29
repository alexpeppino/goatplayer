package mpv

import (
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

func CallMpv(command string) (value string) {
	pr, pw := io.Pipe()
	defer pr.Close()
	cmd := exec.Command("socat", "-", "/tmp/mpvsocket")
	cmd.Stdin = pr
	go func() {
		defer pw.Close()
		fmt.Fprintln(pw, command)
	}()
	v, _ := cmd.Output()
	// trace.Print("CallMpv err: %v", err)

	//!
	// if err != nil {
	// 	log.Print(err)
	// 	dump.TraceAndReportError("CallMpv/socat", false)
	// }

	return string(v)
}

// Returns a value in 1/100 secs.
func GetFileTime() uint32 {
	value := CallMpv(`{ "command": ["get_property", "playback-time"] }`)
	value, _ = strings.CutSuffix(value, "\n")
	// trace.Print("file time value: %s", value)
	if value == "" {
		return 0
	}
	res := tMpvJson{}
	json.Unmarshal([]byte(value), &res)
	// trace.Print("file time res: %v", res)
	return uint32(res.Data * 100)
}

// {"data":8.567633,"request_id":0,"error":"success"}
// str="{ \"command\": [\"get_property\", \"playback-time\"] }"

func SetProperty(p, v string) {
	CallMpv(fmt.Sprintf("{ \"command\": [\"set_property\", \"%s\", \"%s\"] }", p, v))
}

func SetAudioPitchCorrection(b bool, Speed_i uint16) {
	var s string
	var f float32
	if b {
		s = "yes"
		f = 0.001
	} else {
		s = "no"
		f = -0.001
	}
// trace.PrintGreen("set audio-pitch-correction " + s) // t //
	CallMpv("set audio-pitch-correction " + s)
	/// Tilt speed to apply changes
	speed := fmt.Sprintf("%.3f", float32(Speed_i)/100+f)
// trace.PrintGreen(speed) // t //
	SetSpeed(speed)
}

// af remove @e$equ_active_freq
// af add @e$equ_active_freq:equalizer=f=${equ_freqs[equ_active_freq]}:t=q:w=1:g=$(Equ_GetActiveValue)"

func AudiofilterRemove(s string) {
// trace.Print("AudiofilterRemove: %s", s) // t //
	CallMpv("af remove " + s)
}

func AudiofilterRemoveAll() {
// trace.Print("AudiofilterRemoveAll") // t //
	CallMpv("af remove @e0,@e1,@e2,@e3,@e4,@e5,@e6,@e7,@e8,@e9")
}

func AudiofilterAdd(s string) {
// trace.Print("AudiofilterAdd: %s", s) // t //
	CallMpv("af add " + s)
}

func SetSpeed(s string) {
	CallMpv("set speed " + s)
}

func SetVolume(s string) {
	CallMpv("set volume " + s)
}

func Stop() {
// trace.Print("CallMpv: stop") // t //
	CallMpv("stop")
}

func CyclePause() {
// trace.Print("CallMpv: cycle pause") // t //
	CallMpv("cycle pause")
	// Play.isPaused = !Play.isPaused
	// Display.AddLineToPrint(1, 4, "").Flush()
}

func CycleMute() {
// trace.Print("CallMpv: cycle mute") // t //
	CallMpv("cycle mute")
	// Play.isMuted = !Play.isMuted
	// Display.AddLineToPrint(1, 4, "").Flush()
}

func SeekAbsolute(pos float32) {
	s := fmt.Sprintf("seek %.2f absolute", pos)
	CallMpv(s)
}

func SeekP5s() {
	CallMpv("seek 5")
}

func SeekM5s() {
	CallMpv("seek -5")
}

func SeekP30s() {
	CallMpv("seek 30")
}

func SeekM30s() {
	CallMpv("seek -30")
}

type tMpvJson struct {
	Data       float32 `json:"data"`
	Request_id int     `json:"request_id"`
	Error      string  `json:"error"`
}

