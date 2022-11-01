package config

import (
	"context"
	"flag"
	"fmt"

	"github.com/IOPaper/Paper/utils/values"
)

type Args map[string]string

func (a Args) Get(key string) string {
	value, ok := a[key]
	if !ok {
		return ""
	}
	return value
}

const ArgsContextKey string = "CONFIG_ARGS"

func parseArgs(ctx *context.Context) {
	configPath := flag.String("config", "./config.yaml", "config file path")
	flag.Parse()
	*ctx = context.WithValue(*ctx, ArgsContextKey, Args{
		"config": *configPath,
	})
}

func AssertArgs(ctx context.Context) (Args, error) {
	return values.ContextAssertion[Args](ctx, ArgsContextKey)
}

func PrintContextArgs(ctx context.Context) {
	value, err := AssertArgs(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(value)
}
