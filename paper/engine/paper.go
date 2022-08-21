package engine

import (
	"errors"
	"fmt"
	"github.com/IOPaper/Paper/paper/core"
	"github.com/IOPaper/Paper/utils"
)

const PaperDocName = "paper.io"

type Paper struct {
	dataDir string
	Index   *PapersIndex
	Creator *PaperCreate
}

func New(dataDir string) (core.PaperFunc, error) {
	index, err := OpenPapersIndex(dataDir)
	if err != nil {
		return nil, err
	}
	return &Paper{
		dataDir: dataDir,
		Index:   index,
		Creator: NewPaperCreate(),
	}, nil
}

func (p *Paper) Find(path string) (core.PaperAction, error) {
	paper, ok := p.Index.Find(path)
	if !ok {
		return nil, errors.New("paper not found")
	}
	pp, err := paper.Doc(p.dataDir).Open()
	if err != nil {
		return nil, err
	}
	return pp, nil
}

//func (p *Paper) Write(paper *Paper) {
//
//}

func (p *Paper) Writer(path string, paper *core.Paper) error {
	f, err := utils.Open(fmt.Sprintf("%s/%s/%s", p.dataDir, path, PaperDocName))
	if err != nil {
		return err
	}
	defer f.Close()
	p.Index.Put(&PaperIndexDoc{
		dir:  p.dataDir,
		Url:  path,
		Path: path,
	})
	return NewPaperEncode(f).Encode(paper)
}

func (p *Paper) Create(path string) (core.PaperCreateFunc, error) {
	if ok := p.Index.Exist(path); ok {
		return nil, errors.New("paper is occupied")
	}
	p.Index.PutElemTempMapping(path)
	return p.Creator.NewCreateFunc(path, p.Index.DeleteTempMappingElem, p.Writer), nil
}

func (p *Paper) RecoverCreate(paperId int64, path string) (core.PaperCreateFunc, error) {
	if !p.Index.GetTempMappingElem(path) {
		return nil, errors.New("paper not created")
	}
	pp, ok := p.Creator.RecoverCreateFunc(paperId)
	if !ok {
		return nil, errors.New("paper not found")
	}
	if pp.GetPaperPath() != path {
		pp.Close()
		return nil, errors.New("invalid paper path")
	}
	return pp, nil
}

func (p *Paper) List(before, limit uint) (core.PaperBatchAction, error) {
	list, err := p.Index.List(int(before), int(limit))
	if err != nil {
		return nil, err
	}
	papers, err := list.Open(p.dataDir)
	if err != nil {
		return nil, err
	}
	return papers, nil
}

func (p *Paper) Close() error {
	// save memory mapping
	err := p.Index.Write()
	if err != nil {
		return err
	}
	return nil
}
