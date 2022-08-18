package engine

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/IOPaper/Paper/utils"
	"github.com/ugorji/go/codec"
	"os"
)

type PapersIndex struct {
	dir  string
	Size uint `msgpack:"size"`
	// TODO: maybe need replace to PaperMemIndexKey
	Hash []string `msgpack:"hash"`
	Docs []string `msgpack:"docs"`
}

type PaperMemIndexKey string
type PaperMemIndexValue string
type PapersMemIndex map[PaperMemIndexKey]PaperMemIndexValue

type PaperIndexDoc struct {
	Url  string `msgpack:"url"`
	Path string `msgpack:"path"`
}

func NewPapersIndex(dir string) *PapersIndex {
	return &PapersIndex{
		dir:  dir,
		Size: 0,
		Hash: make([]string, 0),
		Docs: make([]string, 0),
	}
}

func NewPapersIndexReader(dir string) (*PapersIndex, error) {
	f, err := os.Open(dir + "/paper_index")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	index := PapersIndex{dir: dir}
	if err = codec.NewDecoder(f, &codec.MsgpackHandle{}).Decode(&index); err != nil {
		return nil, err
	}
	if len(index.Hash) != len(index.Docs) || int(index.Size) != len(index.Hash) {
		return nil, errors.New("index corruption, out of balance")
	}
	return &index, nil
}

// ---------- PapersIndex ---------- //

func (i *PapersIndex) Write() error {
	if len(i.Hash) != len(i.Docs) || int(i.Size) != len(i.Hash) {
		return errors.New("index corruption, out of balance")
	}
	var (
		err error
		buf = &bytes.Buffer{}
	)
	if err = codec.NewEncoder(buf, &codec.MsgpackHandle{}).Encode(i); err != nil {
		return err
	}
	return utils.Write(i.dir+"/paper_index", buf, true)
}

func (i *PapersIndex) put(k, v string) {
	i.Hash = append(i.Hash, k)
	i.Docs = append(i.Docs, v)
	i.Size++
}

func (i *PapersIndex) Put(doc *PaperIndexDoc) error {
	value, err := doc.Value()
	if err != nil {
		return err
	}
	i.put(doc.Key(), value)
	return nil
}

func (i *PapersIndex) ToMap() PapersMemIndex {
	pmi := make(PapersMemIndex)
	for j := 0; j < int(i.Size); j++ {
		pmi.Put(i.Hash[j], i.Docs[j])
	}
	return pmi
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

func (v PaperMemIndexValue) Open(dir string) (*PaperIndexDoc, error) {
	f, err := os.Open(dir + string(v))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var doc PaperIndexDoc
	if err = codec.NewDecoder(f, &codec.MsgpackHandle{}).Decode(&doc); err != nil {
		return nil, err
	}
	return &doc, nil
}

// ---------- PaperIndexDoc ---------- //

func (d *PaperIndexDoc) Key() string {
	s := sha256.Sum256([]byte(d.Path))
	return hex.EncodeToString(s[:])
}

func (d *PaperIndexDoc) Value() (string, error) {
	var (
		err error
		buf = &bytes.Buffer{}
	)
	if err = codec.NewEncoder(buf, &codec.MsgpackHandle{}).Encode(d); err != nil {
		return "", err
	}
	s := sha256.New()
	buf.WriteTo(s)
	buf.Reset()
	hex.NewEncoder(buf).Write(s.Sum(nil))
	return buf.String(), nil
}
