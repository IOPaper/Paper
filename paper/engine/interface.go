package engine

import "io"

type Engine interface {
	Type() string
	GetAttachmentPath(key string) (string, error)
	GetOnePaper(key string) (*Paper, error)
	GetPaperList(before, limit int) (*PaperList, error)
	CheckPaperIndexStatus(index string) bool
	AddOnePaper(index string, verify bool, paperDTO *PaperDTO) error
	AddAttachment(key string, r io.Reader) error
}
