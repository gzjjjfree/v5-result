package plugins

import "github.com/gzjjjfree/v5-result/main/commands/base"

var Plugins []Plugin

type Plugin func(*base.Command) func() error

func RegisterPlugin(plugin Plugin) {
	Plugins = append(Plugins, plugin)
}
