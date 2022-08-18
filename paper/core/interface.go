package core

import "time"

type PaperCreate interface {
}

type PaperEngine interface {
	// Find Generic paper find interface
	Find(url string) (*Paper, error)
	// Write Generic Paper write interface
	// Write(paper *Paper) error

	Create(path string) (PaperCreate, error)
}

type Paper struct {
	Id int64 `msgpack:"id" json:"id" toml:"id"`
	// the Title of the Paper, it can be anything as long as it is a string
	Title string `msgpack:"title" json:"title" toml:"title"`
	// the Content of the Paper, it can be anything as long as it is a string
	Content string `msgpack:"content" json:"content" toml:"content"`
	// Tags for paper, allowing readers to grasp keywords about paper
	Tags []string `msgpack:"tags" json:"tags" toml:"tags"`
	// Attachment of Paper, different modes, different values
	Attachment []string `msgpack:"attachment" json:"attachment" toml:"attachment"`
	// the nickname the Author wants to publish
	Author string `msgpack:"author" json:"author" toml:"author"`
	// each Paper has a unique Sign that never changes
	Sign []byte `msgpack:"sign" json:"sign" toml:"sign"`

	DateCreated time.Time `msgpack:"date-created" json:"date-created" toml:"date-created"`
	// the time when the Paper was last modified, if this value is not empty, it means the Paper has been modified
	DateModified time.Time `msgpack:"date-modified" json:"date-modified" toml:"date-modified"`
}

//type PaperRevising interface {
//	Title(title string)
//}
//
//func (p *Paper) Revising() PaperRevising {
//	return p
//}
