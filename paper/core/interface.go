package core

import (
	"os"
	"reflect"
	"time"
)

type Paper struct {
	Id int64 `msgpack:"id" json:"id" toml:"id" update:""`
	// the Title of the Paper, it can be anything as long as it is a string
	Title string `msgpack:"title" json:"title" toml:"title" update:"ok"`
	// the Content of the Paper, it can be anything as long as it is a string
	Content string `msgpack:"content" json:"content" toml:"content" update:"ok"`
	// Tags for paper, allowing readers to grasp keywords about paper
	Tags []string `msgpack:"tags" json:"tags" toml:"tags" update:"ok"`
	// Attachment of Paper, different modes, different values
	Attachment []string `msgpack:"attachment" json:"attachment" toml:"attachment" update:"ok"`
	// the nickname the Author wants to publish
	Author string `msgpack:"author" json:"author" toml:"author" update:""`
	// each Paper has a unique Sign that never changes
	Sign []byte `msgpack:"sign" json:"sign" toml:"sign" update:""`

	DateCreated time.Time `msgpack:"date-created" json:"date-created" toml:"date-created" update:""`
	// the time when the Paper was last modified, if this value is not empty, it means the Paper has been modified
	DateModified time.Time `msgpack:"date-modified" json:"date-modified" toml:"date-modified" update:""`
}

type PaperExport struct {
	Id           string     `json:"paper_id,omitempty"`
	Title        string     `json:"title"`
	Content      string     `json:"content"`
	Tags         []string   `json:"tags,omitempty"`
	Attachment   []string   `json:"attachment,omitempty"`
	Author       string     `json:"author"`
	Sign         []byte     `json:"sign,omitempty"`
	DateCreate   time.Time  `json:"date_create"`
	DateModified *time.Time `json:"date_modified,omitempty"`
}

type Papers struct {
	Ids    []string
	Papers []Paper
}

type PaperCopy Paper
type PaperAction interface {
	Paper() *Paper
	Export() *PaperExport
	Revising() PaperRevising
}
type PaperBatchAction interface {
	Export() []PaperExport
}
type PaperRevising interface {
	CompareDiff(paper *Paper) error
}

// ---------- PaperCopy ---------- //

func (c *PaperCopy) CompareDiff(paper *Paper) error {
	typeof := reflect.TypeOf(*paper)
	src := reflect.ValueOf(*paper)
	dst := reflect.ValueOf(*c)
	for i := 0; i < typeof.NumField(); i++ {
		field := typeof.Field(i)
		if field.Tag.Get("update") != "ok" {
			continue
		}
		// check src value
		if src.Field(i).IsZero() || src.Field(i).Interface() == dst.Field(i).Interface() {
			continue
		}
		// src value is new
		dst.Field(i).Set(dst.Field(i).Elem())
	}
	c.DateModified = time.Now()
	return nil
}

func (c *PaperCopy) Paper() *Paper {
	return (*Paper)(c)
}

// ---------- Paper ---------- //

func (p *Paper) clone() PaperCopy {
	cp := *p
	return (PaperCopy)(cp)
}

func (p *Paper) Revising() PaperRevising {
	cp := p.clone()
	return &cp
}

func (p *Paper) Paper() *Paper {
	return p
}

func (p *Paper) Export() *PaperExport {
	return &PaperExport{
		Title:      p.Title,
		Content:    p.Content,
		Tags:       p.Tags,
		Attachment: p.Attachment,
		Author:     p.Author,
		Sign:       p.Sign,
		DateCreate: p.DateCreated,
		DateModified: func() *time.Time {
			if p.DateModified.IsZero() {
				return nil
			}
			return &p.DateModified
		}(),
	}
}

// ---------- Papers ---------- //

func (p Papers) Export() []PaperExport {
	size := len(p.Papers)
	vv := make([]PaperExport, size)
	for i := 0; i < size; i++ {
		vv[i] = PaperExport{
			Id:         p.Ids[i],
			Title:      p.Papers[i].Title,
			Content:    p.Papers[i].Content,
			Tags:       p.Papers[i].Tags,
			Attachment: p.Papers[i].Attachment,
			Author:     p.Papers[i].Author,
			Sign:       p.Papers[i].Sign,
			DateCreate: p.Papers[i].DateCreated,
			DateModified: func() *time.Time {
				if p.Papers[i].DateModified.IsZero() {
					return nil
				}
				return &p.Papers[i].DateModified
			}(),
		}
	}
	return vv
}

type PaperCreateFunc interface {
	GetPaperId() int64
	GetPaperPath() string
	NewAttachment(name string) *os.File
	SetContentSize(size uint)

	SetTitle(title string)
	SetContent(content string) error
	SetTags(tag ...string)
	SetAuthor(author string)

	Close()
	Done() error
}

type PaperFunc interface {
	// Find Generic paper find interface
	Find(url string) (PaperAction, error)
	// Write Generic Paper write interface
	// Write(paper *Paper) error

	Create(path string) (PaperCreateFunc, error)

	RecoverCreate(paperId int64, path string) (PaperCreateFunc, error)

	List(before, limit uint) (PaperBatchAction, error)

	Close() error
}
