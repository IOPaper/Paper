package global

import (
	"errors"
	"flag"
	"github.com/IOPaper/Paper/boot/ctl"
	"github.com/IOPaper/Paper/global/structs"
	"github.com/pelletier/go-toml/v2"
	"os"
)

type ImplementConfig struct {
	args struct {
		path string
	}
	cReader *os.File
}

func NewConfig() ctl.I {
	return &ImplementConfig{}
}

func (i *ImplementConfig) Create() (err error) {
	flag.StringVar(&i.args.path, "config", "", "config file path")
	flag.Parse()
	if i.args.path == "" {
		return errors.New("config path is empty")
	}
	if i.cReader, err = os.Open(i.args.path); err != nil {
		return
	}
	return nil
}

func (i *ImplementConfig) Destroy() error {
	return nil
}

func (i *ImplementConfig) Start() error {
	defer i.cReader.Close()
	err := toml.NewDecoder(i.cReader).Decode(&Config)
	if err != nil {
		return err
	}
	for _, option := range map[string]structs.ConfigChecker{
		"engine.log-method":    Config.Engine.LogMethod,
		"engine.log-level":     Config.Engine.LogLevel,
		"paper.index-rule":     Config.Paper.IndexRule,
		"paper.storage-format": Config.Paper.StorageFormat,
	} {
		if err = option.Check(); err != nil {
			return err
		}
	}
	return nil
}

func (i *ImplementConfig) IsAsync() bool {
	return false
}
