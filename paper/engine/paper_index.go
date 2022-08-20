package engine

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/IOPaper/Paper/paper/core"
	"github.com/IOPaper/Paper/utils"
	"github.com/vmihailenco/msgpack/v5"
	"os"
)

const PapersIndexName string = "paper_index"

type PapersIndexMapping interface {
	Put(key, value string)
	Get(key string) (value PaperMemIndexValue, status bool)
	Delete(key string)
	Range(f func(key PaperMemIndexKey, value PaperMemIndexValue) error) error
}

type PaperMemIndexKey string
type PaperMemIndexValue string
type PapersMemIndex map[PaperMemIndexKey]PaperMemIndexValue
type PaperMemIndexValues []string

type PapersIndex struct {
	dir         string
	mapping     PapersIndexMapping
	tempMapping map[string]struct{}
	Size        uint `msgpack:"size"`
	// TODO: maybe need replace to PaperMemIndexKey
	Hash []string `msgpack:"hash"`
	Docs []string `msgpack:"docs"`
}

// PaperIndexDoc TODO: should do better
type PaperIndexDoc struct {
	dir  string
	Url  string `msgpack:"url"`
	Path string `msgpack:"path"`
}

func OpenPapersIndex(dir string) (*PapersIndex, error) {
	if utils.IsExist(dir + PapersIndexName) {
		return NewPapersIndexReader(dir)
	}
	return NewPapersIndex(dir), nil
}

func NewPapersIndex(dir string) *PapersIndex {
	return &PapersIndex{
		dir:         dir,
		mapping:     make(PapersMemIndex),
		tempMapping: make(map[string]struct{}),
		Size:        0,
		Hash:        make([]string, 0),
		Docs:        make([]string, 0),
	}
}

func NewPapersIndexReader(dir string) (*PapersIndex, error) {
	f, err := os.Open(dir + PapersIndexName)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	index := PapersIndex{dir: dir, tempMapping: make(map[string]struct{})}
	if err = msgpack.NewDecoder(f).Decode(&index); err != nil {
		return nil, err
	}
	if len(index.Hash) != len(index.Docs) || int(index.Size) != len(index.Hash) {
		return nil, errors.New("index corruption, out of balance")
	}
	index.mapping = index.toMap()
	return &index, nil
}

// ---------- PapersIndex ---------- //

func (i *PapersIndex) toMap() PapersMemIndex {
	pmi := make(PapersMemIndex)
	for j := 0; j < int(i.Size); j++ {
		pmi.Put(i.Hash[j], i.Docs[j])
	}
	return pmi
}

func (i *PapersIndex) put(k, v string) {
	i.Hash = append(i.Hash, k)
	i.Docs = append(i.Docs, v)
	i.mapping.Put(k, v)
	i.Size++
}

func (i *PapersIndex) Write() error {
	if len(i.Hash) != len(i.Docs) || int(i.Size) != len(i.Hash) {
		return errors.New("index corruption, out of balance")
	}
	var (
		err error
		buf = &bytes.Buffer{}
	)
	if err = msgpack.NewEncoder(buf).Encode(i); err != nil {
		return err
	}
	return utils.Write(i.dir+PapersIndexName, buf, true)
}

func (i *PapersIndex) PutElemTempMapping(path string) {
	i.tempMapping[path] = struct{}{}
}

func (i *PapersIndex) DeleteTempMappingElem(path string) {
	delete(i.tempMapping, path)
}

func (i *PapersIndex) GetTempMappingElem(path string) bool {
	_, ok := i.tempMapping[path]
	return ok
}

func (i *PapersIndex) Put(doc *PaperIndexDoc) error {
	if !i.GetTempMappingElem(doc.Path) {
		return errors.New("the paper does not exist in the waiting list")
	}
	defer i.DeleteTempMappingElem(doc.Path)
	i.put(doc.Key(), doc.Value())
	return nil
}

func (i *PapersIndex) Exist(path string) bool {
	_, ok := i.tempMapping[path]
	if ok {
		return true
	}
	_, ok = i.Find(path)
	return ok
}

func (i *PapersIndex) Mapping() PapersIndexMapping {
	return i.mapping
}

func (i *PapersIndex) Find(path string) (PaperMemIndexValue, bool) {
	s := sha256.Sum256([]byte(path))
	return i.mapping.Get(hex.EncodeToString(s[:]))
}

func (i *PapersIndex) List(before, limit int) (PaperMemIndexValues, error) {
	size := len(i.Docs)
	if size < before {
		return nil, errors.New("before is too large")
	}
	if size < before+limit || limit > 10 {
		return nil, errors.New("limit is too large")
	}
	return i.Docs[before : limit+before], nil
}

// ---------- PapersMemIndex ---------- //

func (i PapersMemIndex) Put(k, v string) {
	i[PaperMemIndexKey(k)] = PaperMemIndexValue(v)
}

func (i PapersMemIndex) Get(k string) (PaperMemIndexValue, bool) {
	v, o := i[PaperMemIndexKey(k)]
	return v, o
}

func (i PapersMemIndex) Delete(k string) {
	delete(i, PaperMemIndexKey(k))
}

func (i PapersMemIndex) Range(f func(k PaperMemIndexKey, v PaperMemIndexValue) error) error {
	var err error
	for k, v := range i {
		if err = f(k, v); err != nil {
			return err
		}
	}
	return nil
}

// ---------- PaperMemIndexValue ---------- //

// Doc
// 获取paper索引信息
func (v PaperMemIndexValue) Doc(dir string) *PaperIndexDoc {
	//f, err := os.Open(dir + string(v))
	//if err != nil {
	//	return nil, err
	//}
	//defer f.Close()
	//doc := PaperIndexDoc{dir: dir}
	//if err = msgpack.NewDecoder(f).Decode(&doc); err != nil {
	//	return nil, err
	//}
	//return &doc, nil
	return &PaperIndexDoc{
		dir:  dir,
		Path: string(v),
	}
}

// ---------- PaperIndexDoc ---------- //

func (d *PaperIndexDoc) Key() string {
	s := sha256.Sum256([]byte(d.Url))
	return hex.EncodeToString(s[:])
}

func (d *PaperIndexDoc) Value() string /*(string, error)*/ {
	//var (
	//	err error
	//	buf = &bytes.Buffer{}
	//)
	//if err = msgpack.NewEncoder(buf).Encode(d); err != nil {
	//	return "", err
	//}
	//s := sha256.New()
	//buf.WriteTo(s)
	//buf.Reset()
	//hex.NewEncoder(buf).Write(s.Sum(nil))
	//return buf.String(), nil
	return d.Path
}

func (d *PaperIndexDoc) Open() (*core.Paper, error) {
	f, err := os.Open(fmt.Sprintf("%s%s/%s", d.dir, d.Path, PaperDocName))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	p := core.Paper{}
	return &p, NewPaperDecode(f).Decode(&p)
}

// ---------- PaperMemIndexValues ---------- //

func (v *PaperMemIndexValues) Open(dir string) ([]core.Paper, error) {
	var (
		f      *os.File
		err    error
		p      core.Paper
		papers = make([]core.Paper, len(*v))
	)
	for i, doc := range *v {
		if f, err = os.Open(fmt.Sprintf("%s%s/%s", dir, doc, PaperDocName)); err != nil {
			return nil, err
		}
		if err = NewPaperDecode(f).Decode(&p); err != nil {
			return nil, err
		}
		f.Close()
		papers[i] = p
	}
	return papers, nil
}
