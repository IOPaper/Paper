package config

type HttpLogLevel string
type PaperEngineType string

const (
	HttpLogLevelDebug   HttpLogLevel = "debug"
	HttpLogLevelRelease HttpLogLevel = "release"

	PaperEngineTypeFs     PaperEngineType = "fs"
	PaperEngineTypeSQLite PaperEngineType = "sqlite"
)

func (l HttpLogLevel) Check() bool {
	switch l {
	case HttpLogLevelDebug, HttpLogLevelRelease:
		return true
	default:
		return false
	}
}

func (l HttpLogLevel) String() string {
	return string(l)
}

func (t PaperEngineType) Check() bool {
	switch t {
	case PaperEngineTypeFs, PaperEngineTypeSQLite:
		return true
	default:
		return false
	}
}
