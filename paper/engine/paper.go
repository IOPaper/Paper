package engine

import "github.com/IOPaper/Paper/paper/core"

type Paper struct {
}

func New() core.PaperEngine {
	return &Paper{}
}

func (p *Paper) Find(path string) (*core.Paper, error) {

	return nil, nil
}

//func (p *Paper) Write(paper *Paper) {
//
//}

func (p *Paper) Create(path string) (core.PaperCreate, error) {

	return nil, nil
}
