package main

import (
	"slices"
	"sort"
)

type tDbIndexes struct {
	sl []uint16
}

func (this *tDbIndexes) AppendIndex(dbIndex uint16) {
	this.sl = append(this.sl, dbIndex)
}

func (this *tDbIndexes) InsertIndex(pos int, dbIndex uint16) {
	this.sl = slices.Insert(this.sl, pos, dbIndex)
}

func (this *tDbIndexes) DeleteIndex(dbIndex int) {
	this.sl = slices.Delete(this.sl, dbIndex, dbIndex+1)
}

func (this *tDbIndexes) ClearSlice() {
	this.sl = nil
}

func (this *tDbIndexes) SortSlice(f f_SortSlice) {
	sort.Slice(this.sl, f)
}

func (this *tDbIndexes) MakeSliceByFilter(filter f_AudiofFilter) {
	// trace.Begin("MakeSliceByFilter")
	this.sl = nil
	for i := range AF.sl {
		if filter(&AF.sl[i]) {
			// trace.Print("valid %d %s", i, AF.sl[i].title)
			this.sl = append(this.sl, uint16(i))
		}
	}
	// trace.End()
}

func (this *tDbIndexes) CopyFrom(that *tDbIndexes) {
	this.sl = make([]uint16, len(that.sl))
	copy(this.sl, that.sl)
}

