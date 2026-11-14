package environment

import "github.com/gzjjjfree/v5-result/common/log"

type ConnectionCapabilitySet interface {
	ConnectionLogCapabilitySet
}

type ConnectionEnvironment interface {
	ConnectionCapabilitySet
	doNotImpl()
}

type ConnectionLogCapabilitySet interface {
	RecordConnectionLog(msg log.Message)
}
