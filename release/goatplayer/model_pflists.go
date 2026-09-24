package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"sort"
	"strings"
)

func (this *mPFListsCon) ActivateFilterByTRow(uint16) {
	this.LoadToPlayer()
}

func (this *mPFListsCon) LoadToPlayer() {
	Play.LoadList(this.keys.GetSelected().v)
	Play.isEditor = true
	Play.UpdateDisplay()
// trace.Print("Play.isEditor: %v", Play.isEditor) // t //
	Play.MakeIndexesFromPFList()
	Play.LoadAndPrepare(nil)
	Play.list = "List: " + this.keys.GetSelected().v
	Play.SetStatusLine_o()
	Play.w.PrintWindowAnyway()
	Play.SelectKeyAndLine(0)
	Play.PlaySelected()

	// var sl []string
	// var listname = PlayerWsp.PFLists.keys.GetSelected().v
	// for i := range PFListsData.sl {
	// 	line := &PFListsData.sl[i]
	// 	if line.listname == listname {
	// 		sl = append(sl, line.filename)
	// 	}
	// }
	// Play.dbIndexes.sl = nil
	// for _, filename := range sl {
	// 	Play.dbIndexes.sl = append(Play.dbIndexes.sl, AF.mapFilenameIndex[filename])
	// }
	// Play.LoadAndPrepare(nil)
	// Play.w.PrintWindowAnyway()
}

/* Loads data into the model and prepares lines */
func (this *mPFListsCon) LoadListnames() {
	if len(PFListsData.sl) == 0 {
		return
	}
	this.Reset()
	for _, listname := range PFListsData.listnames {
		this.AppendKey(listname)
		this.w.AppendLine("  "+listname, false, nil)
	}
	//	Empty line
	this.w.SetEmptyLine()
	this.AfterLoad()
}

func (this *mPFListsCon) DeleteSelectedList() {
// trace.Begin_n(this, "DeleteSelectedList") // t //
	PFListsData.DeleteList(this.keys.GetSelected().v)
	this.ReloadAndPrint()
// trace.End() // t //
}

func (this *mPFListsCon) ReloadAndPrint() {
	PFListsData.PrepareListnames()
	this.LoadListnames()
	if this.w.IsInForeground() {
		this.w.PrintWindowAnyway()
	}
}

func (this *mPFListsCon) PrepareLines() {}

func (this *aPFLists) Reset() {
	this.sl = nil
	this.mapNamesN = make(map[string]int)
	this.listnames = nil
}

/* Appends a line to this.sl */
func (this *aPFLists) AppendLine(line *tPFLine) {
	this.sl = append(this.sl, *line)
}

/* Prepares this.mapNamesN and this.listnames */
func (this *aPFLists) PrepareListnames() {
// trace.Begin_n(this, "PrepareListnames") // t //
	this.listnames = nil
	this.mapNamesN = make(map[string]int)
	for i := range this.sl {
		line := &this.sl[i]
		this.mapNamesN[line.listname]++
	}
	for name := range this.mapNamesN {
		this.listnames = append(this.listnames, name)
	}
	sort.Strings(this.listnames)
// trace.Print("listnames: %#v", this.listnames) // t //
// trace.End() // t //
}

func (this *aPFLists) DeleteLineFromList(listname string, iTrack uint16) {
// trace.BeginAdd_n(this, "DeleteLineFromList", "listname: %s, iTrack: %d", listname, iTrack) // t //
	_, i := this.GetLine(listname, iTrack)
	if i != -1 {
		this.sl = slices.Delete(this.sl, i, i+1)
// trace.Print("line %i has been deleted", i) // t //
	}
// trace.End() // t //
}

func (this *aPFLists) InitPFList(name string) {
	this.name = name
}

/* Appends a line at the bottom of the list */
func (this *aPFLists) AppendLineToList(line *tPFLine) {
// trace.Begin_n(this, "AppendLineToList") // t //
	count := this.GetLinesCount(line.listname)
	line.iTrack = count
	this.AppendLine(line)
// trace.End() // t //
}

/* Get number of list lines */
func (this *aPFLists) GetLinesCount(listname string) uint16 {
// trace.Begin_n(this, "GetLinesCount") // t //
	var count uint16
	for i := range this.sl {
		line := &this.sl[i]
		if line.listname == listname {
			count++
		}
	}
// trace.ReturnAdd("count: %d", count) // t //
	return count
}

func (this *aPFLists) GetLine(listname string, iTrack uint16) (*tPFLine, int) {
// trace.Begin_n(this, "GetLine") // t //
	for i := range this.sl {
		line := &this.sl[i]
		if line.listname == listname && line.iTrack == iTrack {
// trace.ReturnAdd("line found, i: %d", i) // t //
			return line, i
		}
	}
// trace.ReturnAdd("line not found") // t //
	return nil, -1
}

/*
Reads the PFLists from the file.
this.rows, this.mapNamesN => updated
*/
func (this *aPFLists) LoadDataFromFile() {
// trace.BeginSilent_n(this, "LoadDataFromFile") // t //
	this.Reset()
	fd, _ := os.Open(FILE_PFLISTS)
	defer fd.Close()
	scanner := bufio.NewScanner(fd)
	var n int
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), "\t")
		this.AppendLine(&tPFLine{
			fields[0],
			utils.str2uint(fields[1]),
			fields[2],
			utils.str2uint32(fields[3]),
			utils.str2uint32(fields[4]),
			utils.str2uint(fields[5])})
		n++
		// trace.Print(scanner.Text())
		this.mapNamesN[fields[0]]++
	}
	for name := range this.mapNamesN {
		this.listnames = append(this.listnames, name)
	}
	sort.Strings(this.listnames)
// trace.AddToPrint("listnames: %v", this.listnames) // t //
	scanner.Err()
	// trace.Print("%+v \n %+v", this.rows, this.mapNamesN)
	// for s, i := range this.namesMapN {
	// 	// trace.Print("namesMapN  %s %d", s, i)
	// }
	// if err := scanner.Err(); err != nil {
	//     log.Fatal(err)
	// }
// trace.End() // t //
}

/*
Deletes listname from the pflists.
this.sl => updated
*/
func (this *aPFLists) DeleteList(listname string) {
// trace.BeginAdd_n(this, "DeleteList", "listname: '%s'", listname) // t //
// trace.Print("sl len: %d ", len(this.sl)) // t //
	list := make([]tPFLine, len(this.sl))
	copy(list, this.sl)
	this.sl = nil
	for i := range list {
		row := &list[i]
		if row.listname != listname {
			this.sl = append(this.sl, *row)
		}
	}
// trace.Print("sl len: %d ", len(this.sl)) // t //
// trace.End() // t //
}

const (
	PFLISTS_TMP    = "__TMP__"
	PFLISTS_PLAYER = "__PLAYER__"
)

/*
Writes the PFLists to the file.
FILE_PFLISTS => modified
*/
func (this *aPFLists) SaveDataToFile() {
// trace.Begin_n(this, "SaveDataToFile") // t //
	/// Updates the activepl tmplist
	this.DeleteList(PFLISTS_TMP)
	this.DeleteList(PFLISTS_PLAYER)
	if !Play.isEmpty {
		for i := range Play.keys.sl {
			k := &Play.keys.sl[i]
			this.AppendLine(&tPFLine{PFLISTS_PLAYER, uint16(i), k.v, 0, 0, 0})
		}
	}
	fd, _ := os.OpenFile(FILE_PFLISTS, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0664)
	defer fd.Close()

	for i := range this.sl {
		row := &this.sl[i]
		fmt.Fprintf(fd, "%s\t%d\t%s\t%d\t%d\t%d\n",
			row.listname,
			row.iTrack,
			row.filename,
			row.start,
			row.end,
			row.speed)
	}
// trace.EndAdd("written %d lines, %d lists", len(this.sl), len(this.mapNamesN)) // t //
	// this.Read()
}

func (this *aPFLists) GetName() string     { return this.name }
func (this *aPFLists) GetLongName() string { return this.name }

type mPFListsCon struct {
	mContainer
}

type aPFLists struct {
	name           string
	sl             []tPFLine
	mapNamesN      map[string]int
	listnames      []string
	activeListname string
}

type tPFLine struct {
	//	Track fields
	listname string // 0
	iTrack   uint16 // 1
	filename string // 2
	start    uint32 // 3
	end      uint32 // 4
	speed    uint16 // 5
}

