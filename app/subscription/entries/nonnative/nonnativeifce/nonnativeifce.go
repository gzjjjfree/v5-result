package nonnativeifce

import (
	"io/fs"

	"github.com/gzjjjfree/v5-result/app/subscription/entries"
)

type NonNativeConverterConstructorT func(fs fs.FS) (entries.Converter, error)

var NewNonNativeConverterConstructor NonNativeConverterConstructorT
