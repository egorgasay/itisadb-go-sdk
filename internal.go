package itisadb

import "github.com/pbnjay/memory"

type RAM struct {
	Total     uint64
	Available uint64
}

type internal struct{}

var Internal internal

func (i *internal) GetRAM() RAM {
	return RAM{
		Total:     memory.TotalMemory() / 1024 / 1024,
		Available: memory.FreeMemory() / 1024 / 1024,
	}
}
