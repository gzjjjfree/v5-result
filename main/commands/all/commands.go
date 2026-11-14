package all

import (
	"github.com/gzjjjfree/v5-result/main/commands/all/api"
	"github.com/gzjjjfree/v5-result/main/commands/all/tls"
	"github.com/gzjjjfree/v5-result/main/commands/base"
)

//go:generate go run v2ray.com/core/common/errors/errorgen

func init() {
	base.RootCommand.Commands = append(
		base.RootCommand.Commands,
		api.CmdAPI,
		cmdLove,
		tls.CmdTLS,
		cmdUUID,
		cmdVerify,

		// documents
		docFormat,
		docMerge,
	)
}
