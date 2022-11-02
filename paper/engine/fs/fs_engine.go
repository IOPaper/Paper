package fs

import (
	"context"
	"github.com/IOPaper/Paper/paper/engine"
)

type Engine struct {
	FileApi
	ctx   context.Context
	root  string
	index PaperIndexActions
}

func NewEngine(ctx context.Context, root string) (engine.Engine, error) {
	fApi := NewFileApi(root)
	if err := fApi.Init(); err != nil {
		return nil, err
	}
	index, err := NewPaperIndex(fApi)
	if err != nil {
		return nil, err
	}
	return &Engine{
		FileApi: fApi,
		ctx:     ctx,
		root:    root,
		index:   index,
	}, nil
}

func (e *Engine) Type() string {
	return "fs"
}

func (e *Engine) GetOnePaper(index string) (*engine.Paper, error) {
	fiApi, err := e.OpenWithDocIndex(index)
	if err != nil {
		return nil, err
	}
	doc, err := fiApi.OpenDOC()
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func (e *Engine) GetPaperList(before, limit int) (*engine.PaperList, error) {
	indexs, err := e.index.GetSlice(before, limit)
	if err != nil {
		return nil, err
	}
	list := engine.PaperList{
		Len:    len(indexs),
		Papers: make([]engine.PaperStore, len(indexs)),
	}
	for key, index := range indexs {
		docs, er := e.OpenWithDocIndex(index)
		if er != nil {
			return nil, er
		}
		doc, der := docs.OpenDOC()
		if der != nil {
			return nil, der
		}
		list.Papers[key] = engine.PaperStore{
			Index: index,
			Paper: doc,
		}
	}
	return &list, nil
}
