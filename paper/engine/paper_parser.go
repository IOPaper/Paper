package engine

import (
	"github.com/IOPaper/Paper/global"
	jsoniter "github.com/json-iterator/go"
	"github.com/pelletier/go-toml/v2"
	"github.com/vmihailenco/msgpack/v5"
	"io"
)

type PaperEncoding interface {
	Encode(v any) error
}

type PaperDecoding interface {
	Decode(v any) error
}

type Encoder func(w io.Writer) PaperEncoding
type Decoder func(r io.Reader) PaperDecoding

func NewPaperEncode(w io.Writer) PaperEncoding {
	// TODO: using sync.Pool
	switch global.Config.Paper.StorageFormat.String() {
	case "msgpack":
		return msgpack.NewEncoder(w)
	case "toml":
		return toml.NewEncoder(w)
	case "json":
		return jsoniter.NewEncoder(w)
	}
	return nil
}

func NewPaperDecode(r io.Reader) PaperDecoding {
	switch global.Config.Paper.StorageFormat.String() {
	case "msgpack":
		return msgpack.NewDecoder(r)
	case "toml":
		return toml.NewDecoder(r)
	case "json":
		return jsoniter.NewDecoder(r)
	}
	return nil
}
