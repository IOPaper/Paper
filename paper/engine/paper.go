package engine

import (
	"errors"
	"github.com/IOPaper/Paper/paper/core"
)

const paperDocName = "paper.io"

type Paper struct {
	dataDir string
	Index   *PapersIndex
}

func New(dataDir string) (core.PaperFunc, error) {
	index, err := OpenPapersIndex(dataDir)
	if err != nil {
		return nil, err
	}
	return &Paper{
		dataDir: dataDir,
		Index:   index,
	}, nil
}

func (p *Paper) Find(path string) (core.PaperAction, error) {
	paper, ok := p.Index.Mapping().Get(path)
	if !ok {
		return nil, errors.New("paper not found")
	}
	doc, err := paper.Doc(p.dataDir)
	if err != nil {
		return nil, err
	}
	pp, err := doc.Open()
	if err != nil {
		return nil, err
	}
	return pp, nil
}

//func (p *Paper) Write(paper *Paper) {
//
//}

func (p *Paper) Create(path string, paper *core.Paper) error {

	return nil
}

func (p *Paper) Close() error {
	// save memory mapping
	err := p.Index.Write()
	if err != nil {
		return err
	}
	return nil
}
