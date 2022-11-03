package fs

import (
	"context"
	"github.com/IOPaper/Paper/paper/engine"
	"github.com/IOPaper/Paper/paper/id"
	"github.com/pkg/errors"
	"time"
)

type Engine struct {
	FileApi
	ctx    context.Context
	root   string
	index  PaperIndexActions
	worker *id.Worker
}

func NewEngine(ctx context.Context, root string) (engine.Engine, error) {
	worker, err := id.NewWorker(1)
	if err != nil {
		return nil, errors.Wrap(err, "init snowflake worker failed")
	}
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
		worker:  worker,
	}, nil
}

func (e *Engine) Type() string {
	return "fs"
}

func (e *Engine) GetOnePaper(index string) (*engine.Paper, error) {
	fiApi, err := e.OpenPaperWithIndex(index)
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
		docs, er := e.OpenPaperWithIndex(index)
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

func (e *Engine) AddOnePaper(verify bool, dto *engine.PaperDTO) error {
	if e.CheckPaperIndexStatus(dto.Index) {
		return errors.New("paper index already exists")
	}
	id := e.worker.GetId()
	if err := e.AddPaper(dto.Index, &engine.Paper{
		Id:          id,
		Title:       dto.Title,
		Body:        dto.Body,
		Tags:        dto.Tags,
		Attachment:  dto.Attachment,
		Author:      dto.Author,
		Sign:        dto.Sign,
		Verify:      verify,
		DateCreated: time.Now(),
	}); err != nil {
		return err
	}
	defer e.index.Write()
	return e.index.Add(dto.Index, &PaperIndexMetadata{
		Id:         id,
		CreateDate: time.Now(),
	})
}
