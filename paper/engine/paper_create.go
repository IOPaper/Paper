package engine

import (
	"errors"
	"github.com/IOPaper/Paper/paper/core"
	"github.com/IOPaper/Paper/paper/paperId"
	"github.com/IOPaper/Paper/utils"
	"log"
	"os"
	"sync"
	"time"
	"unicode/utf8"
)

type PaperCache struct {
	delete      func(key int64)
	deleteTemp  func(path string)
	write       func(path string, paper *core.Paper) error
	path        string
	contentSize uint
	paper       *core.Paper
}

func (c *PaperCache) GetPaperId() int64 {
	return c.paper.Id
}

func (c *PaperCache) GetPaperPath() string {
	return c.path
}

func (c *PaperCache) NewAttachment(name string) *os.File {
	f, err := utils.Open(c.path + name)
	if err != nil {
		return nil
	}
	return f
}

func (c *PaperCache) SetContentSize(size uint) {
	c.contentSize = size
}

func (c *PaperCache) SetTitle(title string) {
	c.paper.Title = title
}

func (c *PaperCache) SetContent(content string) error {
	log.Println(uint(utf8.RuneCountInString(content)))
	if c.contentSize != uint(utf8.RuneCountInString(content)) {
		return errors.New("invalid content size")
	}
	c.paper.Content = content
	return nil
}

func (c *PaperCache) SetTags(tag ...string) {
	c.paper.Tags = tag
}

func (c *PaperCache) SetAuthor(author string) {
	c.paper.Author = author
}

func (c *PaperCache) Close() {
	c.delete(c.paper.Id)
	c.deleteTemp(c.path)
}

func (c *PaperCache) Done() error {
	defer c.Close()
	c.paper.DateCreated = time.Now()
	return c.write(c.path, c.paper)
}

const DefaultWorkerId int64 = 3

type PaperCreate struct {
	id *paperId.SnowWorker
	// map[int64]core.PaperCreateFunc
	mapping *sync.Map
}

func NewPaperCreate() *PaperCreate {
	return &PaperCreate{
		id: func() *paperId.SnowWorker {
			w, _ := paperId.NewSnowWorker(DefaultWorkerId)
			return w
		}(),
		mapping: &sync.Map{},
	}
}

func (c *PaperCreate) put(key int64, value core.PaperCreateFunc) {
	c.mapping.Store(key, value)
}

func (c *PaperCreate) get(key int64) (core.PaperCreateFunc, bool) {
	v, ok := c.mapping.Load(key)
	if !ok {
		return nil, false
	}
	return v.(core.PaperCreateFunc), true
}

func (c *PaperCreate) delete(key int64) {
	c.mapping.Delete(key)
}

func (c *PaperCreate) NewCreateFunc(path string, del func(path string), write func(path string, paper *core.Paper) error) core.PaperCreateFunc {
	p := &PaperCache{
		delete:     c.delete,
		deleteTemp: del,
		write:      write,
		path:       path,
		paper: &core.Paper{
			Id: c.id.GetId(),
		},
	}
	c.put(p.GetPaperId(), p)
	return p
}

func (c *PaperCreate) RecoverCreateFunc(paperId int64) (core.PaperCreateFunc, bool) {
	return c.get(paperId)
}
