package structs

import "errors"

type LogMethod string
type LogLevel string
type IndexRule string
type StorageFormat string

type ConfigChecker interface {
	Check() error
}

type Config struct {
	Http struct {
		Addr string `toml:"addr"`
	} `toml:"http"`
	Engine struct {
		Repo      string    `toml:"repo"`
		LogMethod LogMethod `toml:"log-method"`
		LogLevel  LogLevel  `toml:"log-level"`
	} `toml:"engine"`
	Paper struct {
		StorageFormat         StorageFormat `toml:"storage-format"`
		IndexRule             IndexRule     `toml:"index-rule"`
		AllowAttachment       bool          `toml:"allow-attachment"`
		AllowAttachmentSuffix []string      `toml:"allow-attachment-suffix"`
		AttachmentMaxSize     int           `toml:"attachment-max-size"`
	} `toml:"paper"`
}

// ---------------- LogMethod ---------------- //

func (m LogMethod) Check() error {
	switch m {
	case "file", "console":
		return nil
	default:
		return errors.New("engine.log-method no supported")
	}
}

// ---------------- LogLevel ---------------- //

func (l LogLevel) String() string {
	return (string)(l)
}

func (l LogLevel) Check() error {
	switch l {
	case "debug", "release":
		return nil
	default:
		return errors.New("engine.log-level no supported")
	}
}

// ---------------- IndexRule ---------------- //

func (r IndexRule) Check() error {
	switch r {
	case "none", "hash":
		return nil
	default:
		return errors.New("paper.index-rule no supported")
	}
}

// ---------------- StorageFormat ---------------- //

func (f StorageFormat) Check() error {
	switch f {
	case "msgpack", "toml", "json":
		return nil
	default:
		return errors.New("paper.storage-format no supported")
	}
}
