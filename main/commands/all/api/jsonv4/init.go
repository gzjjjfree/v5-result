package jsonv4

import "github.com/gzjjjfree/v5-result/main/commands/all/api"

func init() {
	api.CmdAPI.Commands = append(api.CmdAPI.Commands,
		cmdAddInbounds,
		cmdAddOutbounds,
		cmdRemoveInbounds,
		cmdRemoveOutbounds)
}
