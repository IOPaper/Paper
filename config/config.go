package config

import (
	"context"
	"io"
	"os"

	"github.com/IOPaper/Paper/utils/values"
	"gopkg.in/yaml.v2"
)

const ContextKey string = "CONFIG"

type Http struct {
	Addr     string       `yaml:"addr"`
	LogLevel HttpLogLevel `yaml:"log-level"`
}

type PaperEngineFs struct {
	Repo string `yaml:"repo"`
}

type PaperEngine struct {
	Engine PaperEngineType `yaml:"engine"`
	Fs     PaperEngineFs   `yaml:"fs,omitempty"`
}

type Config struct {
	Http        `yaml:"http"`
	PaperEngine `yaml:"paper-engine"`
}

func ParseYAML[T any](r io.Reader, v *T) error {
	return yaml.NewDecoder(r).Decode(v)
}

func WithYAMLConfig(ctx *context.Context) error {
	args, err := AssertArgs(*ctx)
	if err != nil {
		return err
	}
	f, err := os.Open(args.Get("config"))
	if err != nil {
		return err
	}
	defer f.Close()
	v := Config{}
	if err = ParseYAML[Config](f, &v); err != nil {
		return err
	}
	*ctx = context.WithValue(*ctx, ContextKey, &v)
	return nil
}

func Assert(ctx context.Context) (*Config, error) {
	return values.ContextAssertion[*Config](ctx, ContextKey)
}
