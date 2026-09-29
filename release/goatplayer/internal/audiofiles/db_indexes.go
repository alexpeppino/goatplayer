package audiofiles

import (
	"slices"
	"sort"
)

type GDbIndexes struct {
	Sl []uint16
}

func (this *GDbIndexes) AppendIndex(dbIndex uint16) {
	this.Sl = append(this.Sl, dbIndex)
}

func (this *GDbIndexes) InsertIndex(pos int, dbIndex uint16) {
	this.Sl = slices.Insert(this.Sl, pos, dbIndex)
}

func (this *GDbIndexes) DeleteIndex(dbIndex int) {
	this.Sl = slices.Delete(this.Sl, dbIndex, dbIndex+1)
}

func (this *GDbIndexes) ClearSlice() {
	this.Sl = nil
}

func (this *GDbIndexes) SortSlice(f f_int_int_bool) {
	sort.Slice(this.Sl, f)
}

func (this *GDbIndexes) MakeSliceByFilter(filter f_audiof_bool) {
	// trace.Begin("MakeSliceByFilter")
	this.Sl = nil
	for i := range Sl {
		if filter(&Sl[i]) {
			// trace.Print("valid %d %s", i, AF.Sl[i].title)
			this.Sl = append(this.Sl, uint16(i))
		}
	}
	// trace.End()
}

func (this *GDbIndexes) CopyFrom(that *GDbIndexes) {
	this.Sl = make([]uint16, len(that.Sl))
	copy(this.Sl, that.Sl)
}

