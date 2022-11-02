package fs

import (
	"github.com/pkg/errors"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

type paperIndexSortSlice []struct {
	Index string
	Id    int64
}

func (p paperIndexSortSlice) Len() int {
	return len(p)
}

func (p paperIndexSortSlice) Less(i, j int) bool {
	return p[i].Id > p[j].Id
}

func (p paperIndexSortSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

type PaperIndexMetadata struct {
	Id         int64     `json:"id" msgpack:"id"`
	CreateDate time.Time `json:"create_date" msgpack:"create_date"`
}
type PaperIndex struct {
	__mu  sync.Mutex
	Len   int64                          `json:"len"`
	Index map[string]*PaperIndexMetadata `json:"index"`
}

func (i *PaperIndex) Get(key string) (*PaperIndexMetadata, bool) {
	value, status := i.Index[key]
	return value, status
}

func (i *PaperIndex) Set(key string, value *PaperIndexMetadata) error {
	i.__mu.Lock()
	if _, status := i.Get(key); status {
		i.__mu.Unlock()
		return errors.New("index already exists")
	}
	i.Index[key] = value
	i.__mu.Unlock()
	atomic.AddInt64(&i.Len, 1)
	return nil
}

func (i *PaperIndex) Del(key string) error {
	i.__mu.Lock()
	if _, status := i.Get(key); !status {
		i.__mu.Unlock()
		return errors.New("index does not exist")
	}
	delete(i.Index, key)
	i.__mu.Unlock()
	atomic.AddInt64(&i.Len, -1)
	return nil
}

type PaperIndexActions interface {
	Get(key string) (*PaperIndexMetadata, bool)
	Add(key string, value *PaperIndexMetadata) error
	Del(key string) error
	GetSlice(before, limit int) ([]string, error)
	Write() error
}

type paperIndex struct {
	mux    sync.Mutex
	writer WriteIndex
	slice  []string
	*PaperIndex
}

func NewPaperIndex(fApi FileApi) (PaperIndexActions, error) {
	pIndex, err := fApi.OpenIndex()
	if err != nil {
		return nil, err
	}
	// bad code :(
	var slice []string
	{
		var sortSlice paperIndexSortSlice
		for key, value := range pIndex.Index {
			sortSlice = append(sortSlice, struct {
				Index string
				Id    int64
			}{Index: key, Id: value.Id})
		}
		sort.Sort(sortSlice)
		for _, value := range sortSlice {
			slice = append(slice, value.Index)
		}
	}
	return &paperIndex{
		writer:     fApi,
		slice:      slice,
		PaperIndex: pIndex,
	}, nil
}

func (i *paperIndex) Write() error {
	return i.writer.WriteIndex(i.PaperIndex)
}

func (i *paperIndex) Add(key string, value *PaperIndexMetadata) error {
	{
		md, status := i.Get(i.slice[len(i.slice)-1])
		if !status {
			panic("index panic, unable to get metadata of last element in index list")
		}
		if md.Id > value.Id {
			return errors.New("invalid paper id")
		}
	}
	if err := i.Set(key, value); err != nil {
		return err
	}
	i.mux.Lock()
	i.slice = append(i.slice, key)
	i.mux.Unlock()
	return nil
}

func (i *paperIndex) GetSlice(_before, _limit int) ([]string, error) {
	before, limit := int64(_before), int64(_limit)
	switch {
	case i.Len < before:
		return nil, errors.New("before is too large")
	case limit > 10:
		return nil, errors.New("limit is too large")
	case i.Len < before+limit:
		limit = i.Len - before
		if limit > 10 {
			limit = 10
		}
	}
	return i.slice[before : limit+before], nil
}
