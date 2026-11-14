package nonnative

import (
	"io/fs"

	"github.com/gzjjjfree/v5-result/app/subscription/entries"
	"github.com/gzjjjfree/v5-result/app/subscription/entries/nonnative/nonnativeifce"
	"github.com/gzjjjfree/v5-result/app/subscription/entries/outbound"
	"github.com/gzjjjfree/v5-result/app/subscription/specs"
	"github.com/gzjjjfree/v5-result/common"
)

type nonNativeConverter struct {
	matcher *DefMatcher
}

func (n *nonNativeConverter) ConvertToAbstractServerConfig(rawConfig []byte, kindHint string) (*specs.SubscriptionServerConfig, error) {
	nonNativeLink := ExtractAllValuesFromBytes(rawConfig)
	nonNativeLink.Values["_kind"] = kindHint
	result, err := n.matcher.ExecuteAll(nonNativeLink)
	if err != nil {
		return nil, newError("failed to find working converting template").Base(err)
	}
	outboundParser := outbound.NewOutboundEntriesParser()
	outboundEntries, err := outboundParser.ConvertToAbstractServerConfig(result, "")
	if err != nil {
		return nil, newError("failed to parse template output as outbound entries").Base(err)
	}
	return outboundEntries, nil
}

func NewNonNativeConverter(fs fs.FS) (entries.Converter, error) {
	matcher := NewDefMatcher()
	if fs == nil {
		err := matcher.LoadEmbeddedDefinitions()
		if err != nil {
			return nil, newError("failed to load embedded definitions").Base(err)
		}
	} else {
		err := matcher.LoadDefinitions(fs)
		if err != nil {
			return nil, newError("failed to load provided definitions").Base(err)
		}
	}
	return &nonNativeConverter{matcher: matcher}, nil
}

func init() {
	common.Must(entries.RegisterConverter("nonnative", common.Must2(NewNonNativeConverter(nil)).(entries.Converter)))
	nonnativeifce.NewNonNativeConverterConstructor = NewNonNativeConverter
}
