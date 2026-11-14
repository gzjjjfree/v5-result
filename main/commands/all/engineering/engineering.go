package engineering

import "github.com/gzjjjfree/v5-result/main/commands/base"

//go:generate go run github.com/gzjjjfree/v5-result/common/errors/errorgen

var cmdEngineering = &base.Command{
	UsageLine: "{{.Exec}} engineering",
	Commands: []*base.Command{
		cmdConvertPb,
		cmdReversePb,
		cmdNonNativeLinkExtract,
		cmdNonNativeLinkExec,
		cmdSubscriptionEntriesExtract,
		cmdEncodeDataURL,
	},
}

func init() {
	base.RegisterCommand(cmdEngineering)
}

func AddCommand(cmd *base.Command) {
	cmdEngineering.Commands = append(cmdEngineering.Commands, cmd)
}
