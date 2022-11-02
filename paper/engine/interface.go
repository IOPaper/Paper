package engine

type Engine interface {
	Type() string
	GetAttachmentPath(key string) (string, error)
	GetOnePaper(key string) (*Paper, error)
	GetPaperList(before, limit int) (*PaperList, error)
}
