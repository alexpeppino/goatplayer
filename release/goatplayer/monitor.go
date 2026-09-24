package main

import (
	"fmt"
	"os"
	"time"
)

type tMonitorItem struct {
	value any
	name  string
}

type cMonitor struct {
	sl      []tMonitorItem
	content string
}

func key(k *cTableKeys) string {
	return fmt.Sprintf("%d/%d (%d)", k.iSelected, k.Len(), k.iOldSelected)
}

var Monitor cMonitor
var chMonitor chan uint16

func Monitor_go(signal chan uint16) {
	f := func(sec uint16) {
		for {
			select {
			case b2 := <-signal:
				if b2 == 0 {
// trace.GoPrint("gMonitor >> stop monitor") // t //
					return
				} else {
// trace.GoPrint("gMonitor >> sec is %d", b2) // t //
					sec = b2
				}
			default:
				Monitor.UpdateFile(sec)
				utils.Sleep(sec)
			}
		}
	}
	for {
// trace.GoPrint("gMonitor >> waiting for true") // t //
		b1 := <-signal
// trace.GoPrint("gMonitor >> signal arrived, sec is %d", b1) // t //
		if b1 > 0 {
// trace.GoPrint("gMonitor >> true signal, start monitor") // t //
			f(b1)
// trace.GoPrint("gMonitor >> f returned") // t //
		}
	}
}

type tMonitorStruct1 struct {
	Br1Name *string
}

var MonitorStruct1 = tMonitorStruct1{
	Br1Name: &BrowserWsp1.name}

func ffff() {
	m := &MonitorStruct1
	m.Br1Name = &BrowserWsp1.name
}

func (this *cMonitor) InitMonitor() {
	browser := &BrowserWsp1
	this.sl = append(this.sl, tMonitorItem{&Set.AudioPitchCorrection.Value, "AudioPitchCorrection"})
	this.sl = append(this.sl, tMonitorItem{&Set.Equalizer.Pos, "Equalizer"})
	this.sl = append(this.sl, tMonitorItem{func() string { return Set.RepeatMode.GetValue() }, "RepeatMode"})
	this.sl = append(this.sl, tMonitorItem{func() string { return key(&View2.keys) }, "ActivePL keys"})
	this.sl = append(this.sl, tMonitorItem{func() string { return key(&browser.GenresCon.keys) }, "Genres keys"})
	this.sl = append(this.sl, tMonitorItem{func() string { return key(&browser.ArtistsCon.keys) }, "Artists keys"})
	this.sl = append(this.sl, tMonitorItem{func() string { return key(&browser.AlbumsCon.keys) }, "Albums keys"})
	this.sl = append(this.sl, tMonitorItem{func() string { return key(&browser.Dirs.keys) }, "Dirs keys"})
	this.sl = append(this.sl, tMonitorItem{func() string { return Term.activeScreen.name }, "activeScreen"})
	this.sl = append(this.sl, tMonitorItem{func() string {
		if Term.activeScreen.activeWindow != nil {
			return Term.activeScreen.activeWindow.name
		} else {
			return styleMargins.EchoStyle("none", 20)
		}
	}, "active window"})
}

func (this *cMonitor) AddLine(s string) {
	this.content += s + "\n"
}

// .
func (this *cMonitor) UpdateFileWithContent() {
// trace.BeginEnd("Monitor.UpdateFileWithContent") // t //
	fd, _ := os.OpenFile("a/monitor", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0664)
	defer fd.Close()
	fmt.Fprintf(fd, "%s", this.content)
	this.content = ""
}

// Called by the goroutine, updates the file each tot secs.
func (this *cMonitor) UpdateFile(sec uint16) {
	fd, _ := os.OpenFile("a/monitor", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0664)
	defer fd.Close()
	fmt.Fprintf(fd, "Monitor - Time is %s, delay is %d\n\n", time.Now().Format("15:04:05"), sec)
	for i := range this.sl {
		item := &this.sl[i]
		fo := "%-30s %-20v ║  "
		switch v := item.value.(type) {
		case *uint16:
			fmt.Fprintf(fd, fo, item.name, *v)
		case *bool:
			fmt.Fprintf(fd, fo, item.name, *v)
		case *string:
			fmt.Fprintf(fd, fo, item.name, *v)
		case func() string:
			fmt.Fprintf(fd, fo, item.name, v())
		}
		if i%4 == 3 {
			fmt.Fprintf(fd, "\n")
		}
	}
}

