package engine

type Engine interface {
	Type() string
	GetOnePaper(key string) (*Paper, error)
}
