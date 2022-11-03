package engine

import "time"

type PaperAttachment struct {
	Name string `msgpack:"name" json:"name"`
	Hash string `msgpack:"hash" json:"hash"`
}

// Paper [data entity]
type Paper struct {
	Id int64 `msgpack:"id" json:"id" toml:"id"`
	// the Title of the Paper, it can be anything as long as it is a string
	Title string `msgpack:"title" json:"title" toml:"title"`
	// the Body of the Paper, it can be anything as long as it is a string
	Body string `msgpack:"body" json:"body" toml:"body"`
	// Tags for paper, allowing readers to grasp keywords about paper
	Tags []string `msgpack:"tags" json:"tags" toml:"tags"`
	// Attachment of Paper, different modes, different values
	Attachment []PaperAttachment `msgpack:"attachment" json:"attachment" toml:"attachment"`
	// the nickname the Author wants to publish
	Author string `msgpack:"author" json:"author" toml:"author"`
	// each Paper has a unique Sign that never changes
	Sign        []byte    `msgpack:"sign" json:"sign" toml:"sign"`
	Verify      bool      `msgpack:"verify" json:"verify" toml:"verify"`
	DateCreated time.Time `msgpack:"date-created" json:"date_created" toml:"date-created"`
	// the time when the Paper was last modified, if this value is not empty, it means the Paper has been modified
	DateModified time.Time `msgpack:"date-modified" json:"date_modified" toml:"date-modified"`
}

// PaperDTO [data transport object]
type PaperDTO struct {
	Title      string            `json:"title"`
	Body       string            `json:"body"`
	Tags       []string          `json:"tags"`
	Attachment []PaperAttachment `json:"attachment,omitempty"`
	Author     string            `json:"author"`
	Sign       []byte            `json:"sign"`
}

type PaperStore struct {
	Index string
	Paper *Paper
}

type PaperList struct {
	Len    int          `json:"len" msgpack:"len"`
	Papers []PaperStore `json:"papers,omitempty" msgpack:"papers,omitempty"`
}
